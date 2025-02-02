package gomod

import (
	"fmt"
	"os"
)

func Update() error {
	file, err := os.OpenFile("go.mod", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open go.mod for update: %v", err)
	}
	defer file.Close()

	kmmLine := "\nreplace github.com/rh-ecosystem-edge/kernel-module-management => github.com/rh-ecosystem-edge/kernel-module-management release-2.2"
	_, err = file.WriteString(kmmLine)
	if err != nil {
		return fmt.Errorf("failed to append line <%s> to file go.mod: %v", kmmLine, err)
	}

	ctrlLine := "\nreplace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.17.3"
	_, err = file.WriteString(ctrlLine)
	if err != nil {
		return fmt.Errorf("failed to append line <%s> to file go.mod: %v", ctrlLine, err)
	}

	ocmLine := "\nreplace open-cluster-management.io/api => open-cluster-management.io/api v0.13.0"
	_, err = file.WriteString(ocmLine)
        if err != nil {
                return fmt.Errorf("failed to append line <%s> to file go.mod: %v", ocmLine, err)
        }

	uberLine := "\nreplace go.uber.org/mock => go.uber.org/mock v0.4.0"
	_, err = file.WriteString(uberLine)
        if err != nil {
                return fmt.Errorf("failed to append line <%s> to file go.mod: %v", uberLine, err)
        }
	
	return nil
}
