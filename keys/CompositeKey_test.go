package keys_test

import (
	"crypto/sha256"

	"github.com/simonhayward/gkeepassxreader/keys"
	"github.com/simonhayward/gkeepassxreader/keys/keysfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompositeKey", func() {

	Describe("Transform seed", func() {
		var (
			regSeed       []byte
			longSeed      []byte
			noRounds      uint64
			defaultRounds uint64
			compositeKey  *keys.CompositeKey
			fakeEncrypt   *keysfakes.FakeEncrypt
		)

		BeforeEach(func() {
			regSeed = make([]byte, keys.TransformSeedSize)
			regSeed = []byte("abcdefghijklmnopqrstuvwxyz123456")
			longSeed = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			noRounds = uint64(0)
			defaultRounds = uint64(1000)
			fakeEncrypt = &keysfakes.FakeEncrypt{}
			compositeKey = &keys.CompositeKey{Encrypter: fakeEncrypt}
		})

		It("encodes successfully", func() {
			out, err := compositeKey.Transform(regSeed, defaultRounds)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(out)).To(Equal(sha256.Size))
			Expect(fakeEncrypt.EncodeCallCount()).To(Equal(2), "Encryption should be called twice")
		})

		It("returns an error when a seed length is exceeded", func() {
			out, err := compositeKey.Transform(longSeed, defaultRounds)

			Expect(err).To(HaveOccurred())
			Expect(out).To(Equal([]byte{}))
			Expect(fakeEncrypt.EncodeCallCount()).To(Equal(0), "Encryption should not be called")
		})

		It("returns an error when the transform rounds are invalid", func() {
			out, err := compositeKey.Transform(regSeed, noRounds)

			Expect(err).To(HaveOccurred())
			Expect(out).To(Equal([]byte{}))
			Expect(fakeEncrypt.EncodeCallCount()).To(Equal(0), "Encryption should not be called")
		})
	})
})
