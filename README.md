# Webcrawler

This is a simple web crawler which crawls a single domain for URL resources both internal and external.

As this project was made with the minimal amount of time and effort to achieve most basic functionality it lacks a lot of aspects of a usable tool. Some of those include:
  - Error handling
  - Performance optimization
      - Fetch page only once (current implementation fetches page twice  - once for crawling resources to other sites and once for static resources)
      - Crawl multiple pages at the same time

No attempt has been made to mitigate these issues as these are beyond the scope of idea.
However, this program should be sufficient if you want construct a simple sitemap of the domain.


## Installing

Installing using go get pulls from the master branch and builds runnable binary.

```bash
go get -u github.com/lukadante/webcrawler
```

## Running

To crawl a domain, run the following command:

```bash
webcrawler <url to crawl>
```

E.g.

```bash
webcrawler example.com
```

## Testing

Currently, only `crawler` package has a testing suite.

To run tests, position yourself to package directory and run the following:

```bash
go test .
```

## Future work

- [x] Create simple one domain web crawler
- [ ] Persistent storage of results
- [ ] Crawling parallelization
- [ ] Distributed crawling


## Requirements

  - Go 1.9.4

