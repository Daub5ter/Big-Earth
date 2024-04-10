package colly

import (
	"testing"

	"github.com/gocolly/colly"
)

func TestCollector(t *testing.T) {
	collector := colly.NewCollector()

	err := collector.Visit("https://go.dev")
	if err != nil {
		t.Errorf("ошибка, colly collector не работает")
		return
	}
}
