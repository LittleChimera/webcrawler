package crawler

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/lukadante/webcrawler/http"
	"golang.org/x/net/html"
)

func hrefLink(elemBody string) string {
	return valueFromAttribute(elemBody, "href")
}

func srcLink(elemBody string) string {
	return valueFromAttribute(elemBody, "src")
}

// Reads a value of an HTML elements' attribute
func valueFromAttribute(elemBody, attrName string) string {
	decoder := html.NewTokenizer(bytes.NewBufferString(elemBody))
	decoder.Next()
	for {
		key, val, more := decoder.TagAttr()
		if string(key) == attrName {
			return string(val)
		} else if !more {
			break
		}
	}
	return ""
}

type Node struct {
	Host string
	Path string
}

func readNode(elemBody string) *Node {
	decoder := html.NewTokenizer(bytes.NewBufferString(elemBody))
	decoder.Next()
	tagName, _ := decoder.TagName()

	var bodyURL string

	switch string(tagName) {
	case "img", "script":
		bodyURL = srcLink(elemBody)
	case "link", "a":
		bodyURL = hrefLink(elemBody)
	}

	if strings.HasPrefix(bodyURL, "/") {
		bodyURL = CrawlHostname + bodyURL
	}

	if !strings.HasPrefix(bodyURL, "http") {
		bodyURL = "http://" + bodyURL
	}

	u, _ := url.Parse(bodyURL)

	return &Node{
		Host: u.Host,
		Path: u.Path,
	}
}

type Crawler struct {
	client *http.Client
}

func NewCrawler(client *http.Client) *Crawler {
	return &Crawler{
		client: client,
	}
}

func (c *Crawler) crawlPageLinks(url string) []Node {
	return c.crawlPage(url, "a")
}

func (c *Crawler) crawlPageResources(url string) []Node {
	return c.crawlPage(url, "script", "link", "img")
}

func (c *Crawler) crawlPage(url string, tags ...string) []Node {
	pageBody := (*c.client).Get(url)
	decoder := html.NewTokenizer(bytes.NewBufferString(pageBody))
	nodeSet := make(map[Node]bool)
	decoder.Next()
	for {
		t := decoder.Next()
		if t == html.ErrorToken {
			break
		}
		elemTag, _ := decoder.TagName()
		if t == html.StartTagToken || t == html.SelfClosingTagToken {
			var matches bool
			for _, tag := range tags {
				matches = matches || (tag == string(elemTag))
			}
			if !matches {
				continue
			}
			nodeSet[*readNode(string(decoder.Raw()))] = true
		}
	}

	var nodes []Node
	for node := range nodeSet {
		nodes = append(nodes, node)
	}

	return nodes
}
