package ssh

import (
	"fmt"
	"os/exec"

	"craftweave/internal/inventory"
)

// RunShellCommand 通过 SSH 执行命令，并返回 CommandResult
func RunShellCommand(h inventory.Host, command string) CommandResult {
	cmd := exec.Command("ssh",
		"-i", h.KeyFile,
		"-p", h.Port,
		fmt.Sprintf("%s@%s", h.User, h.Address),
		command,
	)
	output, err := cmd.CombinedOutput()
	result := CommandResult{
		Host:   h.Name,
		Output: string(output),
	}
	if err != nil {
		result.ReturnMsg = "FAILED"
		result.ReturnCode = 1
	} else {
		result.ReturnMsg = "CHANGED"
		result.ReturnCode = 0
	}
	return result
}

