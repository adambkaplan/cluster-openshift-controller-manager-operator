package operator

import (
	"reflect"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	configv1 "github.com/openshift/api/config/v1"
	configlistersv1 "github.com/openshift/client-go/config/listers/config/v1"
	v1alpha1 "github.com/openshift/cluster-openshift-controller-manager-operator/pkg/apis/openshiftcontrollermanager/v1alpha1"
)

func TestObserveClusterConfig(t *testing.T) {
	tests := []struct {
		name   string
		cm     *corev1.ConfigMap
		expect map[string]interface{}
	}{
		{
			name: "ensure valid configmap is observed and parsed",
			cm: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "openshift-controller-manager-images",
					Namespace: operatorNamespaceName,
				},
				Data: map[string]string{
					"builderImage":  "quay.io/sample/origin-builder:v4.0",
					"deployerImage": "quay.io/sample/origin-deployer:v4.0",
				},
			},
			expect: map[string]interface{}{
				"build": map[string]interface{}{
					"imageTemplateFormat": map[string]interface{}{
						"format": "quay.io/sample/origin-builder:v4.0",
					},
				},
				"deployer": map[string]interface{}{
					"imageTemplateFormat": map[string]interface{}{
						"format": "quay.io/sample/origin-deployer:v4.0",
					},
				},
			},
		},
		{
			name: "check that extraneous configmap fields are ignored",
			cm: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "openshift-controller-manager-images",
					Namespace: operatorNamespaceName,
				},
				Data: map[string]string{
					"builderImage": "quay.io/sample/origin-builder:v4.0",
					"unknown":      "???",
				},
			},
			expect: map[string]interface{}{
				"build": map[string]interface{}{
					"imageTemplateFormat": map[string]interface{}{
						"format": "quay.io/sample/origin-builder:v4.0",
					},
				},
			},
		},
		{
			name: "expect empty result if no image data is found",
			cm: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "openshift-controller-manager-images",
					Namespace: operatorNamespaceName,
				},
				Data: map[string]string{
					"unknownField":  "quay.io/sample/origin-builder:v4.0",
					"unknownField2": "quay.io/sample/origin-deployer:v4.0",
				},
			},
			expect: map[string]interface{}{},
		},
		{
			name: "expect empty result if no configmap is found",
			cm: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "shall-not-be-found",
					Namespace: operatorNamespaceName,
				},
				Data: map[string]string{
					"builderImage":  "quay.io/sample/origin-builder:v4.0",
					"deployerImage": "quay.io/sample/origin-deployer:v4.0",
				},
			},
			expect: map[string]interface{}{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			indexer.Add(tc.cm)

			listers := Listers{
				coreOperatorsConfigMapLister: corelistersv1.NewConfigMapLister(indexer),
			}
			result, err := observeControllerManagerImagesConfig(listers, map[string]interface{}{}, &v1alpha1.OpenShiftControllerManagerOperatorConfig{})
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}

			if !reflect.DeepEqual(result, tc.expect) {
				t.Errorf("expected %v, but got %v", tc.expect, result)
			}
		})
	}
}

func TestObserveRegistryConfig(t *testing.T) {
	const (
		expectedInternalRegistryHostname = "docker-registry.openshift-image-registry.svc.cluster.local:5000"
	)

	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	imageConfig := &configv1.Image{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Status: configv1.ImageStatus{
			InternalRegistryHostname: expectedInternalRegistryHostname,
		},
	}
	indexer.Add(imageConfig)

	listers := Listers{
		imageConfigLister: configlistersv1.NewImageLister(indexer),
	}

	result, err := observeInternalRegistryHostname(listers, map[string]interface{}{}, &v1alpha1.OpenShiftControllerManagerOperatorConfig{})
	if err != nil {
		t.Error("expected err == nil")
	}
	internalRegistryHostname, _, err := unstructured.NestedString(result, "dockerPullSecret", "internalRegistryHostname")
	if err != nil {
		t.Fatal(err)
	}
	if internalRegistryHostname != expectedInternalRegistryHostname {
		t.Errorf("expected internal registry hostname: %s, got %s", expectedInternalRegistryHostname, internalRegistryHostname)
	}
}

