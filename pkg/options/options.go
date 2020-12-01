package options

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	libgocmd "github.com/open-cluster-management/library-e2e-go/pkg/cmd"

	"gopkg.in/yaml.v1"
	"k8s.io/klog"
)

// TestOptions ...
// Define options available for Tests to consume
type TestOptionsT struct {
	Hub               Cluster         `yaml:"hub,omitempty"`
	ManagedClusters   []Cluster       `yaml:"clusters,omitempty"`
	ImageRegistry     ImageRegistry   `yaml:"imageRegistry,omitempty"`
	OCPReleaseVersion string          `yaml:"ocpReleaseVersion,omitempty"`
	IdentityProvider  string          `yaml:"identityProvider,omitempty"`
	CloudConnection   CloudConnection `yaml:"cloudConnection,omitempty"`
}

// Cluster ...
// Define the shape of clusters that may be added under management
type Cluster struct {
	Name        string          `yaml:"name,omitempty"`
	Namespace   string          `yaml:"namespace,omitempty"`
	Tags        map[string]bool `yaml:"tags,omitempty"`
	BaseDomain  string          `yaml:"baseDomain"`
	User        string          `yaml:"user,omitempty"`
	Password    string          `yaml:"password,omitempty"`
	KubeContext string          `yaml:"kubecontext,omitempty"`
	MasterURL   string          `yaml:"masterURL,omitempty"`
	KubeConfig  string          `yaml:"kubeconfig,omitempty"`
}

// ImageRegistry - define the image repo information
type ImageRegistry struct {
	Server   string `yaml:"server,omitemty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitemty"`
}

// CloudConnection struct for bits having to do with Connections
type CloudConnection struct {
	PullSecret    string  `yaml:"pullSecret"`
	SSHPrivateKey string  `yaml:"sshPrivatekey"`
	SSHPublicKey  string  `yaml:"sshPublickey"`
	APIKeys       APIKeys `yaml:"apiKeys,omitempty"`
	// OCPRelease    string  `yaml:"ocpRelease,omitempty"`
}

// APIKeys - define the cloud connection information
type APIKeys struct {
	AWS   AWSAPIKey   `yaml:"aws,omitempty"`
	GCP   GCPAPIKey   `yaml:"gcp,omitempty"`
	Azure AzureAPIKey `yaml:"azure,omitempty"`
}

// AWSAPIKey ...
type AWSAPIKey struct {
	AWSAccessKeyID  string `yaml:"awsAccessKeyID"`
	AWSAccessSecret string `yaml:"awsSecretAccessKeyID"`
	BaseDNSDomain   string `yaml:"baseDnsDomain"`
	Region          string `yaml:"region"`
}

// GCPAPIKey ...
type GCPAPIKey struct {
	ProjectID             string `yaml:"gcpProjectID"`
	ServiceAccountJSONKey string `yaml:"gcpServiceAccountJsonKey"`
	BaseDNSDomain         string `yaml:"baseDnsDomain"`
	Region                string `yaml:"region"`
}

// AzureAPIKey ...
type AzureAPIKey struct {
	BaseDomainRGN  string `yaml:"azureBaseDomainRGN"`
	BaseDNSDomain  string `yaml:"baseDnsDomain"`
	SubscriptionID string `yaml:"subscriptionID"`
	ClientID       string `yaml:"clientID"`
	ClientSecret   string `yaml:"clientSecret"`
	TenantID       string `yaml:"tenantID"`
	Region         string `yaml:"region"`
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"0123456789"

var TestOptions TestOptionsT

//LoadOptions load the options in the following priority:
//1. The provided file path
//2. The OPTIONS environment variable
//3. Default "resources/options.yaml"
func LoadOptions(optionsFile string) error {
	if err := unmarshal(optionsFile); err != nil {
		klog.Errorf("--options error: %v", err)
		return err
	}
	return nil
}

func unmarshal(optionsFile string) error {

	if optionsFile == "" {
		optionsFile = os.Getenv("OPTIONS")
	}
	if optionsFile == "" {
		optionsFile = "resources/options.yaml"
	}

	klog.V(2).Infof("options filename=%s", optionsFile)

	data, err := ioutil.ReadFile(filepath.Clean(optionsFile))
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal([]byte(data), &TestOptions); err != nil {
		return err
	}

	return nil

}

//GetOwner returns the owner in the following priority:
//1. From command-line
//2. From options.yaml
//3. Using the $USER environment variable.
//4. Default: ginkgo
func GetOwner() string {
	// owner is used to help identify who owns deployed resources
	//    If a value is not supplied, the default is OS environment variable $USER
	owner := libgocmd.End2End.Owner
	if owner == "" {
		owner = os.Getenv("USER")
	}
	if owner == "" {
		owner = "ginkgo"
	}
	return owner
}

//GetUID returns the UID in the following priority:
//1. From command-line
//2. From options.yaml
//3. Generate a new one
func GetUID() (string, error) {
	uid := libgocmd.End2End.UID

	if uid == "" {
		var err error
		uid, err = randString(4)
		if err != nil {
			return "", err
		}
	}
	return uid, nil
}

//GetRegion returns the region for the supported cloud providers
func GetRegion(cloud string) (string, error) {
	switch cloud {
	case "aws":
		return TestOptions.CloudConnection.APIKeys.AWS.Region, nil
	case "azure":
		return TestOptions.CloudConnection.APIKeys.Azure.Region, nil
	case "gcp":
		return TestOptions.CloudConnection.APIKeys.GCP.Region, nil
	default:
		return "", fmt.Errorf("Can not find region as the cloud %s is unsuported", cloud)

	}
}

//GetBaseDomain returns the BaseDomain for the supported cloud providers
func GetBaseDomain(cloud string) (string, error) {
	switch cloud {
	case "aws":
		return TestOptions.CloudConnection.APIKeys.AWS.BaseDNSDomain, nil
	case "azure":
		return TestOptions.CloudConnection.APIKeys.Azure.BaseDNSDomain, nil
	case "gcp":
		return TestOptions.CloudConnection.APIKeys.GCP.BaseDNSDomain, nil
	default:
		return "", fmt.Errorf("Can not find the baseDomain as the cloud %s is unsupported", cloud)

	}
}

//StringWithCharset returns a string of the given length and
// componsed from a characters of the provided charset
func StringWithCharset(length int, charset string) (string, error) {
	b := make([]byte, length)
	for i := range b {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[r.Int64()]
	}
	return string(b), nil
}

func randString(length int) (string, error) {
	return StringWithCharset(length, charset)
}
