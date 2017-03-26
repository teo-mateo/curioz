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

	sol := 1500
	for sol <= 1600 {
		response, err := api.Get(sol, cameras.NAVCAM)
		if err != nil {
			fmt.Println(err)
			break
		}
		worker.Fetch(sol, cameras.MAST, response)
		sol++
	}

}
