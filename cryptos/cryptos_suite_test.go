package cryptos_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCryptos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cryptos Suite")
}
