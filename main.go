package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arran4/golang-ical"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

const (
	propertyCategories  = "CATEGORIES"
	propertyDescription = "DESCRIPTION"
	propertySummary     = "SUMMARY"
	propertyID          = "UID"
	propertyStart       = "DTSTART"
	propertyEnd         = "DTEND"
)

type config struct {
	Ical  string `json:"ical"`
	Title string `json:"title"`
	Port  int    `json:"port"`
}

func main() {

	// Get config
	b, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatalln(err)
	}

	if config.Ical == "" {
		log.Fatalln("ICAL URL not found in config")
	}

	// Get ical
	resp, err := http.Get(config.Ical)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	cal, err := ics.ParseCalendar(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Format data
	var events []event
	for _, v := range cal.Components {

		var date time.Time
		var summary string

		for _, v := range v.UnknownPropertiesIANAProperties() {

			if v.IANAToken == propertySummary {
				summary = v.Value
				break
			}

			if v.IANAToken == propertyStart {
				date, err = time.Parse("20060102T150405Z", v.Value)
				break
			}

			fmt.Println(v.IANAToken, v.Value)
		}

		events = append(events, event{Name: summary, Date: date.String()})
	}

	fmt.Println(events)

	// Template
	app := fiber.New(fiber.Config{Views: html.New("./", ".gohtml")})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	err = app.Listen(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatalln(err)
	}
}

type event struct {
	Name string
	Date string
}
