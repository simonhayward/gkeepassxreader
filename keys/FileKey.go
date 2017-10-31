package keys

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	xmlMetaVersion = "1.00"
	//HexSize length
	HexSize = 64
	//KeySize length
	KeySize = 32
)

type xmlMeta struct {
	XMLName xml.Name `xml:"Meta"`
	Version string   `xml:"Version"`
}

type xmlKey struct {
	XMLName xml.Name `xml:"Key"`
	Data    string   `xml:"Data"`
}

type xmlKeyFile struct {
	XMLName xml.Name `xml:"KeyFile"`
	Meta    xmlMeta  `xml:"Meta"`
	Key     xmlKey   `xml:"Key"`
}

//FileKey represents a file key
type FileKey struct {
	Key []byte
}

// RawKey represents key as byte slice
func (fk *FileKey) RawKey() []byte {
	return fk.Key
}

// Load keyfile
func (fk *FileKey) Load(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}

	if fi.Size() == 0 {
		return false
	}

	// try different key file formats
	if !reset(f) {
		return false
	}

	if fk.validXML(f) {
		if !reset(f) {
			return false
		}
		if fk.loadXML(f) {
			return true
		}
	} else {
		if !reset(f) {
			return false
		}
		if fk.loadBinary(f) {
			return true
		}
		if !reset(f) {
			return false
		}
		if fk.loadHex(f) {
			return true
		}
		if !reset(f) {
			return false
		}
		if fk.loadHashed(f) {
			return true
		}
	}

	return false
}

func reset(f *os.File) bool {
	if _, err := f.Seek(0, 0); err != nil {
		return false
	}
	return true
}

func (fk *FileKey) validXML(f *os.File) bool {
	xmlData, err := ioutil.ReadAll(f)
	if err == nil {
		x := xmlKeyFile{}
		err = xml.Unmarshal(xmlData, &x)
		if err == nil {
			return true
		}
	}
	return false
}

func (fk *FileKey) loadXML(f *os.File) bool {
	xmlFile := &xmlKeyFile{}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Errorf("xml read failed: %s", err)
		return false
	}

	err = xml.Unmarshal(b, xmlFile)
	if err != nil {
		log.Errorf("xml unmarshal failed: %s", err)
		return false
	}

	if xmlFile.XMLName.Local != "KeyFile" {
		return false
	}

	// check meta version
	if xmlFile.Meta.Version != xmlMetaVersion {
		return false
	}

	// check key
	if len(xmlFile.Key.Data) == 0 {
		return false
	}
	if fk.loadxmlKey(xmlFile.Key.Data) == nil {
		return true
	}

	return false
}

func (fk *FileKey) loadxmlKey(k string) error {

	data, err := base64.StdEncoding.DecodeString(k)
	if err != nil {
		return fmt.Errorf("base64 decoding failed: %s", err)
	}

	fk.Key = data
	return nil
}

func (fk *FileKey) loadBinary(f *os.File) bool {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Errorf("binary read failed: %s", err)
		return false
	}
	if len(b) != KeySize {
		return false
	}

	fk.Key = append([]byte(nil), b...)
	return true
}

func (fk *FileKey) loadHex(f *os.File) bool {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Errorf("binary read failed: %s", err)
		return false
	}
	if len(b) != HexSize {
		return false
	}

	dst := make([]byte, KeySize)
	n, err := hex.Decode(dst, b)
	if err != nil {
		log.Errorf("hex decode failed: %s", err)
		return false
	}

	if n != KeySize {
		return false
	}

	fk.Key = append([]byte(nil), dst...)
	return true
}

func (fk *FileKey) loadHashed(f *os.File) bool {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Errorf("hashed read failed: %s", err)
		return false
	}

	h := sha256.New()
	h.Write(b)
	fk.Key = h.Sum(nil)
	return true
}
