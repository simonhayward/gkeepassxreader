package streams

import (
	"golang.org/x/crypto/salsa20"
)

// Salsa20Stream represents a symmetric cipher
type Salsa20Stream struct {
	nonce []byte
	key   *[32]byte
}

//NewSalsa20Stream new stream
func NewSalsa20Stream(nonce []byte, key *[32]byte) *Salsa20Stream {

	s := Salsa20Stream{
		nonce: nonce,
		key:   key,
	}

	return &s
}

//ProcessInPlace update slice in place
func (s *Salsa20Stream) ProcessInPlace(data []byte) {
	salsa20.XORKeyStream(data, data, s.nonce, s.key)
}
