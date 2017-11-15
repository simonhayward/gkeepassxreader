package cryptos_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/simonhayward/gkeepassxreader/cryptos"
)

var _ = Describe("Encrypter", func() {

	Describe("Encode", func() {
		var (
			encrypter cryptos.AesEcbEncrypter
			key       []byte
			seed      []byte
			rounds    uint64
			result    *[]byte
		)

		BeforeEach(func() {
			encrypter = cryptos.AesEcbEncrypter{}
			key = []byte{
				0x60, 0x3d, 0xeb, 0x10, 0x15, 0xca, 0x71, 0xbe, 0x2b, 0x73, 0xae, 0xf0, 0x85, 0x7d, 0x77, 0x81,
				0x1f, 0x35, 0x2c, 0x07, 0x3b, 0x61, 0x08, 0xd7, 0x2d, 0x98, 0x10, 0xa3, 0x09, 0x14, 0xdf, 0xf4,
			}
			seed = []byte{
				0x8e, 0x73, 0xb0, 0xf7, 0xda, 0x0e, 0x64, 0x52, 0xc8, 0x10, 0xf3, 0x2b, 0x80, 0x90, 0x79, 0xe5,
				0x62, 0xf8, 0xea, 0xd2, 0x52, 0x2c, 0x6b, 0x7b,
			}
			rounds = uint64(10000)
			result = &[]byte{}
		})

		Context("when given an incorrect seed", func() {
			It("it returns an error", func() {
				err := encrypter.Encode(key, []byte{0x60, 0x3d, 0xeb}, rounds, result)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when given a valid seed", func() {
			It("it returns no error and populates result", func() {
				Expect(*result).Should(BeEmpty())
				err := encrypter.Encode(key, seed, rounds, result)
				Expect(err).ToNot(HaveOccurred())
				Expect(*result).ShouldNot(BeEmpty())
			})
		})

	})
})
