package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/simonhayward/gkeepassxreader/format"
	"github.com/simonhayward/gkeepassxreader/keys"
	"github.com/simonhayward/gkeepassxreader/output"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "0.0.8"
)

var (
	db      = kingpin.Flag("db", "Keepassx database").Required().File()
	keyfile = kingpin.Flag("keyfile", "Key file").Short('k').File()
	debug   = kingpin.Flag("debug", "Enable debug mode").Short('d').Bool()
	history = kingpin.Flag("history", "Include historical entries").Short('h').Bool()

	cmdSearch       = kingpin.Command("search", "Search for an entry")
	searchTerm      = cmdSearch.Arg("term", "Search by title or UUID").Required().String()
	searchChrs      = cmdSearch.Flag("chrs", "Copy selected characters from password [2,6,7..]").Short('c').String()
	searchClipboard = cmdSearch.Flag("clipboard", "Copy to clipboard").Short('x').Bool()

	cmdList = kingpin.Command("list", "List entries")
)

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	level := log.InfoLevel
	if *debug == true {
		level = log.DebugLevel
	}
	log.SetLevel(level)
	log.SetOutput(os.Stdout)

	var entryService format.EntryService
	var password string

	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Fprint(os.Stderr, "Password (press enter for no password): ")
		stdinPassword, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Fprint(os.Stderr, "\n")

		if err != nil {
			log.Fatalf("error reading password from stdin: %s\n", err)
		}

		password = string(stdinPassword)
	}

	reader, err := format.OpenDatabase(keys.MasterKey(password, *keyfile), *db)
	if err != nil {
		log.Fatalf("open database error: %s", err)
	}

	entryService = &format.EntryServiceOp{
		XMLReader:         reader.XMLReader,
		HistoricalEntries: *history,
	}

	switch kingpin.Parse() {
	case cmdSearch.FullCommand():
		entry, err := entryService.SearchByTerm(*searchTerm)
		if err != nil {
			log.Fatalf("search database error: %s", err)
		}

		if entry == nil {
			log.Fatalf("Search term: '%s' not found\n", *searchTerm)
		} else {
			fields := output.NewDefaults()
			fields.Entries([]format.Entry{*entry})

			// Extract characters from password
			if len(*searchChrs) > 0 {
				err := output.Extract(entry, *searchChrs)
				if err != nil {
					log.Fatalf("unable to extract characters: %s", err)
				}
			}

			// Copy password to clipboard
			if *searchClipboard {
				cp := output.GetClipboard()
				if cp == nil {
					log.Fatalf("unable to identify os to copy to clipboard")
				}

				if err := cp.CopyProcess(entry.Password.PlainText); err != nil {
					log.Fatalf("unable to copy password to clipboard: %s", err)
				}

				fmt.Println("password copied to clipboard")
			} else {
				fields.Data[0] = append(fields.Data[0], entry.Password.PlainText)
				fields.Header = append(fields.Header, "Password")
			}

			output.Table(fields.Header, fields.Data)
		}
	case cmdList.FullCommand():
		allEntries, err := entryService.List()
		if err != nil {
			log.Fatalf("list database error: %s", err)
		}

		fields := output.NewDefaults()
		fields.Entries(allEntries)
		output.Table(fields.Header, fields.Data)
	}
}
