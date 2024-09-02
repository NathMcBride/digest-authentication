package structmarshal

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/errors"
	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structinfo"
)

type StructMarshal struct {
}

func (sm *StructMarshal) Marshal(buffer *bytes.Buffer, info *structinfo.Info, val reflect.Value) error {
	for i := range info.Fields {
		finfo := &info.Fields[i]

		var vf reflect.Value
		for _, x := range finfo.Idx {
			vf = val.Field(x)
		}

		if finfo.Flags&structinfo.FMode != 0 {
			if finfo.Flags&structinfo.FOmitEmpty != 0 && isEmptyValue(vf) {
				continue
			}
			if finfo.Flags&structinfo.FUnq != 0 {
				WriteDelim(buffer, ",", i > 0)
				m, err := MarshalSimple(vf)
				if err != nil {
					return err
				}
				WriteParam(buffer, finfo.Name, m, false)
				continue
			}
		}

		WriteDelim(buffer, ",", i > 0)
		m, err := MarshalSimple(vf)
		if err != nil {
			return err
		}
		WriteParam(buffer, finfo.Name, m, true)
	}
	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Interface, reflect.Pointer:
		return v.IsZero()
	}
	return false
}

func MarshalSimple(val reflect.Value) (string, error) {
	switch val.Kind() {
	case reflect.String:
		return val.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil
	}

	return "", errors.MarshalError("marshal simple")
}

func WriteParam(buffer *bytes.Buffer, name string, value string, quoted bool) {
	v := value
	if quoted {
		v = `"` + value + `"`
	}

	buffer.WriteString(name + "=" + v)
}

func WriteDelim(buffer *bytes.Buffer, delim string, include bool) {
	if !include {
		return
	}
	buffer.WriteString(delim + " ")
}
