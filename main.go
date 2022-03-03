package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arran4/golang-ical"
)

type config struct {
	Ical  string `json:"ical"`
	Title string `json:"title"`
}

func main() {

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

	resp, err := http.Get(config.Ical)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	cal, err := ics.ParseCalendar(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(cal)
}
