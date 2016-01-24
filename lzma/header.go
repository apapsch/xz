// Copyright 2015 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import (
	"errors"
	"fmt"
)

// uint32LE reads an uint32 integer from a byte slize
func uint32LE(b []byte) uint32 {
	x := uint32(b[3]) << 24
	x |= uint32(b[2]) << 16
	x |= uint32(b[1]) << 8
	x |= uint32(b[0])
	return x
}

// uint64LE converts the uint64 value stored as little endian to an uint64
// value.
func uint64LE(b []byte) uint64 {
	x := uint64(b[7]) << 56
	x |= uint64(b[6]) << 48
	x |= uint64(b[5]) << 40
	x |= uint64(b[4]) << 32
	x |= uint64(b[3]) << 24
	x |= uint64(b[2]) << 16
	x |= uint64(b[1]) << 8
	x |= uint64(b[0])
	return x
}

// putUint32LE puts an uint32 integer into a byte slice that must have at least
// a lenght of 4 bytes.
func putUint32LE(b []byte, x uint32) {
	b[0] = byte(x)
	b[1] = byte(x >> 8)
	b[2] = byte(x >> 16)
	b[3] = byte(x >> 24)
}

// putUint64LE puts the uint64 value into the byte slice as little endian
// value. The byte slice b must have at least place for 8 bytes.
func putUint64LE(b []byte, x uint64) {
	b[0] = byte(x)
	b[1] = byte(x >> 8)
	b[2] = byte(x >> 16)
	b[3] = byte(x >> 24)
	b[4] = byte(x >> 32)
	b[5] = byte(x >> 40)
	b[6] = byte(x >> 48)
	b[7] = byte(x >> 56)
}

// noHeaderSize defines the value of the length field in the LZMA header.
const noHeaderSize uint64 = 1<<64 - 1

// maximum header length
const headerLen = 13

// Header represents the header of an LZMA file.
type Header struct {
	Properties Properties
	DictCap    int
	// uncompressed size; negative value if no size is given
	Size int64
}

// marshalBinary marshals the header.
func (h *Header) marshalBinary() (data []byte, err error) {
	if err = h.Properties.Verify(); err != nil {
		return nil, err
	}
	if !(0 <= h.DictCap && int64(h.DictCap) <= MaxDictCap) {
		return nil, fmt.Errorf("lzma: DictCap %d out of range",
			h.DictCap)
	}

	data = make([]byte, 13)

	// property byte
	data[0] = h.Properties.Code()

	// dictionary capacity
	putUint32LE(data[1:5], uint32(h.DictCap))

	// uncompressed size
	var s uint64
	if h.Size > 0 {
		s = uint64(h.Size)
	} else {
		s = noHeaderSize
	}
	putUint64LE(data[5:], s)

	return data, nil
}

// unmarshalBinary unmarshals the header.
func (h *Header) unmarshalBinary(data []byte) error {
	if len(data) != headerLen {
		return errors.New("lzma.unmarshalBinary: data has wrong length")
	}

	// properties
	var err error
	if h.Properties, err = PropertiesForCode(data[0]); err != nil {
		return err
	}

	// dictionary capacity
	h.DictCap = int(uint32LE(data[1:]))
	if h.DictCap < 0 {
		return errors.New(
			"LZMA header: dictionary capacity exceeds maximum " +
				"integer")
	}

	// uncompressed size
	s := uint64LE(data[5:])
	if s == noHeaderSize {
		h.Size = -1
	} else {
		h.Size = int64(s)
		if h.Size < 0 {
			return errors.New(
				"LZMA header: uncompressed size " +
					"out of int64 range")
		}
	}

	return nil
}
