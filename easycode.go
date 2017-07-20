package bitcoder

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

//EasyCoder is a bitcoder that uses type switches and/or reflection to get values
type EasyCoder func(...interface{}) uint64

type easyclass uint8

const (
	classNone easyclass = iota
	classMap
	classVals
)

//NewEasyCoder creates a new EasyCoder from the provided bitpacking code
//All letters are capitalized
func NewEasyCoder(code string) EasyCoder {
	code = strings.ToUpper(code)
	_, args := decode(code)
	fast := NewFastCoder(code)
	return func(dat ...interface{}) (r uint64) {
		//Identify easyclass
		class := classNone
		for _, v := range dat {
			isVal, isPos := func() (val bool, pos bool) {
				val = true
				pos = true
			tloop:
				switch v.(type) {
				case int:
					pos = v.(int) >= 0
				case int8:
					pos = v.(int8) >= 0
				case int16:
					pos = v.(int16) >= 0
				case int32:
					pos = v.(int32) >= 0
				case int64:
					pos = v.(int64) >= 0
				case uint:
				case uint8:
				case uint16:
				case uint32:
				case uint64:
				case bool:
				default:
					if reflect.TypeOf(v).Kind() == reflect.Ptr {
						v = reflect.ValueOf(v).Elem().Interface()
						goto tloop
					}
					return false, true
				}
				return
			}()
			if !isPos {
				panic(errors.New("Negative argument"))
			}
			switch class {
			case classNone:
				if isVal {
					class = classVals
				} else {
					class = classMap
				}
			case classVals:
				if !isVal {
					panic(errors.New("easyclass inconsistency"))
				}
			case classMap:
				if isVal {
					panic(errors.New("easyclass inconsistency"))
				}
			}
		}
		switch class {
		case classNone:
			panic(errors.New("No bitpack arguments"))
		case classMap:
			odat := dat
			dat = make([]interface{}, len(args))
			alst := make(map[rune]bool)
			for _, a := range args {
				alst[a.letter] = true
			}
			dmap := make(map[rune]interface{})
			for _, v := range odat {
				vv := reflect.ValueOf(v)
			kloop:
				switch vv.Kind() {
				case reflect.Ptr:
					vv = vv.Elem() //Dereference
					goto kloop
				case reflect.Map:
					for _, key := range vv.MapKeys() {
						k := key.Interface()
						value := vv.MapIndex(key).Interface()
						var kl rune
					mtswitch:
						switch k.(type) {
						case string:
							k = []rune(k.(string))[0]
							goto mtswitch
						case rune:
							kl = unicode.ToUpper(k.(rune))
						default:
							panic(fmt.Errorf("Invalid key type for map value %q", k))
						}
						if !alst[kl] {
							panic(fmt.Errorf("Invalid key rune %q", k.(rune)))
						}
						dmap[kl] = value
						delete(alst, kl)
					}
				case reflect.Struct:
					nf := vv.NumField()
					vt := vv.Type()
					for fi := 0; fi < nf; fi++ {
						f := vt.Field(fi)
						fv := vv.Field(fi)
						fl := []rune(f.Name)[0]
						if !alst[fl] {
							panic(fmt.Errorf("bitpack field %q (%q) already filled/not present", fl, f.Name))
						}
					fkloop:
						switch fv.Kind() {
						case reflect.Ptr: //Automatically dereference
							fv = fv.Elem()
							goto fkloop
						case reflect.Map:
							fallthrough
						case reflect.Struct:
							bp := f.Tag.Get("bitpack") //sub-pack maps/structs
							if bp == "" {
								panic("Substruct without corresponding bitpack")
							}
							dmap[fl] = NewEasyCoder(bp)(vv.FieldByIndex(f.Index).Interface())
						default:
							dmap[fl] = fv.Interface() //passthrough other values
						}
						delete(alst, fl)
					}
				}
			}
			if len(alst) > 0 { //Check for missed inputs
				missed := []int{}
				for c := range alst {
					missed = append(missed, int(c))
				}
				sort.Ints(missed)
				missedc := make([]rune, len(missed))
				for i, c := range missed {
					missedc[i] = rune(c)
				}
				panic(fmt.Errorf("Missed inputs: %q", string(missedc)))
			}
			for i, a := range args {
				dat[i] = dmap[a.letter]
			}
			fallthrough
		case classVals:
			ar := make([]uint64, len(args))
			for i, v := range dat {
			vloop:
				switch v.(type) {
				case int:
					ar[i] = uint64(v.(int))
				case int8:
					ar[i] = uint64(v.(int8))
				case int16:
					ar[i] = uint64(v.(int16))
				case int32:
					ar[i] = uint64(v.(int32))
				case int64:
					ar[i] = uint64(v.(int64))
				case uint:
					ar[i] = uint64(v.(uint))
				case uint8:
					ar[i] = uint64(v.(uint8))
				case uint16:
					ar[i] = uint64(v.(uint16))
				case uint32:
					ar[i] = uint64(v.(uint32))
				case uint64:
					ar[i] = v.(uint64)
				case bool:
					if v.(bool) {
						ar[i] = 1
					} else {
						ar[i] = 0
					}
				default:
					if reflect.TypeOf(v).Kind() == reflect.Ptr {
						v = reflect.ValueOf(v).Elem().Interface()
						goto vloop
					}
					panic(errors.New("This should never happen"))
				}
			}
			return fast(ar...)
		}
		panic(errors.New("This will never happen"))
	}
}
