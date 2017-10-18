package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/simonhayward/gkeepassxreader/format"
	"github.com/simonhayward/gkeepassxreader/input"
	"github.com/simonhayward/gkeepassxreader/keys"
	"github.com/simonhayward/gkeepassxreader/output"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "0.0.1"
)

var (
	db        = kingpin.Flag("db", "Keepassx database").Required().File()
	search    = kingpin.Flag("search", "Search for title").Required().Short('s').String()
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

	var pass string
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Fprint(os.Stderr, "Password (press enter for no password): ")
		stdinPassword, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Fprint(os.Stderr, "\n")

		if err != nil {
			log.Fatalf("error reading password from stdin: %s\n", err)
		}

		pass = string(stdinPassword)
	}

	i := input.Params{
		Db:       *db,
		Keyfile:  *keyfile,
		Password: &pass,
		Search:   *search,
		Chrs:     *chrs,
	}

	reader, err := OpenDatabase(i)
	if err != nil {
		log.Fatalf("open database error: %s", err)
	}

	entry, err := reader.SearchDatabase(i.Search, i.Chrs)
	if err != nil {
		log.Fatalf("search database error: %s", err)
	}

	if entry == nil {
		log.Fatalf("Search term: '%s' not found\n", i.Search)
	} else {
		fmt.Printf("Found\nTitle: %s\n", entry.Title)
		if *clipboard {
			cp := output.GetClipboard()
			if err := cp.CopyProcess(entry.PlainTextPassword); err != nil {
				log.Fatalf("unable to copy password to clipboard: %s", err)
			}
			fmt.Println("password copied to clipboard")
		} else {
			fmt.Printf("Password: %s\n", entry.PlainTextPassword)
		}
	}
}

// OpenDatabase with input params
func OpenDatabase(i input.Provider) (*format.KeePass2Reader, error) {
	masterKey := databaseKey(i)

	k := format.NewKeePass2Reader()
	err := k.ReadDatabase(i.GetDb(), masterKey)
	if err != nil {
		return nil, fmt.Errorf("read database error: %s", err.Error())
	}

	return k, nil
}

func databaseKey(i input.Provider) *keys.CompositeKey {

	var masterKey keys.CompositeKey

	if len(*i.GetPassword()) > 0 {
		pk := &keys.PasswordKey{}
		pk.SetPassword(*i.GetPassword())
		masterKey.AddKey(pk)
	}

	if i.GetKeyfile() != nil {
		kf := &keys.FileKey{}
		if !kf.Load(i.GetKeyfile()) {
			log.Warn("unable to load key file")
			return &keys.CompositeKey{}
		}
		masterKey.AddKey(kf)
	}

	return &masterKey
}
