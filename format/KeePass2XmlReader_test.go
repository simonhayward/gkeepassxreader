package format_test

import (
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/simonhayward/gkeepassxreader/format"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Xml", func() {

	Describe("Unmarshal XML", func() {

		var (
			db *os.File
			v  format.KeePass2XmlFile
		)

		BeforeEach(func() {
			v = format.KeePass2XmlFile{}
		})

		AfterEach(func() {
			db.Close()
		})

		Context("when given valid xml", func() {
			It("it returns no error and sets the first groups entry values", func() {

				xmlBody, err := ioutil.ReadFile("test_data/XmlBody.xml")
				Expect(err).ToNot(HaveOccurred())

				err = xml.Unmarshal(xmlBody, &v)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v.Root.Groups)).To(Equal(1))
				Expect(v.Root.Groups[0].UUID).To(Equal("96EgMmATevhdT3G8zV4eaA=="))
				Expect(len(v.Root.Groups[0].Entry[0].StringEntry)).To(Equal(5))

				for _, e := range v.Root.Groups[0].Entry[0].StringEntry {

					switch e.Key {
					case "Password":
						Expect(e.Value.Protected).To(Equal("True"))
						Expect(e.Value.Data).To(Equal("WxH8XrA5iZ/GUdw="))
					case "Title":
						Expect(e.Value.Protected).To(Equal(""))
						Expect(e.Value.Data).To(Equal("my secret entry"))
					case "Notes":
						Expect(e.Value.Protected).To(Equal(""))
						Expect(e.Value.Data).To(Equal(""))
					case "URL":
						Expect(e.Value.Protected).To(Equal(""))
						Expect(e.Value.Data).To(Equal(""))
					case "UserName":
						Expect(e.Value.Protected).To(Equal(""))
						Expect(e.Value.Data).To(Equal("my username"))
					}
				}

				Expect(v.Meta.HeaderHash).To(Equal("xbYFxGYATFmL2Z8Jxj7mRDTMYEQyQH9/yixknbuJSWw="))
			})
		})

		Context("when given valid xml group with historical entries", func() {
			It("returns no error and sets the entries historical values", func() {
				xmlBody, err := ioutil.ReadFile("test_data/XmlGroupBody.xml")
				Expect(err).ToNot(HaveOccurred())
				err = xml.Unmarshal(xmlBody, &v)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(v.Root.Groups)).To(Equal(1))
				Expect(len(v.Root.Groups[0].Groups)).To(Equal(1))
				Expect(len(v.Root.Groups[0].Entry[0].HistoryEntries)).To(Equal(4))

				historyEntriesPasswords := []string{
					"",
					"0i3rbm8tHxSow0paE7mo9FqN0kP21ThHUuiN8Cy2fbTBtJbu5fB+UgzKK3l6SWS8KxfghKORmoEjwCgurzBOEqvvaohmvD8A1wqcLg2+0V6fkvXgUeOXWD03M+qrb7oAx3bm3WZZl4rT/2k9Zx6uKWqFPeip1LYSvzrW4tbvRCF+tsAyBcaW/M5vuY4AD+bKep+iJtT4+OGhRKRGfOi5J8YQOOFjrkE5Ch2v9QqT8f3d0pv+jJPZZkDU0AsKncxgmug=",
					"j/4iLP+KgA1/R18V6+V1X4ajxjwxshwgiiCDqHDxIHCdxP2FvLfNXUuUP0Jm+CbmKv6RhGiai7STMScb8oMnZERFE8gVOUYH4wvBLw0Q55RNWBLlffIOLJiIzJaIscsKLEw54tOoMLoe6BQQXE7Wfm+WcYkSkZREMTO/Ar11HhdENvhGqoPz2d6yQLTlobAlNUD1n7lxbFySz4wUFKS6CBwSpEXy9xQ43lpMt23emHv0ITVDP3SpTWG/lFoInfpYmWPiZs0wjtET6ziJTzdIK7FoYgYYsNI31jVJKUs+lgJYTfY+7BNFSkfmtX5KkvrbLhMtfNN6YECd3YpzIP47N5tAa1x9fuXVSyrlc0sSBjLeLZ3xSFH6rUy9o9c4Nu2UUJqQWPEdxW/MeL0nUWuS/Vm+0B0HS5/0eT3YkzHthHHbsg436RhIa97Z7vR6Gik3legFR+JFmosA+gWb8QAlxEC59kG4uS5zS/gXfmix7TCvj297kPW7TUxWmA==",
					"Y5V5fMzlOexbeeheQv/W5+oqoEqc9p6fYUn7x5Otvnj12E9vXFXqoNGAtN1w/cmjMZzLE4R4kh5C6yCXtllk1MLxUwnMXDY5LNY7BvZLgP7pN7IckOMUoItXJIYABAND1eApVVgaFqtF74JH8MwLMNmrAcEzjn7CplXRlzjTI48wPOvGHmfqV+ah3kfdgwZ56u+qTr6CcCOTMSkVpHRBKUaNFxFNK4j5LVRtr79VOa1FnwAAjI5tzMhfyapeO+dNP50AUvrGfQVV2G0ORhdng0bXmWtYnLG1RacRVJGMO5tJRJ1xd/VXYxv3LeYKKBauaby4pJx4dtnZyPr1/Pb3olqEDnlYExcHVd+2f0SAcFZr2/7g173w3QMyjOhST23mMqYut3BorTJVLgINUBBUwn3BX3ELnx0akcBOpxOTcZOBcKJr3Mtc2AndWGFBqAlSQ8v5Q6QoU7IpqAAilxTldRASQyFYcRlRK4Dgw0/CRKwKigOFl5PUYFdtYSn7cB3yH3CArviBLxQ3xs40zeeGP7Lnit0YvAJQ+o4ka+C9MyMiEha/mW+LmBxpCmloXLqgF6EmV3GrOwd9C1EylK9lpBN95qN46wWMFg/VA8/MOd4SRKTlFoO/9ukzMATsoRWKWhzVayebmqnAYb5KRvfj7w3f3TymgVVwsOuA45tOcw0wRXYmP1gnVn9MxUXBBb34oKFriyudd94vW6EBUWPtjaWtPth94Be6MjJu7GqNgtxw",
				}

				for i, entry := range v.Root.Groups[0].Entry[0].HistoryEntries {
					for _, e := range entry.StringEntry {
						if e.Key == "Password" {
							Expect(e.Value.Protected).To(Equal("True"))
							Expect(e.Value.Data).To(Equal(historyEntriesPasswords[i]))
						}
					}
				}
			})
		})

		Context("when given a single valid xml group with historical entries", func() {
			It("returns no error and sets the entries historical values", func() {
				xmlBody, err := ioutil.ReadFile("test_data/History.xml")
				Expect(err).ToNot(HaveOccurred())
				err = xml.Unmarshal(xmlBody, &v)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(v.Root.Groups)).To(Equal(1))
				Expect(len(v.Root.Groups[0].Groups)).To(Equal(1))
				Expect(v.Root.Groups[0].Groups[0].Name).To(Equal("Emails"))
				Expect(len(v.Root.Groups[0].Groups[0].Entry[0].HistoryEntries)).To(Equal(3))

				historyEntriesPasswords := []string{
					"ecOhbRkpfOSDhTTCaUlkdg==",
					"9FmeO/ozrBmRrkvz/z4aiA==",
					"58MDXimsUWKtEIP4JEW+qg==",
				}

				for i, entry := range v.Root.Groups[0].Groups[0].Entry[0].HistoryEntries {
					for _, e := range entry.StringEntry {
						if e.Key == "Password" {
							Expect(e.Value.Protected).To(Equal("True"))
							Expect(e.Value.Data).To(Equal(historyEntriesPasswords[i]))
						}
					}
				}
			})
		})

		Context("when given a valid (format 200) xml database", func() {
			It("returns no error and sets the expected values", func() {
				xmlBody, err := ioutil.ReadFile("test_data/XmlFormat200.xml")
				Expect(err).ToNot(HaveOccurred())
				err = xml.Unmarshal(xmlBody, &v)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(v.Root.Groups)).To(Equal(1))
				Expect(v.Root.Groups[0].UUID).To(Equal("llmpEu56QkCMkm5zu2GzkQ=="))
				Expect(len(v.Root.Groups[0].Groups)).To(Equal(6))

				for _, e := range v.Root.Groups[0].Entry[0].StringEntry {

					switch e.Key {
					case "Password":
						Expect(e.Value.Protected).To(Equal(""))
						Expect(e.Value.Data).To(Equal("Password"))
					case "Title":
						Expect(e.Value.Protected).To(Equal(""))
						Expect(e.Value.Data).To(Equal("Sample Entry"))
					case "Notes":
						Expect(e.Value.Protected).To(Equal(""))
						Expect(e.Value.Data).To(Equal("Notes"))
					case "URL":
						Expect(e.Value.Protected).To(Equal("True"))
						Expect(e.Value.Data).To(Equal("YtgAIKYL4ggH0nzaFP4srWV+GC8A1B0I"))
					case "UserName":
						Expect(e.Value.Protected).To(Equal("True"))
						Expect(e.Value.Data).To(Equal("VK3dahCqvgR8"))
					}
				}

				Expect(v.Meta.HeaderHash).To(Equal(""))
			})
		})

	})

	Describe("Groups XML", func() {

		var (
			v       format.KeePass2XmlFile
			err     error
			xmlBody []byte
		)

		BeforeEach(func() {
			v = format.KeePass2XmlFile{}
			xmlBody, err = ioutil.ReadFile("test_data/NewDatabase.xml")
			Expect(err).ToNot(HaveOccurred())
			err = xml.Unmarshal(xmlBody, &v)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when opening a new database", func() {
			It("succeeds and extracts the root group", func() {
				Expect(v.Root.Groups[0].UUID).To(Equal("lmU+9n0aeESKZvcEze+bRg=="))
				Expect(v.Root.Groups[0].Name).To(Equal("NewDatabase"))
				Expect(len(v.Root.Groups[0].Groups)).To(Equal(3))
				Expect(len(v.Root.Groups[0].Entry)).To(Equal(2))
			})

			It("succesfully extracts group1", func() {
				Expect(v.Root.Groups[0].Groups[0].UUID).To(Equal("AaUYVdXsI02h4T1RiAlgtg=="))
				Expect(v.Root.Groups[0].Groups[0].Name).To(Equal("General"))
			})

			It("succesfully extracts group2", func() {
				Expect(v.Root.Groups[0].Groups[1].UUID).To(Equal("1h4NtL5DK0yVyvaEnN//4A=="))
				Expect(v.Root.Groups[0].Groups[1].Name).To(Equal("Windows"))

				Expect(v.Root.Groups[0].Groups[1].Groups[0].UUID).To(Equal("HoYE/BjLfUSW257pCHJ/eA=="))
				Expect(v.Root.Groups[0].Groups[1].Groups[0].Name).To(Equal("Subsub"))

				Expect(v.Root.Groups[0].Groups[1].Groups[0].Entry[0].UUID).To(Equal("GZpdQvGXOU2kaKRL/IVAGg=="))

				for _, e := range v.Root.Groups[0].Groups[1].Groups[0].Entry[0].StringEntry {
					switch e.Key {
					case "Title":
						Expect(e.Value.Protected).To(Equal(""))
						Expect(e.Value.Data).To(Equal("Subsub Entry"))
					default:
						Expect(e.Value.Protected).To(Equal(""))
						Expect(e.Value.Data).To(Equal(""))
					}
				}
			})

		})
	})

	Describe("Historical entries XML", func() {

		Context("when opening an xml database with historical entries", func() {
			It("succeeds marking those entries as historical", func() {
				s := "f174ffa5ae3b7914ce301bf9841386c66b23ec435d71810df3c5c3b6ded86ae9"
				decoded, err := hex.DecodeString(s)
				Expect(err).ToNot(HaveOccurred())

				var randomKey [32]byte
				copy(randomKey[:], decoded)

				xmlDevice, err := os.Open("test_data/History.xml")
				Expect(err).ToNot(HaveOccurred())

				XMLReader, err := format.NewKeePass2XmlReader(xmlDevice, &randomKey)
				Expect(err).ToNot(HaveOccurred())

				entries := []format.Entry{}
				randomBytesOffset := 0
				err = XMLReader.ReadGroups(&entries, XMLReader.KeePass2XmlFile.Root.Groups, &randomBytesOffset)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(entries)).To(Equal(4))

				Expect(entries[0].Historical).To(Equal(false))

				for _, e := range entries[1:] {
					Expect(e.Historical).To(Equal(true))
				}
			})
		})
	})
})
