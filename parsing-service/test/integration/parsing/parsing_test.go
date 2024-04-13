package parsing

import (
	"testing"

	"github.com/gocolly/colly"
)

func TestParsing(t *testing.T) {
	collector := colly.NewCollector()

	url := "https://go.dev"
	t.Logf("отправка запроса на посещение сайта %s ...", url)
	err := collector.Visit(url)
	if err != nil {
		t.Errorf("ошибка, colly collector не работает: %v", err)
		return
	}
}
