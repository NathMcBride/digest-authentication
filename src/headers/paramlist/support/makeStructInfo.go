package support

import "github.com/NathMcBride/digest-authentication/src/headers/paramlist/structinfo"

type MakeStructInfo struct {
	Info structinfo.Info
}

func NewMakeStructInfo() *MakeStructInfo {
	return &MakeStructInfo{
		Info: structinfo.Info{
			Fields: []structinfo.FieldInfo{},
		}}
}

func (ms *MakeStructInfo) WithFUnqFlag() *MakeStructInfo {
	return ms.AddField("field", structinfo.FUnq)
}

func (ms *MakeStructInfo) WithFOmitEmptyFlag() *MakeStructInfo {
	return ms.AddField("field", structinfo.FOmitEmpty)
}

func (ms *MakeStructInfo) WithAllFlags() *MakeStructInfo {
	return ms.AddField("field", structinfo.FUnq|structinfo.FOmitEmpty)
}

func (ms *MakeStructInfo) WithNoFlags() *MakeStructInfo {
	return ms.AddField("field", 0)
}

func (ms *MakeStructInfo) AddField(name string, flags structinfo.FieldFlags) *MakeStructInfo {
	idx := 0
	length := len(ms.Info.Fields)
	if length > 0 {
		idx = length
	}

	ms.Info.Fields = append(ms.Info.Fields,
		structinfo.FieldInfo{
			Idx:   []int{idx},
			Name:  name,
			Flags: flags})
	return ms
}

func (ms *MakeStructInfo) Build() structinfo.Info {
	return ms.Info
}

type TestStruct struct {
	Field string
}
