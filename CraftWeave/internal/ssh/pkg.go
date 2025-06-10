package ssh

import (
	"fmt"

	"craftweave/internal/inventory"
)

// InstallAptPackage installs a package using apt-get install
func InstallAptPackage(h inventory.Host, name string) CommandResult {
	if name == "" {
		return CommandResult{Host: h.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "package name required"}
	}
	cmd := fmt.Sprintf("sudo apt-get update -y && sudo apt-get install -y %s", name)
	return RunShellCommand(h, cmd)
}

// InstallYumPackage installs a package using yum install
func InstallYumPackage(h inventory.Host, name string) CommandResult {
	if name == "" {
		return CommandResult{Host: h.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "package name required"}
	}
	cmd := fmt.Sprintf("sudo yum install -y %s", name)
	return RunShellCommand(h, cmd)
}
