package code_templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yevgeny-shnaidman/gpu-operator-template/internal/config"
)

type TemplatesData struct {
	Vendor                   string `yaml:"vendor"`
	Repo                     string `yaml:"repo"`
	APIVersion               string `yaml:"apiVersion"`
	Group                    string `yaml:"group"`
	PCIVendorID              string `yaml:"pciVendorID"`
	KernelModuleName         string `yaml:"kernelModuleName"`
	DefaultDevicePluginImage string `yaml:"defaultDevicePluginImage"`
	ImageFirmwarePath        string `yaml:"imageFirmwarePath"`
	DefaultDriverVersion     string `yaml:"defaultDriverVersion"`
	DefaultNodeLabellerImage string `yaml:"defaultNodeLabellerImage"`
	NodeMetricsImage         string `yaml:"nodeMetricsImage"`
}

const (
	apiVersionDir = "API_VERSION"
)

var (
	//go:embed templates
	templatesFS embed.FS
)

func RunTemplates(config *config.TemplaterConfig) error {
	err := fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// get the path in the target
		targetPath := getTargetPath(path, config)
		if d.IsDir() {
			os.Mkdir(targetPath, 0750)
		} else {
			targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("failed to create file %s: error %v", targetPath, err)
			}
			tmpl, err := template.ParseFS(templatesFS, path)
			if err != nil {
				return fmt.Errorf("failed to parse file %s: %w", path, err)
			}
			err = tmpl.ExecuteTemplate(targetFile, filepath.Base(path)/*templateFile*/, *config)
			if err != nil {
				return fmt.Errorf("failed to parse templates for file %s: err %v", path, err)
			}
			targetFile.Close()
		}
		return nil
	})
	return err
}

func getTargetPath(sourcePath string, values *config.TemplaterConfig) string {
	trimmedSourcePath := strings.TrimPrefix(sourcePath, "templates/")
	if trimmedSourcePath == "Dockerfile.skipper-repo" {
		return strings.Replace(trimmedSourcePath, "skipper-repo", values.RepoName, 1) + "-build"
	} 
	// replace API_VERSION with the real api version
	versionedSourcePath := strings.Replace(trimmedSourcePath, apiVersionDir, values.API.Version, 1)
	return strings.Replace(versionedSourcePath, "gotmpl", "go", 1)
}
