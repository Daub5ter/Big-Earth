package parsing

import (
	"testing"

	"github.com/gocolly/colly"
	"github.com/stretchr/testify/suite"
)

// ClientSuite - структура для тестов.
type ClientSuite struct {
	suite.Suite
	collector *colly.Collector
}

// SetupSuite настраивает тесты
// (включается перед тестами) .
func (c *ClientSuite) SetupSuite() {
	c.collector = colly.NewCollector()
}

// TestParsing запускает тесты.
func TestParsing(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

// TestVisit проверяет работу посещение сайта/сервиса парсером.
func (c *ClientSuite) TestVisit() {
	url := "https://go.dev"
	err := c.collector.Visit(url)
	c.NoError(err)
	c.T().Log("библиотека парсинга может посещать сайты/сервисы")
}
