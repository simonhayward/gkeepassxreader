package output

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/simonhayward/gkeepassxreader/format"
)

//Extract specific characters from string
func Extract(entry *format.Entry, chrs string) error {
	var buffer bytes.Buffer
	idxs := strings.Split(chrs, ",")
	for _, strIdx := range idxs {
		idx, err := strconv.Atoi(strIdx)
		if err != nil {
			return err
		}
		idx--
		buffer.WriteByte(entry.PlainTextPassword[idx])
	}
	entry.PlainTextPassword = buffer.String()

	return nil
}
