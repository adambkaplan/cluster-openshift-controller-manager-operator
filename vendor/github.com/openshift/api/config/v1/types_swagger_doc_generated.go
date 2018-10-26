package v1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-generated-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_AdmissionPluginConfig = map[string]string{
	"":              "AdmissionPluginConfig holds the necessary configuration options for admission plugins",
	"location":      "Location is the path to a configuration file that contains the plugin's configuration",
	"configuration": "Configuration is an embedded configuration object to be used as the plugin's configuration. If present, it will be used instead of the path to the configuration file.",
}

func (AdmissionPluginConfig) SwaggerDoc() map[string]string {
	return map_AdmissionPluginConfig
}

var map_AuditConfig = map[string]string{
	"":                         "AuditConfig holds configuration for the audit capabilities",
	"enabled":                  "If this flag is set, audit log will be printed in the logs. The logs contains, method, user and a requested URL.",
	"auditFilePath":            "All requests coming to the apiserver will be logged to this file.",
	"maximumFileRetentionDays": "Maximum number of days to retain old log files based on the timestamp encoded in their filename.",
	"maximumRetainedFiles":     "Maximum number of old log files to retain.",
	"maximumFileSizeMegabytes": "Maximum size in megabytes of the log file before it gets rotated. Defaults to 100MB.",
	"policyFile":               "PolicyFile is a path to the file that defines the audit policy configuration.",
	"policyConfiguration":      "PolicyConfiguration is an embedded policy configuration object to be used as the audit policy configuration. If present, it will be used instead of the path to the policy file.",
	"logFormat":                "Format of saved audits (legacy or json).",
	"webHookKubeConfig":        "Path to a .kubeconfig formatted file that defines the audit webhook configuration.",
	"webHookMode":              "Strategy for sending audit events (block or batch).",
}

func (AuditConfig) SwaggerDoc() map[string]string {
	return map_AuditConfig
}

var map_Build = map[string]string{
	"":     "Build holds cluster-wide information on how to handle builds. The canonical name is `cluster`",
	"spec": "Spec holds user-settable values for the build controller configuration",
}

func (Build) SwaggerDoc() map[string]string {
	return map_Build
}

var map_BuildDefaults = map[string]string{
	"gitHTTPProxy":  "GitHTTPProxy is the location of the HTTPProxy for Git source",
	"gitHTTPSProxy": "GitHTTPSProxy is the location of the HTTPSProxy for Git source",
	"gitNoProxy":    "GitNoProxy is the list of domains for which the proxy should not be used",
	"env":           "Env is a set of default environment variables that will be applied to the build if the specified variables do not exist on the build",
	"sourceStrategyDefaults": "SourceStrategyDefaults are default values that apply to builds using the source strategy.",
	"imageLabels":            "ImageLabels is a list of docker labels that are applied to the resulting image. User can override a default label by providing a label with the same name in their Build/BuildConfig.",
	"nodeSelector":           "NodeSelector is a selector which must be true for the build pod to fit on a node",
	"annotations":            "Annotations are annotations that will be added to the build pod",
	"resources":              "Resources defines resource requirements to execute the build.",
}

func (BuildDefaults) SwaggerDoc() map[string]string {
	return map_BuildDefaults
}

var map_BuildOverrides = map[string]string{
	"forcePull":    "ForcePull indicates whether the build strategy should always be set to ForcePull=true",
	"imageLabels":  "ImageLabels is a list of docker labels that are applied to the resulting image. If user provided a label in their Build/BuildConfig with the same name as one in this list, the user's label will be overwritten.",
	"nodeSelector": "NodeSelector is a selector which must be true for the build pod to fit on a node",
	"annotations":  "Annotations are annotations that will be added to the build pod",
	"tolerations":  "Tolerations is a list of Tolerations that will override any existing tolerations set on a build pod.",
}

func (BuildOverrides) SwaggerDoc() map[string]string {
	return map_BuildOverrides
}

var map_BuildSpec = map[string]string{
	"additionalTrustedCA": "AdditionalTrustedCA is a reference to a ConfigMap containing additional CAs that should be trusted for image pushes and pulls during builds.",
	"buildDefaults":       "BuildDefaults controls the default information for Builds",
	"buildOverrides":      "BuildOverrides controls override settings for builds",
}

func (BuildSpec) SwaggerDoc() map[string]string {
	return map_BuildSpec
}

var map_CertInfo = map[string]string{
	"":         "CertInfo relates a certificate with a private key",
	"certFile": "CertFile is a file containing a PEM-encoded certificate",
	"keyFile":  "KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile",
}

func (CertInfo) SwaggerDoc() map[string]string {
	return map_CertInfo
}

