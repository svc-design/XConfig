# CraftWeave Design Overview

This document summarizes the current codebase layout and the main functions available in each module. The goal is to provide an easy reference when extending CraftWeave with new features (e.g. copy, command, package modules).

## Directory Structure

```
CraftWeave/
├── cmd/           # Cobra subcommands
├── core/          # Executor and parser logic
├── internal/      # SSH runner and inventory parser
├── example/       # Sample inventory, playbooks and roles
└── main.go        # CLI entry point
```

### Subcommands (`cmd/`)
- **root.go** – sets up the CLI and registers subcommands.
- **ansible.go** – ad-hoc task execution. Supports `shell` and `script` modules.
- **playbook.go** – runs a YAML playbook with `shell`, `script` and `template` tasks.
- **vault.go** – placeholder for future encryption features.
- **cmdb.go** – placeholder to export topology graphs.
- **plugin.go** – placeholder for loading plugins.

### Core Modules (`core/`)
- **executor/playbook.go** – iterates over plays and tasks and executes them using the SSH helpers. Supports roles by loading tasks from `roles/<role>/tasks/main.yaml`.
- **parser/parser.go** – defines YAML structures (`Play`, `Task`, `Template`) and parses playbooks.
- **cmdb/** – reserved for CMDB graph generation (not yet implemented).

### Internal Libraries (`internal/`)
- **inventory/inventory.go** – parses INI style inventory files and returns a list of `Host` structures.
- **ssh/runner.go** – runs remote shell commands via SSH with key or password authentication.
- **ssh/script.go** – uploads a local script using base64 and executes it remotely.
- **ssh/template.go** – renders a Go template and uploads the result to the remote host.
- **ssh/formatter.go** – utilities for aggregated output.
- **ssh/result.go** – defines the `CommandResult` struct used across the executor.

