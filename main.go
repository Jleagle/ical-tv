package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arran4/golang-ical"
)

const url = ""

func main() {

	resp, err := http.Get(url)
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
