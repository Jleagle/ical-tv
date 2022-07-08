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

type event struct {
	ID     string `json:"id"`
	Name   string `json:"title"`
	Date   string `json:"start"`
	AllDay bool   `json:"allDay"`
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
	for k, v := range cal.Components {

		var date time.Time
		var summary string

		for _, vv := range v.UnknownPropertiesIANAProperties() {

			if vv.IANAToken == propertySummary {
				summary = vv.Value
			}

			if vv.IANAToken == propertyStart {
				date, err = time.Parse("20060102T150405Z", vv.Value)
			}
		}

		events = append(events, event{
			ID:     fmt.Sprint(k),
			Name:   summary,
			Date:   date.Format("2006-01-02"),
			AllDay: true,
		})
	}

	// Template
	app := fiber.New(fiber.Config{
		Views:                 html.New("./", ".gohtml"),
		DisableStartupMessage: true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":  config.Title,
			"Events": events,
		})
	})

	err = app.Listen(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatalln(err)
	}
}
