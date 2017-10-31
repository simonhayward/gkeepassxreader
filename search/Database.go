package search

import (
	"fmt"
	"sort"
	"strings"

	"github.com/simonhayward/gkeepassxreader/format"
)

//ByTitle for searching
type ByTitle struct{ format.Entries }

//Less for comparisons
func (e ByTitle) Less(i, j int) bool { return e.Entries[i].Title < e.Entries[j].Title }

// Database searches the xml database for specified search term
func Database(xmlReader *format.KeePass2XmlReader, searchTerm string) (*format.Entry, error) {

	entries := []format.Entry{}
	randomBytesOffset := 0
	if err := xmlReader.ReadGroups(&entries, xmlReader.KeePass2XmlFile.Root.Groups, &randomBytesOffset); err != nil {
		return nil, fmt.Errorf("unable to read groups: %s", err)
	}

	sort.Sort(ByTitle{entries})

	var titles []string
	for _, e := range entries {
		titles = append(titles, e.Title)
	}

	idx := searchTitles(searchTerm, titles)
	if idx < len(titles) {
		plaintext, err := xmlReader.KeePass2RandomStream.Process(entries[idx].RandomOffset, []byte(entries[idx].CipherText))
		if err != nil {
			return nil, fmt.Errorf("unable to decode password: %s", err)
		}
		entries[idx].PlainTextPassword = string(plaintext)
		return &entries[idx], nil
	}

	return nil, nil
}

func searchTitles(searchTerm string, titles []string) int {
	p1, p2 := make(chan int, 1), make(chan int, 1)
	l := len(titles)

	go searchExact(searchTerm, titles, p1)
	go searchLowerCase(searchTerm, titles, p2)

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
