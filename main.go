package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// used for comparing data to todays date to get 3 days
const d float64 = -72

// used for weekend
const w float64 = -168
const sun int = 1
const mon int = 2
const tue int = 3
const wed int = 4
const thu int = 5
const fri int = 6
const sat int = 7

// R1 is top level struct
type R1 struct {
	D1 D1 `json:"metcheckData"`
}

// D1 is level 2 struct
type D1 struct {
	Location F1 `json:"forecastLocation"`
}

// F1 is level 3
type F1 struct {
	Forecast []Res `json:"forecast"`
}

// Res represents the things we actually want from the json response
type Res struct {
	Temp   string `json:"temperature"`
	Chance string `json:"chanceofrain"`
	Rain   string `json:"rain"`
	Wind   string `json:"windgustspeed"`
	Humid  string `json:"humidity"`
	Day    string `json:"dayOfWeek"` // 1 (sunday) to 7 (saturday)
	Utc    string `json:"utcTime"`
	DayN   string `json:"weekday"`
}

// check is a basic error checker
func check(e error) {
	if e != nil {
		fmt.Println("an error occured... panic!!!")
		panic(e)
	}
}

// getInfo calls endpoint and returns []byte
func getInfo(url string) []byte {

	res, err := http.Get("http://ws1.metcheck.com/ENGINE/v9_0/json.asp?lat=53.9&lon=-1.6&lid=67633&Fc=No")
	check(err)
	info, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	check(err)
	//fmt.Println(info)
	return info

}

// unmarshal sorts through json to filter what is needed
func unmarshal(info []byte) {
	// to struct
	var data R1
	err := json.Unmarshal(info, &data)
	check(err)
	fCast := data.D1.Location.Forecast
	forecast(fCast)
}

// forecast checks through forecast for next 7 days
func forecast(fCast []Res) {
	loc, _ := time.LoadLocation("UTC")
	layout := "2006-01-02T15:04:05"

	now := time.Now().In(loc).Truncate(24 * time.Hour)

	for _, t := range fCast {
		// attempt to parse utc to time
		tt, err := time.Parse(layout, t.Utc)
		check(err)
		diff := now.Sub(tt.Truncate(24 * time.Hour))

		switch {
		case diff.Hours() == 0: // today
			fmt.Println("today")
		case diff.Hours() == -24: // tomorrow
			fmt.Println("tomorrow")
		case diff.Hours() == -48: //day after
			fmt.Println("the day after")
		case diff.Hours() > w && t.Day == "7": // saturday
			fmt.Println("saturday")
		case diff.Hours() > w && t.Day == "1": // sunday
			fmt.Println("sunday")

		}

		/* filter down to next d days
		if diff.Hours() > d {
			fmt.Println("temp = ", t.Temp, "degrees celcius")
			fmt.Println("chance of rain = ", t.Chance, "%")
			fmt.Println("humidity = ", t.Humid, "%")
			fmt.Println("ammount of rain = ", t.Rain, "mm per hour")
			fmt.Println("day = ", t.Day)
			fmt.Println("utc = ", t.Utc)
			fmt.Println("day = ", t.DayN)
		} */
	}

}

func main() {
	url := "http://ws1.metcheck.com/ENGINE/v9_0/json.asp?lat=53.9&lon=-1.6&lid=67633&Fc=No"
	info := getInfo(url)
	unmarshal(info) // unmarshall to struct

	// getting date for main title of page
	day := time.Now().Weekday()
	fmt.Println(day)
}
