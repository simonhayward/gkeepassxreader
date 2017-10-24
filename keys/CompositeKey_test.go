package keys_test

import (
	"crypto/sha256"

	"github.com/simonhayward/gkeepassxreader/keys/keysfakes"

	"github.com/simonhayward/gkeepassxreader/cryptos/cryptosfakes"
	"github.com/simonhayward/gkeepassxreader/keys"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompositeKey", func() {

	Describe("Transformation", func() {
		var (
			regSeed       []byte
			longSeed      []byte
			noRounds      uint64
			defaultRounds uint64
			compositeKey  *keys.CompositeKey
			fakeEncrypt   *cryptosfakes.FakeEncrypt
		)

		BeforeEach(func() {
			regSeed = []byte("abcdefghijklmnopqrstuvwxyz123456")
			longSeed = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			noRounds = uint64(0)
			defaultRounds = uint64(1000)
			fakeEncrypt = &cryptosfakes.FakeEncrypt{}
			compositeKey = keys.NewCompositeKey()
		})

		Context("when given a valid seed with vaid transformation rounds", func() {
			It("it returns no error and will call the encryption", func() {
				compositeKey.Encrypter = fakeEncrypt
				out, err := compositeKey.Transform(regSeed, defaultRounds)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(out)).To(Equal(sha256.Size))
				Expect(fakeEncrypt.EncodeCallCount()).To(Equal(2), "Encryption should be called twice")
			})
		})

		Context("when given a invalid seed length and valid transformation rounds", func() {
			It("returns an error and encryption is not called", func() {
				compositeKey.Encrypter = fakeEncrypt
				out, err := compositeKey.Transform(longSeed, defaultRounds)

				Expect(err).To(HaveOccurred())
				Expect(out).To(Equal([]byte{}))
				Expect(fakeEncrypt.EncodeCallCount()).To(Equal(0), "Encryption should NOT be called")
			})
		})

		Context("when given a valid seed length and invalid transformation rounds", func() {
			It("returns an error and encryption is not called", func() {
				compositeKey.Encrypter = fakeEncrypt
				out, err := compositeKey.Transform(regSeed, noRounds)

				Expect(err).To(HaveOccurred())
				Expect(out).To(Equal([]byte{}))
				Expect(fakeEncrypt.EncodeCallCount()).To(Equal(0), "Encryption should NOT be called")
			})
		})
	})

	Describe("RawKeys", func() {
		var (
			compositeKey *keys.CompositeKey
			passwordKey  *keys.PasswordKey
			fakeFileKey  *keysfakes.FakeKey
		)

		BeforeEach(func() {
			compositeKey = keys.NewCompositeKey()
			passwordKey = &keys.PasswordKey{}
			fakeFileKey = &keysfakes.FakeKey{}
		})
		Context("when given a password key", func() {
			It("it returns the checksum value", func() {
				p := "my password"
				passwordKey.SetPassword(p)
				compositeKey.AddKey(passwordKey)
				expected := sha256.Sum256(passwordKey.RawKey())

				out := compositeKey.RawKey()
				Expect(len(out)).To(Equal(sha256.Size))
				Expect(out).To(Equal(expected[:]))
			})
		})

		Context("when given multiple keys", func() {
			It("it returns the checksum of all keys combined", func() {
				p := "my password"
				passwordKey.SetPassword(p)
				compositeKey.AddKey(passwordKey)

				myFilePassword := sha256.Sum256([]byte("file password\n"))
				fakeFileKey.RawKeyReturns(myFilePassword[:])
				compositeKey.AddKey(fakeFileKey)

				h := sha256.New()
				h.Write(passwordKey.RawKey())
				h.Write(fakeFileKey.RawKey())
				expected := h.Sum(nil)

				out := compositeKey.RawKey()
				Expect(len(out)).To(Equal(sha256.Size))
				Expect(out).To(Equal(expected[:]))
			})
		})
	})

})
