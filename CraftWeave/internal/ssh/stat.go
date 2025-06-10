package ssh

import (
	"fmt"

	"craftweave/internal/inventory"
)

// Stat checks whether the given path exists on the remote host.
// It returns a boolean indicating existence along with a CommandResult
// describing the check.
func Stat(h inventory.Host, path string) (bool, CommandResult) {
	cmd := fmt.Sprintf("test -e '%s'", path)
	res := RunShellCommand(h, cmd)

	exists := res.ReturnCode == 0
	res.ReturnCode = 0
	res.ReturnMsg = "OK"
	if exists {
		res.Output = fmt.Sprintf("%s exists", path)
	} else {
		res.Output = fmt.Sprintf("%s does not exist", path)
	}

	return exists, res
}