var map_ClientConnectionOverrides = map[string]string{
	"acceptContentTypes": "acceptContentTypes defines the Accept header sent by clients when connecting to a server, overriding the default value of 'application/json'. This field will control all connections to the server used by a particular client.",
	"contentType":        "contentType is the content type used when sending data to the server from this client.",
	"qps":                "qps controls the number of queries per second allowed for this connection.",
	"burst":              "burst allows extra queries to accumulate when a client is exceeding its rate.",
}

func (ClientConnectionOverrides) SwaggerDoc() map[string]string {
	return map_ClientConnectionOverrides
}

var map_ConfigMapReference = map[string]string{
	"":         "ConfigMapReference references the location of a configmap.",
	"filename": "Key allows pointing to a specific key/value inside of the configmap.  This is useful for logical file references.",
}

func (ConfigMapReference) SwaggerDoc() map[string]string {
	return map_ConfigMapReference
}

var map_EtcdConnectionInfo = map[string]string{
	"":     "EtcdConnectionInfo holds information necessary for connecting to an etcd server",
	"urls": "URLs are the URLs for etcd",
	"ca":   "CA is a file containing trusted roots for the etcd server certificates",
}

func (EtcdConnectionInfo) SwaggerDoc() map[string]string {
	return map_EtcdConnectionInfo
}

var map_EtcdStorageConfig = map[string]string{
	"storagePrefix": "StoragePrefix is the path within etcd that the OpenShift resources will be rooted under. This value, if changed, will mean existing objects in etcd will no longer be located.",
}

func (EtcdStorageConfig) SwaggerDoc() map[string]string {
	return map_EtcdStorageConfig
}

var map_GenericAPIServerConfig = map[string]string{
	"":                   "GenericAPIServerConfig is an inline-able struct for aggregated apiservers that need to store data in etcd",
	"servingInfo":        "ServingInfo describes how to start serving",
	"corsAllowedOrigins": "CORSAllowedOrigins",
	"auditConfig":        "AuditConfig describes how to configure audit information",
	"storageConfig":      "StorageConfig contains information about how to use",
}

func (GenericAPIServerConfig) SwaggerDoc() map[string]string {
	return map_GenericAPIServerConfig
}

var map_HTTPServingInfo = map[string]string{
	"": "HTTPServingInfo holds configuration for serving HTTP",
	"maxRequestsInFlight":   "MaxRequestsInFlight is the number of concurrent requests allowed to the server. If zero, no limit.",
	"requestTimeoutSeconds": "RequestTimeoutSeconds is the number of seconds before requests are timed out. The default is 60 minutes, if -1 there is no limit on requests.",
}

func (HTTPServingInfo) SwaggerDoc() map[string]string {
	return map_HTTPServingInfo
}

var map_Image = map[string]string{
	"":         "Image holds cluster-wide information about how to handle images.  The canonical name is `cluster`",
	"metadata": "Standard object's metadata.",
	"spec":     "spec holds user settable values for configuration",
	"status":   "status holds observed values from the cluster. They may not be overridden.",
}

func (Image) SwaggerDoc() map[string]string {
	return map_Image
}

var map_ImageLabel = map[string]string{
	"name":  "Name defines the name of the label. It must have non-zero length.",
	"value": "Value defines the literal value of the label.",
}

func (ImageLabel) SwaggerDoc() map[string]string {
	return map_ImageLabel
}

var map_ImageList = map[string]string{
	"metadata": "Standard object's metadata.",
}

func (ImageList) SwaggerDoc() map[string]string {
	return map_ImageList
}

var map_ImageSpec = map[string]string{
	"allowedRegistriesForImport": "AllowedRegistriesForImport limits the docker registries that normal users may import images from. Set this list to the registries that you trust to contain valid Docker images and that you want applications to be able to import from. Users with permission to create Images or ImageStreamMappings via the API are not affected by this policy - typically only administrators or system integrations will have those permissions.",
	"externalRegistryHostname":   "ExternalRegistryHostname sets the hostname for the default external image registry. The external hostname should be set only when the image registry is exposed externally. The value is used in 'publicDockerImageRepository' field in ImageStreams. The value must be in \"hostname[:port]\" format.",
	"additionalTrustedCA":        "AdditionalTrustedCA is a reference to a ConfigMap containing additional CAs that should be trusted during imagestream import.",
}

func (ImageSpec) SwaggerDoc() map[string]string {
	return map_ImageSpec
}

var map_ImageStatus = map[string]string{
	"internalRegistryHostname": "this value is set by the image registry operator which controls the internal registry hostname InternalRegistryHostname sets the hostname for the default internal image registry. The value must be in \"hostname[:port]\" format. For backward compatibility, users can still use OPENSHIFT_DEFAULT_REGISTRY environment variable but this setting overrides the environment variable.",
}

func (ImageStatus) SwaggerDoc() map[string]string {
	return map_ImageStatus
}

