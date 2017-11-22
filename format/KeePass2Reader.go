package format

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/simonhayward/gkeepassxreader/core"
	"github.com/simonhayward/gkeepassxreader/keys"
	"github.com/simonhayward/gkeepassxreader/streams"
	log "github.com/sirupsen/logrus"
)

const (
	keepass1Signature1 uint32 = 0x9AA2D903
	keepass1Signature2 uint32 = 0xB54BFB65

	keepass2Signature1 uint32 = 0x9AA2D903
	keepass2Signature2 uint32 = 0xB54BFB67

	keepass2FileVersion                    = 0x00030001
	keepass2FileVersionMin                 = 0x00020000
	keepass2FileVersionCriticalMask uint32 = 0xFFFF0000

	// HeaderFieldID
	keepass2EndOfHeader         = 0
	keepass2Comment             = 1
	keepass2CipherID            = 2
	keepass2CompressionFlags    = 3
	keepass2MasterSeed          = 4
	keepass2TransformSeed       = 5
	keepass2TransformRounds     = 6
	keepass2EncryptionIV        = 7
	keepass2ProtectedStreamKey  = 8
	keepass2StreamStartBytes    = 9
	keepass2InnerRandomStreamID = 10

	// ProtectedStreamAlgo
	keepass2Salsa20 = 2
)

//KeePass2Reader represents a KeePass2Reader
type KeePass2Reader struct {
	Db                 *core.Database
	headerEnd          bool
	XMLReader          *KeePass2XmlReader
	masterSeed         []byte
	transformSeed      []byte
	encryptionIV       []byte
	streamStartBytes   []byte
	protectedStreamKey []byte
	headerStoredData   []byte
}

//NewKeePass2Reader with default values
func NewKeePass2Reader() *KeePass2Reader {
	return &KeePass2Reader{
		Db: core.NewDatabase(),
	}
}

// OpenDatabase with key
func OpenDatabase(masterKey *keys.CompositeKey, dbFile *os.File) (*KeePass2Reader, error) {
	k := NewKeePass2Reader()
	err := k.ReadDatabase(dbFile, masterKey)
	if err != nil {
		return nil, fmt.Errorf("read database error: %s", err.Error())
	}

	return k, nil
}

