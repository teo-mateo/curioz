package main

import (
	"curioz/api"
	"curioz/cameras"
	"fmt"
)

type camera struct {
	ID   string
	Name string
}

func main() {
	_, err := api.Get(1000, cameras.FHAZ)
	if err != nil {
		fmt.Println(err)
	}

}
