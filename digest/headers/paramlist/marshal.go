package paramlist

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
)

type typeInfo struct {
	Fields []fieldInfo
}

type fieldInfo struct {
	Idx   []int
	Name  string
	Flags fieldFlags
}

type fieldFlags int

const (
	fElement fieldFlags = 1 << iota
	fUnq
	fOmitEmpty

	fMode = fElement | fUnq | fOmitEmpty
)

type Marshaler struct{}

func (m *Marshaler) Marshal(v any) ([]byte, error) {
	val := reflect.ValueOf(v)
	if !val.IsValid() {
		return nil, MarshalError("invalid value passed to marshal")
	}

	tinfo := GetTypeInfo(val.Type())
	buffer := new(bytes.Buffer)

	err := MarshalStruct(buffer, tinfo, val)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func GetTypeInfo(typ reflect.Type) *typeInfo {
	tinfo := typeInfo{}
	if typ.Kind() == reflect.Struct {
		n := typ.NumField()
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if (!f.IsExported() && !f.Anonymous) || f.Tag.Get("hparam") == "-" {
				continue
			}

			finfo := StructFieldInfo(&f)
			tinfo.Fields = append(tinfo.Fields, *finfo)
		}
	}
	return &tinfo
}

func StructFieldInfo(f *reflect.StructField) *fieldInfo {
	finfo := &fieldInfo{Idx: f.Index}
	tag := f.Tag.Get("hparam")

	tokens := strings.Split(tag, ",")
	if len(tokens) == 1 {
		finfo.Flags = fElement
	} else {
		tag = tokens[0]
		for _, flag := range tokens[1:] {
			switch flag {
			case "unq":
				finfo.Flags |= fUnq
			case "omitempty":
				finfo.Flags |= fOmitEmpty
			}
		}
	}

	if tag == "" {
		finfo.Name = f.Name
		return finfo
	}
	finfo.Name = tag

	return finfo
}

func MarshalStruct(buffer *bytes.Buffer, tinfo *typeInfo, val reflect.Value) error {
	for i := range tinfo.Fields {
		finfo := &tinfo.Fields[i]

		var vf reflect.Value
		for _, x := range finfo.Idx {
			vf = val.Field(x)
		}

		if finfo.Flags&fMode != 0 {
			if finfo.Flags&fOmitEmpty != 0 && isEmptyValue(vf) {
				continue
			}
			if finfo.Flags&fUnq != 0 {
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

	return "", MarshalError("marshal simple")
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
