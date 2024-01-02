package parse

import "github.com/gocolly/colly"

type Parsing struct {
	Collector *colly.Collector
}

type Data struct {
	Events []Event
}

type Event struct {
	Name  string
	Image string
	Link  string
}
