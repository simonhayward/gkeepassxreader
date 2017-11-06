package format

// Entry representation
type Entry struct {
	Title             string
	Group             string
	Username          string
	Password          string
	PasswordProtected bool
	PlainTextPassword string
	UUID              string
	URL               string
	Notes             string
	RandomOffset      int
	CipherText        []byte
}

//Entries is a slice of Entry
type Entries []Entry
