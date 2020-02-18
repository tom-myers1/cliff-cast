package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

// Final is used to calc the final values to display
type Final struct {
	Day        string
	RainTotal  float64
	RainChance []int64
	Temp       []int64
	Wind       []int64
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

	var tDay string
	var d1, d2, d3 Final
	loc, _ := time.LoadLocation("UTC")
	layout := "2006-01-02T15:04:05"

	now := time.Now().In(loc).Truncate(24 * time.Hour)

	for _, t := range fCast {
		// attempt to parse utc to time
		tt, err := time.Parse(layout, t.Utc)
		check(err)
		diff := now.Sub(tt.Truncate(24 * time.Hour))

		xTemp, err := strconv.ParseInt(t.Temp, 10, 10)
		check(err)
		xRChance, err := strconv.ParseInt(t.Chance, 10, 10)
		check(err)
		xRTotal, err := strconv.ParseFloat(t.Rain, 10)
		check(err)
		xWind, err := strconv.ParseInt(t.Wind, 10, 10)
		check(err)

		switch {
		case diff.Hours() == 0: // today

			tDay = "today"
			d1.Day = tDay
			d1.RainChance = append(d1.RainChance, xRChance)
			d1.RainTotal = d1.RainTotal + xRTotal
			d1.Wind = append(d1.Wind, xWind)
			d1.Temp = append(d1.Temp, xTemp)
		case diff.Hours() == -24: // tomorrow

			tDay = "tomorrow"
			d2.Day = tDay
			d2.RainChance = append(d2.RainChance, xRChance)
			d2.RainTotal = d2.RainTotal + xRTotal
			d2.Wind = append(d2.Wind, xWind)
			d2.Temp = append(d2.Temp, xTemp)
		case diff.Hours() == -48: //day after

			tDay = t.DayN
			d3.Day = tDay
			d3.RainChance = append(d3.RainChance, xRChance)
			d3.RainTotal = d3.RainTotal + xRTotal
			d3.Wind = append(d3.Wind, xWind)
			d3.Temp = append(d3.Temp, xTemp)
		default:
			// do nothing
		}

	}
	fmt.Println("today...")
	fmt.Println(d1)
	fmt.Println("tomorrow")
	fmt.Println(d2)
	fmt.Println("day after")
	fmt.Println(d3)
}

func main() {
	url := "http://ws1.metcheck.com/ENGINE/v9_0/json.asp?lat=53.9&lon=-1.6&lid=67633&Fc=No"
	info := getInfo(url)
	unmarshal(info) // unmarshall to struct


}
