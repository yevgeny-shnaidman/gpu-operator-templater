package operator_sdk

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed binaries/operator-sdk
var sdkBinary embed.FS

func InitializeRepo() error {
	// Extract the operator-sdk binary
    	sdkPath := filepath.Join(os.TempDir(), "operator-sdk")
    	err := os.WriteFile(sdkPath, readBinaryFile("binaries/operator-sdk"), 0755)
    	if err != nil {
		fmt.Errorf("failed to  extract operator-sdk: %v", err)
    	}

	err = runInit(sdkPath)
	if err != nil {
		return fmt.Errorf("failed to run runInit: %v", err)
	}

	return runCreateAPI(sdkPath)
}

func runInit(operatorSDKPath string) error {
	// initialize the repo
        cmd := exec.Command(operatorSDKPath, "init", "--domain=sigs.x-k8s.io", "--repo=github.com/yevgeny-shnaidman/test-gpu-operator", "--skip-go-version-check")
        output, err := cmd.CombinedOutput()
        if err != nil {
                return fmt.Errorf("failed running operator-sdk init, output %s: %v", output, err)
        }
	return nil
}

func runCreateAPI(operatorSDKPath string) error {
	cmd := exec.Command(operatorSDKPath, "create", "api", "--controller=false",  "--group=amd", "--kind=DeviceConfig", "--resource=true", "--version=v1alpha1")
	output, err := cmd.CombinedOutput()
        if err != nil {
		return fmt.Errorf("failed running operator-sdk create api, output %s: %v", output, err)
        }
	return nil
}

// Hhelper function to read binary file
func readBinaryFile(name string) []byte {
    data, err := fs.ReadFile(sdkBinary, name)
    if err != nil {
        fmt.Println("Error reading binary file:", err)
        os.Exit(1)
    }
    return data
}
