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
