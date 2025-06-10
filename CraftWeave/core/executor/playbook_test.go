package executor_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"craftweave/core/executor"
	"craftweave/core/parser"
)

// helper to load example playbook from path relative to this package
func loadPlaybook(t *testing.T, relPath string) []parser.Play {
	t.Helper()
	path := filepath.Join("..", "..", "example", relPath)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", relPath, err)
	}
	var plays []parser.Play
	if err := yaml.Unmarshal(data, &plays); err != nil {
		t.Fatalf("failed to parse playbook %s: %v", relPath, err)
	}
	return plays
}

func TestParseExampleRun(t *testing.T) {
	plays := loadPlaybook(t, "run_example")
	if len(plays) != 1 {
		t.Fatalf("expected 1 play, got %d", len(plays))
	}
	if len(plays[0].Tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(plays[0].Tasks))
	}
	if plays[0].Tasks[0].Shell == "" {
		t.Errorf("first task should be shell")
	}
	if plays[0].Tasks[1].Script == "" {
		t.Errorf("second task should be script")
	}
	if plays[0].Tasks[2].Template == nil {
		t.Errorf("third task should be template")
	}
}

func TestExecutePlaybookCheckMode(t *testing.T) {
	plays := loadPlaybook(t, "run_example")
	inventory := filepath.Join("..", "..", "example", "inventory")
	// enable check mode to avoid real SSH connections
	executor.CheckMode = true
	executor.AggregateOutput = false
	executor.ExecutePlaybook(plays, inventory, filepath.Join("..", "..", "example"), nil)
}

func TestExampleScripts(t *testing.T) {
	scripts := []string{"echo.sh", "id.sh", "nproc.sh", "uname.sh"}
	for _, s := range scripts {
		path := filepath.Join("..", "..", "example", "scripts", s)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("script %s missing: %v", s, err)
		}
		if err := os.Chmod(path, 0755); err != nil { // ensure executable
			t.Fatalf("chmod failed for %s: %v", s, err)
		}
		cmd := exec.Command("bash", path)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("script %s failed: %v, output: %s", s, err, string(out))
		}
		if len(out) == 0 {
			t.Fatalf("script %s produced no output", s)
		}
	}
}
