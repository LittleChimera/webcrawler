package crawler

import (
	"bytes"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func hrefLink(elemBody string) string {
	return valueFromAttribute(elemBody, "href")
}

func imageLink(elemBody string) string {
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
		bodyURL = crawlHostname + bodyURL
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
