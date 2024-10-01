package code_templates

import (
	"embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplatesData struct {
	Vendor      string `yaml:"vendor"`
	Repo        string `yaml:"repo"`
	APIVersion  string `yaml:"apiversion"`
	Group       string `yaml:"group"`
	PCIVendorID string `yaml:"pcivendorid"`
}

var (
	//go:embed templates
	templatesFS embed.FS
	tmpl        = template.Must(
		template.ParseFS(templatesFS, "templates/*/*.gotmpl", "templates/*/*/*.gotmpl"),
	)
)

func RunTemplates(valuesFilePath string) error {
	var substValues TemplatesData
	yamlFile, err := os.ReadFile(valuesFilePath)
	if err != nil {
		return fmt.Errorf("failed to access values file %s: error %v", valuesFilePath, err)
	}
	err = yaml.Unmarshal(yamlFile, &substValues)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the values file %s into struct: error %v", valuesFilePath, err)
	}
	err = fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// get the path in the target
		targetPath := strings.Replace(strings.TrimPrefix(path, "templates/"), "gotmpl", "go", 1)
		if d.IsDir() {
			os.Mkdir(targetPath, 0750)
		} else {
			targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("failed to create file %s: error %v", targetPath, err)
			}
			templateFile := filepath.Base(path)
			err = tmpl.ExecuteTemplate(targetFile, templateFile, substValues)
			if err != nil {
				return fmt.Errorf("failed to parse templates for file %s: err %v", path, err)
			}
			targetFile.Close()
		}
		return nil
	})
	return err
}
