package crawler

import (
	"bytes"
	"fmt"
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

func (n *Node) URL() string {
	return n.Host + n.Path
}

func newNode(link string) *Node {
	if strings.HasPrefix(link, "/") {
		link = CrawlHostname + link
	}

	if !strings.HasPrefix(link, "http") {
		link = "http://" + link
	}

	u, _ := url.Parse(link)

	return &Node{
		Host: u.Host,
		Path: u.Path,
	}
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

	return newNode(bodyURL)
}

type Crawler struct {
	client       *http.Client
	visitedPages map[Node]bool
}

func NewCrawler(client *http.Client) *Crawler {
	return &Crawler{
		client:       client,
		visitedPages: make(map[Node]bool),
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

type Page struct {
	Path          Node
	LinkNodes     []Node
	ResourceNodes []Node
}

func (c *Crawler) CrawlSite(url string) []Page {
	return c.CrawlNode(newNode(url))
}

func (c *Crawler) CrawlNode(pageNode *Node) []Page {
	fmt.Println(pageNode.URL())
	var pages []Page
	if c.visitedPages[*pageNode] {
		return pages
	}
	c.visitedPages[*pageNode] = true
	pageResources := c.crawlPageResources(pageNode.URL())
	pageLinks := c.crawlPageLinks(pageNode.URL())
	for _, linkNode := range pageLinks {
		pages = append(pages, c.CrawlNode(&linkNode)...)
	}

	pages = append(pages, Page{
		Path:          *pageNode,
		LinkNodes:     pageLinks,
		ResourceNodes: pageResources,
	})

	return pages
}
