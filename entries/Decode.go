package entries

import (
	"github.com/pkg/errors"
	"github.com/simonhayward/gkeepassxreader/format"
)

func decodeEntryValue(xmlReader *format.KeePass2XmlReader, eValue *format.EntryValue) error {
	plaintext, err := xmlReader.KeePass2RandomStream.Process(eValue.RandomOffset, []byte(eValue.CipherText))
	if err != nil {
		return errors.Wrap(err, "entry value decode failed")
	}
	eValue.PlainText = string(plaintext)
	return nil
}

func decodeEntries(xmlReader *format.KeePass2XmlReader, entries []format.Entry, decodePassword bool) error {
	for idx := range entries {
		var entryValues []*format.EntryValue

		entryValues = []*format.EntryValue{
			entries[idx].Notes,
			entries[idx].URL,
			entries[idx].Username,
			entries[idx].Title,
		}

		if decodePassword {
			entryValues = append(entryValues, entries[idx].Password)
		}

		for _, ev := range entryValues {
			if ev != nil && ev.Protected {
				err := decodeEntryValue(xmlReader, ev)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
