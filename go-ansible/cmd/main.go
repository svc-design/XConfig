package main

import (
	"flag"
	"fmt"
	"os"

	"go-ansible/inventory"
	"go-ansible/ssh"
)

func main() {
	inventoryFile := flag.String("i", "", "Path to the inventory file")
	module := flag.String("m", "", "Module to run (e.g., 'ping')")
	group := flag.String("g", "all", "Host group to target")
	flag.Parse()

	if *inventoryFile == "" || *module == "" {
		fmt.Println("Usage: go-ansible -i <inventory-file> -m <module>")
		os.Exit(1)
	}

	hosts, err := inventory.LoadHosts(*inventoryFile)
	if err != nil {
		fmt.Printf("Failed to load inventory: %v\n", err)
		os.Exit(1)
	}

	if *module == "ping" {
		for _, host := range hosts {
			if err := ssh.Ping(host.Address); err != nil {
				fmt.Printf("Failed to ping host %s: %v\n", host.Name, err)
			}
		}
	} else {
		fmt.Printf("Module %s is not supported.\n", *module)
		os.Exit(1)
	}
}
