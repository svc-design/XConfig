package ssh

import (
	"fmt"
	"os/exec"

	"craftweave/internal/inventory"
)

func RunShellCommand(h inventory.Host, command string) {
	cmd := exec.Command("ssh",
		"-i", h.KeyFile,
		"-p", h.Port,
		fmt.Sprintf("%s@%s", h.User, h.Address),
		command,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s | FAILED | %v\n", h.Name, err)
	} else {
		fmt.Printf("%s | CHANGED | rc=0 >>\n%s\n", h.Name, string(output))
	}
}

