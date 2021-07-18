package typemeta

import (
	"encoding/json"
	"reflect"
)

// UnmarshalValue unmarshals data, sets default values, and returns an error if unsuccessful.
// It is much, much slower than `json.Unmarshal`. Not sure why you would even use this tbh.
func UnmarshalValue(t interface{}, data []byte) (interface{}, error) {
	rv, err := unmarshalValue(t, data)
	if err != nil {
		return rv.Interface(), err
	}
	return rv.Interface(), nil
}

func unmarshalValue(t interface{}, data []byte) (reflect.Value, error) {
	tm := Get(t)
	nonPtrKind := NonPtr(tm).Kind()
	switch nonPtrKind {
	case reflect.Array, reflect.Slice:
		rv := reflect.New(reflect.TypeOf([]interface{}{}))
		err := json.Unmarshal(data, rv.Interface())
		if err != nil {
			return rv.Elem(), err
		}
		rv, err = ConvertValue(rv.Elem(), tm)
		if err != nil {
			return rv, err
		}
		return rv, nil
	case reflect.Map:
		rv := reflect.New(reflect.TypeOf(map[string]interface{}{}))
		err := json.Unmarshal(data, rv.Interface())
		if err != nil {
			return rv.Elem(), err
		}
		rv, err = ConvertValue(rv.Elem(), tm)
		if err != nil {
			return rv, err
		}
		return rv, nil
	case reflect.Struct:
		if tm.Primitive() {
			break
		}
		rv := reflect.New(reflect.TypeOf(map[string]interface{}{}))
		err := json.Unmarshal(data, rv.Interface())
		if err != nil {
			return rv.Elem(), err
		}
		rv, err = ConvertValue(rv.Elem(), tm)
		if err != nil {
			return rv, err
		}
		return rv, nil
	}
	rv := reflect.New(tm.Type())
	err := json.Unmarshal(data, rv.Interface())
	if err != nil {
		return rv.Elem(), err
	}
	return rv.Elem(), nil
}
