package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Res represents the things we want from the json response
type Res struct {
	Temp   string `json:"temperature"`
	Chance string `json:"chanceofrain"`
	Rain   string `json:"rain"`
	Wind   string `json:"windgustspeed"`
	Humid  string `json:"humidity"`
	Day    string `json:"dayOfWeek"`
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

	return info

}

func test1(info []byte) {
	m := make(map[string]interface{})
	err := json.Unmarshal(info, &m)
	check(err)
	//fmt.Print(m)
	test := m["metcheckData"]
	fmt.Println("test1")
	fmt.Println("          ")
	fmt.Println(test)

	test2(test)
}

func test2(info interface{}) {

	fmt.Println("test2")
	fmt.Println("          ")
	gob.Register(map[string]interface{}{})
	var buf bytes.Buffer // stand in for network
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(info)
	check(err)

	fmt.Print(buf.Bytes())

}

func main() {
	url := "http://ws1.metcheck.com/ENGINE/v9_0/json.asp?lat=53.9&lon=-1.6&lid=67633&Fc=No"
	info := getInfo(url)
	test1(info) // json unmarshal
	//test2(info) // gob
	//test3(info) // something else
}
