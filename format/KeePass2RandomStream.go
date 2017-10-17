package format

import (
	"github.com/simonhayward/gkeepassxreader/streams"
)

//KeePass2RandomStream represents a random stream
type KeePass2RandomStream struct {
	cipherStream *streams.Salsa20Stream
	offset       int
	buffer       []byte
}

//NewKeePass2RandomStream returns a new random stream pointer
func NewKeePass2RandomStream(nonce []byte, key *[32]byte) *KeePass2RandomStream {

	s := KeePass2RandomStream{
		cipherStream: streams.NewSalsa20Stream(nonce, key),
	}
	return &s
}

func (r *KeePass2RandomStream) randomBytes(offset, length int) []byte {
	r.buffer = make([]byte, offset+length)
	r.cipherStream.ProcessInPlace(r.buffer)
	return r.buffer[offset : offset+length]
}

// Process request
func (r *KeePass2RandomStream) Process(offset int, ciphertext []byte) ([]byte, error) {

	randomData := r.randomBytes(offset, len(ciphertext))
	result := make([]byte, len(ciphertext))

	for i := 0; i < len(ciphertext); i++ {
		result[i] = ciphertext[i] ^ randomData[i]
	}

	return result, nil
}
