package typemeta

import (
	"reflect"
	"testing"
)

func TestConvertValue(t *testing.T) {
	type StructA struct {
		Name   string
		Active bool
	}
	t.Run("slice to elem", func(t *testing.T) {
		cd := []interface{}{
			[]StructA{{"Test", true}},
			[]map[string]interface{}{{"name": "Test"}},
			[]interface{}{1},
			[]string{"hello", "there"},
			[]bool{true},
		}
		for _, cd := range cd {
			elem := Get(cd).(*Slice).Elem
			_, err := ConvertValue(reflect.ValueOf(cd), elem)
			if err != nil {
				t.Error("failed converting slice to elem: " + err.Error())
			}
		}
	})
	type MapType map[string]interface{}
	t.Run("map to map type", func(t *testing.T) {
		cd := []interface{}{
			map[string]interface{}{"Hello": 123},
			map[int]interface{}{123: "456"},
		}
		for _, cd := range cd {
			_, err := ConvertValue(reflect.ValueOf(cd), MapType{})
			if err != nil {
				t.Error("failed converting map type: " + err.Error())
			}
		}
	})
}
