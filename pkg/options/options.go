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
	ConfigDir       string `yaml:"configDir,omitempty"`
	Owner           string `yaml:"owner,omitempty"`
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

	klog.V(2).Infof("options filename=%s", optionsFile)

	data, err := ioutil.ReadFile(filepath.Clean(optionsFile))
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
	uidPostfix, err := randString(4)
	if err != nil {
		return "", err
	}
	if libgocmd.End2End.UID != "" {
		uidPostfix = libgocmd.End2End.UID
	}
	switch cloud {
	case "aws", "gcp", "azure":
		return ownerPrefix + "-" + cloud + "-" + uidPostfix, nil
	default:
		return "", fmt.Errorf("Can not generate cluster name as the cloud %s is unsupported", cloud)
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
		return "", fmt.Errorf("Can not find region as the cloud %s is unsuported", cloud)

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
		return "", fmt.Errorf("Can not find the baseDomain as the cloud %s is unsupported", cloud)

	}
}

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