//ReadDatabase reads the input database
func (k *KeePass2Reader) ReadDatabase(db *os.File, compositeKey *keys.CompositeKey) error {

	if err := k.CheckSignature(db); err != nil {
		return errors.Wrap(err, "Signature check failed")
	}

	version, err := k.CheckVersion(db)
	if err != nil {
		return errors.Wrap(err, "Version check failed")
	}

	for {
		continueLoop, err := k.ReadHeaders(db)
		if err != nil {
			return errors.Wrap(err, "Reading headers failed")
		}
		if continueLoop == false {
			break
		}
	}

	if err := k.CheckHeaders(); err != nil {
		return errors.Wrap(err, "Header check failed")
	}

	if err := k.Db.SetKey(compositeKey, k.transformSeed); err != nil {
		return errors.Wrap(err, "Unable to calculate master key")
	}

	h := sha256.New()
	h.Write(k.masterSeed)
	h.Write(k.Db.TransformedMasterKey)
	finalKey := h.Sum(nil)

	block, err := aes.NewCipher(finalKey)
	if err != nil {
		return errors.Wrap(err, "New AES Cipher error")
	}

	cipherStream, err := streams.NewSymmetricCipherStream(block, k.encryptionIV, db, streams.DirectionDecrypt)
	if err != nil {
		return errors.Wrap(err, "Cipher stream error")
	}

	var realStart []byte
	cipherStream.ReadData(&realStart, 32)

	if !bytes.Equal(realStart, k.streamStartBytes) {
		return errors.New("Wrong key or database file is corrupt")
	}

	/*
		Hashed stream
	*/
	hashBlock := streams.NewHashedBlock(cipherStream.BlockMode, cipherStream)
	var result []byte
	var bytesRead int
	byteChunks := 65500

	for {
		bytesRead, err = hashBlock.ReadData(&result, byteChunks)
		if err != nil {
			return err
		}
		if bytesRead == 0 {
			break
		}
	}

	var xmlDevice io.Reader
	if k.Db.CompressionAlgo == core.CompressionNone {
		log.Debugf("no compression set")
		xmlDevice = bytes.NewReader(result)
	} else {
		log.Debugf("compression set")

		buf := bytes.NewBuffer(result)
		zr, err := gzip.NewReader(buf)

		if err != nil {
			return errors.Wrap(err, "gzip new reader failed")
		}

		b, err := ioutil.ReadAll(zr)
		if err != nil {
			return errors.Wrap(err, "xml error")
		}
		xmlDevice = bytes.NewReader(b)
	}

	randomKey := sha256.Sum256(k.protectedStreamKey)
	k.XMLReader, err = NewKeePass2XmlReader(xmlDevice, &randomKey)
	if err != nil {
		return errors.Wrap(err, "keepass2xml reader creation failed")
	}

	xmlHeaderHash, err := k.XMLReader.HeaderHash()
	if err != nil {
		return errors.Wrap(err, "xml header hash error")
	}

	if !(version < keepass2FileVersion || len(xmlHeaderHash) > 0) {
		return errors.New("xml header hash error")
	}

	if len(xmlHeaderHash) > 0 {
		hh := sha256.New()
		hh.Write(k.headerStoredData)
		headerHash := hh.Sum(nil)

		if !bytes.Equal(headerHash, xmlHeaderHash) {
			return errors.New("header doesn't match hash")
		}
	}

	return nil
}

//CheckSignature inspects to see if this is a valid keepass database
func (k *KeePass2Reader) CheckSignature(db *os.File) error {

	signature1Bytes := make([]byte, 4)
	_, err := db.Read(signature1Bytes)

	if err != nil {
		return errors.Wrap(err, "unable to read signature1")
	}

	k.headerStoredData = append(k.headerStoredData, signature1Bytes...)

	var signature1 uint32
	buf1 := bytes.NewReader(signature1Bytes)
	if err := binary.Read(buf1, binary.LittleEndian, &signature1); err != nil {
		return errors.Wrap(err, "signature1 read failed")
	}

	if signature1 != keepass2Signature1 {
		return errors.New("not a KeePass database")
	}

	signature2Bytes := make([]byte, 4)
	_, err = db.Read(signature2Bytes)

	if err != nil {
		return errors.Wrap(err, "unable to read signature2")
	}

	k.headerStoredData = append(k.headerStoredData, signature2Bytes...)

	var signature2 uint32
	buf2 := bytes.NewReader(signature2Bytes)
	if err := binary.Read(buf2, binary.LittleEndian, &signature2); err != nil {
		return errors.Wrap(err, "signature2 read failed")
	}

	if signature2 == keepass1Signature2 {
		return errors.New("the selected file is an old KeePass 1 database (.kdb)")
	} else if signature2 != keepass2Signature2 {
		return errors.New("not a KeePass database")
	}

	return nil
}

//CheckVersion validates the keepass version supported
func (k *KeePass2Reader) CheckVersion(db *os.File) (uint32, error) {
	versionBytes := make([]byte, 4)
	_, err := db.Read(versionBytes)

	if err != nil {
		return 0, errors.Wrap(err, "unable to read version")
	}

	k.headerStoredData = append(k.headerStoredData, versionBytes...)

	buf := bytes.NewReader(versionBytes)
	var version uint32

	if err := binary.Read(buf, binary.LittleEndian, &version); err != nil {
		return 0, errors.Wrap(err, "binary.Read failed")
	}

	version = version & keepass2FileVersionCriticalMask

	var maxVersion = keepass2FileVersion & keepass2FileVersionCriticalMask

	log.Debugf("checking versions. min: %d max: %d", keepass2FileVersionMin, maxVersion)

	if (version < keepass2FileVersionMin) || (version > maxVersion) {
		return 0, errors.New("unsupported KeePass database version")
	}

	log.Debugf("version: %d", version)

	return version, nil
}

