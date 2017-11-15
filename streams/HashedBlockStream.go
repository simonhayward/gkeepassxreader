package streams

import (
	"bytes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// HashedBlock represents a hashed block
type HashedBlock struct {
	buffer       []byte
	blockIndex   uint32
	bufferPos    int
	mode         cipher.BlockMode
	cipherStream *SymmetricCipherStream
	eof          bool
}

//NewHashedBlock create new hashed block
func NewHashedBlock(mode cipher.BlockMode, stream *SymmetricCipherStream) *HashedBlock {
	return &HashedBlock{
		mode:         mode,
		cipherStream: stream,
	}
}

// ReadData in hashed blocks
func (hb *HashedBlock) ReadData(data *[]byte, maxSize int) (int, error) {

	bytesRemaining := maxSize
	var offset, bytesToCopy int

	if hb.eof {
		return 0, nil
	}

	for bytesRemaining > 0 {
		if hb.bufferPos == len(hb.buffer) {
			res, err := hb.readHashedBlock()
			if res == false {
				if err != nil {
					return 0, err
				}
				return maxSize - bytesRemaining, nil
			}
		}

		bytesToCopy = len(hb.buffer) - hb.bufferPos
		if bytesRemaining < bytesToCopy {
			bytesToCopy = bytesRemaining
		}

		if len(*data) < offset+bytesToCopy {
			newSlice := make([]byte, offset+bytesToCopy)
			copy(newSlice, *data)
			*data = newSlice
		}

		copy((*data)[offset:offset+bytesToCopy], hb.buffer[hb.bufferPos:hb.bufferPos+bytesToCopy])

		offset += bytesToCopy
		hb.bufferPos += bytesToCopy
		bytesRemaining -= bytesToCopy
	}

	return maxSize, nil
}

func (hb *HashedBlock) readHashedBlock() (bool, error) {
	indexBytes := make([]byte, 4)
	_, err := hb.cipherStream.ReadData(&indexBytes, 4)

	if err != nil {
		return false, fmt.Errorf("unable to read block index: %s", err)
	}

	var index uint32
	bufIndex := bytes.NewReader(indexBytes)
	if err := binary.Read(bufIndex, binary.LittleEndian, &index); err != nil {
		return false, fmt.Errorf("index read failed: %s", err)
	}

	if index != hb.blockIndex {
		return false, fmt.Errorf("invalid block index: %d -> %d", index, hb.blockIndex)
	}

	hash := make([]byte, 32)
	_, err = hb.cipherStream.ReadData(&hash, 32)

	if err != nil {
		return false, fmt.Errorf("unable to read hash: %s", err)
	}

	if len(hash) != 32 {
		return false, fmt.Errorf("invalid hash size: %d %v ", len(hash), hash)
	}

	blockSizeBytes := make([]byte, 4)
	_, err = hb.cipherStream.ReadData(&blockSizeBytes, 4)

	if err != nil {
		return false, fmt.Errorf("unable to read block size: %s", err)
	}

	var blockSize uint32
	bufBlockSize := bytes.NewReader(blockSizeBytes)
	if err := binary.Read(bufBlockSize, binary.LittleEndian, &blockSize); err != nil {
		return false, fmt.Errorf("block size read failed: %s", err)
	}

	if blockSize < 0 {
		return false, fmt.Errorf("invalid block size")
	}

	if blockSize == 0 {
		if bytes.Count(hash, []byte{0}) != 32 {
			return false, fmt.Errorf("invalid hash of final block")
		}

		// EOF
		hb.eof = true
		return false, nil
	}

	_, err = hb.cipherStream.ReadData(&hb.buffer, int(blockSize))

	if err != nil {
		return false, fmt.Errorf("unable to buffer: %s", err)
	}

	if len(hb.buffer) != int(blockSize) {
		return false, fmt.Errorf("block too short")
	}

	h := sha256.New()
	h.Write(hb.buffer)
	bufferHash := h.Sum(nil)

	if !bytes.Equal(hash, bufferHash) {
		return false, fmt.Errorf("mismatch between hash and data")
	}

	hb.bufferPos = 0
	hb.blockIndex++

	return true, nil
}
