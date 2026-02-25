package structs_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nyttikord/gokord/internal/structs"
)

func validMap(t *testing.T, check, res map[string]any) {
	if len(check) != len(res) {
		t.Errorf("invalid value: got %v, wanted %v", check, res)
	}
	for k, val := range check {
		if !reflect.DeepEqual(res[k], val) {
			t.Errorf("invalid value at %s, got %v, wanted %v", k, val, res[k])
		}
	}
}

type test1 struct {
	A string `json:"a"`
	B string `json:",omitempty"`
}

func TestMarshalToMap_Simple(t *testing.T) {
	h := test1{"aaa", ""}
	v := structs.MarshalToMap(h)
	validMap(t, v, map[string]any{"a": "aaa"})

	h = test1{"", "bbb"}
	v = structs.MarshalToMap(h)
	validMap(t, v, map[string]any{"a": "", "B": "bbb"})
}

type test2 struct {
	Hey *test1 `json:"hey,omitempty"`
	Key string `json:"key"`
}

func TestMarshalToMap_Nested(t *testing.T) {
	h := test2{&test1{"aaa", ""}, "key"}
	v := structs.MarshalToMap(h)
	validMap(t, v, map[string]any{"hey": map[string]any{"a": "aaa"}, "key": "key"})

	h = test2{nil, "k"}
	v = structs.MarshalToMap(h)
	validMap(t, v, map[string]any{"key": "k"})
}

func (t test1) String() string {
	return fmt.Sprintf("%s:%s", t.A, t.B)
}

type test3 struct {
	Default int `json:"default"`
	ConvDef int `json:"conv,string"`
}

func TestMarshalToMap_String(t *testing.T) {
	h := test3{1, 1}
	v := structs.MarshalToMap(h)
	validMap(t, v, map[string]any{"default": 1, "conv": "1"})
}

type test4 struct {
	A int
}

func (t test4) MarshalMap() map[string]any {
	return map[string]any{"a": t.A + 5}
}

func TestMarshalToMap_CustomMarshal(t *testing.T) {
	h := test4{0}
	v := structs.MarshalToMap(h)
	validMap(t, v, map[string]any{"a": 5})
}
