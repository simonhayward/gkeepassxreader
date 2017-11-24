package output

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/simonhayward/gkeepassxreader/format"
)

const whitespace = " \t"

//Extract specific characters from string
func Extract(entry *format.Entry, chrs string) error {
	var buffer bytes.Buffer

	chrs = strings.Trim(chrs, whitespace)
	chrs = strings.Trim(chrs, ",")
	idxs := strings.Split(chrs, ",")

	for _, strIdx := range idxs {
		strIdx = strings.Trim(strIdx, whitespace)

		idx, err := strconv.Atoi(strIdx)
		if err != nil {
			return errors.Wrap(err, "failure to convert string to int")
		}

		if idx <= 0 || idx > len(entry.Password.PlainText) {
			return errors.Errorf("index out of range: %d max value allowed is: %d", idx, len(entry.Password.PlainText))
		}

		idx--

		buffer.WriteByte(entry.Password.PlainText[idx])
	}

	entry.Password.PlainText = buffer.String()

	return nil
}
