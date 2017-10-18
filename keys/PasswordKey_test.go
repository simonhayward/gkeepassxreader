package keys_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/simonhayward/gkeepassxreader/keys"
)

func TestSetPassword(t *testing.T) {

	password := "my secret password"
	passwordHexDec := "528eb8404a697e419a740ec1d02908739efce4a04f05a7b90d9f877131d13a44"
	pk := &keys.PasswordKey{}
	pk.SetPassword(password)
	expected, err := hex.DecodeString(passwordHexDec)

	if err != nil {
		t.Errorf("decode failed: %s", err)
	}

	if !bytes.Equal(pk.RawKey(), expected) {
		t.Errorf("password unequal, expected: %x received: %x", expected, pk.RawKey())
	}
}