//CheckHeaders checks if all required headers were present
func (k *KeePass2Reader) CheckHeaders() error {
	if len(k.masterSeed) == 0 || len(k.transformSeed) == 0 || len(k.encryptionIV) == 0 ||
		len(k.streamStartBytes) == 0 || len(k.protectedStreamKey) == 0 ||
		len(k.Db.Cipher.Data) == 0 {
		return errors.New("missing database headers")
	}
	return nil
}

// ReadHeaders extracts the headers of the database
func (k *KeePass2Reader) ReadHeaders(db *os.File) (bool, error) {
	headerEnd := false

	fieldIDArray := make([]byte, 1)
	_, err := db.Read(fieldIDArray)

	if err != nil {
		return false, errors.Wrap(err, "unable to read fieldIDArray")
	}

	k.headerStoredData = append(k.headerStoredData, fieldIDArray...)

	if len(fieldIDArray) != 1 {
		return false, errors.New("invalid header id size")
	}
	var fieldID = fieldIDArray[0]
	log.Debugf("header field id: %d", fieldID)

	var fieldLen uint16
	if err := binary.Read(db, binary.LittleEndian, &fieldLen); err != nil {
		return false, errors.Wrap(err, "invalid header field length")
	}

	var h, l uint8 = uint8(fieldLen >> 8), uint8(fieldLen & 0xff)
	k.headerStoredData = append(k.headerStoredData, []byte{l, h}...)

	log.Debugf("header field length: %d", fieldLen)

	var fieldData []byte
	if fieldLen != 0 {
		fieldData = make([]byte, int(fieldLen))
		n, err := db.Read(fieldData)
		if err != nil {
			return false, errors.New("unable to read field length")
		}
		if n != int(fieldLen) {
			return false, errors.New("invalid header data length")
		}
	}

	k.headerStoredData = append(k.headerStoredData, fieldData...)

	switch fieldID {
	case keepass2EndOfHeader:
		headerEnd = true
		log.Debugf("end of header: %d", fieldID)
	case keepass2CipherID:
		log.Debugf("setting cipher: FieldID: %d fieldData len: %d", fieldID, len(fieldData))
		if err = k.setCipher(fieldData); err != nil {
			return false, errors.Wrap(err, "cipher not set")
		}
	case keepass2CompressionFlags:
		log.Debugf("setting compression flags: fieldID: %d", fieldID)
		if err = k.setCompressionFlags(fieldData); err != nil {
			return false, errors.Wrap(err, "compression flags not set")
		}
	case keepass2MasterSeed:
		log.Debugf("setting master seed: %d", fieldID)
		if err = k.setMasterSeed(fieldData); err != nil {
			return false, errors.Wrap(err, "master seed not set")
		}
	case keepass2TransformSeed:
		log.Debugf("setting transform seed: %d", fieldID)
		if err = k.setTransformSeed(fieldData); err != nil {
			return false, errors.Wrap(err, "transform seed not set")
		}
	case keepass2TransformRounds:
		log.Debugf("setting transform rounds: %d", fieldID)
		if err = k.setTransformRounds(fieldData); err != nil {
			return false, errors.Wrap(err, "transform rounds not set")
		}
	case keepass2EncryptionIV:
		log.Debugf("set setEncryptionIV: %d", fieldID)
		if err = k.setEncryptionIV(fieldData); err != nil {
			return false, errors.Wrap(err, "encryptionIV not set")
		}
	case keepass2ProtectedStreamKey:
		log.Debugf("setting protected stream key: %d", fieldID)
		if err = k.setProtectedStreamKey(fieldData); err != nil {
			return false, errors.Wrap(err, "protected stream key not set")
		}
	case keepass2StreamStartBytes:
		log.Debugf("setting StreamStartBytes: %d", fieldID)
		if err = k.setStreamStartBytes(fieldData); err != nil {
			return false, errors.Wrap(err, "stream start bytes not set")
		}
	case keepass2InnerRandomStreamID:
		log.Debugf("setting InnerRandomStreamID: %d", fieldID)
		if err = k.setInnerRandomStreamID(fieldData); err != nil {
			return false, errors.Wrap(err, "innerRandomStreamID not set")
		}
	default:
		log.Errorf("unknown header field read: id=%d", fieldID)
		return false, errors.Wrapf(err, "unknown header field: %d", fieldID)
	}

	return !headerEnd, nil
}

