package ssh

import (
	"fmt"
	"os/exec"
)

func Ping(host string) error {
	cmd := exec.Command("ping", "-c", "1", host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ping failed: %s\n%s", err, output)
	}
	fmt.Printf("Ping to %s succeeded.\n", host)
	return nil
}
