package output

import (
	"fmt"
	"sort"
	"strings"

	"craftweave/internal/ssh"
)

// Collector defines interface for task result collection
// Flush prints or handles aggregated results if needed
// Flush may be a no-op for simple collectors
//go:generate mockgen -destination=collector_mock.go -package=output . Collector

type Collector interface {
	Collect(res ssh.CommandResult)
	Flush()
}

type StdoutCollector struct{}

func (StdoutCollector) Collect(res ssh.CommandResult) {
	fmt.Printf("%s | %s | rc=%d >>\n%s\n", res.Host, res.ReturnMsg, res.ReturnCode, res.Output)
}

func (StdoutCollector) Flush() {}

type AggregateCollector struct {
	results []ssh.CommandResult
}

func (a *AggregateCollector) Collect(res ssh.CommandResult) {
	a.results = append(a.results, res)
}

func (a *AggregateCollector) Flush() {
	grouped := make(map[string][]ssh.CommandResult)
	for _, r := range a.results {
		key := fmt.Sprintf("%s-%d-%s", r.ReturnMsg, r.ReturnCode, r.Output)
		grouped[key] = append(grouped[key], r)
	}
	for _, grp := range grouped {
		if len(grp) == 0 {
			continue
		}
		hosts := make([]string, len(grp))
		for i, r := range grp {
			hosts[i] = r.Host
		}
		sort.Strings(hosts)
		fmt.Printf("%s | %s | rc=%d >>\n%s\n", strings.Join(hosts, ","), grp[0].ReturnMsg, grp[0].ReturnCode, grp[0].Output)
	}
	a.results = nil
}