func (k *KeePass2Reader) setCipher(b []byte) error {

	if len(b) != core.UUIDLength {
		return errors.Errorf("invalid cipher uuid length: %d expected: %d", len(b), core.UUIDLength)
	}

	cipher := core.UUID{
		Data: b,
	}

	if !bytes.Equal(b, core.Keepass2CipherAes) {
		return errors.New("unsupported cipher")
	}

	k.Db.Cipher = cipher

	return nil
}

func (k *KeePass2Reader) setCompressionFlags(b []byte) error {
	if len(b) != 4 {
		return errors.New("invalid compression flags length")
	}

	var id uint32
	buf := bytes.NewBuffer(b)
	if err := binary.Read(buf, binary.LittleEndian, &id); err != nil {
		return errors.Wrap(err, "binary.Read failed")
	}

	if id > core.CompressionAlgorithmMax {
		return errors.New("unsupported compression algorithm")
	}

	k.Db.CompressionAlgo = id
	return nil
}

func (k *KeePass2Reader) setMasterSeed(b []byte) error {
	if len(b) != 32 {
		return errors.New("invalid master seed size")
	}

	k.masterSeed = b
	return nil
}

func (k *KeePass2Reader) setTransformSeed(b []byte) error {
	if len(b) != 32 {
		return errors.New("invalid transform seed size")
	}

	k.transformSeed = b
	return nil
}

func (k *KeePass2Reader) setTransformRounds(b []byte) error {
	if len(b) != 8 {
		return errors.New("invalid transform rounds size")
	}

	var rounds uint64
	buf := bytes.NewBuffer(b)
	if err := binary.Read(buf, binary.LittleEndian, &rounds); err != nil {
		return errors.Wrap(err, "binary.Read failed")
	}

	if k.Db.TransformRounds != rounds {
		log.Debugf("updating transform rounds from: %d to: %d", k.Db.TransformRounds, rounds)
		k.Db.TransformRounds = rounds
	}

	return nil
}

func (k *KeePass2Reader) setEncryptionIV(b []byte) error {
	if len(b) != 16 {
		return errors.New("invalid encryption iv size")
	}

	k.encryptionIV = b
	return nil
}

func (k *KeePass2Reader) setProtectedStreamKey(b []byte) error {
	if len(b) != 32 {
		return errors.New("invalid stream key size")
	}

	k.protectedStreamKey = b
	return nil
}

func (k *KeePass2Reader) setStreamStartBytes(b []byte) error {
	if len(b) != 32 {
		return errors.New("invalid start bytes size")
	}

	k.streamStartBytes = b
	return nil
}

func (k *KeePass2Reader) setInnerRandomStreamID(b []byte) error {
	if len(b) != 4 {
		return errors.New("invalid random stream id size")
	}

	var id uint32
	buf := bytes.NewBuffer(b)
	if err := binary.Read(buf, binary.LittleEndian, &id); err != nil {
		return errors.Wrap(err, "binary.Read failed")
	}

	if id != keepass2Salsa20 {
		return errors.New("unsupported random stream algorithm")
	}

	log.Debugf("setting inner random stream id: %d", id)
	return nil
}
