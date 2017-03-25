package worker

import (
	"curioz/api"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

//Fetch fetches all images from the given API response
func Fetch(response api.Response) (err error) {

	for _, photo := range response.Photos {
		wg.Add(1)
		go fetchImage(photo.ImgSrc)
	}

	wg.Wait()
	return nil // no error
}

func fetchImage(url string) (err error) {

	//wait for me!
	defer wg.Done()

	fmt.Println(url)

	webResponse, er := http.Get(url)
	if er != nil {
		return er
	}

	//defer closing of stream
	defer webResponse.Body.Close()

	bytes, er := ioutil.ReadAll(webResponse.Body)
	if er != nil {
		return er
	}

	fmt.Printf("worker fetched %d bytes", len(bytes))
	return nil
}
