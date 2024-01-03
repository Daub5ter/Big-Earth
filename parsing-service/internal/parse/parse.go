package parse

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"parsing-service/internal/db"
	"parsing-service/pkg/random"
	"time"
)

const url = "https://afisha.yandex.ru/krasnodar"

func NewParsing() *Parsing {
	return &Parsing{Collector: colly.NewCollector()}
}

func NewData() *Data {
	var data Data
	data.Events = make([]Event, 0, 32)
	return &data
}

func (p *Parsing) Parse(r db.Request) {
	data := NewData()

	for counts := 0; len(data.Events) == 0; counts++ {
		log.Println("try to parse date...")

		d, err := p.parseRussia()
		if err != nil {
			log.Println(err)
		} else {
			data.Events = d.Events
		}

		number, err := random.CreateNumber(1, 4)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(number * time.Second)

		if counts > 10 {
			counts = 0
			time.Sleep(12 * time.Second)
		}
	}

	log.Println(data)
}

func (p *Parsing) parseRussia() (*Data, error) {
	var data *Data
	p.Collector.AllowURLRevisit = true

	p.Collector.OnHTML("div.rubric-featured__container", func(e *colly.HTMLElement) {

		if e.ChildText("div.rubric-featured__preview") != "" {
			if data == nil {
				data = NewData()
			}

			var name string
			if e.ChildText("div.rubric-featured__top") != "" {
				name = e.ChildText("div.rubric-featured__top") + " "
			}
			name += e.ChildText("div.rubric-featured__title") + " " +
				e.ChildText("div.rubric-featured__preview")

			event := Event{
				Name:  name,
				Image: e.ChildAttrs("img", "src")[0],
				Link:  "https://afisha.yandex.ru" + e.ChildAttrs("a", "href")[0],
			}

			data.Events = append(data.Events, event)
		}
	})

	err := p.Collector.Visit(url)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error visit", url, err))
	}

	if data == nil {
		return nil, errors.New("error captcha/protect")
	}

	return data, nil
}
