package ssh

import (
	"strings"

	"craftweave/internal/inventory"
)

// GatherFacts collects basic system information from the target host.
// Currently it returns ansible_architecture and ansible_pkg_mgr.
func GatherFacts(h inventory.Host) map[string]string {
	facts := make(map[string]string)

	archRes := RunShellCommand(h, "uname -m")
	if archRes.ReturnCode == 0 {
		facts["ansible_architecture"] = strings.TrimSpace(archRes.Output)
	} else {
		facts["ansible_architecture"] = "unknown"
	}

	pkgCmd := "if command -v apt-get >/dev/null 2>&1; then echo apt; " +
		"elif command -v yum >/dev/null 2>&1; then echo yum; else echo unknown; fi"
	pkgRes := RunShellCommand(h, pkgCmd)
	if pkgRes.ReturnCode == 0 {
		facts["ansible_pkg_mgr"] = strings.TrimSpace(pkgRes.Output)
	} else {
		facts["ansible_pkg_mgr"] = "unknown"
	}

	return facts
}
