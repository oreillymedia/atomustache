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

  Describe("#RenderView", func() {
  
    It("should render atomic templates into view", func() {
      renderer := atomustache.New("test_templates")
      result := renderer.RenderView("folder/oneview", map[string]string{"name":"Rune"})
      Expect(result).To(Equal("This is show Rune: This is organism Rune: This is molecules Rune. "))
    })

    It("should render partials no matter their order", func() {
      renderer := atomustache.New("test_templates")
      result := renderer.RenderView("folder/anotherview", map[string]string{"name":"Rune"})
      Expect(result).To(Equal("This is show Rune: This is molecules with: This is molecules Rune."))
    })

  })

  Describe("#RenderViewInLayout", func() {

    It("should render view in layout with atomic templates", func() {
      renderer := atomustache.New("test_templates")
      result := renderer.RenderViewInLayout("folder/oneview", "test", map[string]string{"name":"Rune"})
      Expect(result).To(Equal("Before This is show Rune: This is organism Rune: This is molecules Rune.  After"))
    })

  })

})