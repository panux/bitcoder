package bitcoder

import "testing"

func TestOverShiftRegression(t *testing.T) {
	code := "1AA"
	_, vars := decode(code)
	if vars[0].shift != 0 {
		t.Errorf("Overshift regression - expected shift 0 but got %d", vars[0].shift)
	}
}

func TestFastCoder(t *testing.T) {
	code := "1AA"
	A := uint64(3)
	expected := uint64(7)
	actual := NewFastCoder(code)(A)
	if expected != actual {
		t.Log(decode(code)) //Debugging info
		t.Errorf("Expected result %d from code %q but got %d", expected, code, actual)
	}
}

func TestValueCoder(t *testing.T) {
	code := "1ABBCCCC"
	A := true
	B := 3
	C := uint8(15)
	expected := uint64(255)
	actual := NewEasyCoder(code)(A, &B, C)
	if expected != actual {
		t.Errorf("Expected result %d from code %q but got %d", expected, code, actual)
	}
}

type flagtest struct {
	A, B, C bool
}

type reftest struct {
	Flags flagtest `bitpack:"1ABC"`
	Iptr  ***int
}

func TestReflectCoder(t *testing.T) {
	code := "FFFFIIB0"
	expected := NewFastCoder("10101100")()
	three := 3
	pthree := &three
	ppthree := &pthree
	pppthree := &ppthree
	actual := NewEasyCoder(code)(
		reftest{
			Flags: flagtest{
				A: false,
				B: true,
				C: false,
			},
			Iptr: pppthree,
		},
		struct {
			B bool
		}{
			B: false,
		},
	)
	if expected != actual {
		t.Errorf("Expected result %d from code %q but got %d", expected, code, actual)
	}
}
