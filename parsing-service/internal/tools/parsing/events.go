package parsing

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"parsing-service/internal/models"

	"github.com/gocolly/colly"
)

// parseEventRussia парсит данные о событиях в стране Россия.
func (p parser) parseEventRussia(urlPlace string, events *[]*models.Event) error {
	var err error
	maxEvents := 10
	eventCounter := 1

	p.collector.OnHTML("div.whitespace-nowrap.mr-3", func(e *colly.HTMLElement) {
		parts := strings.Split(e.Text, "/")
		if len(parts) < 2 {
			err = errors.New("no necessary data")
			return
		}

		maxEvents, err = strconv.Atoi(parts[1])
	})

	p.collector.OnHTML("div.swiper-slide", func(e *colly.HTMLElement) {
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

		*events = append(*events,
			&models.Event{
				Name:  name,
				Image: image,
				Link:  link,
			})

		eventCounter++
	})

	errVisit := p.collector.Visit(urlPlace)
	if errVisit != nil {
		return errors.New(fmt.Sprintf("Error visit %s, %s", urlPlace, errVisit.Error()))
	}

	if err != nil {
		return errors.New(fmt.Sprintf("Error sraping %s, %s", urlPlace, err.Error()))
	}

	return nil
}
