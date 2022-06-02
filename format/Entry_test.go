package format_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/simonhayward/gkeepassxreader/format"
	"github.com/simonhayward/gkeepassxreader/keys"
)

var _ = Describe("Search", func() {

	var (
		entryService    *format.EntryServiceOp
		entryValues     []format.EntryValue
		entriesExpected []format.Entry
	)

	BeforeEach(func() {
		entryService = &format.EntryServiceOp{}
		entryValues = []format.EntryValue{
			format.EntryValue{
				PlainText: "Entry Number 1",
			},
			format.EntryValue{
				PlainText: "Entry Number 2",
			},
			format.EntryValue{
				PlainText: "QDKlkWFWQ6UJWbKyNPQs51Z8RcA5bYgU",
			},
			format.EntryValue{
				PlainText: "Entry Number 4",
			},
			format.EntryValue{
				PlainText: "A Title Repeated",
			},
			format.EntryValue{
				PlainText: "a title repeated",
			},
		}

		entriesExpected = []format.Entry{
			format.Entry{
				Title: &entryValues[0],
				UUID:  "u3uqCWOYDQ5UEpGGWoEjZLAMbWLzFkdC",
			},
			format.Entry{
				Title: &entryValues[1],
				UUID:  "wba62vPxT7aQfocEhJjyu3lVVRXN6GX7",
			},
			format.Entry{
				Title: &entryValues[2],
				UUID:  "qAcei1Z5fbcMTklrcikE4GccEYtlZx8p",
			},
			format.Entry{
				Title: &entryValues[3],
				UUID:  "QDKlkWFWQ6UJWbKyNPQs51Z8RcA5bYgU",
			},
			format.Entry{
				Title: &entryValues[4],
				UUID:  "YPNgv5b5mmWd1aAhNTofGdruRRefoqyJ",
			},
			format.Entry{
				Title: &entryValues[5],
				UUID:  "qXEZG88PvPiZd3abw1fBin8CL8fJV2aQ",
			},
		}
	})

	Context("when searching for a match by title", func() {
		It("succeeds and returns the correct match", func() {
			idx := entryService.Search("Entry Number 1", entriesExpected)
			Expect(idx).To(Equal(0))
			idx = entryService.Search("ENTRY number 1", entriesExpected)
			Expect(idx).To(Equal(0))
			idx = entryService.Search("Entry Number 2", entriesExpected)
			Expect(idx).To(Equal(1))
			idx = entryService.Search("entry number 2", entriesExpected)
			Expect(idx).To(Equal(1))
		})
	})

	Context("when searching for a match by uuid", func() {
		It("succeeds and returns the correct match", func() {
			idx := entryService.Search("u3uqCWOYDQ5UEpGGWoEjZLAMbWLzFkdC", entriesExpected)
			Expect(idx).To(Equal(0))
			idx = entryService.Search("wba62vPxT7aQfocEhJjyu3lVVRXN6GX7", entriesExpected)
			Expect(idx).To(Equal(1))
		})
	})

	Context("when searching for a match which does not exist", func() {
		It("returns an integer equal to the length of the entries", func() {
			idx := entryService.Search("Does not Exist", entriesExpected)
			Expect(idx).To(Equal(len(entriesExpected)))
		})
	})

	Context("when searching for a match uuid has precedence over title", func() {
		It("succeeds and returns the correct match", func() {
			idx := entryService.Search("QDKlkWFWQ6UJWbKyNPQs51Z8RcA5bYgU", entriesExpected)
			Expect(idx).To(Equal(3))
		})
	})

	Context("when searching for a match exact has precedence", func() {
		It("succeeds and returns the correct match", func() {
			idx := entryService.Search("A Title Repeated", entriesExpected)
			Expect(idx).To(Equal(4))

			idx = entryService.Search("a title repeated", entriesExpected)
			Expect(idx).To(Equal(5))
		})
	})

	Context("when searching for multiple possible lower cased matches", func() {
		It("succeeds and returns the first match", func() {
			idx := entryService.Search("A TITLE REPEATED", entriesExpected)
			Expect(idx).To(Equal(4))
		})
	})

	Context("when searching for an entry excluding historical entries", func() {
		It("does not return matching historical entries", func() {

			db, err := os.Open("test_data/HistoryTitle.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "password"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())

			entryService.XMLReader = reader.XMLReader
			entry, err := entryService.SearchByTerm("Mynew email address")
			Expect(err).ToNot(HaveOccurred())
			Expect(entry).Should(BeNil())
		})
	})

	Context("when searching for an entry including historical entries", func() {
		It("succeeds and returns the matching historical entry", func() {

			db, err := os.Open("test_data/HistoryTitle.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "password"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())

			entryService.HistoricalEntries = true
			entryService.XMLReader = reader.XMLReader
			entry, err := entryService.SearchByTerm("Mynew email address")
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.Historical).To(Equal(true))

			Expect(entry.Title.Protected).To(Equal(false))
			Expect(entry.Title.Data).To(Equal("Mynew email address"))
			Expect(entry.Title.PlainText).To(Equal("Mynew email address"))

			Expect(entry.Password.Data).To(Equal("DfjoK9i3kcH0JUggt6LX9Q=="))
			Expect(entry.Password.PlainText).To(Equal("3MAVouuiK2g6Qi5Q"))
			Expect(entry.Password.Protected).To(Equal(true))

			Expect(entry.URL.Data).To(Equal(""))
			Expect(entry.URL.PlainText).To(Equal(""))
			Expect(entry.URL.Protected).To(Equal(false))

			Expect(entry.Username.Data).To(Equal("test@test.com"))
			Expect(entry.Username.PlainText).To(Equal("test@test.com"))
			Expect(entry.Username.Protected).To(Equal(false))
		})
	})
})

var _ = Describe("Entries", func() {

	var (
		entryService *format.EntryServiceOp
		db           *os.File
	)

	BeforeEach(func() {
		entryService = &format.EntryServiceOp{}
	})

	AfterEach(func() {
		db.Close()
	})

	Context("when requesting a list of all entries", func() {
		It("succeeds and returns all entries", func() {

			db, err := os.Open("test_data/ProtectedStrings.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "masterpw"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())

			entryService.XMLReader = reader.XMLReader
			entryService.HistoricalEntries = true
			listEntries, err := entryService.List()
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
					Group:      "Protected",
					Title:      &entryValues[0],
					Username:   &entryValues[1],
					Password:   &entryValues[2],
					URL:        &entryValues[3],
					Notes:      &entryValues[4],
					UUID:       "a8370aa88afd3c4593ce981eafb789c8",
					Historical: false,
				},
				format.Entry{
					Group:      "Protected",
					Title:      &entryValues[5],
					Username:   &entryValues[6],
					Password:   &entryValues[7],
					URL:        &entryValues[8],
					Notes:      &entryValues[9],
					UUID:       "a8370aa88afd3c4593ce981eafb789c8",
					Historical: true,
				},
			}

			Expect(listEntries).To(Equal(expectedEntries))
		})
	})

	Context("when requesting a list of example entries", func() {
		It("succeeds and returns all entries", func() {

			db, err := os.Open("test_data/Example.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "password"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())

			entryService.XMLReader = reader.XMLReader
			listEntries, err := entryService.List()
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

			db, err := os.Open("test_data/Format200.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "a"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())

			entryService.XMLReader = reader.XMLReader
			listEntries, err := entryService.List()
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

	Context("when requesting a list of entries excluding historical entries", func() {
		It("succeeds and returns a single entry excluding the historical entries", func() {

			db, err := os.Open("test_data/History.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "password"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())

			entryService.XMLReader = reader.XMLReader
			listEntries, err := entryService.List()
			Expect(err).ToNot(HaveOccurred())

			Expect(len(listEntries)).To(Equal(1))
			Expect(listEntries[0].Group).To(Equal("Emails"))

			Expect(listEntries[0].Historical).To(Equal(false))

			Expect(listEntries[0].Title.Protected).To(Equal(false))
			Expect(listEntries[0].Title.Data).To(Equal("My email address"))
			Expect(listEntries[0].Title.PlainText).To(Equal("My email address"))

			Expect(listEntries[0].Password.Data).To(Equal("tzAeMlvTMXTm6Ebq2B6p+A=="))
			Expect(listEntries[0].Password.PlainText).To(Equal(""))
			Expect(listEntries[0].Password.Protected).To(Equal(true))

			Expect(listEntries[0].URL.Data).To(Equal(""))
			Expect(listEntries[0].URL.PlainText).To(Equal(""))
			Expect(listEntries[0].URL.Protected).To(Equal(false))

			Expect(listEntries[0].Username.Data).To(Equal("test@test.com"))
			Expect(listEntries[0].Username.PlainText).To(Equal("test@test.com"))
			Expect(listEntries[0].Username.Protected).To(Equal(false))
		})
	})

	Context("when requesting a list of entries including historical entries", func() {
		It("succeeds and returns all entries including the historical entries", func() {

			db, err := os.Open("test_data/History.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "password"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())

			entryService.XMLReader = reader.XMLReader
			entryService.HistoricalEntries = true

			listEntries, err := entryService.List()
			Expect(err).ToNot(HaveOccurred())

			Expect(len(listEntries)).To(Equal(4))
			Expect(listEntries[0].Group).To(Equal("Emails"))

			Expect(listEntries[0].Historical).To(Equal(false))
			Expect(listEntries[1].Historical).To(Equal(true))
			Expect(listEntries[2].Historical).To(Equal(true))
			Expect(listEntries[3].Historical).To(Equal(true))
		})
	})
})
