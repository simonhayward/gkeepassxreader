package keys_test

import (
	"encoding/hex"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/simonhayward/gkeepassxreader/keys"
)

var _ = Describe("PasswordKey", func() {

	It("sets the expected value", func() {

		password := "my secret password"
		passwordHexDec := "528eb8404a697e419a740ec1d02908739efce4a04f05a7b90d9f877131d13a44"
		pk := &keys.PasswordKey{}
		pk.SetPassword(password)
		expected, err := hex.DecodeString(passwordHexDec)

		Expect(err).ToNot(HaveOccurred())
		Expect(pk.RawKey()).To(Equal(expected), "password is expected to be equal")
	})

})
