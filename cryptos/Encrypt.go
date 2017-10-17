package cryptos

import (
	"crypto/aes"
	"fmt"
)

//Encrypt to transform key
type Encrypt interface {
	Encode(key []byte, seed []byte, rounds uint64, result *[]byte) error
}

//AesEcbEncrypter AES/ECB encrypter
type AesEcbEncrypter struct{}

//Encode is the implementation of AES/ECB encryption
func (a *AesEcbEncrypter) Encode(key []byte, seed []byte, rounds uint64, result *[]byte) error {
	cipherBlock, err := aes.NewCipher(seed)
	if err != nil {
		return fmt.Errorf("unable to create cipher block: %s", err)
	}

	dst := make([]byte, len(key))
	copy(dst, key)

	encrypter := NewECBEncrypter(cipherBlock)

	for i := uint64(0); i < rounds; i++ {
		encrypter.CryptBlocks(dst, dst)
	}

	*result = dst
	return nil
}
