package format

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"
)

var (
	innerStreamSalsa20Iv = []byte{0xE8, 0x30, 0x09, 0x4B, 0x97, 0x20, 0x5D, 0x2A}
)

type value struct {
	Data      string `xml:"Value"`
	Protected string `xml:"Protected,attr"`
}

type stringEntry struct {
	XMLName xml.Name `xml:"String"`
	Key     string   `xml:"Key"`
	Value   string   `xml:"Value"`
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

// Entry representation
type Entry struct {
	Title             string
	Password          string
	PlainTextPassword string
	Protected         bool
	UUID              string
	RandomOffset      int
	CipherText        []byte
}

//Entries contains entries
type Entries []Entry

//Len length of entries
func (e Entries) Len() int { return len(e) }

//Swap entries
func (e Entries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

//ByTitle for searching
type ByTitle struct{ Entries }

//Less for comparisons
func (e ByTitle) Less(i, j int) bool { return e.Entries[i].Title < e.Entries[j].Title }

func (k *KeePass2XmlReader) readEntries(entries *[]Entry, rEntries []entry, randomBytesOffset *int) error {
	for _, entry := range rEntries {
		var title, password string
		for _, sEntry := range entry.StringEntry {
			if sEntry.Key == "Title" {
				title = sEntry.Value
			} else if sEntry.Key == "Password" {
				password = sEntry.Value
			}
		}

		cipherText, err := base64.StdEncoding.DecodeString(string(password))
		if err != nil {
			return fmt.Errorf("ciphertext decode err: %s", err)
		}

		e := Entry{
			UUID:         entry.UUID,
			Title:        title,
			Password:     password,
			CipherText:   cipherText,
			RandomOffset: *randomBytesOffset,
		}

		*randomBytesOffset += len(cipherText)

		*entries = append(*entries, e)

		if len(entry.HistoryEntries) > 0 {
			if err := k.readEntries(entries, entry.HistoryEntries, randomBytesOffset); err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *KeePass2XmlReader) readGroups(entries *[]Entry, groups []group, randomBytesOffset *int) error {
	for _, group := range groups {
		if len(group.Entry) > 0 {
			if err := k.readEntries(entries, group.Entry, randomBytesOffset); err != nil {
				return err
			}
		}

		if len(group.Groups) > 0 {
			if err := k.readGroups(entries, group.Groups, randomBytesOffset); err != nil {
				return err
			}
		}
	}
	return nil
}

//Search reader
func (k *KeePass2XmlReader) Search(searchTerm string) (*Entry, error) {

	entries := []Entry{}
	randomBytesOffset := 0
	if err := k.readGroups(&entries, k.KeePass2XmlFile.Root.Groups, &randomBytesOffset); err != nil {
		return nil, fmt.Errorf("unable to read groups: %s", err)
	}

	sort.Sort(ByTitle{entries})

	var titles []string
	for _, e := range entries {
		titles = append(titles, e.Title)
	}

	idx := searchTitles(searchTerm, titles)
	if idx < len(titles) {
		plaintext, err := k.KeePass2RandomStream.Process(entries[idx].RandomOffset, []byte(entries[idx].CipherText))
		if err != nil {
			return nil, fmt.Errorf("unable to decode password: %s", err)
		}
		entries[idx].PlainTextPassword = string(plaintext)
		return &entries[idx], nil
	}

	return nil, nil
}

func searchTitles(searchTerm string, titles []string) int {
	s1, s2 := make(chan int, 1), make(chan int, 1)
	l := len(titles)

	go searchExact(searchTerm, titles, s1)
	go searchLowerCase(searchTerm, titles, s2)

	// priority
	i := <-s1
	if i < l {
		return i
	}

	// less priority
	i = <-s2
	if i < l {
		return i
	}

	return l
}

func searchExact(searchTerm string, terms []string, c chan int) {
	defer close(c)

	for i, v := range terms {
		if searchTerm == v {
			c <- i
			break
		}
	}
	c <- len(terms)
}

func searchLowerCase(searchTerm string, terms []string, c chan int) {
	defer close(c)

	searchTerm = strings.ToLower(searchTerm)
	for i, v := range terms {
		if searchTerm == strings.ToLower(v) {
			c <- i
			break
		}
	}
	c <- len(terms)
}

//HeaderHash return hashed header
func (k *KeePass2XmlReader) HeaderHash() ([]byte, error) {

	data, err := base64.StdEncoding.DecodeString(k.KeePass2XmlFile.Meta.HeaderHash)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %s", err)
	}

	return data, nil
}
