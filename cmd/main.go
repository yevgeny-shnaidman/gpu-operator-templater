package main

import (
    "fmt"
    "os"
    "github.com/yevgeny-shnaidman/gpu-operator-template/internal/operator_sdk"
)


func main() {
    err := operator_sdk.InitializeRepo()
    if err != nil {
	fmt.Errorf("failed to initialize repo with operator-sdk, error %s\n", err)
	os.Exit(-1)
    }
}
