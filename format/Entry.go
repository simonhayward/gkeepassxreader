package format

//EntryValue individual entry value
type EntryValue struct {
	Data         string
	Protected    bool
	PlainText    string
	RandomOffset int
	CipherText   []byte
}

// Entry representation
type Entry struct {
	Group    string
	Title    *EntryValue
	Username *EntryValue
	Password *EntryValue
	URL      *EntryValue
	Notes    *EntryValue
	UUID     string
}

//Entries is a slice of Entry
type Entries []Entry
