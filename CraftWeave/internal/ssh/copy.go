package ssh

import (
	"encoding/base64"
	"fmt"
	"os"

	"craftweave/internal/inventory"
)

// CopyFile uploads a local file to a remote host with optional permission mode.
func CopyFile(h inventory.Host, src, dest, mode string) CommandResult {
	content, err := os.ReadFile(src)
	if err != nil {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("read file failed: %v", err),
		}
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	script := fmt.Sprintf("echo \"%s\" | base64 -d > %s", encoded, dest)
	if mode != "" {
		script += fmt.Sprintf(" && chmod %s %s", mode, dest)
	}

	return RunShellCommand(h, script)
}
