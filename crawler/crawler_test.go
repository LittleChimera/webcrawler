package crawler

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	tutil "github.com/lukadante/webcrawler/test_util"

	"github.com/lukadante/webcrawler/http"
	"github.com/lukadante/webcrawler/http/mock_http"
)

const (
	hrefElemTemplate = `<a href="%v">lipsum.com</a>`
)

var httpClient http.Client

func init() {
	CrawlHostname = "lipsum.com"
}

func TestReadHrefLink(t *testing.T) {
	url := "lipsum.com"
	hrefElement := fmt.Sprintf(hrefElemTemplate, url)

	tutil.Assert(t, url, hrefLink(hrefElement))
}
func TestReadSrcLink(t *testing.T) {
	url := "lipsum.com/image.png"
	hrefElement := fmt.Sprintf(`<img src="%v" />`, url)

	tutil.Assert(t, url, srcLink(hrefElement))
}
func TestCompareEqualNodes(t *testing.T) {
	node1 := readNode(fmt.Sprintf(hrefElemTemplate, "http://lipsum.com/some/page"))
	node2 := readNode(fmt.Sprintf(hrefElemTemplate, "lipsum.com/some/page"))
	node3 := readNode(fmt.Sprintf(hrefElemTemplate, "/some/page"))

	tutil.Assert(t, *node1, *node2)
	tutil.Assert(t, *node1, *node3)
	tutil.Assert(t, *node2, *node3)
}

func mockClient(t *testing.T) (http.Client, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockClient := mock_http.NewMockClient(ctrl)
	mockClient.EXPECT().Get("lipsum.com").Return(htmlMockSource)
	return mockClient, ctrl
}

func TestSinglePageCrawlLinksResults(t *testing.T) {
	client, ctrl := mockClient(t)
	defer ctrl.Finish()

	const host = "lipsum.com"
	crawler := NewCrawler(&client)
	crawledNodes := crawler.crawlPageLinks(host)

	var internalLinksCount int
	for _, node := range crawledNodes {
		if node.Host == host {
			internalLinksCount++
		}
	}
	tutil.Assert(t, 4, internalLinksCount)
	tutil.Assert(t, 51, len(crawledNodes)-internalLinksCount)
}
