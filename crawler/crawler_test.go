package crawler

import (
	"fmt"
	"testing"

	tutil "../test_util"
)

func init() {
	crawlHostname = "example.com"
}

func TestReadHrefLink(t *testing.T) {
	url := "example.com"
	hrefElement := fmt.Sprintf(`<a href="%v">Example.com</a>`, url)

	tutil.Assert(t, url, hrefLink(hrefElement))
}
func TestReadImageLink(t *testing.T) {
	url := "example.com/image.png"
	hrefElement := fmt.Sprintf(`<img src="%v" />`, url)

	tutil.Assert(t, url, imageLink(hrefElement))
}
func TestCompareEqualNodes(t *testing.T) {
	hrefElemTemplate := `<a href="%v">Example.com</a>`
	node1 := readNode(fmt.Sprintf(hrefElemTemplate, "http://example.com/some/page"))
	node2 := readNode(fmt.Sprintf(hrefElemTemplate, "example.com/some/page"))
	node3 := readNode(fmt.Sprintf(hrefElemTemplate, "/some/page"))

	tutil.Assert(t, *node1, *node2)
	tutil.Assert(t, *node1, *node3)
	tutil.Assert(t, *node2, *node3)
}
