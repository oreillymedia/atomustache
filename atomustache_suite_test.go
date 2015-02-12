package atomustache_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAtomustache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Atomustache Suite")
}
