package operator_sdk

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yevgeny-shnaidman/gpu-operator-template/internal/config"
)

//go:embed binaries/operator-sdk
var sdkBinary embed.FS

func InitializeRepo(config *config.TemplaterConfig) error {
	// Extract the operator-sdk binary
	sdkPath := filepath.Join(os.TempDir(), "operator-sdk")
	err := os.WriteFile(sdkPath, readBinaryFile("binaries/operator-sdk"), 0755)
	if err != nil {
		fmt.Errorf("failed to  extract operator-sdk: %v", err)
	}

	err = runInit(sdkPath, config)
	if err != nil {
		return fmt.Errorf("failed to run runInit: %v", err)
	}
	err = runCreateAPI(sdkPath, config)
	if err != nil {
		return fmt.Errorf("failed to run runCreateAPI: %v", err)
	}
	return cleanup()
}

func runInit(operatorSDKPath string, config *config.TemplaterConfig) error {
	// initialize the repo
	params := []string{
		"init",
		"--domain=" + config.API.Domain,
		"--repo=" + config.API.CodeRepo,
		"--skip-go-version-check",
	}
	cmd := exec.Command(operatorSDKPath, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed running operator-sdk init, output %s: %v", output, err)
	}
	return nil
}

func runCreateAPI(operatorSDKPath string, config *config.TemplaterConfig) error {
	params := []string{
		"create",
		"api",
		"--controller=false",
		"--group=" + config.API.Group,
		"--kind=DeviceConfig",
		"--resource=true",
		"--version=" + config.API.Version,
	}
	cmd := exec.Command(operatorSDKPath, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed running operator-sdk create api, output %s: %v", output, err)
	}
	return nil
}

func cleanup() error {
	if err := os.Remove("cmd/main.go");err != nil {
		return fmt.Errorf("failed to remove cmd/main.go")
	}
	if err := os.Remove(".golangci.yml");err != nil {
                return fmt.Errorf("failed to remove .golangci.yml")
        }
	if err := os.RemoveAll("bin");err != nil {
		return fmt.Errorf("failed to remove bin directory")
	}

	return os.RemoveAll("test")
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
