package format_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/simonhayward/gkeepassxreader/core"
	"github.com/simonhayward/gkeepassxreader/entries"
	"github.com/simonhayward/gkeepassxreader/format"
	"github.com/simonhayward/gkeepassxreader/keys"
)

var _ = Describe("Databases", func() {

	var (
		db *os.File
	)

	AfterEach(func() {
		db.Close()
	})

	Context("when searching a protected strings test database with compression", func() {
		It("succeeds and returns the correct entry", func() {
			db, err := os.Open("test_data/ProtectedStrings.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "masterpw"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())
			Expect(reader.Db.CompressionAlgo).To(Equal(core.CompressionGzip))

			searchTerm := "Sample Entry"
			entry, err := entries.SearchByTerm(reader.XMLReader, searchTerm)
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.Group).To(Equal("Protected"))
			Expect(entry.UUID).To(Equal("a8370aa88afd3c4593ce981eafb789c8"))
			Expect(entry.Username.PlainText).To(Equal("Protected User Name"))
			Expect(entry.URL.PlainText).To(Equal("http://www.somesite.com/"))
			Expect(entry.Notes.PlainText).To(Equal("Notes"))
			Expect(entry.Password.PlainText).To(Equal("ProtectedPassword"))

		})
	})

	Context("when trying to open database whose protected stream key has been modified in the header", func() {
		It("returns an error", func() {
			db, err := os.Open("test_data/BrokenHeaderHash.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := ""
			_, err = format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when opening a databse with non ascii password without compression", func() {
		It("succeeds", func() {
			db, err := os.Open("test_data/NonAscii.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "\xce\x94\xc3\xb6\xd8\xb6"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())
			Expect(reader.Db.CompressionAlgo).To(Equal(core.CompressionNone))
		})
	})

	Context("when opening a databse with with compression", func() {
		It("succeeds", func() {
			db, err := os.Open("test_data/Compressed.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := ""
			masterKey := keys.NewCompositeKey()
			pk := &keys.PasswordKey{}
			pk.SetPassword(password)
			masterKey.AddKey(pk)

			reader, err := format.OpenDatabase(masterKey, db)
			Expect(err).ToNot(HaveOccurred())
			Expect(reader.Db.CompressionAlgo).To(Equal(core.CompressionGzip))
		})
	})

	Context("when opening a format 200 database", func() {
		It("succeeds and returns the correct entry", func() {
			db, err := os.Open("test_data/Format200.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "a"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())
			Expect(reader.Db.CompressionAlgo).To(Equal(core.CompressionGzip))

			searchTerm := "Sample Entry"
			entry, err := entries.SearchByTerm(reader.XMLReader, searchTerm)
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.Group).To(Equal("Format200"))
			Expect(entry.Title.PlainText).To(Equal("Sample Entry"))
			Expect(entry.Password.PlainText).To(Equal("Password"))
			Expect(entry.Username.PlainText).To(Equal("User Name"))
			Expect(entry.URL.PlainText).To(Equal("http://www.somesite.com/"))
			Expect(entry.Notes.PlainText).To(Equal("Notes"))

		})
	})

	Context("when opening a format 300 database", func() {
		It("succeeds", func() {
			db, err := os.Open("test_data/Format300.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "a"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())
			Expect(reader.Db.CompressionAlgo).To(Equal(core.CompressionGzip))

			searchTerm := "Sample Entry"
			entry, err := entries.SearchByTerm(reader.XMLReader, searchTerm)
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.Title.PlainText).To(Equal("Sample Entry"))
			Expect(entry.Password.PlainText).To(Equal("Password"))
			Expect(entry.Username.PlainText).To(Equal("User Name"))
			Expect(entry.URL.PlainText).To(Equal("http://www.somesite.com/"))
			Expect(entry.Notes.PlainText).To(Equal("Notes"))
		})
	})

	Context("when opening an example database", func() {
		It("succeeds", func() {
			db, err := os.Open("test_data/Example.kdbx")
			Expect(err).ToNot(HaveOccurred())

			password := "password"
			reader, err := format.OpenDatabase(keys.MasterKey(password, nil), db)
			Expect(err).ToNot(HaveOccurred())
			Expect(reader.Db.CompressionAlgo).To(Equal(core.CompressionGzip))

			searchTerm := "Sample Entry"
			entry, err := entries.SearchByTerm(reader.XMLReader, searchTerm)
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.UUID).To(Equal("640c38611c3ea4489ced361f54e43dbe"))
			Expect(entry.Title.PlainText).To(Equal("Sample Entry"))
			Expect(entry.Password.PlainText).To(Equal("Password"))
			Expect(entry.Username.PlainText).To(Equal("User Name"))
			Expect(entry.URL.PlainText).To(Equal("http://keepass.info/"))
			Expect(entry.Notes.PlainText).To(Equal("Notes"))
		})
	})
})
