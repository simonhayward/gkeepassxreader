package input

import "os"

//Provider represents the input
type Provider interface {
	GetDb() *os.File
	GetPassword() *string
	GetKeyfile() *os.File
	GetSearch() string
	GetChrs() string
}

// Params represents the input parameters
type Params struct {
	Db       *os.File
	Password *string
	Keyfile  *os.File
	Search   string
	Chrs     string
}

func (i Params) GetDb() *os.File {
	return i.Db
}
func (i Params) GetPassword() *string {
	return i.Password
}
func (i Params) GetKeyfile() *os.File {
	return i.Keyfile
}
func (i Params) GetSearch() string {
	return i.Search
}
func (i Params) GetChrs() string {
	return i.Chrs
}
