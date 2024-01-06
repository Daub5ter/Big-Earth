package data

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
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

func (p *Parsing) Parse(place Place) (*PlaceInformation, error) {
	placeInformation := newPlaceInformation(24)
	p.Collector.AllowURLRevisit = true

	err := p.parseRussia(placeInformation)
	if err != nil {
		return nil, err
	}

	return placeInformation, nil
}

func (p *Parsing) parseRussia(placeInformation *PlaceInformation) error {
	var err error
	maxEvents := 10
	eventCounter := 1

	p.Collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
	})

	p.Collector.OnHTML("div.whitespace-nowrap.mr-3", func(e *colly.HTMLElement) {
		parts := strings.Split(e.Text, "/")
		if len(parts) < 2 {
			err = errors.New("no necessary data")
			return
		}

		maxEvents, err = strconv.Atoi(parts[1])
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

	errVisit := p.Collector.Visit(url)
	if errVisit != nil {
		return errors.New(fmt.Sprintf("Error visit %s, %s", url, errVisit.Error()))
	}

	if err != nil {
		return errors.New(fmt.Sprintf("Error sraping %s, %s", url, err.Error()))
	}

	return nil
}
