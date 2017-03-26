package worker

import (
	"curioz/api"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

//Fetch fetches all images from the given API response
func Fetch(sol int, camera string, response api.Response) (err error) {

	fmt.Println("SOL: ", sol)

	//channel to use as a semaphore.
	//its capacity of 6 will block until a task finishes and removes from it
	//
	sem := make(chan bool, 6)

	if len(response.Photos) == 0 {
		fmt.Printf("Sol %d: no pictures.\n", sol)
		return nil
	}

	for _, photo := range response.Photos {
		sem <- true
		wg.Add(1)
		go func(s int, c string, i string) {
			defer func() { <-sem }()
			fetchImage(s, c, i)
		}(sol, camera, photo.ImgSrc)
	}

	return nil
}

func fetchImage(sol int, camera string, url string) (err error) {

	//wait for me!
	defer wg.Done()

	fmt.Println(url)

	//build path of output file;
	index := strings.LastIndex(url, "/") + 1
	fn := "sol_" + strconv.Itoa(sol) + "_" + camera + "_" + url[index:]
	cwd, er := os.Getwd()
	if er != nil {
		return er
	}
	//full file name
	fullfn := cwd + "\\..\\img_result\\" + fn

	//check if file already exists; if yes, return
	_, er = os.Stat(fullfn)
	if er == nil {
		fmt.Printf("File exists: %s\n", fn)
		return nil
	}

	webResponse, er := http.Get(url)
	if er != nil {
		return er
	}

	//defer closing of stream
	defer webResponse.Body.Close()

	//read response into byte array
	bytes, er := ioutil.ReadAll(webResponse.Body)
	if er != nil {
		return er
	}

	fmt.Printf("worker fetched %d bytes\n", len(bytes))

	//create new file
	f, er := os.Create(fullfn)
	if er != nil {
		return er
	}

	//write bytes to file
	_, er = f.Write(bytes)
	if er != nil {
		return er
	}

	f.Close()

	return nil
}
