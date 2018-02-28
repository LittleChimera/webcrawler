package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lukadante/webcrawler/crawler"
	"github.com/lukadante/webcrawler/http"
)

func main() {
	crawler.CrawlHostname = os.Args[1]
	client := http.SimpleClient{}
	crawlerClient := crawler.NewCrawler(client)
	result := crawlerClient.CrawlSite(crawler.CrawlHostname)
	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
}
