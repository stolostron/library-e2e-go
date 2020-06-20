package options

type TestOptionsContainer struct {
	Options TestOptions `yaml:"options"`
}

// Define options available for Tests to consume
type TestOptions struct {
	Hub              Hub             `yaml:"hub,omitempty"`
	ManagedClusters  ManagedClusters `yaml:"managedClusters,omitempty"`
	IdentityProvider int             `yaml:"identityProvider,omitempty"`
	Connection       CloudConnection `yaml:"cloudConnection,omitempty"`
	Headless         string          `yaml:"headless,omitempty"`
	OwnerPrefix      string          `yaml:"ownerPrefix,omitempty"`
}

// Define the shape of clusters that may be added under management
type Hub struct {
	ConfigDir  string `yaml:"configDir,omitempty"`
	BaseDomain string `yaml:"baseDomain"`
	User       string `yaml:"user,omitempty"`
	Password   string `yaml:"password,omitempty"`
}

// Define the shape of clusters that may be added under management
type ManagedClusters struct {
	ConfigDir string `yaml:"configDir,omitempty"`
}

// CloudConnection struct for bits having to do with Connections
type CloudConnection struct {
	PullSecret    string  `yaml:"pullSecret"`
	SSHPrivateKey string  `yaml:"sshPrivatekey"`
	SSHPublicKey  string  `yaml:"sshPublickey"`
	Keys          APIKeys `yaml:"apiKeys,omitempty"`
	OCPRelease    string  `yaml:"ocpRelease,omitempty"`
}

type APIKeys struct {
	AWS   AWSAPIKey   `yaml:"aws,omitempty"`
	GCP   GCPAPIKey   `yaml:"gcp,omitempty"`
	Azure AzureAPIKey `yaml:"azure,omitempty"`
}

type AWSAPIKey struct {
	AWSAccessID     string `yaml:"awsAccessKeyID"`
	AWSAccessSecret string `yaml:"awsSecretAccessKeyID"`
	BaseDnsDomain   string `yaml:"baseDnsDomain"`
	Region          string `yaml:"region"`
}

type GCPAPIKey struct {
	ProjectID             string `yaml:"gcpProjectID"`
	ServiceAccountJsonKey string `yaml:"gcpServiceAccountJsonKey"`
	BaseDnsDomain         string `yaml:"baseDnsDomain"`
	Region                string `yaml:"region"`
}

type AzureAPIKey struct {
	BaseDnsDomain        string `yaml:"baseDnsDomain"`
	BaseDomainRGN        string `yaml:"azureBaseDomainRGN"`
	ServicePrincipalJson string `yaml:"azureServicePrincipalJson"`
	Region               string `yaml:"region"`
}
