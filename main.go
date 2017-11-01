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
	db        = kingpin.Flag("db", "Keepassx database").Required().File()
	term      = kingpin.Flag("term", "Search term").Required().Short('s').String()
	chrs      = kingpin.Flag("chrs", "Copy characters from password [2,6,7..]").Short('c').String()
	keyfile   = kingpin.Flag("keyfile", "Key file").Short('k').File()
	debug     = kingpin.Flag("debug", "Enable debug mode").Short('d').Bool()
	clipboard = kingpin.Flag("clipboard", "Copy to clipboard").Short('x').Bool()
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

	entry, err := search.Database(reader.XMLReader, *term)
	if err != nil {
		log.Fatalf("search database error: %s", err)
	}

	if entry == nil {
		log.Fatalf("Search term: '%s' not found\n", *term)
	} else {
		// fields
		dataHeader := []string{"UUID", "Title", "Username", "URL", "Notes"}
		dataEntry := []string{entry.UUID, entry.Title, entry.Username, entry.URL, entry.Notes}

		// Extract x characters from password
		if len(*chrs) > 0 {
			output.Extract(entry, *chrs)
		}

		// Copy password to clipboard
		if *clipboard {
			cp := output.GetClipboard()
			if cp == nil {
				log.Fatalf("unable to identify os to copy to clipboard")
			}

			if err := cp.CopyProcess(entry.PlainTextPassword); err != nil {
				log.Fatalf("unable to copy password to clipboard: %s", err)
			}

			fmt.Println("password copied to clipboard")
		} else {
			dataEntry = append(dataEntry, entry.PlainTextPassword)
			dataHeader = append(dataHeader, "Password")
		}

		output.Table(dataHeader, [][]string{dataEntry})
	}
}
