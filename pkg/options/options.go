package options

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

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
	Owner     string `yaml:"owner,omitempty"`
	//TODO: add OCP image, as an array in order to test sequentially
	//or a single value and launch concurrently multiple tests with different options
	ImageSetRefName string `yaml:"imageSetRefName,omitempty"`
}

// CloudConnection struct for bits having to do with Connections
type CloudConnection struct {
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
	AWSAccessKeyID  string `yaml:"aws_access_key_id"`
	AWSAccessSecret string `yaml:"aws_secret_access_key"`
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

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

var TestOptions TestOptionsT

func LoadOptions(optionsFile string) error {
	err := unmarshal(optionsFile)
	if err != nil {
		klog.Errorf("--options error: %v", err)
		return err
	}
	return nil
}

func unmarshal(optionsFile string) error {

	if optionsFile == "" {
		optionsFile = os.Getenv("OPTIONS")
		if optionsFile == "" {
			optionsFile = "resources/options.yaml"
		}
	}

	klog.V(1).Infof("options filename=%s", optionsFile)

	data, err := ioutil.ReadFile(optionsFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(data), &TestOptions)
	if err != nil {
		return err
	}

	return nil

}

func GetClusterName(cloud string) (string, error) {
	var ownerPrefix string
	var uidPostfix string
	// OwnerPrefix is used to help identify who owns deployed resources
	//    If a value is not supplied, the default is OS environment variable $USER
	if libgocmd.End2End.Owner != "" {
		ownerPrefix = libgocmd.End2End.Owner
	} else {
		if TestOptions.ManagedClusters.Owner == "" {
			ownerPrefix = os.Getenv("USER")
			if ownerPrefix == "" {
				ownerPrefix = "ginkgo"
			}
		} else {
			ownerPrefix = TestOptions.ManagedClusters.Owner
		}
	}
	klog.V(1).Infof("ownerPrefix=%s", ownerPrefix)
	uidPostfix = randString(4)
	if libgocmd.End2End.UID != "" {
		uidPostfix = libgocmd.End2End.UID
	}
	switch cloud {
	case "aws", "gcp", "azure":
		return ownerPrefix + "-" + cloud + "-" + uidPostfix, nil
	default:
		return "", fmt.Errorf("Unsupporter cloud %s", cloud)
	}
}

func GetRegion(cloud string) (string, error) {
	switch cloud {
	case "aws":
		return TestOptions.Connection.Keys.AWS.Region, nil
	case "azure":
		return TestOptions.Connection.Keys.Azure.Region, nil
	case "gcp":
		return TestOptions.Connection.Keys.GCP.Region, nil
	default:
		return "", fmt.Errorf("Unsupporter cloud %s", cloud)

	}
}

func GetBaseDomain(cloud string) (string, error) {
	switch cloud {
	case "aws":
		return TestOptions.Connection.Keys.AWS.BaseDnsDomain, nil
	case "azure":
		return TestOptions.Connection.Keys.Azure.BaseDnsDomain, nil
	case "gcp":
		return TestOptions.Connection.Keys.GCP.BaseDnsDomain, nil
	default:
		return "", fmt.Errorf("Unsupporter cloud %s", cloud)

	}
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randString(length int) string {
	return StringWithCharset(length, charset)
}
