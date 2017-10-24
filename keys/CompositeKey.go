package keys

import (
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/simonhayward/gkeepassxreader/cryptos"
)

const (
	//TransformSeedSize length
	TransformSeedSize = 32
)

//CompositeKey holds keys
type CompositeKey struct {
	keys      []Key
	Encrypter cryptos.Encrypt
}

//NewCompositeKey defaults
func NewCompositeKey() *CompositeKey {
	return &CompositeKey{
		Encrypter: &cryptos.AesEcbEncrypter{},
	}
}

//AddKey appends Key to slice
func (c *CompositeKey) AddKey(k Key) {
	c.keys = append(c.keys, k)
}

//Transform the composite key by performing encryption and returning a checksum
func (c *CompositeKey) Transform(seed []byte, rounds uint64) ([]byte, error) {
	if len(seed) != TransformSeedSize {
		return []byte{}, fmt.Errorf("seed size error, expected: %d received: %d", TransformSeedSize, len(seed))
	}

	if rounds == uint64(0) {
		return []byte{}, fmt.Errorf("rounds error, expected greater than zero")
	}

	var resultLeft, resultRight []byte
	var wg sync.WaitGroup

	errc := make(chan error, 2)
	splitKey := len(c.RawKey()) / 2

	wg.Add(2)
	go c.Encrypt(c.RawKey()[:splitKey], seed, rounds, &resultLeft, errc, &wg)
	go c.Encrypt(c.RawKey()[splitKey:], seed, rounds, &resultRight, errc, &wg)
	wg.Wait()
	close(errc)

	for e := range errc {
		if e != nil {
			return nil, e
		}
	}

	transformed := resultLeft
	transformed = append(transformed, resultRight...)

	h := sha256.New()
	h.Write(transformed)

	return h.Sum(nil), nil
}

//Encrypt performs the encryption
func (c *CompositeKey) Encrypt(key []byte, seed []byte, rounds uint64, result *[]byte, errc chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	errc <- c.Encrypter.Encode(key, seed, rounds, result)
}

//RawKey returns the checksum of all keys combined
func (c *CompositeKey) RawKey() []byte {

	h := sha256.New()
	for _, key := range c.keys {
		h.Write(key.RawKey())
	}
	return h.Sum(nil)
}
