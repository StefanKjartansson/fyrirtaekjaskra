package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	BufferSize = 256
	cc         = make(chan Company, BufferSize)
)

func main() {

	streets, err := ImportStreets("./gotuskra.csv")
	if err != nil {
		fmt.Println(err)
	}

	total := 0
	hasWritten := false

	go func() {
		for _, s := range streets {
			if len(strings.Split(s, " ")) == 1 {
				ScrapeStreet(s, cc)
			}
		}
	}()

	go func() {
		for {
			select {
			case ev := <-cc:
				total++
				if hasWritten {
					os.Stdout.Write([]byte(","))
				}
				b, _ := json.MarshalIndent(ev, "", "  ")
				os.Stdout.Write(b)
				hasWritten = true
			case <-time.After(2 * time.Second):
				fmt.Printf(",{\"total\": %d}\n", total)
				os.Exit(1)
			}
		}
	}()

	select {}
}
