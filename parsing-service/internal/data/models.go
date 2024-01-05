package data

import "github.com/gocolly/colly"

type Request struct {
	Country string
	City    string
}

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
