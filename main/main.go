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

	sol := 1
	for sol <= 500 {
		response, err := api.Get(sol, cameras.FHAZ)
		if err != nil {
			fmt.Println(err)
			break
		}
		worker.Fetch(sol, response)
		sol++
	}

}
