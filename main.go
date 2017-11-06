package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/simonhayward/gkeepassxreader/format"
	"github.com/simonhayward/gkeepassxreader/keys"
	"github.com/simonhayward/gkeepassxreader/output"
	"github.com/simonhayward/gkeepassxreader/search"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "0.0.1"
)

var (
	db      = kingpin.Flag("db", "Keepassx database").Required().File()
	keyfile = kingpin.Flag("keyfile", "Key file").Short('k').File()
	debug   = kingpin.Flag("debug", "Enable debug mode").Short('d').Bool()

	cmdSearch       = kingpin.Command("search", "Search for an entry")
	searchTerm      = cmdSearch.Arg("term", "Search by title or UUID").Required().String()
	searchChrs      = cmdSearch.Flag("chrs", "Copy selected characters from password [2,6,7..]").Short('c').String()
	searchClipboard = cmdSearch.Flag("clipboard", "Copy to clipboard").Short('x').Bool()
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

	switch kingpin.Parse() {
	case cmdSearch.FullCommand():
		entry, err := search.Database(reader.XMLReader, *searchTerm)
		if err != nil {
			log.Fatalf("search database error: %s", err)
		}

		if entry == nil {
			log.Fatalf("Search term: '%s' not found\n", *searchTerm)
		} else {
			fields := output.NewDefaults()
			fields.Entries([]*format.Entry{entry})

			// Extract x characters from password
			if len(*searchChrs) > 0 {
				output.Extract(entry, *searchChrs)
			}

			// Copy password to clipboard
			if *searchClipboard {
				cp := output.GetClipboard()
				if cp == nil {
					log.Fatalf("unable to identify os to copy to clipboard")
				}

				if err := cp.CopyProcess(entry.PlainTextPassword); err != nil {
					log.Fatalf("unable to copy password to clipboard: %s", err)
				}

				fmt.Println("password copied to clipboard")
			} else {
				fields.Data[0] = append(fields.Data[0], entry.PlainTextPassword)
				fields.Header = append(fields.Header, "Password")
			}

			output.Table(fields.Header, fields.Data)
		}
	}
}
