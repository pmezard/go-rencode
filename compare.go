package rencode

//
// go-rencode v0.1.0 - Go implementation of rencode - fast (basic)
//                  object serialization similar to bencode
// Copyright (C) 2015 gdm85 - https://github.com/gdm85/go-rencode/

// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

import (
	"bytes"
)

// this hack allows fetching keys by either string or byte slice type
// no other special matching is performed with other slice types
func deepEqual(a, b interface{}) bool {
	switch a.(type) {
	case string:
		switch b.(type) {
		case []byte:
			return bytes.Compare([]byte(a.(string)), b.([]byte)) == 0
		case string:
			return a == b
		default:
			return false
		}
	case []byte:
		switch b.(type) {
		case []byte:
			return bytes.Compare(a.([]byte), b.([]byte)) == 0
		case string:
			return bytes.Compare([]byte(b.(string)), a.([]byte)) == 0
		default:
			return false
		}
	case Dictionary:
		switch b.(type) {
		case Dictionary:
			d1 := a.(Dictionary)
			d2 := b.(Dictionary)
			return d1.Equals(&d2)
		}
		return false
	case List:
		switch b.(type) {
		case List:
			l1 := a.(List)
			l2 := b.(List)
			return l1.Equals(&l2)
		}
		return false
	}

	return a == b
}
