package output_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/simonhayward/gkeepassxreader/format"
	"github.com/simonhayward/gkeepassxreader/output"
)

var _ = Describe("Extract", func() {

	var (
		entry *format.Entry
	)

	BeforeEach(func() {
		entry = &format.Entry{
			Password: &format.EntryValue{PlainText: "this is my password"},
		}
	})

	Context("when trying to extract characters from a string", func() {
		It("succeeds and updates the password", func() {
			err := output.Extract(entry, "1,2,3")
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("thi"))
		})
	})

	Context("when trying to extract repeated characters from a string", func() {
		It("succeeds and updates the password", func() {
			err := output.Extract(entry, "19,1,2,3,19")
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("dthid"))
		})
	})

	Context("when trying to extract characters from a string with whitespace", func() {
		It("succeeds and updates the password", func() {
			err := output.Extract(entry, " 1, 2 ,	3	")
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("thi"))
		})
	})

	Context("when trying to extract characters from a string with trailing comma", func() {
		It("succeeds and updates the password", func() {
			err := output.Extract(entry, "1,2,")
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("th"))
		})
	})

	Context("when trying to extract characters from a string with leading comma", func() {
		It("succeeds and updates the password", func() {
			err := output.Extract(entry, ",1,2")
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("th"))
		})
	})

	Context("when trying to extract characters from a string with whitespace and comma", func() {
		It("succeeds and updates the password", func() {
			err := output.Extract(entry, " , 1,2")
			Expect(err).ToNot(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("th"))
		})
	})

	Context("when trying to extract characters out of range for the string", func() {
		It("returns an error and the password remains unchanged", func() {
			err := output.Extract(entry, "0,2,3")
			Expect(err).To(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("this is my password"))

			err = output.Extract(entry, "1, 2, 199")
			Expect(err).To(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("this is my password"))
		})
	})

	Context("when trying to extract negative integers", func() {
		It("returns an error and the password remains unchanged", func() {
			err := output.Extract(entry, "1,2,-3")
			Expect(err).To(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("this is my password"))
		})
	})

	Context("when trying to extract non character integers", func() {
		It("returns an error and the password remains unchanged", func() {
			err := output.Extract(entry, "1,2,P")
			Expect(err).To(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("this is my password"))

			err = output.Extract(entry, "1,A,2")
			Expect(err).To(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("this is my password"))
		})
	})

	Context("when trying to extract an empty string", func() {
		It("returns an error and the password remains unchanged", func() {
			err := output.Extract(entry, "")
			Expect(err).To(HaveOccurred())

			Expect(entry.Password.PlainText).To(Equal("this is my password"))
		})
	})

})
