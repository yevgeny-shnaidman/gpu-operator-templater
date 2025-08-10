package main

import (
	"flag"
	"fmt"
	"github.com/yevgeny-shnaidman/gpu-operator-template/internal/code_templates"
	"github.com/yevgeny-shnaidman/gpu-operator-template/internal/config"
	"github.com/yevgeny-shnaidman/gpu-operator-template/internal/gomod"
	"github.com/yevgeny-shnaidman/gpu-operator-template/internal/operator_sdk"
	"os"
)

func main() {
	valuesFilePath := flag.String("f", "", "files with values constituations")
	flag.Parse()

	configData, err := config.InitConfigData(*valuesFilePath)
	if err != nil {
		fmt.Printf("failed to initialize configuration: %s\n", err)
		os.Exit(-1)
	}
	err = operator_sdk.InitializeRepo(configData)
	if err != nil {
		fmt.Printf("failed to initialize repo with operator-sdk: %s\n", err)
		os.Exit(-1)
	}
	err = code_templates.RunTemplates(configData)
	if err != nil {
		fmt.Printf("failed to run templates and create the files in the target repo, error %s\n", err)
		os.Exit(-1)
	}

	err = gomod.Update()
	if err != nil {
		fmt.Printf("failed to update the go.mod, error %s\n", err)
		os.Exit(-1)
	}
}
