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
	"fmt"
	"math"
	"math/big"
)

// Encode is the generic encoder method that will encode any of the following supported types:
// * big.Int
// * List
// * Dictionary
// * bool
// * float32, float34
// * []byte, string (all strings are stored as byte slices anyway)
// * int8, int16, int32, int64, int
// * uint8, uint16, uint32, uint64, uint
func (r *Encoder) Encode(data interface{}) error {
	if data == nil {
		return r.EncodeNone()
	}
	switch data.(type) {
	case big.Int:
		x := data.(big.Int)
		s := x.String()
		if len(s) > MAX_INT_LENGTH {
			return fmt.Errorf("Number is longer than %d characters", MAX_INT_LENGTH)
		}
		return r.EncodeBigNumber(s)
	case List:
		x := data.(List)
		if x.Length() < LIST_FIXED_COUNT {
			_, err := r.buffer.Write([]byte{byte(LIST_FIXED_START + x.Length())})
			if err != nil {
				return err
			}
			for _, v := range x.Values() {
				err = r.Encode(v)
				if err != nil {
					return err
				}
			}
			return nil
		}
		_, err := r.buffer.Write([]byte{byte(CHR_LIST)})
		if err != nil {
			return err
		}

		for _, v := range x.Values() {
			err = r.Encode(v)
			if err != nil {
				return err
			}
		}

		_, err = r.buffer.Write([]byte{byte(CHR_TERM)})
		return err
	case Dictionary:
		x := data.(Dictionary)
		if x.Length() < DICT_FIXED_COUNT {
			_, err := r.buffer.Write([]byte{byte(DICT_FIXED_START + x.Length())})
			if err != nil {
				return err
			}
			keys := x.Keys()
			for i, v := range x.Values() {
				err = r.Encode(keys[i])
				if err != nil {
					return err
				}
				err = r.Encode(v)
				if err != nil {
					return err
				}
			}
			return nil
		}
		_, err := r.buffer.Write([]byte{byte(CHR_DICT)})
		if err != nil {
			return err
		}
		keys := x.Keys()
		for i, v := range x.Values() {
			err = r.Encode(keys[i])
			if err != nil {
				return err
			}
			err = r.Encode(v)
			if err != nil {
				return err
			}
		}

		_, err = r.buffer.Write([]byte{byte(CHR_TERM)})
		return err
	case bool:
		return r.EncodeBool(data.(bool))
	case float32:
		return r.EncodeFloat32(data.(float32))
	case float64:
		return r.EncodeFloat64(data.(float64))
	case []byte:
		return r.EncodeBytes(data.([]byte))
	case string:
		// all strings will be treated as byte arrays
		return r.EncodeBytes([]byte(data.(string)))
	case int8:
		return r.EncodeInt8(data.(int8))
	case int:
		x := data.(int)
		if math.MinInt8 <= x && x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		if math.MinInt16 <= x && x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}
		if math.MinInt32 <= x && x <= math.MaxInt32 {
			return r.EncodeInt32(int32(x))
		}
		return r.EncodeInt64(int64(x))
	case uint8:
		x := data.(uint8)
		if x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
	case uint16:
		x := data.(uint16)
		if x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		if x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}
	case int16:
		x := data.(int16)
		if math.MinInt8 <= x && x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		return r.EncodeInt16(int16(x))
	case uint32:
		x := data.(uint32)
		if x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		if x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}
		if x <= math.MaxInt32 {
			return r.EncodeInt32(int32(x))
		}
	case int32:
		x := data.(int32)
		if math.MinInt8 <= x && x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		if math.MinInt16 <= x && x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}
		return r.EncodeInt32(int32(x))
	case int64:
		x := data.(int64)
		if math.MinInt8 <= x && x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		if math.MinInt16 <= x && x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}
		if math.MinInt32 <= x && x <= math.MaxInt32 {
			return r.EncodeInt32(int32(x))
		}
		return r.EncodeInt64(int64(x))
	case uint64, uint:
		s := fmt.Sprintf("%d", data)
		if len(s) > MAX_INT_LENGTH {
			return fmt.Errorf("Number is longer than %d characters", MAX_INT_LENGTH)
		}
		return r.EncodeBigNumber(s)
	default:
		return fmt.Errorf("could not encode data of type %T", data)
	}
	panic("unexpected fallthrough")
}