func TestObserveBuildControllerConfig(t *testing.T) {
	memLimit, err := resource.ParseQuantity("1G")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name        string
		buildConfig *configv1.Build
		expectError bool
	}{
		{
			name: "no build config",
		},
		{
			name: "valid build config",
			buildConfig: &configv1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster",
				},
				Spec: configv1.BuildSpec{
					BuildDefaults: configv1.BuildDefaults{
						GitHTTPProxy:  "http://my-proxy",
						GitHTTPSProxy: "https://my-proxy",
						GitNoProxy:    "https://no-proxy",
						Env: []corev1.EnvVar{
							{
								Name:  "FOO",
								Value: "BAR",
							},
						},
						ImageLabels: []configv1.ImageLabel{
							{
								Name:  "build.openshift.io",
								Value: "test",
							},
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: memLimit,
							},
						},
					},
					BuildOverrides: configv1.BuildOverrides{
						ImageLabels: []configv1.ImageLabel{
							{
								Name:  "build.openshift.io",
								Value: "teset2",
							},
						},
						NodeSelector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"foo": "bar",
							},
						},
						Tolerations: []corev1.Toleration{
							{
								Key:      "somekey",
								Operator: corev1.TolerationOpExists,
								Effect:   corev1.TaintEffectNoSchedule,
							},
						},
					},
				},
			},
		},
		{
			name: "match expressions",
			buildConfig: &configv1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster",
				},
				Spec: configv1.BuildSpec{
					BuildOverrides: configv1.BuildOverrides{
						NodeSelector: metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "mylabel",
									Values:   []string{"foo", "bar"},
									Operator: metav1.LabelSelectorOpIn,
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			if test.buildConfig != nil {
				indexer.Add(test.buildConfig)
			}
			listers := Listers{
				buildConfigLister: configlistersv1.NewBuildLister(indexer),
			}
			config := map[string]interface{}{}
			operatorConfig := &v1alpha1.OpenShiftControllerManagerOperatorConfig{}
			observed, err := observeBuildControllerConfig(listers, config, operatorConfig)
			if err != nil {
				if !test.expectError {
					t.Fatalf("unexpected error observing build controller config: %v", err)
				}
			}
			if test.expectError {
				if err == nil {
					t.Error("expected error to be thrown, but was not")
				}
				if len(observed) > 0 {
					t.Error("expected returned config to be empty")
				}
				return
			}
			if test.buildConfig == nil {
				if len(observed) > 0 {
					t.Errorf("expected empty observed config, got %v", observed)
				}
				return
			}

			testNestedField(observed, test.buildConfig.Spec.BuildDefaults.GitHTTPProxy, "build.buildDefaults.gitHTTPProxy", false, t)
			testNestedField(observed, test.buildConfig.Spec.BuildDefaults.GitHTTPSProxy, "build.buildDefaults.gitHTTPSProxy", false, t)
			testNestedField(observed, test.buildConfig.Spec.BuildDefaults.GitNoProxy, "build.buildDefaults.gitNoProxy", false, t)
			testNestedField(observed, test.buildConfig.Spec.BuildDefaults.Env, "build.buildDefaults.env", false, t)
			testNestedField(observed, test.buildConfig.Spec.BuildDefaults.ImageLabels, "build.buildDefaults.imageLabels", false, t)
			testNestedField(observed, test.buildConfig.Spec.BuildOverrides.ImageLabels, "build.buildOverrides.imageLabels", false, t)
			testNestedField(observed, test.buildConfig.Spec.BuildOverrides.NodeSelector.MatchLabels, "build.buildOverrides.nodeSelector", false, t)
			testNestedField(observed, test.buildConfig.Spec.BuildOverrides.Tolerations, "build.buildOverrides.tolerations", false, t)
		})
	}
}

