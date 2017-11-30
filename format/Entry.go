package format

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

//EntryValue represents an individual entry value
type EntryValue struct {
	Data         string
	Protected    bool
	PlainText    string
	RandomOffset int
	CipherText   []byte
}

// Entry represents a single Entry
type Entry struct {
	Group      string
	Title      *EntryValue
	Username   *EntryValue
	Password   *EntryValue
	URL        *EntryValue
	Notes      *EntryValue
	UUID       string
	Historical bool
}

//Entries represents a collection of Entry
type Entries []Entry

//EntryServiceOp handles the interactions with individual entries
type EntryServiceOp struct {
	XMLReader         *KeePass2XmlReader
	HistoricalEntries bool
}

//EntryService is an interface for interfacing with individual entries
type EntryService interface {
	List() ([]Entry, error)
	SearchByTerm(searchTerm string) (*Entry, error)
	Search(searchTerm string, entries []Entry) int
}

var _ EntryService = &EntryServiceOp{}

//List all entries
func (s *EntryServiceOp) List() ([]Entry, error) {
	entries := []Entry{}
	randomBytesOffset := 0
	if err := s.XMLReader.ReadGroups(&entries, s.XMLReader.KeePass2XmlFile.Root.Groups, &randomBytesOffset); err != nil {
		return nil, err
	}

	if !s.HistoricalEntries {
		entries = removeHistorical(entries)
	}

	if err := decodeEntries(s.XMLReader, entries, false); err != nil {
		return nil, err
	}

	return entries, nil
}

// SearchByTerm searches the xml database for specified search term
func (s *EntryServiceOp) SearchByTerm(searchTerm string) (*Entry, error) {

	entries := []Entry{}
	randomBytesOffset := 0
	if err := s.XMLReader.ReadGroups(&entries, s.XMLReader.KeePass2XmlFile.Root.Groups, &randomBytesOffset); err != nil {
		return nil, fmt.Errorf("unable to read groups: %s", err)
	}

	if !s.HistoricalEntries {
		entries = removeHistorical(entries)
	}

	// decode all titles
	for _, e := range entries {
		if e.Title.Protected {
			err := decodeEntryValue(s.XMLReader, e.Title)
			if err != nil {
				return nil, err
			}
		}
	}

	idx := s.Search(searchTerm, entries)
	if idx < len(entries) {
		if err := decodeEntries(s.XMLReader, []Entry{entries[idx]}, true); err != nil {
			return nil, err
		}

		return &entries[idx], nil
	}

	return nil, nil
}

//Search queries the given entries to find a match
func (s *EntryServiceOp) Search(searchTerm string, entries []Entry) int {
	var titles, uuids []string
	for _, e := range entries {
		titles = append(titles, e.Title.PlainText)
		uuids = append(uuids, e.UUID)
	}

	p1, p2, p3 := make(chan int, 1), make(chan int, 1), make(chan int, 1)
	l := len(entries)

	go searchExact(searchTerm, uuids, p1)
	go searchExact(searchTerm, titles, p2)
	go searchLowerCase(searchTerm, titles, p3)

	// priority #1
	i := <-p1
	if i < l {
		return i
	}

	// priority #2
	i = <-p2
	if i < l {
		return i
	}

	// priority #3
	i = <-p3
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
			return
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
			return
		}
	}
	c <- len(terms)
}

func decodeEntries(xmlReader *KeePass2XmlReader, entries []Entry, decodePassword bool) error {
	for idx := range entries {
		var entryValues []*EntryValue

		entryValues = []*EntryValue{
			entries[idx].Notes,
			entries[idx].URL,
			entries[idx].Username,
			entries[idx].Title,
		}

		if decodePassword {
			entryValues = append(entryValues, entries[idx].Password)
		}

		for _, ev := range entryValues {
			if ev != nil && ev.Protected {
				err := decodeEntryValue(xmlReader, ev)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func decodeEntryValue(xmlReader *KeePass2XmlReader, eValue *EntryValue) error {
	plaintext, err := xmlReader.KeePass2RandomStream.Process(eValue.RandomOffset, []byte(eValue.CipherText))
	if err != nil {
		return errors.Wrap(err, "entry value decode failed")
	}
	eValue.PlainText = string(plaintext)
	return nil
}

func removeHistorical(entries []Entry) []Entry {
	var notHistorical []Entry
	for _, e := range entries {
		if !e.Historical {
			notHistorical = append(notHistorical, e)
		}
	}
	return notHistorical
}
