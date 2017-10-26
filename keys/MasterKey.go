package keys

import (
	"os"

	log "github.com/sirupsen/logrus"
)

//MasterKey from password and file key
func MasterKey(password string, keyFile *os.File) *CompositeKey {

	masterKey := NewCompositeKey()

	if len(password) > 0 {
		pk := &PasswordKey{}
		pk.SetPassword(password)
		masterKey.AddKey(pk)
	}

	if keyFile != nil {
		kf := &FileKey{}
		if !kf.Load(keyFile) {
			log.Warn("unable to load key file")
			return &CompositeKey{}
		}
		masterKey.AddKey(kf)
	}

	return masterKey
}
