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
	J     map[rune]int `bitpack:"XYZ1"`
}

func TestReflectCoder(t *testing.T) {
	code := "FFFFIIB0JJJJ"
	expected := NewFastCoder("101011001111")()
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
			J: map[rune]int{
				'X': 1,
				'Y': 1,
				'Z': 1,
			},
		},
		map[string]uint{
			"b": uint(0),
		},
	)
	if expected != actual {
		t.Errorf("Expected result %d from code %q but got %d", expected, code, actual)
	}
}

var vals = struct {
	A int
	B int8
	C int16
	D int32
	E int64
	F uint
	G uint8
	H uint16
	I uint32
	J uint64
	K bool
}{
	A: 1,
	B: 0,
	C: 1,
	D: 0,
	E: 1,
	F: 0,
	G: 1,
	H: 0,
	I: 1,
	J: 0,
	K: true,
}

var valstest = NewEasyCoder("abcdefghijk")
var valsexpected = NewFastCoder("10101010101")()

func TestTypes(t *testing.T) {
	v := vals
	actual := valstest(&v)
	if valsexpected != actual {
		t.Errorf("Expected result %d but got %d", valsexpected, actual)
	}
}

func TestErrNegInt(t *testing.T) {
	v := vals
	v.A = -1
	errtest(t, "negint", "Negative argument", func() {
		valstest(v)
	})
}
func TestErrNegInt8(t *testing.T) {
	v := vals
	v.B = -1
	errtest(t, "negint8", "Negative argument", func() {
		valstest(v)
	})
}
func TestErrNegInt16(t *testing.T) {
	v := vals
	v.C = -1
	errtest(t, "negint16", "Negative argument", func() {
		valstest(v)
	})
}
func TestErrNegInt32(t *testing.T) {
	v := vals
	v.D = -1
	errtest(t, "negint32", "Negative argument", func() {
		valstest(v)
	})
}
func TestErrNegInt64(t *testing.T) {
	v := vals
	v.E = -1
	errtest(t, "negint64", "Negative argument", func() {
		valstest(v)
	})
}

func TestErrInconsistentVM(t *testing.T) {
	errtest(t, "inconsistent-value-map", "easyclass inconsistency", func() {
		NewEasyCoder("")(1, struct{}{})
	})
}
func TestErrInconsistentMV(t *testing.T) {
	errtest(t, "inconsistent-map-value", "easyclass inconsistency", func() {
		NewEasyCoder("")(struct{}{}, 1)
	})
}
func TestErrNothing(t *testing.T) {
	errtest(t, "nothing", "No bitpack arguments", func() {
		NewEasyCoder("A")()
	})
}

func TestErrInvalidMap(t *testing.T) {
	errtest(t, "bad-map-key", "Invalid key type for map value '\\x01'", func() { //Meh, this is OK
		NewEasyCoder("A")(map[int]int{1: 1})
	})
}
func TestErrInvalidMapRune(t *testing.T) {
	errtest(t, "bad-map-rune", "Invalid key rune 'B'", func() { //Meh, this is OK
		NewEasyCoder("A")(map[rune]int{'B': 1})
	})
}
func TestErrInvalidVal(t *testing.T) {
	v := vals
	vtest := NewEasyCoder("abcdefghjk")
	errtest(t, "invalid-val", "bitpack field 'I' (\"I\") already filled/not present", func() {
		vtest(v)
	})
}
func TestErrMissingTag(t *testing.T) {
	errtest(t, "missing-tag", "Substruct without corresponding bitpack", func() {
		NewEasyCoder("A")(struct{ A struct{} }{A: struct{}{}})
	})
}
func TestErrMissingValue(t *testing.T) {
	errtest(t, "missing-value", "Missed inputs: \"A\"", func() {
		NewEasyCoder("A")(struct{}{})
	})
}

func TestTooMany(t *testing.T) {
	errtest(t, "too-many", "Too many args to FastCoder", func() {
		NewFastCoder("A")(1, 2)
	})
}
func TestTooFew(t *testing.T) {
	errtest(t, "too-few", "Too few args to FastCoder", func() {
		NewFastCoder("AB")(1)
	})
}
func TestOversize(t *testing.T) {
	errtest(t, "oversize-arg", "Oversized argument 0 - should be 1 bits but is 1 bits", func() {
		NewFastCoder("A")(2)
	})
}

func TestOver64(t *testing.T) {
	errtest(t, "oversize-bitpack", "Bitpacking code must be less than 64 bits", func() {
		decode("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	})
}
func TestIllegalRune(t *testing.T) {
	errtest(t, "illegal-rune", "Illegal character ';' in bitpacking code", func() {
		decode(";")
	})
}

func errtest(t *testing.T, name string, goal string, f func()) {
	tst := testerr(f)
	if tst == nil {
		t.Fatalf("Test %q should have thrown an error but did not", name)
	} else if tst.Error() != goal {
		t.Fatalf("Test %q should have thrown error %q but threw %q", name, goal, tst.Error())
	} else {
		t.Logf("Error test %q succeeded with error %q", name, tst.Error())
	}
}

func testerr(f func()) (err error) {
	defer func() {
		e := recover()
		if e == nil {
			err = nil
		} else {
			err = e.(error)
		}
	}()
	f()
	panic(nil)
}
