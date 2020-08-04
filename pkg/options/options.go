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

// Define options available for Tests to consume
type TestOptionsT struct {
	Hub              Hub             `yaml:"hub,omitempty"`
	ManagedClusters  ManagedClusters `yaml:"managedClusters,omitempty"`
	IdentityProvider int             `yaml:"identityProvider,omitempty"`
	Connection       CloudConnection `yaml:"cloudConnection,omitempty"`
	Headless         string          `yaml:"headless,omitempty"`
	Owner            string          `yaml:"owner,omitempty"`
	UID              string          `yaml:"uid,omitempty"`
}

// Define the shape of clusters that may be added under management
type Hub struct {
	ConfigDir  string `yaml:"configDir,omitempty"`
	BaseDomain string `yaml:"baseDomain"`
	//The hub kubeconfig path
	KubeConfigPath string `yaml:"kubeconfig,omitempty"`
}

// Define the shape of clusters that may be added under management
type ManagedClusters struct {
	ConfigDir string `yaml:"configDir,omitempty"`
	//The ocp imageset name to use while deploying the cluster
	ImageSetRefName string `yaml:"imageSetRefName,omitempty"`
	//TODO: Create image set named <owner>-<tag-of-image>-<uid>
	//OCPImageRelease will use it to create an imageSet named <owner>-<tag-of-image>-<uid>
	//if the ImageSetRefName is empty
	//example quay.io/openshift-release-dev/ocp-release:4.3.28-x86_64
	OCPImageRelease string `yaml:"OCPImageRelease,omitempty"`
}

// CloudConnection struct for bits having to do with Connections
type CloudConnection struct {
	SSHPrivateKey string  `yaml:"sshPrivatekey"`
	SSHPublicKey  string  `yaml:"sshPublickey"`
	Keys          APIKeys `yaml:"apiKeys,omitempty"`
	// OCPRelease    string  `yaml:"ocpRelease,omitempty"`
}

type APIKeys struct {
	AWS   AWSAPIKey   `yaml:"aws,omitempty"`
	GCP   GCPAPIKey   `yaml:"gcp,omitempty"`
	Azure AzureAPIKey `yaml:"azure,omitempty"`
}

type AWSAPIKey struct {
	AWSAccessKeyID  string `yaml:"awsAccessKeyID"`
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
	BaseDnsDomain  string `yaml:"baseDnsDomain"`
	BaseDomainRGN  string `yaml:"azureBaseDomainRGN"`
	ClientId       string `yaml:"clientId"`
	ClientSecret   string `yaml:"clientSecret"`
	TenantId       string `yaml:"tenantId"`
	SubscriptionId string `yaml:"subscriptionId"`
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
		owner = TestOptions.Owner
	}
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
		uid = TestOptions.UID
	}
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
		return TestOptions.Connection.Keys.AWS.Region, nil
	case "azure":
		return TestOptions.Connection.Keys.Azure.Region, nil
	case "gcp":
		return TestOptions.Connection.Keys.GCP.Region, nil
	default:
		return "", fmt.Errorf("Can not find region as the cloud %s is unsuported", cloud)

	}
}

//GetBaseDomain returns the BaseDomain for the supported cloud providers
func GetBaseDomain(cloud string) (string, error) {
	switch cloud {
	case "aws":
		return TestOptions.Connection.Keys.AWS.BaseDnsDomain, nil
	case "azure":
		return TestOptions.Connection.Keys.Azure.BaseDnsDomain, nil
	case "gcp":
		return TestOptions.Connection.Keys.GCP.BaseDnsDomain, nil
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
