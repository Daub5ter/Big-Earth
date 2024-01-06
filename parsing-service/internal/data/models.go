package data

import "github.com/gocolly/colly"

type Place struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

type Parsing struct {
	Collector *colly.Collector
}

type PlaceInformation struct {
	Events []Event `json:"events"`
}

type Event struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Link  string `json:"link"`
}
