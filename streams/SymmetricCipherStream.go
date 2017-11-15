package streams

import (
	"crypto/cipher"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	//DirectionEncrypt for encryption
	DirectionEncrypt = iota

	//DirectionDecrypt for decryption
	DirectionDecrypt = iota
)

// SymmetricCipherStream represents a symmetric cipher
type SymmetricCipherStream struct {
	buffer        []byte
	bufferPos     int
	bufferFilling bool
	Block         cipher.Block
	BlockMode     cipher.BlockMode
	db            *os.File
}

//NewSymmetricCipherStream new stream
func NewSymmetricCipherStream(block cipher.Block, encryptionIV []byte, db *os.File, direction int) (*SymmetricCipherStream, error) {
	var blockMode cipher.BlockMode

	if direction == DirectionEncrypt {
		blockMode = cipher.NewCBCEncrypter(block, encryptionIV)
	} else {
		blockMode = cipher.NewCBCDecrypter(block, encryptionIV)
	}

	s := SymmetricCipherStream{
		Block:         block,
		BlockMode:     blockMode,
		bufferFilling: false,
		bufferPos:     0,
		db:            db,
	}

	return &s, nil
}

//ReadData read data from stream
func (s *SymmetricCipherStream) ReadData(data *[]byte, maxSize int) (int, error) {

	log.Debugf("[SymmetricCipherStream::ReadData] maxSize: %d", maxSize)
	bytesRemaining := maxSize
	offset := 0
	var bytesToCopy int

	for bytesRemaining > 0 {
		if (s.bufferPos == len(s.buffer)) || s.bufferFilling {
			res, err := s.readBlock()
			if res == false {
				if err != nil {
					return 0, err
				}
				return maxSize - bytesRemaining, nil
			}
		}

		bytesToCopy = len(s.buffer) - s.bufferPos
		if bytesRemaining < bytesToCopy {
			bytesToCopy = bytesRemaining
		}

		log.Debugf("[SymmetricCipherStream::ReadData] bytesToCopy: %d", bytesToCopy)

		if len(*data) < offset+bytesToCopy {
			newSlice := make([]byte, offset+bytesToCopy)
			copy(newSlice, *data)
			*data = newSlice
		}

		copy((*data)[offset:offset+bytesToCopy], s.buffer[s.bufferPos:s.bufferPos+bytesToCopy])
		log.Debugf("[SymmetricCipherStream::ReadData] offset: %d bufferPos: %d", offset, s.bufferPos)

		offset += bytesToCopy
		s.bufferPos += bytesToCopy
		bytesRemaining -= bytesToCopy
	}

	return maxSize, nil
}

func (s *SymmetricCipherStream) readBlock() (bool, error) {

	var newData []byte

	if s.bufferFilling {
		newData = make([]byte, s.Block.BlockSize()-len(s.buffer))
	} else {
		s.buffer = nil
		newData = make([]byte, s.Block.BlockSize())
	}

	readResult, err := io.ReadAtLeast(s.db, newData, len(newData))

	log.Debugf("[SymmetricCipherStream::ReadBlock] readResult: %d", readResult)

	if err != nil {
		return false, err
	}

	s.buffer = append(s.buffer, newData...)
	log.Debugf("[SymmetricCipherStream::readDataBlock] buffer APPEND: %d", len(s.buffer))

	if len(s.buffer) != s.Block.BlockSize() {
		s.bufferFilling = true
		return false, nil
	}

	log.Debugf("[SymmetricCipherStream::readDataBlock] processInPlace: %d", len(s.buffer))
	s.BlockMode.CryptBlocks(s.buffer, s.buffer)

	s.bufferPos = 0
	s.bufferFilling = false

	return true, nil
}
