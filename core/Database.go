package core

import (
	"github.com/simonhayward/gkeepassxreader/keys"
)

//UUID represents the unqiue identifier
type UUID struct {
	Data []byte
}

const (
	//UUIDLength is fixed
	UUIDLength int = 16

	// CompressionNone value
	CompressionNone = uint32(0)
	//CompressionGzip value
	CompressionGzip = uint32(1)

	// CompressionAlgorithmMax algo
	CompressionAlgorithmMax = CompressionGzip

	defaultTransformRounds = uint64(100000)
)

var (
	// Keepass2CipherAes == 31c1f2e6bf714350be5805216afc5aff
	Keepass2CipherAes = []byte{49, 193, 242, 230, 191, 113, 67, 80, 190, 88, 5, 33, 106, 252, 90, 255}
)

//DatabaseData is data
type DatabaseData struct {
	Cipher               UUID
	CompressionAlgo      uint32
	TransformSeed        []byte
	TransformRounds      uint64
	TransformedMasterKey []byte
	Key                  *keys.CompositeKey
	HasKey               bool
}

// Database is database
type Database struct {
	Data DatabaseData
}

// NewDatabase with default values
func NewDatabase() *Database {
	u := UUID{Data: Keepass2CipherAes}

	return &Database{
		Data: DatabaseData{
			Cipher:          u,
			CompressionAlgo: CompressionGzip,
			TransformRounds: defaultTransformRounds,
			HasKey:          false,
		},
	}
}

//SetKey sets up key
func (d *Database) SetKey(key *keys.CompositeKey, transformSeed []byte, updateChangedTime bool) error {

	var transformedMasterKey []byte

	transformedMasterKey, err := key.Transform(transformSeed, d.Data.TransformRounds)

	if err != nil {
		return err
	}

	d.Data.Key = key
	d.Data.TransformSeed = transformSeed
	d.Data.TransformedMasterKey = transformedMasterKey
	d.Data.HasKey = true

	return nil
}