func testNestedField(obj map[string]interface{}, expectedVal interface{}, field string, existIfEmpty bool, t *testing.T) {
	nestedField := strings.Split(field, ".")
	switch expected := expectedVal.(type) {
	case string:
		value, found, err := unstructured.NestedString(obj, nestedField...)
		if err != nil {
			t.Fatalf("failed to read nested string %s: %v", field, err)
		}
		if expected != value {
			t.Errorf("expected field %s to be %s, got %s", field, expectedVal, value)
		}
		if existIfEmpty && !found {
			t.Errorf("expected field %s to exist, even if empty", field)
		}
	case map[string]string:
		value, found, err := unstructured.NestedStringMap(obj, nestedField...)
		if err != nil {
			t.Fatal(err)
		}
		if !equality.Semantic.DeepEqual(value, expected) {
			t.Errorf("expected %s to be %v, got %v", field, expected, value)
		}
		if existIfEmpty && !found {
			t.Errorf("expected field %s to exist, even if empty", field)
		}
	case []corev1.EnvVar:
		value, found, err := unstructured.NestedSlice(obj, nestedField...)
		if err != nil {
			t.Fatal(err)
		}
		rawExpected, err := ConvertJSON(expected)
		if err != nil {
			t.Fatalf("unable to test field %s: %v", field, err)
		}
		if !equality.Semantic.DeepEqual(value, rawExpected) {
			t.Errorf("expected %s to be %v, got %v", field, expected, value)
		}
		if existIfEmpty && !found {
			t.Errorf("expected field %s to exist, even if empty", field)
		}
	case []corev1.Toleration:
		value, found, err := unstructured.NestedSlice(obj, nestedField...)
		if err != nil {
			t.Fatal(err)
		}
		rawExpected, err := ConvertJSON(expected)
		if err != nil {
			t.Fatalf("unable to test field %s: %v", field, err)
		}
		if !equality.Semantic.DeepEqual(value, rawExpected) {
			t.Errorf("expected %s to be %v, got %v", field, expected, value)
		}
		if existIfEmpty && !found {
			t.Errorf("expected field %s to exist, even if empty", field)
		}
	case []configv1.ImageLabel:
		value, found, err := unstructured.NestedSlice(obj, nestedField...)
		if err != nil {
			t.Fatal(err)
		}
		rawExpected, err := ConvertJSON(expected)
		if err != nil {
			t.Fatalf("unable to test field %s: %v", field, err)
		}
		if !equality.Semantic.DeepEqual(value, rawExpected) {
			t.Errorf("expected %s to be %v, got %v", field, expected, value)
		}
		if existIfEmpty && !found {
			t.Errorf("expected field %s to exist, even if empty", field)
		}
	case []interface{}:
		value, found, err := unstructured.NestedSlice(obj, nestedField...)
		if err != nil {
			t.Fatalf("unable to test field %s: %v", field, err)
		}
		rawExpected, err := ConvertJSON(expected)
		if err != nil {
			t.Fatalf("unable to test field %s: %v", field, err)
		}
		if !equality.Semantic.DeepEqual(value, rawExpected) {
			t.Errorf("expected %s to be %v, got %v", field, expected, value)
		}
		if existIfEmpty && !found {
			t.Errorf("expected field %s to exist, even if empty", field)
		}
	default:
		value, found, err := unstructured.NestedFieldCopy(obj, nestedField...)
		if err != nil {
			t.Fatalf("unable to test field %s: %v", field, err)
		}
		rawExpected, err := ConvertJSON(expected)
		if err != nil {
			t.Fatalf("unable to test field %s: %v", field, err)
		}
		if !equality.Semantic.DeepEqual(rawExpected, value) {
			t.Errorf("expected %s to be %v; got %v", field, expected, value)
		}
		if existIfEmpty && !found {
			t.Errorf("expected field %s to exist, even if empty", field)
		}
	}
}

