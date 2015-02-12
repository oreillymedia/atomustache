package atomustache_test

import (
  "github.com/oreillymedia/atomustache"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestAtomustache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Atomustache Suite")
}

var _ = Describe("Atomustache", func() {

  It("#RenderView", func() {
    renderer := atomustache.New("test_templates")
    result := renderer.RenderView("topics/show", map[string]string{"name":"Rune"})
    Expect(result).To(Equal("This is show Rune: This is organism Rune: This is molecules Rune. "))
  })

  It("#RenderViewInLayout", func() {
    renderer := atomustache.New("test_templates")
    result := renderer.RenderViewInLayout("topics/show", "test", map[string]string{"name":"Rune"})
    Expect(result).To(Equal("Before This is show Rune: This is organism Rune: This is molecules Rune.  After"))
  })

})