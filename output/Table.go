package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/simonhayward/gkeepassxreader/format"
)

//Data header/data
type Data struct {
	Header []string
	Data   [][]string
}

//NewDefaults entries
func NewDefaults() *Data {
	return &Data{
		Header: []string{"UUID", "Group", "Title", "Username", "URL", "Notes"},
	}
}

//Entries fields to display
func (d *Data) Entries(entries []format.Entry) {
	for _, entry := range entries {

		var url, notes string
		if entry.Notes != nil {
			notes = entry.Notes.PlainText
		}
		if entry.URL != nil {
			url = entry.URL.PlainText
		}

		d.Data = append(d.Data, []string{
			entry.UUID,
			entry.Group,
			entry.Title.PlainText,
			entry.Username.PlainText,
			url,
			notes,
		})
	}
}

//Table entries in ascii table
func Table(dataHeader []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(dataHeader)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
