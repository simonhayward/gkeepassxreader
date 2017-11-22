package entries_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEntries(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Entries Suite")
}
