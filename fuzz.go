// +build gofuzz

package rencode

import "bytes"

func Fuzz(data []byte) int {
	d := NewDecoder(bytes.NewBuffer(data))
	for {
		_, err := d.DecodeNext()
		if err != nil {
			break
		}
	}
	return 1
}
