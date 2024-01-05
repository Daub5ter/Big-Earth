package data

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"parsing-service/pkg/random"
	"strings"
	"time"
)

const url = "https://afisha.yandex.ru/krasnodar"

func NewParsing() *Parsing {
	return &Parsing{Collector: colly.NewCollector()}
}

func newPlaceInformation() *PlaceInformation {
	var data PlaceInformation
	data.Events = make([]Event, 0, 32)
	return &data
}

func (p *Parsing) Parse(place Place) *PlaceInformation {
	var data *PlaceInformation
	p.Collector.AllowURLRevisit = true

	for counts := 0; data == nil; counts++ {
		log.Println("try to parse date...")

		err := p.parseRussia(data)
		if err != nil {
			log.Println(err)
		}

		dur, err := random.CreateDuration(1, 4)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(dur * time.Second)

		if counts > 10 {
			counts = 0
			time.Sleep(20 * time.Second)
		}
	}

	log.Println(data)

	return data
}

func (p *Parsing) parseRussia(data *PlaceInformation) error {
	p.Collector.OnHTML("div.rubric-featured__container", func(e *colly.HTMLElement) {
		if e.ChildText("div.rubric-featured__preview") != "" {
			if data == nil {
				data = newPlaceInformation()
			}

			var name strings.Builder
			if e.ChildText("div.rubric-featured__top") != "" {
				name.WriteString(e.ChildText("div.rubric-featured__top"))
				name.WriteString(" ")
			}
			name.WriteString(e.ChildText("div.rubric-featured__title"))
			name.WriteString(" ")
			name.WriteString(e.ChildText("div.rubric-featured__preview"))

			nameCorrect := strings.ReplaceAll(name.String(), "    ", " ")

			event := Event{
				Name:  nameCorrect,
				Image: e.ChildAttrs("img", "src")[0],
				Link:  "https://afisha.yandex.ru" + e.ChildAttrs("a", "href")[0],
			}

			data.Events = append(data.Events, event)
		}
	})

	err := p.Collector.Visit(url)
	if err != nil {
		return errors.New(fmt.Sprintf("Error visit %s, %s", url, err.Error()))
	}

	if data == nil {
		return errors.New("error captcha/protect")
	}

	return nil
}
