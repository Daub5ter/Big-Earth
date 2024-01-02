package parse

import (
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
	data.Events = make([]Event, 0, 24)
	return &data
}

func (p *Parsing) Parse(r db.Request) {
	data := NewData()

	for counts := 0; len(data.Events) == 0; counts++ {
		log.Println("try to parse date...")
		data = p.parseRussia()

		number, err := random.CreateNumber(1, 4)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(number * time.Second)

		if counts > 10 {
			counts = 0
			time.Sleep(15 * time.Second)
		}
	}

	log.Println(data)
}

func (p *Parsing) parseRussia() *Data {
	data := NewData()
	p.Collector.AllowURLRevisit = true

	p.Collector.OnHTML("div.rubric-featured__container", func(e *colly.HTMLElement) {

		if e.ChildText("div.rubric-featured__preview") != "" {

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
		log.Println("Error visit", url, err)
	}

	return data
}
