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

	Context("when requesting a list of example entries", func() {
		It("succeeds and returns all entries", func() {

			db, err := os.Open("../format/test_data/Example.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "password"
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
					Data:         "User Name",
					Protected:    false,
					PlainText:    "User Name",
					RandomOffset: 0,
					CipherText:   nil,
				},
				format.EntryValue{
					Data:         "xzlvdrx7Cd8=",
					Protected:    true,
					PlainText:    "",
					RandomOffset: 0,
					CipherText:   []byte{199, 57, 111, 118, 188, 123, 9, 223},
				},
				format.EntryValue{
					Data:         "http://keepass.info/",
					Protected:    false,
					PlainText:    "http://keepass.info/",
					RandomOffset: 0,
					CipherText:   nil,
				},
				format.EntryValue{
					Data: "Notes", Protected: false, PlainText: "Notes", RandomOffset: 0, CipherText: nil,
				},
				format.EntryValue{
					Data:         "Sample Entry #2",
					Protected:    false,
					PlainText:    "Sample Entry #2",
					RandomOffset: 0,
					CipherText:   nil,
				},
				format.EntryValue{
					Data:         "Michael321",
					Protected:    false,
					PlainText:    "Michael321",
					RandomOffset: 0,
					CipherText:   nil,
				},
				format.EntryValue{
					Data:         "cQPYAOM=",
					Protected:    true,
					PlainText:    "",
					RandomOffset: 8,
					CipherText:   []byte{113, 3, 216, 0, 227},
				},
				format.EntryValue{
					Data:         "http://keepass.info/help/kb/testform.html",
					Protected:    false,
					PlainText:    "http://keepass.info/help/kb/testform.html",
					RandomOffset: 0,
					CipherText:   nil,
				},
			}

			expectedEntries := []format.Entry{
				format.Entry{
					Group:    "example",
					Title:    &entryValues[0],
					Username: &entryValues[1],
					Password: &entryValues[2],
					URL:      &entryValues[3],
					Notes:    &entryValues[4],
					UUID:     "640c38611c3ea4489ced361f54e43dbe",
				},
				format.Entry{
					Group:    "example",
					Title:    &entryValues[5],
					Username: &entryValues[6],
					Password: &entryValues[7],
					URL:      &entryValues[8],
					Notes:    nil,
					UUID:     "db8e52f8c86d7d468ecd53d4c2fe0a31",
				},
			}

			Expect(listEntries).To(Equal(expectedEntries))
		})
	})

	Context("when requesting a list of entries which are protected", func() {
		It("succeeds and returns all entries", func() {

			db, err := os.Open("../format/test_data/Format200.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "a"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())

			listEntries, err := entries.List(reader.XMLReader)
			Expect(err).ToNot(HaveOccurred())

			Expect(listEntries[0].Group).To(Equal("Format200"))

			Expect(listEntries[0].Title.Protected).To(Equal(false))
			Expect(listEntries[0].Title.Data).To(Equal("Sample Entry"))
			Expect(listEntries[0].Title.PlainText).To(Equal("Sample Entry"))

			Expect(listEntries[0].Password.Data).To(Equal("Password"))
			Expect(listEntries[0].Password.PlainText).To(Equal("Password"))
			Expect(listEntries[0].Password.Protected).To(Equal(false))

			Expect(listEntries[0].URL.Data).To(Equal("YtgAIKYL4ggH0nzaFP4srWV+GC8A1B0I"))
			Expect(listEntries[0].URL.PlainText).To(Equal("http://www.somesite.com/"))
			Expect(listEntries[0].URL.Protected).To(Equal(true))

			Expect(listEntries[0].Username.Data).To(Equal("VK3dahCqvgR8"))
			Expect(listEntries[0].Username.PlainText).To(Equal("User Name"))
			Expect(listEntries[0].Username.Protected).To(Equal(true))
		})
	})
})
