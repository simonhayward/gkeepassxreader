package format

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

//Entries is a slice of Entry
type Entries []Entry

//Len length of entries
func (e Entries) Len() int { return len(e) }

//Swap entries
func (e Entries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
