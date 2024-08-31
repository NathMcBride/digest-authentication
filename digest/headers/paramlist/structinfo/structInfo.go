package structinfo

import (
	"reflect"
	"strings"
)

type Info struct {
	Fields []FieldInfo
}

type FieldInfo struct {
	Idx   []int
	Name  string
	Flags FieldFlags
}

type FieldFlags int

const (
	FElement FieldFlags = 1 << iota
	FUnq
	FOmitEmpty

	FMode = FElement | FUnq | FOmitEmpty
)

type StructInfo struct {
}

func (si *StructInfo) GetTypeInfo(typ reflect.Type) *Info {
	info := Info{}
	if typ.Kind() == reflect.Struct {
		n := typ.NumField()
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if (!f.IsExported() && !f.Anonymous) || f.Tag.Get("hparam") == "-" {
				continue
			}

			finfo := si.FieldInfo(&f)
			info.Fields = append(info.Fields, *finfo)
		}
	}
	return &info
}

func (si *StructInfo) FieldInfo(f *reflect.StructField) *FieldInfo {
	finfo := &FieldInfo{Idx: f.Index}
	tag := f.Tag.Get("hparam")

	tokens := strings.Split(tag, ",")
	if len(tokens) == 1 {
		finfo.Flags = FElement
	} else {
		tag = tokens[0]
		for _, flag := range tokens[1:] {
			switch flag {
			case "unq":
				finfo.Flags |= FUnq
			case "omitempty":
				finfo.Flags |= FOmitEmpty
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
