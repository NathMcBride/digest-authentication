package paramlist

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/NathMcBride/web-authentication/digest/headers/paramlist/errors"
	"github.com/NathMcBride/web-authentication/digest/headers/paramlist/structinfo"
	"github.com/NathMcBride/web-authentication/digest/parsers"
)

// Test
func Unmarshal(data []byte, v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Pointer {
		return errors.UnmarshalError("not a pointer")
	}

	if val.IsNil() {
		return errors.UnmarshalError("nil pointer")
	}

	parsed, err := parsers.ParseDigestAuth(string(data))
	if err != nil {
		return err
	}

	val = val.Elem()
	typ := val.Type()
	s := structinfo.StructInfo{}
	tinfo := s.GetTypeInfo(typ)

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
