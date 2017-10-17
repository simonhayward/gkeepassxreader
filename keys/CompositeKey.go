package keys

import (
	"crypto/aes"
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/simonhayward/gkeepassxreader/cryptos"
)

//CompositeKey holds keys
type CompositeKey struct {
	keys []Key
}

//AddKey appends Key to slice
func (c *CompositeKey) AddKey(k Key) {
	c.keys = append(c.keys, k)
}

//Transform the composite key return sha 256 sum
func (c *CompositeKey) Transform(seed []byte, rounds uint64) ([]byte, error) {
	if len(seed) != 32 {
		return nil, fmt.Errorf("seed size error, expected: %d received: %d", 32, len(seed))
	}

	if rounds <= uint64(0) {
		return nil, fmt.Errorf("rounds error, expected greater than: %d", rounds)
	}

	var resultLeft, resultRight []byte
	var wg sync.WaitGroup

	errc := make(chan error, 2)
	splitKey := len(c.RawKey()) / 2

	wg.Add(2)
	go c.transformKeyRaw(c.RawKey()[:splitKey], seed, rounds, &resultLeft, errc, &wg)
	go c.transformKeyRaw(c.RawKey()[splitKey:], seed, rounds, &resultRight, errc, &wg)
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

//RawKey returns the 256 sum of keys
func (c *CompositeKey) RawKey() []byte {

	h := sha256.New()

	for _, key := range c.keys {
		h.Write(key.RawKey())
	}

	return h.Sum(nil)
}

func (c *CompositeKey) transformKeyRaw(key []byte, seed []byte, rounds uint64, result *[]byte, errc chan error, wg *sync.WaitGroup) {

	cipherBlock, err := aes.NewCipher(seed)
	if err != nil {
		errc <- err
		return
	}

	dst := make([]byte, len(key))
	copy(dst, key)

	encrypter := cryptos.NewECBEncrypter(cipherBlock)

	for i := uint64(0); i < rounds; i++ {
		encrypter.CryptBlocks(dst, dst)
	}

	*result = dst
	errc <- nil
	wg.Done()
}
