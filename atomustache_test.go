package atomustache_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/oreillymedia/atomustache"
)

func TestAtomustache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Atomustache Suite")
}

var _ = Describe("Atomustache", func() {

	Describe("#RenderView", func() {

		It("should render atomic templates into view", func() {
			renderer, _ := atomustache.New(
				"./test_templates/styleguide",
				"./test_templates/layouts",
				"./test_templates/views",
				".mustache",
			)
			result, err := renderer.RenderView("folder/oneview", map[string]string{"name": "Rune"})
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("This is show Rune: This is organism Rune: This is molecules Rune. "))
		})

		It("should render partials no matter their order", func() {
			renderer, _ := atomustache.New(
				"./test_templates/styleguide",
				"./test_templates/layouts",
				"./test_templates/views",
				".mustache",
			)
			result, err := renderer.RenderView("folder/anotherview", map[string]string{"name": "Rune"})
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("This is show Rune: This is molecules with: This is molecules Rune."))
		})

	})

	Describe("#RenderViewInLayout", func() {

		It("should render view in layout with atomic templates", func() {
			renderer, _ := atomustache.New(
				"./test_templates/styleguide",
				"./test_templates/layouts",
				"./test_templates/views",
				".mustache",
			)
			result, err := renderer.RenderViewInLayout("folder/oneview", "test", map[string]string{"name": "Rune"})
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("Before This is show Rune: This is organism Rune: This is molecules Rune.  After"))
		})

	})

})
