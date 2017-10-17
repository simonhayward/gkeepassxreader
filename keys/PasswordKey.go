package keys

import (
	"crypto/sha256"
)

type PasswordKey struct {
	key []byte
}

func (p *PasswordKey) RawKey() []byte {
	return p.key
}

func (p *PasswordKey) SetPassword(password string) {
	b := sha256.Sum256([]byte(password))
	p.key = b[:]
}
