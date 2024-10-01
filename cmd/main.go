package main

import (
	"flag"
	"fmt"
	"github.com/yevgeny-shnaidman/gpu-operator-template/internal/code_templates"
	"github.com/yevgeny-shnaidman/gpu-operator-template/internal/operator_sdk"
	"os"
)

func main() {
	valuesFilePath := flag.String("f", "", "files with values constituations")
	flag.Parse()
	err := operator_sdk.InitializeRepo()
	if err != nil {
		fmt.Printf("failed to initialize repo with operator-sdk, error %s\n", err)
		os.Exit(-1)
	}

	err = code_templates.RunTemplates(*valuesFilePath)
	if err != nil {
		fmt.Printf("failed to run templates and create the files in the target repo, error %s\n", err)
		os.Exit(-1)
	}
}
