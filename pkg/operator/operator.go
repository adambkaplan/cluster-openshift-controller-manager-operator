package operator

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/client-go/util/workqueue"

	operatorsv1 "github.com/openshift/api/operator/v1"
	operatorconfigclientv1 "github.com/openshift/cluster-openshift-controller-manager-operator/pkg/generated/clientset/versioned/typed/openshiftcontrollermanager/v1"
	operatorconfiginformerv1 "github.com/openshift/cluster-openshift-controller-manager-operator/pkg/generated/informers/externalversions/openshiftcontrollermanager/v1"
	"github.com/openshift/library-go/pkg/operator/events"
)

const (
	kubeAPIServerNamespaceName         = "openshift-kube-apiserver"
	targetNamespaceName                = "openshift-controller-manager"
	operatorNamespaceName              = "openshift-controller-manager-operator"
	userSpecifiedGlobalConfigNamespace = "openshift-config"
	workQueueKey                       = "key"
	workloadFailingCondition           = "WorkloadFailing"
)

type OpenShiftControllerManagerOperator struct {
	targetImagePullSpec  string
	operatorConfigClient operatorconfigclientv1.OpenshiftcontrollermanagerV1Interface

	kubeClient kubernetes.Interface

	// queue only ever has one item, but it has nice error handling backoff/retry semantics
	queue workqueue.RateLimitingInterface

	rateLimiter flowcontrol.RateLimiter
	recorder    events.Recorder
}

func NewOpenShiftControllerManagerOperator(
	targetImagePullSpec string,
	operatorConfigInformer operatorconfiginformerv1.OpenShiftControllerManagerOperatorConfigInformer,
	kubeInformersForOpenshiftControllerManager informers.SharedInformerFactory,
	operatorConfigClient operatorconfigclientv1.OpenshiftcontrollermanagerV1Interface,
	kubeClient kubernetes.Interface,
	recorder events.Recorder,
) *OpenShiftControllerManagerOperator {
	c := &OpenShiftControllerManagerOperator{
		targetImagePullSpec:  targetImagePullSpec,
		operatorConfigClient: operatorConfigClient,
		kubeClient:           kubeClient,
		queue:                workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "KubeApiserverOperator"),
		rateLimiter:          flowcontrol.NewTokenBucketRateLimiter(0.05 /*3 per minute*/, 4),
		recorder:             recorder,
	}

	operatorConfigInformer.Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenshiftControllerManager.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenshiftControllerManager.Core().V1().ServiceAccounts().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenshiftControllerManager.Core().V1().Services().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForOpenshiftControllerManager.Apps().V1().Deployments().Informer().AddEventHandler(c.eventHandler())

	// we only watch some namespaces
	kubeInformersForOpenshiftControllerManager.Core().V1().Namespaces().Informer().AddEventHandler(c.namespaceEventHandler())

	return c
}

func (c OpenShiftControllerManagerOperator) sync() error {
	operatorConfig, err := c.operatorConfigClient.OpenShiftControllerManagerOperatorConfigs().Get("cluster", metav1.GetOptions{})
	if err != nil {
		return err
	}
	switch operatorConfig.Spec.ManagementState {
	case operatorsv1.Unmanaged:
		return nil

	case operatorsv1.Removed:
		// TODO probably need to watch until the NS is really gone
		if err := c.kubeClient.CoreV1().Namespaces().Delete(targetNamespaceName, nil); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
		// TODO report that we are removing?
		return nil
	}

	forceRequeue, err := syncOpenShiftControllerManager_v311_00_to_latest(c, operatorConfig)
	if forceRequeue && err != nil {
		c.queue.AddRateLimited(workQueueKey)
	}

	return err
}

// Run starts the openshift-controller-manager and blocks until stopCh is closed.
func (c *OpenShiftControllerManagerOperator) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	glog.Infof("Starting OpenShiftControllerManagerOperator")
	defer glog.Infof("Shutting down OpenShiftControllerManagerOperator")

	// doesn't matter what workers say, only start one.
	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
}

func (c *OpenShiftControllerManagerOperator) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *OpenShiftControllerManagerOperator) processNextWorkItem() bool {
	dsKey, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(dsKey)

	// before we call sync, we want to wait for token.  We do this to avoid hot looping.
	c.rateLimiter.Accept()

	err := c.sync()
	if err == nil {
		c.queue.Forget(dsKey)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("%v failed with : %v", dsKey, err))
	c.queue.AddRateLimited(dsKey)

	return true
}

// eventHandler queues the operator to check spec and status
func (c *OpenShiftControllerManagerOperator) eventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { c.queue.Add(workQueueKey) },
		UpdateFunc: func(old, new interface{}) { c.queue.Add(workQueueKey) },
		DeleteFunc: func(obj interface{}) { c.queue.Add(workQueueKey) },
	}
}

// this set of namespaces will include things like logging and metrics which are used to drive
var interestingNamespaces = sets.NewString(targetNamespaceName)

func (c *OpenShiftControllerManagerOperator) namespaceEventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if !ok {
				c.queue.Add(workQueueKey)
			}
			if ns.Name == targetNamespaceName {
				c.queue.Add(workQueueKey)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			ns, ok := old.(*corev1.Namespace)
			if !ok {
				c.queue.Add(workQueueKey)
			}
			if ns.Name == targetNamespaceName {
				c.queue.Add(workQueueKey)
			}
		},
		DeleteFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
					return
				}
				ns, ok = tombstone.Obj.(*corev1.Namespace)
				if !ok {
					utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Namespace %#v", obj))
					return
				}
			}
			if ns.Name == targetNamespaceName {
				c.queue.Add(workQueueKey)
			}
		},
	}
}
