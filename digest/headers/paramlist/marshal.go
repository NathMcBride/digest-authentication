package paramlist

import (
	"bytes"
	"reflect"

	"github.com/NathMcBride/web-authentication/digest/headers/paramlist/errors"
	"github.com/NathMcBride/web-authentication/digest/headers/paramlist/structinfo"
)

type StructInfoer interface {
	GetTypeInfo(typ reflect.Type) *structinfo.Info
	FieldInfo(f *reflect.StructField) *structinfo.FieldInfo
}

type StructMarshaler interface {
	Marshal(buffer *bytes.Buffer, info *structinfo.Info, val reflect.Value) error
}

type Marshaler struct {
	StructInfoer    StructInfoer
	StructMarshaler StructMarshaler
}

func (m *Marshaler) Marshal(v any) ([]byte, error) {
	val := reflect.ValueOf(v)
	if !val.IsValid() {
		return nil, errors.MarshalError("invalid value passed to marshal")
	}

	tinfo := m.StructInfoer.GetTypeInfo(val.Type())
	buffer := new(bytes.Buffer)

	err := m.StructMarshaler.Marshal(buffer, tinfo, val)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
