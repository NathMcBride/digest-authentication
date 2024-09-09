package fakes

import (
	"reflect"

	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structinfo"
)

type FakeStructInfoer struct {
	getTypeInfoCallCount   int
	getTypeInfoReturns     struct{ info structinfo.Info }
	getTypeInfoArgsForCall []struct{ typ reflect.Type }
}

func (fs *FakeStructInfoer) GetTypeInfoCallCount() int {
	return fs.getTypeInfoCallCount
}

func (fs *FakeStructInfoer) GetTypeInfoReturns(info structinfo.Info) {
	fs.getTypeInfoReturns = struct{ info structinfo.Info }{info}
}

func (fs *FakeStructInfoer) GetTypeInfoArgsForCall(i int) reflect.Type {
	return fs.getTypeInfoArgsForCall[0].typ
}

func (fs *FakeStructInfoer) GetTypeInfo(typ reflect.Type) *structinfo.Info {
	fs.getTypeInfoCallCount++
	fs.getTypeInfoArgsForCall = append(fs.getTypeInfoArgsForCall, struct{ typ reflect.Type }{typ})
	return &fs.getTypeInfoReturns.info
}

func (fs *FakeStructInfoer) FieldInfo(f *reflect.StructField) *structinfo.FieldInfo { return nil }
