package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
)

format(b byte) {
	
}

func check() {
	res, err := http.Get("http://ws1.metcheck.com/ENGINE/v9_0/json.asp?lat=53.9&lon=-1.6&lid=67633&Fc=No")
	if err != nil {
		log.Fatal(err)
	}
	info, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", info)

	format(info)
}

func main() {
	check()
}
