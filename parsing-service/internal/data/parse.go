package data

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

const url = "https://krd.kassir.ru"

func NewParsing() *Parsing {
	return &Parsing{Collector: colly.NewCollector()}
}

func newPlaceInformation(allocationEvent int) *PlaceInformation {
	var placeInformation PlaceInformation
	placeInformation.Events = make([]Event, 0, allocationEvent)
	return &placeInformation
}

func (p *Parsing) Parse(place Place) *PlaceInformation {
	placeInformation := newPlaceInformation(24)
	p.Collector.AllowURLRevisit = true

	err := p.parseRussia(placeInformation)
	if err != nil {
		log.Println(err)
	}

	return placeInformation
}

func (p *Parsing) parseRussia(placeInformation *PlaceInformation) error {
	maxEvents := 10
	eventCounter := 1

	p.Collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
	})

	p.Collector.OnHTML("div.whitespace-nowrap.mr-3", func(e *colly.HTMLElement) {
		var err error
		parts := strings.Split(e.Text, "/")
		maxEvents, err = strconv.Atoi(parts[1])
		if err != nil {
			log.Println(err)
		}
	})

	p.Collector.OnHTML("div.swiper-slide", func(e *colly.HTMLElement) {
		if eventCounter == maxEvents {
			return
		}

		name := e.ChildText("h2.line-clamp-2")
		if name == "" {
			return
		}

		image := e.ChildAttrs("source", "srcset")[0]
		if image == "" {
			return
		}

		var link string
		if len(e.ChildAttrs("div.cursor-pointer", "href")) != 0 {
			link = e.ChildAttrs("div.cursor-pointer", "href")[0]
		} else {
			link = e.ChildAttrs("a.cursor-pointer", "href")[0]
		}

		event := Event{
			Name:  name,
			Image: image,
			Link:  link,
		}

		placeInformation.Events = append(placeInformation.Events, event)

		eventCounter++
	})

	err := p.Collector.Visit(url)
	if err != nil {
		log.Println(errors.New(fmt.Sprintf("Error visit %s, %s", url, err.Error())))
	}

	return nil
}
