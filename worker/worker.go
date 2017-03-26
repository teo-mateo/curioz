package worker

import (
	"curioz/api"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

// CONCURRENCY is the number of threads
const CONCURRENCY = 6

var wg sync.WaitGroup

//Fetch fetches all images from the given API response
func Fetch(sol int, camera string, response api.Response, dest string) (err error) {

	fmt.Println("SOL: ", sol)

	//channel to use as a semaphore.
	//its capacity of 6 will block until a task finishes and removes from it
	//
	sem := make(chan bool, CONCURRENCY)

	if len(response.Photos) == 0 {
		fmt.Printf("Sol %d: no pictures.\n", sol)
		return nil
	}

	for _, photo := range response.Photos {
		sem <- true
		wg.Add(1)
		go func(imgSrc string) {
			defer func() { <-sem }()

			fullfn := getFullFN(sol, camera, imgSrc, dest)
			if fullfn == "" {
				return
			}

			bytes, err := fetchImage(sol, camera, imgSrc)
			if err != nil {
				return
			}

			err = writeFile(bytes, fullfn)
			if err != nil {
				os.Exit(0)
			}

		}(photo.ImgSrc)
	}

	return nil
}

func getFullFN(sol int, camera string, url string, dest string) (fullfn string) {

	//build path of output file;
	index := strings.LastIndex(url, "/") + 1
	fn := "sol_" + strconv.Itoa(sol) + "_" + camera + "_" + url[index:]

	fullfn = path.Join(dest, fn)

	//check if file already exists; if yes, return
	_, er := os.Stat(fullfn)
	if er == nil {
		fmt.Printf("File exists: %s\n", fn)
		return ""
	}

	return fullfn
}

func writeFile(bytes []byte, fullfn string) (err error) {
	//create new file
	f, er := os.Create(fullfn)
	if er != nil {
		return er
	}

	//write bytes to file
	_, err = f.Write(bytes)
	if err != nil {
		return
	}

	f.Close()

	return nil
}

func fetchImage(sol int, camera string, url string) (bytes []byte, err error) {

	//wait for me!
	defer wg.Done()

	fmt.Println(url)

	//get response
	webResponse, er := http.Get(url)
	if er != nil {
		return nil, er
	}

	//defer closing of stream
	defer webResponse.Body.Close()

	//read response into byte array
	bytes, er = ioutil.ReadAll(webResponse.Body)
	if er != nil {
		return nil, er
	}

	fmt.Printf("worker fetched %d bytes\n", len(bytes))

	return bytes, nil
}
