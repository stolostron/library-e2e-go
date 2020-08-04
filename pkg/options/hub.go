package options

import (
	"os"
	"path/filepath"
)

//GetHubKubeConfig returns the hub kubeconfig path if provided
// if not it search for a kubeocnfig.yaml in the configDir if exists and empty string if not
//The kubeconfig file for the hub is supposed to be in <configDir>/kubeconfig
func GetHubKubeConfig(configDir, kubeConfigPath string) string {
	kubeConfigFilePath := kubeConfigPath
	if kubeConfigFilePath == "" {
		if configDir == "" {
			return ""
		}
		kubeConfigPath = filepath.Join(configDir, "kubeconfig.yaml")
		if _, err := os.Stat(kubeConfigPath); os.IsNotExist(err) {
			kubeConfigPath = ""
		}
	}
	return kubeConfigPath
}
