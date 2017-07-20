package bitcoder

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type bcentry struct {
	shift  uint8
	size   uint8
	letter rune
}

func (e bcentry) checkSize(val uint64) bool {
	return (val >> e.size) == 0
}

func decode(code string) (uint64, []bcentry) {
	code = strings.Replace(code, " ", "", -1)
	if len(code) > 64 {
		panic(errors.New("Bitpacking code must be less than 64 bits"))
	}
	cdat := uint64(0)
	prev := '0'
	args := []bcentry{}
	sh := uint8(len(code) - 1)
	for _, v := range []rune(code) {
		if v == '0' || v == '1' {
			val := uint64(v - '0')
			cdat |= val << sh
		} else if unicode.IsLetter(v) {
			if v == prev {
				args[len(args)-1].size++
				args[len(args)-1].shift = sh
			} else {
				args = append(args, bcentry{
					shift:  sh,
					size:   1,
					letter: v,
				})
			}
		} else {
			panic(fmt.Errorf("Illegal character %q in bitpacking code", v))
		}
		prev = v
		sh--
	}
	return cdat, args
}
