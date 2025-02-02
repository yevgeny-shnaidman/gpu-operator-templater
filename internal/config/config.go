package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

/*
type TemplaterConfig struct {
        Vendor                   string `yaml:"vendor"`
        CodeRepo                 string `yaml:"codeRepo"`
        APIVersion               string `yaml:"apiVersion"`
	Domain                   string `yaml:"domain"`
        Group                    string `yaml:"group"`
        PCIVendorID              string `yaml:"pciVendorID"`
        KernelModuleName         string `yaml:"kernelModuleName"`
        DefaultDevicePluginImage string `yaml:"defaultDevicePluginImage"`
        ImageFirmwarePath        string `yaml:"imageFirmwarePath"`
        DefaultDriverVersion     string `yaml:"defaultDriverVersion"`
        DefaultNodeLabellerImage string `yaml:"defaultNodeLabellerImage"`
        NodeMetricsImage         string `yaml:"nodeMetricsImage"`
	OperatorImage            string `yaml:"operatorImage"`
	RepoName                 string

}
*/

type TemplaterConfig struct {
	API struct {
		Vendor                   string `yaml:"vendor"`
		CodeRepo                 string `yaml:"codeRepo"`
		Version               string `yaml:"version"`
		Domain                   string `yaml:"domain"`
		Group                    string `yaml:"apiGroup"`

	} `yaml:"api"`
	KMM struct {
		PCIVendorID              string `yaml:"pciVendorID"`
		KernelModuleName         string `yaml:"kernelModuleName"`
		EnableDevicePlugin       bool `yaml:"enableDevicePlugin"`
		DevicePluginImage string `yaml:"devicePluginImage"`
		EnableFirmware           bool `yaml:"enableFirmware"`
		ImageFirmwarePath        string `yaml:"imageFirmwarePath"`
		DriverVersion     string `yaml:"driverVersion"`
	} `yaml:"kmm"`
	NodeLabeller struct {
		Enable  bool `yaml:"enable"`
		Image string `yaml:"image"`
	} `yaml:"nodeLabeller"`
	NodeMetrics struct {
		Enable  bool `yaml:"enable"`
		Image string `yaml:"image"`
	} `yaml:"nodeMetrics"`
	
	OperatorImage            string `yaml:"operatorImage"`
	RepoName                 string
}

func InitConfigData(configFilePath string) (*TemplaterConfig, error) {
	var configData TemplaterConfig
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file %s: error %v", configFilePath, err)
	}

	err = yaml.Unmarshal(yamlFile, &configData)
        if err != nil {
                return nil, fmt.Errorf("failed to unmarshal the values file %s into struct: error %v", configFilePath, err)
        }

	configData.RepoName = filepath.Base(configData.API.CodeRepo)

	return &configData, nil
}
