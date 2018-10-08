package cloud

import (
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"
)

// DevSpaceCloudConfigPath holds the path to the cloud config file
const DevSpaceCloudConfigPath = ".devspace/clouds.yaml"

// DevSpaceKubeContextName is the name for the kube config context
const DevSpaceKubeContextName = "devspace"

// ProviderConfig holds all the different providers and their configuration
type ProviderConfig map[string]*Provider

// Provider describes the struct to hold the cloud configuration
type Provider struct {
	Name        string `yaml:"name,omitempty"`
	KubeContext string `yaml:"kubecontext,omitempty"`
	Login       string `yaml:"login,omitempty"`
	GetConfig   string `yaml:"getConfig,omitempty"`
	Token       string `yaml:"token,omitempty"`
}

// DevSpaceCloudProviderName is the name of the default devspace-cloud provider
const DevSpaceCloudProviderName = "devspace-cloud"

// DevSpaceCloudProviderConfig holds the information for the devspace-cloud
var DevSpaceCloudProviderConfig = &Provider{
	Login:       "https://cloud.devspace.covexo.com/login",
	GetConfig:   "https://cloud.devspace.covexo.com/clusterConfig",
	KubeContext: DevSpaceKubeContextName,
}

// ParseCloudConfig parses the cloud configuration and returns a map containing the configurations
func ParseCloudConfig() (ProviderConfig, error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath.Join(homedir, DevSpaceCloudConfigPath))
	if os.IsNotExist(err) {
		return ProviderConfig{
			DevSpaceCloudProviderName: DevSpaceCloudProviderConfig,
		}, nil
	}

	cloudConfig := make(ProviderConfig)
	err = yaml.Unmarshal(data, cloudConfig)
	if err != nil {
		return nil, err
	}

	if _, ok := cloudConfig[DevSpaceCloudProviderName]; ok {
		cloudConfig[DevSpaceCloudProviderName].GetConfig = DevSpaceCloudProviderConfig.GetConfig
		cloudConfig[DevSpaceCloudProviderName].Login = DevSpaceCloudProviderConfig.Login
	} else {
		cloudConfig[DevSpaceCloudProviderName] = DevSpaceCloudProviderConfig
	}

	for configName, config := range cloudConfig {
		config.Name = configName
	}

	return cloudConfig, nil
}

// SaveCloudConfig saves the provider configuration to file
func SaveCloudConfig(config ProviderConfig) error {
	homedir, err := homedir.Dir()
	if err != nil {
		return err
	}

	cfgPath := filepath.Join(homedir, DevSpaceCloudConfigPath)

	for name, provider := range config {
		provider.Name = ""

		if name == DevSpaceCloudProviderName {
			provider.Login = ""
			provider.GetConfig = ""
		}
	}

	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(cfgPath), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgPath, out, 0600)
}
