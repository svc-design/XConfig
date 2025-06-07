package ssh

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"

	"craftweave/internal/inventory"
)

// RenderTemplate renders the given template file with data and uploads it to the remote host
func RenderTemplate(h inventory.Host, src, dest string, data map[string]string) CommandResult {
	content, err := os.ReadFile(src)
	if err != nil {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("read template failed: %v", err),
		}
	}

	t, err := template.New("tpl").Parse(string(content))
	if err != nil {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("parse template failed: %v", err),
		}
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("execute template failed: %v", err),
		}
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	script := fmt.Sprintf("echo \"%s\" | base64 -d > %s", encoded, dest)

	return RunShellCommand(h, script)
}
