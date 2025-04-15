package inventory

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Host struct {
	Name    string
	Address string
	User    string
}

func LoadHosts(inventoryPath string) ([]Host, error) {
	file, err := os.Open(inventoryPath)
	if err != nil {
		return nil, fmt.Errorf("could not open inventory: %w", err)
	}
	defer file.Close()

	var hosts []Host
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Simple parser for "host ansible_host=x.x.x.x ansible_user=y"
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		host := Host{Name: parts[0]}
		for _, part := range parts[1:] {
			if strings.HasPrefix(part, "ansible_host=") {
				host.Address = strings.TrimPrefix(part, "ansible_host=")
			}
			if strings.HasPrefix(part, "ansible_user=") {
				host.User = strings.TrimPrefix(part, "ansible_user=")
			}
		}
		hosts = append(hosts, host)
	}

	return hosts, scanner.Err()
}
