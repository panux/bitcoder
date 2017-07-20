package bitcoder

import (
	"errors"
	"fmt"
	"math"
)

//FastCoder is the fastest bitcoder - it takes data in as variadic args
type FastCoder func(...uint64) uint64

//NewFastCoder decodes a bitpacking code and returns a FastCoder
func NewFastCoder(code string) FastCoder {
	cnst, args := decode(code)
	return func(dat ...uint64) (r uint64) {
		if len(dat) > len(args) {
			panic(errors.New("Too many args to FastCoder"))
		} else if len(dat) < len(args) {
			panic(errors.New("Too few args to FastCoder"))
		}
		r = cnst
		for i, a := range args {
			v := dat[i]
			if !a.checkSize(v) {
				bits := int(math.Ceil(math.Log2(float64(v))))
				panic(fmt.Errorf("Oversized argument %d - should be %d bits but is %d bits", i, a.size, bits))
			}
			r |= v << a.shift
		}
		return
	}
}
