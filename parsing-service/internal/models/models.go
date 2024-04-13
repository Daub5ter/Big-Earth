package models

// Place - место на планете.
type Place struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

// PlaceInformation - информация о месте.
type PlaceInformation struct {
	Text   string   `json:"text"`
	Photos []string `json:"photos"`
	Videos []string `json:"videos"`
}

// Event - событие в месте.
type Event struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Link  string `json:"link"`
}
