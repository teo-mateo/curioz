package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//CURIOSITY Mars rover API root
const CURIOSITY = "https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos"

// Get JSON response for the given sol (Mars day) and camera
func Get(sol int, camera string) (r Response, err error) {

	var j string

	//build url
	key := os.Getenv("NASADEV")
	url := CURIOSITY + "?sol=" + strconv.Itoa(sol) + "&camera=" + strings.ToLower(camera) + "&api_key=" + key
	fmt.Println(url)

	//make http request
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	//read body, convert to string
	defer resp.Body.Close()
	bytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return
	}
	j = string(bytes)

	//decode json
	dec := json.NewDecoder(strings.NewReader(j))
	dec.Decode(&r)
	fmt.Printf("%d bytes decoded", len(bytes))

	//return response
	return
}
