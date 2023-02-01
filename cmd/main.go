package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/walkergriggs/tsgo"
)

func main() {
	b := make([]byte, 188)

	for {
		n, err := os.Stdin.Read(b)
		if n == 0 && err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		h, err := tsgo.ParseHeader(b)
		if err != nil {
			log.Fatal(err)
		}

		s, err := json.MarshalIndent(h, "", "\t")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(s))
	}
}
