package project_package

import (
	"reflect"
	"testing"
)

func TestGetTime(t *testing.T){
	var intTime interface{} = GetTime()
	_, ok := intTime.(int)
	if !ok {
		t.Errorf("GetTime return value type is %v, want int", reflect.TypeOf(intTime))
	}
}
