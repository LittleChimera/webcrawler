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
	bodyURL := hrefLink(elemBody)
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
	client       *http.Client
	visitedNodes map[Node]bool
}

func NewCrawler(client *http.Client) *Crawler {
	return &Crawler{
		client: client,
	}
}

func (c *Crawler) crawlPageLinks(url string) []Node {
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
		if t == html.StartTagToken && string(elemTag) == "a" {
			nodeSet[*readNode(string(decoder.Raw()))] = true
		}
	}

	var nodes []Node
	for node := range nodeSet {
		nodes = append(nodes, node)
	}

	return nodes
}
