package entries_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/simonhayward/gkeepassxreader/entries"
	"github.com/simonhayward/gkeepassxreader/format"
	"github.com/simonhayward/gkeepassxreader/keys"
)

var _ = Describe("Entries", func() {

	var (
		db *os.File
	)

	AfterEach(func() {
		db.Close()
	})

	Context("when requesting a list of all entries", func() {
		It("succeeds and returns all entries", func() {

			db, err := os.Open("../format/test_data/ProtectedStrings.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "masterpw"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())

			listEntries, err := entries.List(reader.XMLReader)
			Expect(err).ToNot(HaveOccurred())

			entryValues := []format.EntryValue{
				format.EntryValue{
					Data:         "Sample Entry",
					Protected:    false,
					PlainText:    "Sample Entry",
					RandomOffset: 0,
					CipherText:   nil,
				},
				format.EntryValue{
					Data:         "Protected User Name",
					Protected:    false,
					PlainText:    "Protected User Name",
					RandomOffset: 0,
					CipherText:   nil,
				},
				format.EntryValue{
					Data:         "ZapxLOhEQgWELdKbLMgCBp4=",
					Protected:    true,
					PlainText:    "",
					RandomOffset: 0,
					CipherText:   []byte{101, 170, 113, 44, 232, 68, 66, 5, 132, 45, 210, 155, 44, 200, 2, 6, 158},
				},
				format.EntryValue{
					Data:         "http://www.somesite.com/",
					Protected:    false,
					PlainText:    "http://www.somesite.com/",
					RandomOffset: 0,
					CipherText:   nil,
				},
				format.EntryValue{
					Data: "Notes", Protected: false, PlainText: "Notes", RandomOffset: 0, CipherText: nil,
				},
				format.EntryValue{
					Data:         "Sample Entry",
					Protected:    false,
					PlainText:    "Sample Entry",
					RandomOffset: 0,
					CipherText:   nil,
				},
				format.EntryValue{
					Data:         "Protected User Name",
					Protected:    false,
					PlainText:    "Protected User Name",
					RandomOffset: 0,
					CipherText:   nil,
				},
				format.EntryValue{
					Data:         "u8PhlyS8ep0VjyRUP8Su88c=",
					Protected:    true,
					PlainText:    "",
					RandomOffset: 17,
					CipherText:   []byte{187, 195, 225, 151, 36, 188, 122, 157, 21, 143, 36, 84, 63, 196, 174, 243, 199},
				},
				format.EntryValue{
					Data:         "http://www.somesite.com/",
					Protected:    false,
					PlainText:    "http://www.somesite.com/",
					RandomOffset: 0,
					CipherText:   nil,
				},
				format.EntryValue{
					Data: "Notes", Protected: false, PlainText: "Notes", RandomOffset: 0, CipherText: nil,
				},
			}

			expectedEntries := []format.Entry{
				format.Entry{
					Group:    "Protected",
					Title:    &entryValues[0],
					Username: &entryValues[1],
					Password: &entryValues[2],
					URL:      &entryValues[3],
					Notes:    &entryValues[4],
					UUID:     "a8370aa88afd3c4593ce981eafb789c8",
				},
				format.Entry{
					Group:    "Protected",
					Title:    &entryValues[5],
					Username: &entryValues[6],
					Password: &entryValues[7],
					URL:      &entryValues[8],
					Notes:    &entryValues[9],
					UUID:     "a8370aa88afd3c4593ce981eafb789c8",
				},
			}

			Expect(listEntries).To(Equal(expectedEntries))
		})
	})
})
