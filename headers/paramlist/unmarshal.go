package paramlist

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/NathMcBride/web-authentication/parsers"
)

func Unmarshal(data []byte, v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Pointer {
		return UnmarshalError("not a pointer")
	}

	if val.IsNil() {
		return UnmarshalError("nil pointer")
	}

	parsed, err := parsers.ParseDigestAuth(string(data))
	if err != nil {
		return err
	}

	val = val.Elem()
	typ := val.Type()
	tinfo := GetTypeInfo(typ)

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
				return UnmarshalError("unable to parse bool")
			}
			vf.SetBool(value)
		}

	}

	return nil
}
