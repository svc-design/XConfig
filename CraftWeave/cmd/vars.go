// cmd/vars.go
package cmd

var (
	AggregateOutput bool              // --aggregate / -A
	CheckMode       bool              // --check / -C
	InventoryPath   string            // --inventory / -i
	ExtraVars       map[string]string // --extra-vars / -e
)
