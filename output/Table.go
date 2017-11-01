package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

//Table entries in ascii table
func Table(dataHeader []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(dataHeader)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
