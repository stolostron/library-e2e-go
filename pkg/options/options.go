package options

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	libgocmd "github.com/open-cluster-management/library-e2e-go/pkg/cmd"

	"k8s.io/klog"
)

type TestOptionsContainer struct {
	Options TestOptionsT `json:"options"`
}

// TestOptions ...
// Define options available for Tests to consume
type TestOptionsT struct {
	Hub               Cluster         `json:"hub"`
	ManagedClusters   []Cluster       `json:"clusters"`
	ImageRegistry     ImageRegistry   `json:"imageRegistry,omitempty"`
	OCPReleaseVersion string          `json:"ocpReleaseVersion,omitempty"`
	IdentityProvider  string          `json:"identityProvider,omitempty"`
	CloudConnection   CloudConnection `json:"cloudConnection,omitempty"`
}

// Cluster ...
// Define the shape of clusters that may be added under management
type Cluster struct {
	Name        string          `json:"name,omitempty"`
	Namespace   string          `json:"namespace,omitempty"`
	Tags        map[string]bool `json:"tags,omitempty"`
	BaseDomain  string          `json:"baseDomain"`
	User        string          `json:"user,omitempty"`
	Password    string          `json:"password,omitempty"`
	KubeContext string          `json:"kubecontext,omitempty"`
	MasterURL   string          `json:"masterURL,omitempty"`
	KubeConfig  string          `json:"kubeconfig,omitempty"`
}

// ImageRegistry - define the image repo information
type ImageRegistry struct {
	Server   string `json:"server"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// CloudConnection struct for bits having to do with Connections
type CloudConnection struct {
	PullSecret    string  `json:"pullSecret"`
	SSHPrivateKey string  `json:"sshPrivatekey"`
	SSHPublicKey  string  `json:"sshPublickey"`
	APIKeys       APIKeys `json:"apiKeys,omitempty"`
	// OCPRelease    string  `yaml:"ocpRelease,omitempty"`
}

// APIKeys - define the cloud connection information
type APIKeys struct {
	AWS       AWSAPIKey       `json:"aws,omitempty"`
	GCP       GCPAPIKey       `json:"gcp,omitempty"`
	Azure     AzureAPIKey     `json:"azure,omitempty"`
	BareMetal BareMetalAPIKey `json:"baremetal,omitempty"`
}

// AWSAPIKey ...
type AWSAPIKey struct {
	AWSAccessKeyID  string `json:"awsAccessKeyID"`
	AWSAccessSecret string `json:"awsSecretAccessKeyID"`
	BaseDNSDomain   string `json:"baseDnsDomain"`
	Region          string `json:"region"`
}

// GCPAPIKey ...
type GCPAPIKey struct {
	ProjectID             string `json:"gcpProjectID"`
	ServiceAccountJSONKey string `json:"gcpServiceAccountJsonKey"`
	BaseDNSDomain         string `json:"baseDnsDomain"`
	Region                string `json:"region"`
}

// AzureAPIKey ...
type AzureAPIKey struct {
	BaseDomainRGN  string `json:"azureBaseDomainRGN"`
	BaseDNSDomain  string `json:"baseDnsDomain"`
	SubscriptionID string `json:"subscriptionID"`
	ClientID       string `json:"clientID"`
	ClientSecret   string `json:"clientSecret"`
	TenantID       string `json:"tenantID"`
	Region         string `json:"region"`
}

// BareMetalAPIKey ...
type BareMetalAPIKey struct {
	ClusterName                  string   `json:"clusterName"`
	BaseDNSDomain                string   `json:"baseDnsDomain"`
	LibvirtURI                   string   `json:"libvirtURI"`
	ProvisioningNetworkCIDR      string   `json:"provisioningNetworkCIDR"`
	ProvisioningNetworkInterface string   `json:"provisioningNetworkInterface"`
	ProvisioningBridge           string   `json:"provisioningBridge"`
	ExternalBridge               string   `json:"externalBridge"`
	APIVIP                       string   `json:"apiVIP"`
	IngressVIP                   string   `json:"ingressVIP"`
	SSHKnownHostsList            []string `json:"sshKnownHostsList"`
	ImageRegistryMirror          string   `json:"imageRegistryMirror"`
	BootstrapOSImage             string   `json:"bootstrapOSImage"`
	ClusterOSImage               string   `json:"clusterOSImage"`
	TrustBundle                  string   `json:"trustBundle"`
	Hosts                        []Hosts  `json:"hosts"`
}

// Hosts ...
type Hosts struct {
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	Role            string `json:"role"`
	Bmc             BMC    `json:"bmc"`
	BootMACAddress  string `json:"bootMACAddress"`
	HardwareProfile string `json:"hardwareProfile"`
}

// BMC ...
type BMC struct {
	Address                        string `json:"address"`
	DisableCertificateVerification bool   `json:"disableCertificateVerification"`
	Username                       string `json:"username"`
	Password                       string `json:"password"`
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"0123456789"

var TestOptions TestOptionsContainer

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

	// if owner == "" {
	// 	owner = os.Getenv("USER")
	// }
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
		return TestOptions.Options.CloudConnection.APIKeys.AWS.Region, nil
	case "azure":
		return TestOptions.Options.CloudConnection.APIKeys.Azure.Region, nil
	case "gcp":
		return TestOptions.Options.CloudConnection.APIKeys.GCP.Region, nil
	default:
		return "", fmt.Errorf("Can not find region as the cloud %s is unsuported", cloud)

	}
}

//GetBaseDomain returns the BaseDomain for the supported cloud providers
func GetBaseDomain(cloud string) (string, error) {
	switch cloud {
	case "aws":
		return TestOptions.Options.CloudConnection.APIKeys.AWS.BaseDNSDomain, nil
	case "azure":
		return TestOptions.Options.CloudConnection.APIKeys.Azure.BaseDNSDomain, nil
	case "gcp":
		return TestOptions.Options.CloudConnection.APIKeys.GCP.BaseDNSDomain, nil
	case "baremetal":
		return TestOptions.Options.CloudConnection.APIKeys.BareMetal.BaseDNSDomain, nil
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
