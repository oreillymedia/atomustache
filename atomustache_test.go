package atomustache_test

import (
	"github.com/oreillymedia/atomustache"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Atomustache", func() {

  It("should render HTML string from atomic design folder structure", func() {
    renderer := atomustache.New("test_templates")
    result := renderer.RenderView("topics/show", map[string]string{"name":"Rune"})
    Expect(result).To(Equal("This is show Rune: This is organism Rune: This is molecules Rune. "))
  })

})
