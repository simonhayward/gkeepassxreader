package keys

import (
	"crypto/sha256"
)

//PasswordKey represents a password key
type PasswordKey struct {
	key []byte
}

//RawKey returns the key in bytes
func (p *PasswordKey) RawKey() []byte {
	return p.key
}

//SetPassword creates the key from a checksum hash
func (p *PasswordKey) SetPassword(password string) {
	b := sha256.Sum256([]byte(password))
	p.key = b[:]
}