var map_KubeClientConfig = map[string]string{
	"kubeConfig":          "kubeConfig is a .kubeconfig filename for going to the owning kube-apiserver.  Empty uses an in-cluster-config",
	"connectionOverrides": "connectionOverrides specifies client overrides for system components to loop back to this master.",
}

func (KubeClientConfig) SwaggerDoc() map[string]string {
	return map_KubeClientConfig
}

var map_LeaderElection = map[string]string{
	"":              "LeaderElection provides information to elect a leader",
	"disable":       "disable allows leader election to be suspended while allowing a fully defaulted \"normal\" startup case.",
	"namespace":     "namespace indicates which namespace the resource is in",
	"name":          "name indicates what name to use for the resource",
	"leaseDuration": "leaseDuration is the duration that non-leader candidates will wait after observing a leadership renewal until attempting to acquire leadership of a led but unrenewed leader slot. This is effectively the maximum duration that a leader can be stopped before it is replaced by another candidate. This is only applicable if leader election is enabled.",
	"renewDeadline": "renewDeadline is the interval between attempts by the acting master to renew a leadership slot before it stops leading. This must be less than or equal to the lease duration. This is only applicable if leader election is enabled.",
	"retryPeriod":   "retryPeriod is the duration the clients should wait between attempting acquisition and renewal of a leadership. This is only applicable if leader election is enabled.",
}

func (LeaderElection) SwaggerDoc() map[string]string {
	return map_LeaderElection
}

var map_NamedCertificate = map[string]string{
	"":      "NamedCertificate specifies a certificate/key, and the names it should be served for",
	"names": "Names is a list of DNS names this certificate should be used to secure A name can be a normal DNS name, or can contain leading wildcard segments.",
}

func (NamedCertificate) SwaggerDoc() map[string]string {
	return map_NamedCertificate
}

var map_RegistryLocation = map[string]string{
	"":           "RegistryLocation contains a location of the registry specified by the registry domain name. The domain name might include wildcards, like '*' or '??'.",
	"domainName": "DomainName specifies a domain name for the registry In case the registry use non-standard (80 or 443) port, the port should be included in the domain name as well.",
	"insecure":   "Insecure indicates whether the registry is secure (https) or insecure (http) By default (if not specified) the registry is assumed as secure.",
}

func (RegistryLocation) SwaggerDoc() map[string]string {
	return map_RegistryLocation
}

var map_RemoteConnectionInfo = map[string]string{
	"":    "RemoteConnectionInfo holds information necessary for establishing a remote connection",
	"url": "URL is the remote URL to connect to",
	"ca":  "CA is the CA for verifying TLS connections",
}

func (RemoteConnectionInfo) SwaggerDoc() map[string]string {
	return map_RemoteConnectionInfo
}

var map_ServingInfo = map[string]string{
	"":                  "ServingInfo holds information about serving web pages",
	"bindAddress":       "BindAddress is the ip:port to serve on",
	"bindNetwork":       "BindNetwork is the type of network to bind to - defaults to \"tcp4\", accepts \"tcp\", \"tcp4\", and \"tcp6\"",
	"clientCA":          "ClientCA is the certificate bundle for all the signers that you'll recognize for incoming client certificates",
	"namedCertificates": "NamedCertificates is a list of certificates to use to secure requests to specific hostnames",
	"minTLSVersion":     "MinTLSVersion is the minimum TLS version supported. Values must match version names from https://golang.org/pkg/crypto/tls/#pkg-constants",
	"cipherSuites":      "CipherSuites contains an overridden list of ciphers for the server to support. Values must match cipher suite IDs from https://golang.org/pkg/crypto/tls/#pkg-constants",
}

func (ServingInfo) SwaggerDoc() map[string]string {
	return map_ServingInfo
}

var map_SourceStrategyDefaults = map[string]string{
	"incremental": "Incremental indicates if s2i build strategies should perform an incremental build or not",
}

func (SourceStrategyDefaults) SwaggerDoc() map[string]string {
	return map_SourceStrategyDefaults
}

var map_StringSource = map[string]string{
	"": "StringSource allows specifying a string inline, or externally via env var or file. When it contains only a string value, it marshals to a simple JSON string.",
}

func (StringSource) SwaggerDoc() map[string]string {
	return map_StringSource
}

var map_StringSourceSpec = map[string]string{
	"":        "StringSourceSpec specifies a string value, or external location",
	"value":   "Value specifies the cleartext value, or an encrypted value if keyFile is specified.",
	"env":     "Env specifies an envvar containing the cleartext value, or an encrypted value if the keyFile is specified.",
	"file":    "File references a file containing the cleartext value, or an encrypted value if a keyFile is specified.",
	"keyFile": "KeyFile references a file containing the key to use to decrypt the value.",
}

func (StringSourceSpec) SwaggerDoc() map[string]string {
	return map_StringSourceSpec
}

// AUTO-GENERATED FUNCTIONS END HERE
