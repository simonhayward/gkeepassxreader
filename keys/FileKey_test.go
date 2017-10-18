package keys_test

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"github.com/simonhayward/gkeepassxreader/keys"
)

func createTempFile(t *testing.T, name string) *os.File {
	tmpFile, err := ioutil.TempFile("", name)

	if err != nil {
		t.Fatalf("Tempfile creation err: %s", err)
	}

	return tmpFile
}

func TestLoadXmlValid(t *testing.T) {
	tmpFile := createTempFile(t, "_xml_valid_")
	defer os.Remove(tmpFile.Name())

	xmlValidFile := `<?xml version="1.0" encoding="UTF-8"?>
	<KeyFile>
		<Meta>
			<Version>1.00</Version>
		</Meta>
		<Key>
			<Data>VbM4cH69dgdelucJwa+u6g038mMrxCTHbUr1a9haLBY=</Data>
		</Key>
	</KeyFile>`

	if _, err := tmpFile.Write([]byte(xmlValidFile)); err != nil {
		t.Fatalf("Write failed: %s", err)
	}

	fk := &keys.FileKey{}
	if fk.Load(tmpFile) == false {
		t.Errorf("Xml is valid but returns false")
	}
}

func TestLoadXmlInvalidMeta(t *testing.T) {
	tmpFile := createTempFile(t, "_xml_invalid_meta_")
	defer os.Remove(tmpFile.Name())

	xmlInvalidFile := `<?xml version="1.0" encoding="UTF-8"?>
	<KeyFile>
		<Meta>
		</Meta>
		<Key>
			<Data>VbM4cH69dgdelucJwa+u6g038mMrxCTHbUr1a9haLBY=</Data>
		</Key>
	</KeyFile>`

	if _, err := tmpFile.Write([]byte(xmlInvalidFile)); err != nil {
		t.Fatalf("Write failed: %s", err)
	}

	fk := &keys.FileKey{}
	if fk.Load(tmpFile) == true {
		t.Errorf("Xml is invalid but returns true")
	}
}

func TestLoadXmlInvalidKey(t *testing.T) {
	tmpFile := createTempFile(t, "_xml_invalid_key_")
	defer os.Remove(tmpFile.Name())

	xmlInvalidFile := `<?xml version="1.0" encoding="UTF-8"?>
	<KeyFile>
		<Meta>
			<Version>1.00</Version>
		</Meta>
		<Key>
		</Key>
	</KeyFile>`

	if _, err := tmpFile.Write([]byte(xmlInvalidFile)); err != nil {
		t.Fatalf("Write failed: %s", err)
	}

	fk := &keys.FileKey{}
	if fk.Load(tmpFile) == true {
		t.Errorf("Xml is invalid but returns true")
	}
}
func TestLoadXmlInvalidKeyBase64(t *testing.T) {
	tmpFile := createTempFile(t, "_xml_invalid_key_base64_")
	defer os.Remove(tmpFile.Name())

	xmlInvalidFile := `<?xml version="1.0" encoding="UTF-8"?>
	<KeyFile>
		<Meta>
			<Version>1.00</Version>
		</Meta>
		<Key>
			<Data>1yuu-(%this</Data>
		</Key>
	</KeyFile>`

	if _, err := tmpFile.Write([]byte(xmlInvalidFile)); err != nil {
		t.Fatalf("Write failed: %s", err)
	}

	fk := &keys.FileKey{}
	if fk.Load(tmpFile) == true {
		t.Errorf("Xml is invalid but returns true")
	}
}

func TestLoadBinaryValid(t *testing.T) {
	tmpFile := createTempFile(t, "_binary_valid_")
	defer os.Remove(tmpFile.Name())

	var binaryValidFile []byte
	for i := 0; i < keys.KeySize; i++ {
		binaryValidFile = append(binaryValidFile, byte(1))
	}

	if _, err := tmpFile.Write(binaryValidFile); err != nil {
		t.Fatalf("Write failed: %s", err)
	}

	fk := &keys.FileKey{}
	if fk.Load(tmpFile) == false {
		t.Errorf("Binary is valid but returns false")
	}

	if !bytes.Equal(fk.Key, binaryValidFile) {
		t.Errorf("Binary is valid but key not equal")
	}
}

func TestLoadBinaryInvalid(t *testing.T) {
	tmpFile := createTempFile(t, "_binary_invalid_")
	defer os.Remove(tmpFile.Name())

	var binaryInvalidFile []byte
	for i := 0; i < keys.KeySize; i++ {
		binaryInvalidFile = append(binaryInvalidFile, byte(1))
	}
	binaryInvalidFile = append(binaryInvalidFile, byte(1)) //excess

	if _, err := tmpFile.Write(binaryInvalidFile); err != nil {
		t.Fatalf("Write failed: %s", err)
	}

	fk := &keys.FileKey{}
	if bytes.Equal(fk.Key, binaryInvalidFile) {
		t.Errorf("Binary is invalid but key is equal")
	}
}

func TestLoadHexValid(t *testing.T) {
	tmpFile := createTempFile(t, "_hex_valid_")
	defer os.Remove(tmpFile.Name())

	var key []byte
	for i := 0; i < keys.KeySize; i++ {
		key = append(key, byte(1))
	}

	hexValidFile := make([]byte, keys.HexSize)
	hex.Encode(hexValidFile, key)

	if _, err := tmpFile.Write(hexValidFile); err != nil {
		t.Fatalf("Write failed: %s", err)
	}

	fk := &keys.FileKey{}
	if fk.Load(tmpFile) == false {
		t.Errorf("Hex is valid but returns false")
	}
}

func TestLoadHexInvalid(t *testing.T) {
	tmpFile := createTempFile(t, "_hex_invalid_")
	defer os.Remove(tmpFile.Name())

	var key []byte
	for i := 0; i < keys.KeySize; i++ {
		key = append(key, byte(1))
	}
	key = append(key, byte(1)) //excessive key length

	hexInvalidFile := make([]byte, hex.EncodedLen(len(key)))
	hex.Encode(hexInvalidFile, key)

	if _, err := tmpFile.Write(hexInvalidFile); err != nil {
		t.Fatalf("Write failed: %s", err)
	}

	fk := &keys.FileKey{}
	if bytes.Equal(fk.Key, hexInvalidFile) {
		t.Errorf("Hex is invalid but key is equal")
	}
}

func TestLoadHashedValid(t *testing.T) {
	tmpFile := createTempFile(t, "_hashed_valid_")
	defer os.Remove(tmpFile.Name())

	text := "Valid hashed text"
	s := "6aa6fbec3065ac2b23bd3357a1ad23e7ccf1eca0865f658f1736a882525d744d"
	expected, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("Decode failed: %s", err)
	}

	if _, err := tmpFile.Write([]byte(text)); err != nil {
		t.Fatalf("Write failed: %s", err)
	}

	fk := &keys.FileKey{}
	if fk.Load(tmpFile) == false {
		t.Errorf("Hash is valid but returns false")
	}

	if !bytes.Equal(fk.Key, expected) {
		t.Errorf("Expected: %s received: %s", hex.EncodeToString(expected), hex.EncodeToString(fk.Key))
	}
}
func TestEmptyFileInvalid(t *testing.T) {
	tmpFile := createTempFile(t, "_empty_")
	defer os.Remove(tmpFile.Name())

	fk := &keys.FileKey{}
	if fk.Load(tmpFile) == true {
		t.Errorf("File is empty but returns true")
	}
}
