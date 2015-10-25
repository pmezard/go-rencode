package main

import (
	"fmt"
	"os"

	"github.com/gdm85/go-rencode"
)

func main() {
	d := rencode.Dictionary{}
	d.Add("int", int64(42))
	d.Add("float", float64(3.14))
	d.Add("string", "some string")
	d.Add("bool", true)

	nested := rencode.Dictionary{}
	list := rencode.List{}
	list.Add(2)
	list.Add("another string")
	list.Add(false)
	d.Add("nested", nested)
	empty := rencode.List{}
	d.Add("empty", empty)
	enc := rencode.Encoder{}
	err := enc.Encode(d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(enc.Bytes())
}