func TestObserveBuildAdditionalCA(t *testing.T) {
	const caMountPoint = "/var/run/configmaps/additional-ca/additional-ca.crt"
	tests := []struct {
		name        string
		inputCA     *corev1.ConfigMap
		currentCA   *corev1.ConfigMap
		expectedCA  string
		expectError bool
	}{
		{
			name: "has certificate",
			inputCA: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ca",
					Namespace: openshiftConfigNamespaceName,
				},
				Data: map[string]string{
					"ca.crt": dummyCA,
				},
			},
			expectedCA: caMountPoint,
		},
		{
			name: "no certificate",
		},
		{
			name: "update certificate",
			inputCA: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ca",
					Namespace: openshiftConfigNamespaceName,
				},
				Data: map[string]string{
					"ca.crt": dummyCA,
				},
			},
			currentCA: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "other-ca",
					Namespace: openshiftConfigNamespaceName,
				},
				Data: map[string]string{
					"ca.crt": dummyCA + dummyCA,
				},
			},
			expectedCA: caMountPoint,
		},
		{
			name: "bad certificate",
			inputCA: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "other-ca",
					Namespace: openshiftConfigNamespaceName,
				},
				Data: map[string]string{
					"ca.crt": "THIS IS A BAD CERT",
				},
			},
			expectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buildIndex := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			cmIndex := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
			expectedHash := ""
			opConfig := &v1alpha1.OpenShiftControllerManagerOperatorConfig{
				Status: v1alpha1.OpenShiftControllerManagerOperatorConfigStatus{},
			}
			if tc.inputCA != nil {
				expectedHash = hashCAMap(tc.inputCA)
				cmIndex.Add(tc.inputCA)
				bc := &configv1.Build{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.BuildSpec{
						AdditionalTrustedCA: configv1.ConfigMapReference{
							Name:      tc.inputCA.Name,
							Namespace: tc.inputCA.Namespace,
						},
					},
				}
				buildIndex.Add(bc)
			}
			if tc.currentCA != nil {
				cmIndex.Add(tc.currentCA)
				trustedCA := &v1alpha1.AdditionalTrustedCA{
					SHA1Hash:      hashCAMap(tc.currentCA),
					ConfigMapName: tc.currentCA.Name,
				}
				opConfig.Spec.AdditionalTrustedCA = trustedCA
				opConfig.Status.AdditionalTrustedCA = trustedCA
			}
			listers := Listers{
				buildConfigLister:            configlistersv1.NewBuildLister(buildIndex),
				clusterConfigConfigMapLister: corelistersv1.NewConfigMapLister(cmIndex),
			}
			result, err := observeBuildAdditionalCA(listers, map[string]interface{}{}, opConfig)
			if tc.expectError {
				if err == nil {
					t.Error("expected error to occur")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			additionalCA, _, err := unstructured.NestedString(result, "build", "additionalTrustedCA")
			if err != nil {
				t.Fatal(err)
			}
			if additionalCA != tc.expectedCA {
				t.Errorf("expected additional CA to be set to %s, got %s", tc.expectedCA, additionalCA)
			}
			if len(expectedHash) > 0 {
				if opConfig.Spec.AdditionalTrustedCA == nil {
					t.Error("expected operator config to have spec.additionalTrustedCA set, got nil")
					return
				}
				if expectedHash != opConfig.Spec.AdditionalTrustedCA.SHA1Hash {
					t.Errorf("expected CA hash %s, got %s", expectedHash, opConfig.Spec.AdditionalTrustedCA.SHA1Hash)
				}
				if tc.inputCA.Name != opConfig.Spec.AdditionalTrustedCA.ConfigMapName {
					t.Errorf("expected operator config spec to reference %s/%s, got %s/%s",
						tc.inputCA.Namespace,
						tc.inputCA.Name,
						openshiftConfigNamespaceName,
						opConfig.Spec.AdditionalTrustedCA.ConfigMapName)
				}
			}
		})
	}
}
