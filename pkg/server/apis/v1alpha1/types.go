package v1alpha1

import (
	core_v1 "k8s.io/api/core/v1"
	cmd_api "k8s.io/client-go/tools/clientcmd/api"
)

// Tenant contains information about tenant
type Tenant struct {
	// Metadata contains metadata information about tenant
	Metadata Metadata `json:"metadata"`
	// Spec contains tenant spec
	Spec TenantSpec `json:"spec"`
}

// Metadata contains metadata information
type Metadata struct {
	// Name is the name of the resource
	Name string `json:"name"`
	// Description describes the resource
	Description string `json:"description,omitempty"`
	// CreationTime records the time of the resource creation
	CreationTime string `json:"creationTime"`
}

// TenantSpec contains the tenant spec information
type TenantSpec struct {
	// PersistentVolumeClaim describes information about persistent volume claim
	PersistentVolumeClaim PersistentVolumeClaim `json:"persistentVolumeClaim"`

	// ResourceQuota describes the resource quota of the namespace,
	// eg map[core_v1.ResourceName]string{"cpu": "2", "memory": "4Gi"}
	ResourceQuota map[core_v1.ResourceName]string `json:"resourceQuota"`
}

// PersistentVolumeClaim describes information about pvc belongs to a tenant
type PersistentVolumeClaim struct {
	// Name is the pvc name specified by user, and if Name is not nil, cyclone will
	// use this pvc and not to create another one.
	// Name string `json:"name"`

	// StorageClass represents the strorageclass used to create pvc
	StorageClass string `json:"storageClass"`

	// Size represents the capacity of the pvc, unit supports 'Gi' or 'Mi'
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#capacity
	Size string `json:"size"`
}

// Integration contains information about external systems
type Integration struct {
	// Metadata contains metadata information about integration
	Metadata Metadata `json:"metadata"`
	// Spec contains integration spec
	Spec IntegrationSpec `json:"spec"`
}

// IntegrationType defines the type of integration
type IntegrationType string

const (
	// SonarQube is the SonarQube integration
	SonarQube IntegrationType = "SonarQube"
	// DockerRegistry is the DockerRegistry integration
	DockerRegistry IntegrationType = "DockerRegistry"
	// SCM is the SCM integration
	SCM IntegrationType = "SCM"
	// Cluster is the Cluster integration
	Cluster IntegrationType = "Cluster"
	// GeneralIntegration is the General integration
	GeneralIntegration IntegrationType = "General"
)

// IntegrationSpec contains the integration spec information
type IntegrationSpec struct {
	// Type of integration
	Type IntegrationType `json:"type"`
	// The actual info about various external systems.
	IntegrationSource `json:",inline"`
}

// IntegrationSource contains various external systems.
// exactly one of its members must be set, and the member must equal with the integration's type.
type IntegrationSource struct {
	// SonarQube describes info about external system sonar qube, and is used for code scanning in CI.
	SonarQube *SonarQubeSource `json:"sonarQube"`

	// DockerRegistry describes info about external system docker registry, and is used to manager containers.
	DockerRegistry *DockerRegistrySource `json:"dockerRegistry"`

	// SCM describes info about external Source Code Management system, and is used to manager code.
	SCM *SCMSource `json:"scm"`

	// Cluster contains information about clusters.
	// Users can define which cluster will be used to run workload,
	// and clusters integrated here can be used to deploy application in CD tasks.
	Cluster *ClusterSource `json:"cluster"`

	// General contains parameters defined by users.
	General []ParameterItem `json:"general"`
}

// ParameterItem defines a parameter
type ParameterItem struct {
	// Name of the parameter
	Name string `json:"name"`
	// Value of the parameter
	Value string `json:"value"`
}

// SonarQubeSource represents a code scanning tool for CI.
type SonarQubeSource struct {
	// Server represents the server address of sonar qube .
	Server string `json:"server"`
	// Token is the credential to access sonar server.
	Token string `json:"token"`
}

// DockerRegistrySource represents a docker registry to manager containers.
type DockerRegistrySource struct {
	// Server represents the domain of docker registry.
	Server string `json:"server"`
	// User is a user of registry.
	User string `json:"user"`
	// Password is the password of the corresponding user.
	Password string `json:"password"`
}

// SCMType defines the type of Source Code Management
type SCMType string

const (
	// GitLab is the Gitlab scm
	GitLab SCMType = "GitLab"
	// GitHub is the GitHub scm
	GitHub = "GitHub"
	// SVN is the SVN scm
	SVN = "SVN"
)

// SCMSource represents Source Code Management to manage code.
type SCMSource struct {
	// Type is the type of scm, e.g. GitLab, GitHub, SVN
	Type SCMType `json:"type"`
	// Server represents the domain of docker registry.
	Server string `json:"server"`
	// User is a user of the SCM.
	User string `json:"user"`
	// Password is the password of the corresponding user.
	Password string `json:"password"`
	// Token is the credential to access SCM.
	Token string `json:"token"`
}

// ClusterSource contains info about clusters.
type ClusterSource struct {
	// Credential is the credential info of the cluster
	Credential ClusterCredential `json:"credential"`
	// IsControlCluster describes whether the cluster is the control cluster itself
	IsControlCluster bool `json:"isControlCluster,omitempty"`
	// IsWorkerCluster defines whether this cluster can be used to run workflow.
	// True, will create namespace and pvc associated with tenant in the cluster.
	// False, will delete namespace and pvc associated with tenant in the cluster.
	IsWorkerCluster bool `json:"isWorkerCluster"`
	// Namespace is the namespace where workflow will run in.
	// If set, cyclone will use it directly, otherwise a new one will be created.
	// It's used when 'IsWorkerCluster' is True.
	Namespace string `json:"namespace"`
	// PVC is the pvc name specified by user, and if this is not nil, cyclone will
	// use this pvc and not to create another one.
	// It's used when 'IsWorkerCluster' is True.
	PVC string `json:"pvc"`
}

// ClusterCredential contains credential info about cluster
type ClusterCredential struct {
	// Server represents the address of cluster.
	Server string `json:"server"`
	// User is a user of the cluster.
	User string `json:"user"`
	// Password is the password of the corresponding user.
	Password string `json:"password"`
	// BearerToken is the credential to access cluster.
	BearerToken string `json:"bearerToken"`
	// TLSClientConfig is the config about TLS
	TLSClientConfig *TLSClientConfig `json:"tlsClientConfig,omitempty"`
	// KubeConfig is the config about kube config
	KubeConfig *cmd_api.Config `json:"kubeConfig,omitempty"`
}

// TLSClientConfig contains settings to enable transport layer security
type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	Insecure bool `json:"insecure,omitempty" bson:"insecure"`

	// CAFile is the trusted root certificates for server
	CAFile string `json:"caFile,omitempty" bson:"caFile"`

	// CAData holds PEM-encoded bytes (typically read from a root certificates bundle).
	// CAData takes precedence over CAFile
	CAData []byte `json:"caData,omitempty" bson:"caData"`
}
