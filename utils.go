package fyrirtaekjaskra

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ReadOrGetURL(filename string, url string) (content []byte, err error) {

	if _, err := os.Stat(filename); err == nil {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(fmt.Sprintf("unable to read a file: %s", filename))
		}
		return content, nil
	}

	log.Printf("Fetching: %s\n", url)

	res, err := http.Get(url)
	if err != nil {
		return
	}
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filename, content, 0644)

	return content, nil

}
