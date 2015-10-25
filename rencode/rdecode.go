package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/gdm85/go-rencode"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	d := rencode.NewDecoder(bytes.NewBuffer(data))
	for {
		_, err := d.DecodeNext()
		if err != nil {
			break
		}
	}
}
