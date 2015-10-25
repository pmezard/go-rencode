/*
 * go-rencode v0.1.0 - Go implementation of rencode - fast (basic)
                    object serialization similar to bencode
 * Copyright (C) 2015 gdm85 - https://github.com/gdm85/go-rencode/

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/
package rencode

//go:generate go run --tags=generate generate.go

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	DEFAULT_FLOAT_BITS = 32 // Default number of bits for serialized floats, either 32 or 64 (also a parameter for dumps()).
	MAX_INT_LENGTH     = 64 // Maximum length of integer when written as base 10 string.
	// The bencode 'typecodes' such as i, d, etc have been extended and relocated on the base-256 character set.
	CHR_LIST    = 59
	CHR_DICT    = 60
	CHR_INT     = 61
	CHR_INT1    = 62
	CHR_INT2    = 63
	CHR_INT4    = 64
	CHR_INT8    = 65
	CHR_FLOAT32 = 66
	CHR_FLOAT64 = 44
	CHR_TRUE    = 67
	CHR_FALSE   = 68
	CHR_NONE    = 69
	CHR_TERM    = 127
	// Positive integers with value embedded in typecode.
	INT_POS_FIXED_START = 0
	INT_POS_FIXED_COUNT = 44
	// Dictionaries with length embedded in typecode.
	DICT_FIXED_START = 102
	DICT_FIXED_COUNT = 25
	// Negative integers with value embedded in typecode.
	INT_NEG_FIXED_START = 70
	INT_NEG_FIXED_COUNT = 32
	// Strings with length embedded in typecode.
	STR_FIXED_START = 128
	STR_FIXED_COUNT = 64
	// Lists with length embedded in typecode.
	LIST_FIXED_START = STR_FIXED_START + STR_FIXED_COUNT
	LIST_FIXED_COUNT = 64
)

type Encoder struct {
	buffer bytes.Buffer
}

func (r *Encoder) Bytes() []byte {
	return r.buffer.Bytes()
}

func (r *Encoder) EncodeInt8(x int8) error {
	if 0 <= x && x < INT_POS_FIXED_COUNT {
		_, err := r.buffer.Write([]byte{byte(INT_POS_FIXED_START + x)})
		return err
	}
	if -INT_NEG_FIXED_COUNT <= x && x < 0 {
		_, err := r.buffer.Write([]byte{byte(INT_NEG_FIXED_START - 1 - x)})
		return err
	}
	if -128 < x && x <= 127 {
		_, err := r.buffer.Write([]byte{CHR_INT1, byte(x)})
		return err
	}
	panic("impossible just happened")
}

func (r *Encoder) EncodeBool(b bool) error {
	var data byte
	if b {
		data = CHR_TRUE
	} else {
		data = CHR_FALSE
	}

	_, err := r.buffer.Write([]byte{data})
	return err
}

func (r *Encoder) EncodeInt16(x int16) error {
	_, err := r.buffer.Write([]byte{CHR_INT2})
	if err != nil {
		return err
	}
	return binary.Write(&r.buffer, binary.BigEndian, x)
}

func (r *Encoder) EncodeInt32(x int32) error {
	_, err := r.buffer.Write([]byte{CHR_INT4})
	if err != nil {
		return err
	}
	return binary.Write(&r.buffer, binary.BigEndian, x)
}

func (r *Encoder) EncodeInt64(x int64) error {
	_, err := r.buffer.Write([]byte{CHR_INT8})
	if err != nil {
		return err
	}
	return binary.Write(&r.buffer, binary.BigEndian, x)
}

func (r *Encoder) EncodeBigNumber(s string) error {
	_, err := r.buffer.Write([]byte{CHR_INT})
	if err != nil {
		return err
	}
	_, err = r.buffer.Write([]byte(s))
	if err != nil {
		return err
	}
	_, err = r.buffer.Write([]byte{CHR_TERM})
	return err
}

func (r *Encoder) EncodeNone() error {
	_, err := r.buffer.Write([]byte{CHR_NONE})
	return err
}

func (r *Encoder) EncodeBytes(b []byte) error {
	if len(b) < STR_FIXED_COUNT {
		_, err := r.buffer.Write([]byte{byte(STR_FIXED_START + len(b))})
		if err != nil {
			return err
		}
		_, err = r.buffer.Write(b)
		return err
	}

	prefix := []byte(fmt.Sprintf("%d:", len(b)))

	_, err := r.buffer.Write(prefix)
	if err != nil {
		return err
	}

	_, err = r.buffer.Write(b)
	return err
}

func (r *Encoder) EncodeFloat32(f float32) error {
	_, err := r.buffer.Write([]byte{CHR_FLOAT32})
	if err != nil {
		return err
	}
	return binary.Write(&r.buffer, binary.BigEndian, f)
}

func (r *Encoder) EncodeFloat64(f float64) error {
	_, err := r.buffer.Write([]byte{CHR_FLOAT64})
	if err != nil {
		return err
	}
	return binary.Write(&r.buffer, binary.BigEndian, f)
}
