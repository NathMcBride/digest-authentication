package paramlist

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/errors"
)

type Parser interface {
	Parse(auth string) (map[string]string, error)
}

type UnMarshaler struct {
	StructInfoer StructInfoer
	Parser       Parser
}

func (um *UnMarshaler) Unmarshal(data []byte, v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Pointer {
		return errors.UnmarshalError("not a pointer")
	}

	if val.IsNil() {
		return errors.UnmarshalError("nil pointer")
	}

	parsed, err := um.Parser.Parse(string(data))
	if err != nil {
		return err
	}

	val = val.Elem()
	typ := val.Type()

	tinfo := um.StructInfoer.GetTypeInfo(typ)
	for i := range tinfo.Fields {
		finfo := &tinfo.Fields[i]

		var vf reflect.Value
		for _, x := range finfo.Idx {
			vf = val.Field(x)
		}

		src := parsed[finfo.Name]
		switch vf.Kind() {
		case reflect.String:
			vf.SetString(src)
		case reflect.Bool:
			if len(src) == 0 {
				vf.SetBool(false)
				return nil
			}
			value, err := strconv.ParseBool(strings.TrimSpace(string(src)))
			if err != nil {
				return errors.UnmarshalError("unable to parse bool")
			}
			vf.SetBool(value)
		}
	}

	return nil
}
