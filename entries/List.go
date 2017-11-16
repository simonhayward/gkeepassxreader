package entries

import (
	"github.com/simonhayward/gkeepassxreader/format"
)

//List all entries
func List(xmlReader *format.KeePass2XmlReader) ([]format.Entry, error) {
	entries := []format.Entry{}
	randomBytesOffset := 0
	if err := xmlReader.ReadGroups(&entries, xmlReader.KeePass2XmlFile.Root.Groups, &randomBytesOffset); err != nil {
		return nil, err
	}

	return entries, nil
}
