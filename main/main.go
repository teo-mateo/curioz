package main

import (
	"curioz/api"
	"curioz/cameras"
	"curioz/worker"
	"fmt"
)

type camera struct {
	ID   string
	Name string
}

func main() {
	response, err := api.Get(1000, cameras.FHAZ)
	if err != nil {
		fmt.Println(err)
	}

	worker.Fetch(response)

}
