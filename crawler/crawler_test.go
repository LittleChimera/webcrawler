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
func TestCompareEqualLinkNodes(t *testing.T) {
	node1 := readNode(fmt.Sprintf(hrefElemTemplate, "http://lipsum.com/some/page"))
	node2 := readNode(fmt.Sprintf(hrefElemTemplate, "lipsum.com/some/page"))
	node3 := readNode(fmt.Sprintf(hrefElemTemplate, "/some/page"))

	tutil.Assert(t, *node1, *node2)
	tutil.Assert(t, *node1, *node3)
	tutil.Assert(t, *node2, *node3)
}

func TestCompareEqualResourceNodes(t *testing.T) {
	node1 := readNode(fmt.Sprintf(`<img src="%v" />`, "http://lipsum.com/some/image.png"))
	node2 := readNode(fmt.Sprintf(`<link href="%v" />`, "lipsum.com/some/image.png"))
	node3 := readNode(fmt.Sprintf(`<script src="%v" />`, "/some/image.png"))

	tutil.Assert(t, *node1, *node2)
	tutil.Assert(t, *node1, *node3)
	tutil.Assert(t, *node2, *node3)
}

func mockClient(t *testing.T, response string) (http.Client, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockClient := mock_http.NewMockClient(ctrl)
	mockClient.EXPECT().Get("lipsum.com").Return(response)
	return mockClient, ctrl
}

func TestSinglePageCrawlLinksResults(t *testing.T) {
	client, ctrl := mockClient(t, htmlMockSource)
	defer ctrl.Finish()

	const host = "lipsum.com"
	crawler := NewCrawler(client)
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

func TestSinglePageCrawlResourcesResults(t *testing.T) {
	client, ctrl := mockClient(t, htmlMockSource)
	defer ctrl.Finish()

	const host = "lipsum.com"
	crawler := NewCrawler(client)
	crawledNodes := crawler.crawlPageResources(host)

	var internalResourcesCount int
	for _, node := range crawledNodes {
		if node.Host == host {
			internalResourcesCount++
		}
	}
	tutil.Assert(t, 9, internalResourcesCount)
	tutil.Assert(t, 1, len(crawledNodes)-internalResourcesCount)
}

func TestCrawlSiteLoop(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_http.NewMockClient(ctrl)
	mockClient.EXPECT().Get("lipsum.com").Return(htmlLoopMockSources)
	mockClient.EXPECT().Get("lipsum.com").Return(htmlLoopMockSources)
	mockClient.EXPECT().Get("lipsum.com/about").Return(htmlLoopMockSources)
	mockClient.EXPECT().Get("lipsum.com/about").Return(htmlLoopMockSources)
	mockClient.EXPECT().Get("lipsum.com/generate").Return(htmlLoopMockSources)
	mockClient.EXPECT().Get("lipsum.com/generate").Return(htmlLoopMockSources)
	defer ctrl.Finish()

	const host = "lipsum.com"

	var client http.Client
	client = mockClient
	crawler := NewCrawler(client)
	crawledPages := crawler.CrawlSite(host)

	tutil.Assert(t, 3, len(crawledPages))
}
