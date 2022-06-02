package keys_test

import (
	"encoding/hex"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/simonhayward/gkeepassxreader/keys"
)

var _ = Describe("FileKey", func() {

	var (
		tmpFile *os.File
		err     error
		fk      *keys.FileKey
	)

	BeforeEach(func() {
		tmpFile, err = ioutil.TempFile("", "filekey")
		Expect(err).ToNot(HaveOccurred())
		fk = &keys.FileKey{}
	})

	AfterEach(func() {
		os.Remove(tmpFile.Name())
	})

	Context("when given an xml file which is valid", func() {
		It("returns true", func() {
			file := `<?xml version="1.0" encoding="UTF-8"?>
			<KeyFile>
				<Meta>
					<Version>1.00</Version>
				</Meta>
				<Key>
					<Data>VbM4cH69dgdelucJwa+u6g038mMrxCTHbUr1a9haLBY=</Data>
				</Key>
			</KeyFile>`

			_, err = tmpFile.Write([]byte(file))
			Expect(err).ToNot(HaveOccurred())
			Expect(fk.Load(tmpFile)).To(Equal(true), "Xml is valid but returns false")
		})
	})

	Context("when given an invalid xml file", func() {
		It("returns false", func() {
			file := `<?xml version="1.0" encoding="UTF-8"?>
			<KeyFile>
				<Meta>
				</Meta>
				<Key>
					<Data>VbM4cH69dgdelucJwa+u6g038mMrxCTHbUr1a9haLBY=</Data>
				</Key>
			</KeyFile>`

			_, err = tmpFile.Write([]byte(file))
			Expect(err).ToNot(HaveOccurred())
			Expect(fk.Load(tmpFile)).To(Equal(false), "Xml is invalid but returns true")
		})
	})

	Context("when given an xml file with an invalid key", func() {
		It("returns false", func() {
			file := `<?xml version="1.0" encoding="UTF-8"?>
		<KeyFile>
			<Meta>
				<Version>1.00</Version>
			</Meta>
			<Key>
			</Key>
		</KeyFile>`

			_, err := tmpFile.Write([]byte(file))
			Expect(err).ToNot(HaveOccurred())
			Expect(fk.Load(tmpFile)).To(Equal(false), "Xml is invalid but returns true")
		})
	})

	Context("when given an xml file with an invalid base 64 key", func() {
		It("returns false", func() {
			file := `<?xml version="1.0" encoding="UTF-8"?>
		<KeyFile>
			<Meta>
				<Version>1.00</Version>
			</Meta>
			<Key>
				<Data>1yuu-(%this</Data>
			</Key>
		</KeyFile>`

			_, err := tmpFile.Write([]byte(file))
			Expect(err).ToNot(HaveOccurred())
			Expect(fk.Load(tmpFile)).To(Equal(false), "Xml is invalid but returns true")
		})
	})

	Context("when given a valid binary file", func() {
		It("returns true and to match the key", func() {
			var binaryValidFile []byte
			for i := 0; i < keys.KeySize; i++ {
				binaryValidFile = append(binaryValidFile, byte(1))
			}

			_, err = tmpFile.Write(binaryValidFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(fk.Load(tmpFile)).To(Equal(true), "Binary is valid but returns false")
			Expect(fk.Key).To(Equal(binaryValidFile), "Binary is valid but key not equal")
		})
	})

	Context("when given an invalid binary file", func() {
		It("expects the key to not match", func() {
			var binaryInvalidFile []byte
			for i := 0; i < keys.KeySize; i++ {
				binaryInvalidFile = append(binaryInvalidFile, byte(1))
			}
			binaryInvalidFile = append(binaryInvalidFile, byte(1)) //excess

			_, err = tmpFile.Write(binaryInvalidFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(fk.Key).ToNot(Equal(binaryInvalidFile), "Binary is invalid but key is equal")
		})
	})

	Context("when given a valid hex file", func() {
		It("returns true", func() {
			var key []byte
			for i := 0; i < keys.KeySize; i++ {
				key = append(key, byte(1))
			}

			hexValidFile := make([]byte, keys.HexSize)
			hex.Encode(hexValidFile, key)

			_, err = tmpFile.Write(hexValidFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(fk.Load(tmpFile)).To(Equal(true), "Hex is valid but returns false")
		})
	})

	Context("when given an invalid hex file", func() {
		It("expects the key to not match", func() {
			var key []byte
			for i := 0; i < keys.KeySize; i++ {
				key = append(key, byte(1))
			}
			key = append(key, byte(1)) //excessive key length

			hexInvalidFile := make([]byte, hex.EncodedLen(len(key)))
			hex.Encode(hexInvalidFile, key)

			_, err = tmpFile.Write(hexInvalidFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(fk.Key).ToNot(Equal(hexInvalidFile), "Hex is invalid but key is equal")
		})
	})

	Context("when given a valid hashed file", func() {
		It("returns true and the key should match", func() {

			text := "Valid hashed text"
			s := "6aa6fbec3065ac2b23bd3357a1ad23e7ccf1eca0865f658f1736a882525d744d"
			expected, err := hex.DecodeString(s)
			Expect(err).ToNot(HaveOccurred())

			_, err = tmpFile.Write([]byte(text))
			Expect(err).ToNot(HaveOccurred())
			Expect(fk.Load(tmpFile)).To(Equal(true), "Hash is valid but returns false")

			Expect(fk.Key).To(Equal(expected), "hash is valid and should match key")
		})
	})

	Context("when given a empty file", func() {
		It("returns false", func() {
			Expect(fk.Load(tmpFile)).To(Equal(false), "Empty file is invalid but returns true")

		})
	})
})
