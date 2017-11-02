package search

import (
	"fmt"
	"strings"

	"github.com/simonhayward/gkeepassxreader/format"
)

// Database searches the xml database for specified search term
func Database(xmlReader *format.KeePass2XmlReader, searchTerm string) (*format.Entry, error) {

	entries := []format.Entry{}
	randomBytesOffset := 0
	if err := xmlReader.ReadGroups(&entries, xmlReader.KeePass2XmlFile.Root.Groups, &randomBytesOffset); err != nil {
		return nil, fmt.Errorf("unable to read groups: %s", err)
	}

	idx := search(searchTerm, entries)
	if idx < len(entries) {
		plaintext, err := xmlReader.KeePass2RandomStream.Process(entries[idx].RandomOffset, []byte(entries[idx].CipherText))
		if err != nil {
			return nil, fmt.Errorf("unable to decode password: %s", err)
		}
		entries[idx].PlainTextPassword = string(plaintext)
		return &entries[idx], nil
	}

	return nil, nil
}

func search(searchTerm string, entries []format.Entry) int {
	var titles, uuids []string
	for _, e := range entries {
		titles = append(titles, e.Title)
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
