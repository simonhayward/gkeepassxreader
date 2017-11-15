package format

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

var (
	innerStreamSalsa20Iv = []byte{0xE8, 0x30, 0x09, 0x4B, 0x97, 0x20, 0x5D, 0x2A}
)

type value struct {
	Data      string `xml:",chardata"`
	Protected string `xml:"Protected,attr"`
}

type stringEntry struct {
	XMLName xml.Name `xml:"String"`
	Key     string   `xml:"Key"`
	Value   value    `xml:"Value"`
}

type entry struct {
	XMLName        xml.Name      `xml:"Entry"`
	UUID           string        `xml:"UUID"`
	StringEntry    []stringEntry `xml:"String"`
	HistoryEntries []entry       `xml:"History>Entry"`
}

type group struct {
	XMLName xml.Name `xml:"Group"`
	UUID    string   `xml:"UUID"`
	Name    string   `xml:"Name"`
	Entry   []entry  `xml:"Entry"`
	Groups  []group  `xml:"Group"`
}

type root struct {
	XMLName xml.Name `xml:"Root"`
	Groups  []group  `xml:"Group"`
}

type meta struct {
	XMLName    xml.Name `xml:"Meta"`
	HeaderHash string   `xml:"HeaderHash"`
}

//KeePass2XmlFile represents the xml file
type KeePass2XmlFile struct {
	XMLName xml.Name `xml:"KeePassFile"`
	Root    root     `xml:"Root"`
	Meta    meta     `xml:"Meta"`
}

//KeePass2XmlReader reads the xml
type KeePass2XmlReader struct {
	KeePass2XmlFile      KeePass2XmlFile
	KeePass2RandomStream *KeePass2RandomStream
}

//NewKeePass2XmlReader creates a new reader
func NewKeePass2XmlReader(xmlDevice io.Reader, key *[32]byte) (*KeePass2XmlReader, error) {

	data, err := ioutil.ReadAll(xmlDevice)
	if err != nil {
		return nil, fmt.Errorf("read xml error: %s", err)
	}

	f := KeePass2XmlFile{}
	err = xml.Unmarshal(data, &f)
	if err != nil {
		return nil, fmt.Errorf("unmarshal xml error: %s", err)
	}

	return &KeePass2XmlReader{
		KeePass2XmlFile:      f,
		KeePass2RandomStream: NewKeePass2RandomStream(innerStreamSalsa20Iv, key),
	}, nil
}

//HeaderHash return hashed header
func (k *KeePass2XmlReader) HeaderHash() ([]byte, error) {

	data, err := base64.StdEncoding.DecodeString(k.KeePass2XmlFile.Meta.HeaderHash)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %s", err)
	}

	return data, nil
}

func (k *KeePass2XmlReader) newEntryValue(sEntry stringEntry, randomBytesOffset *int) (*EntryValue, error) {

	if sEntry.Value.Protected == "True" && len(sEntry.Value.Data) > 0 {
		cipherText, err := base64.StdEncoding.DecodeString(sEntry.Value.Data)
		if err != nil {
			return nil, fmt.Errorf("ciphertext decode err: %s", err)
		}

		e := &EntryValue{
			Data:         sEntry.Value.Data,
			Protected:    true,
			CipherText:   cipherText,
			RandomOffset: *randomBytesOffset,
		}

		*randomBytesOffset += len(cipherText)

		return e, nil
	}

	return &EntryValue{
		Data:      sEntry.Value.Data,
		Protected: false,
		PlainText: sEntry.Value.Data,
	}, nil

}

func (k *KeePass2XmlReader) readEntries(entries *[]Entry, rEntries []entry, entryGroup group, randomBytesOffset *int) error {

	for _, entry := range rEntries {
		var title, password, username, url, notes *EntryValue
		var err error

		for _, sEntry := range entry.StringEntry {
			switch sEntry.Key {
			case "Notes":
				notes, err = k.newEntryValue(sEntry, randomBytesOffset)
				if err != nil {
					return fmt.Errorf("Notes entry value failed: %s", err)
				}
			case "Password":
				password, err = k.newEntryValue(sEntry, randomBytesOffset)
				if err != nil {
					return fmt.Errorf("Password entry value failed: %s", err)
				}
			case "Title":
				title, err = k.newEntryValue(sEntry, randomBytesOffset)
				if err != nil {
					return fmt.Errorf("Title entry value failed: %s", err)
				}
			case "URL":
				url, err = k.newEntryValue(sEntry, randomBytesOffset)
				if err != nil {
					return fmt.Errorf("URL entry value failed: %s", err)
				}
			case "UserName":
				username, err = k.newEntryValue(sEntry, randomBytesOffset)
				if err != nil {
					return fmt.Errorf("UserName entry value failed: %s", err)
				}
			}
		}

		uuid, err := base64.StdEncoding.DecodeString(entry.UUID)
		if err != nil {
			return fmt.Errorf("base64 decode for uuid failed: %s", err)
		}

		e := Entry{
			UUID:     hex.EncodeToString(uuid),
			Title:    title,
			Group:    entryGroup.Name,
			Password: password,
			Username: username,
			URL:      url,
			Notes:    notes,
		}

		*entries = append(*entries, e)

		// history entries are needed for randomBytes offset in random stream
		if len(entry.HistoryEntries) > 0 {
			if err := k.readEntries(entries, entry.HistoryEntries, entryGroup, randomBytesOffset); err != nil {
				return err
			}
		}
	}
	return nil
}

//ReadGroups iterates over database groups
func (k *KeePass2XmlReader) ReadGroups(entries *[]Entry, groups []group, randomBytesOffset *int) error {
	for _, group := range groups {
		if len(group.Entry) > 0 {
			if err := k.readEntries(entries, group.Entry, group, randomBytesOffset); err != nil {
				return err
			}
		}

		if len(group.Groups) > 0 {
			if err := k.ReadGroups(entries, group.Groups, randomBytesOffset); err != nil {
				return err
			}
		}
	}
	return nil
}
