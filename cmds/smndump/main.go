// Copyright 2010-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/memio"
)

// awkward!

// newInt constructs a UintN with the specified value and bits.
func newInt(val uint64, bits int) memio.UintN {
	switch bits {
	case 8:
		val := memio.Uint8(int8(val))
		return &val
	case 16:
		val := memio.Uint16(uint16(val))
		return &val
	case 32:
		val := memio.Uint32(uint32(val))
		return &val
	case 64:
		val := memio.Uint64(uint64(val))
		return &val
	default:
		panic(fmt.Sprintf("invalid number of bits %d", bits))
	}
}

/*
fn snmr(a: u32) -> u32 {
    // the smn device is at (0)
    unsafe {
        outl(0xcf8, 0x800000b8);
    outl(0xcfc, a);
    outl(0xcf8, 0x800000bc);
    inl(0xcfc)
    }
}*/
var s = []struct {
	a uint64
	z uint64
}{
	{0x02800000, 0x40},
}

func main() {
	var xindex = newInt(0x8000005c, 32)
	var index = newInt(0x80000060, 32)
	var data = newInt(0x80000064, 32)
	flag.Parse()
	a := flag.Args()
	if len(a) == 0 {
		a = []string{"2800000/40"}
	}
	for _, el := range a {
		b := strings.Split(el, "/")
		if len(b) != 2 {
			log.Fatalf("Could not split %v", el)
		}
		addr, err := strconv.ParseUint(b[0], 16, 36)
		if err != nil {
			log.Fatal(err)
		}
		amt, err := strconv.ParseUint(b[1], 16, 32)
		if err != nil {
			log.Fatal(err)
		}
		for i := addr; i < addr + amt; i += 4 {
			high := newInt((i>>32)+i, 32)
			low := newInt(i, 32)
			if err := memio.ArchOut(0xcf8, xindex); err != nil {
				log.Fatal(err)
			}
			if err := memio.ArchOut(0xcfc, high); err != nil {
				log.Fatal(err)
			}
			if err := memio.ArchOut(0xcf8, index); err != nil {
				log.Fatal(err)
			}
			if err := memio.ArchOut(0xcfc, low); err != nil {
				log.Fatal(err)
			}
			if err := memio.ArchOut(0xcf8, data); err != nil {
				log.Fatal(err)
			}
			var v = newInt(0, 32)
			if err := memio.ArchIn(0xcfc, v); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%#09x: %s\n", i, v)
		}
	}
}
