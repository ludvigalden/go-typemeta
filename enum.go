package typemeta

import (
	"fmt"
	"reflect"
	"strings"
)

// Enum is type meta for an enum
type Enum struct {
	name     string
	typ      reflect.Type
	values   map[string]interface{}
	firstKey string
}

// NewEnum creates a new enum with the specified name. The value can be an array of values, which will be converted to a map[string]string (using fmt.Sprint if the values are not strings),
// or a map of values such as map[string]interface{}. The value can also be nil, and in that case values should obviously be set later using the `Enum.SetValue` method.
func NewEnum(name string, v interface{}) *Enum {
	var typ reflect.Type
	values := map[string]interface{}{}
	var firstKey string
	if strarr, ok := v.([]string); ok {
		for _, stri := range strarr {
			values[stri] = stri
		}
	} else if iarr, ok := v.([]interface{}); ok {
		for index, iv := range iarr {
			strv := fmt.Sprint(iv)
			if index == 0 {
				firstKey = strv
			}
			values[strv] = strv
			if typ == nil {
				typ = reflect.TypeOf(typ)
			}
		}
	} else if imap, ok := v.(map[string]interface{}); ok {
		values = imap
	} else {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < rv.Len(); i++ {
				iv := rv.Index(i)
				if iv.Kind() == reflect.Ptr {
					iv = iv.Elem()
				}
				if typ == nil {
					typ = iv.Type()
				}
				str := iv.String()
				if i == 0 {
					firstKey = str
				}
				values[str] = str
			}
		default:
			panic("Unable to create enum " + name + " from value " + fmt.Sprint(v))
		}
	}
	if typ == nil {
		typ = reflect.TypeOf(string(""))
	}
	return &Enum{name, typ, values, firstKey}
}

// SetValue sets a value of the enum
func (e Enum) SetValue(name string, value interface{}) Enum {
	if e.values == nil {
		e.values = map[string]interface{}{name: value}
	} else {
		e.values[name] = value
	}
	return e
}

// IterateValues iterates the values of the enum
func (e Enum) IterateValues(iteratee func(string, interface{})) {
	if e.values == nil {
		return
	}
	for name, value := range e.values {
		iteratee(name, value)
	}
}

// FirstEntry returns the first entry of the enum. If the enum was specified as a slice or array, the first
// entry is guaranteed to be for first item in the slice or array. If it was specified as a map, the first entry is random (lol).
func (e Enum) FirstEntry() (string, interface{}) {
	if e.values == nil {
		return "", nil
	}
	if e.firstKey != "" {
		return e.firstKey, e.values[e.firstKey]
	}
	for name, value := range e.values {
		return name, value
	}
	return "", nil
}

// Parse iterates through the values of the enum and whether it was found
func (e Enum) Parse(value interface{}) (interface{}, bool) {
	if e.values == nil {
		return nil, false
	}
	// reflectValue := reflect.ValueOf(value)
	// if reflectValue.Kind() == reflect.Ptr {
	// 	reflectValue = reflectValue.Elem()
	// }
	// if reflectValue.Kind() == reflect.String {
	// 	stringValue := reflectValue.String()
	// 	for name, enumValue := range e.values {
	// 		if name == stringValue {

	// if !value.IsValid() {
	// 	return value, errors.New("typemeta: Received invalid value to parse")
	// } else if value.Type() == s.typ {
	// 	return value, nil
	// }, true
	// 		} else if enumReflectValue := reflect.ValueOf(enumValue); enumReflectValue.Kind() == reflect.String && enumReflectValue.String() == stringValue {

	// if !value.IsValid() {
	// 	return value, errors.New("typemeta: Received invalid value to parse")
	// } else if value.Type() == s.typ {
	// 	return value, nil
	// }, true
	// 		}
	// 	}
	// } else {
	// 	for _, enumValue := range e.values {
	// 		if enumReflectValue := reflect.ValueOf(enumValue); enumReflectValue.Kind() == reflectValue.Kind() {
	// 			switch enumReflectValue.Kind() {
	// 			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
	// 				if reflectValue.Int() == enumReflectValue.Int() {

	// if !value.IsValid() {
	// 	return value, errors.New("typemeta: Received invalid value to parse")
	// } else if value.Type() == s.typ {
	// 	return value, nil
	// }, true
	// 				}
	// 				return nil, false
	// 			default:
	// 				println("Warning: Unable to compare values " + reflectValue.String() + " and " + enumReflectValue.String() + " for enum " + e.name)
	// 				return nil, false
	// 			}
	// 		}
	// 	}
	// }
	return nil, false
}

// Has returns whether the enum matches a value
func (e Enum) Has(value interface{}) bool {
	_, ok := e.Parse(value)
	return ok
}

func (e Enum) String() string {
	return "Enum(" + e.innerString() + ")"
}

func (e Enum) innerString() string {
	strParts := []string{}
	if e.name != "" {
		strParts = append(strParts, e.name)
	}
	if e.values != nil && len(e.values) > 0 {
		valuesStrParts := []string{}
		for key, value := range e.values {
			if key == value {
				valuesStrParts = append(valuesStrParts, key)
			} else {
				valuesStrParts = append(valuesStrParts, key+": "+fmt.Sprint(value))

			}
		}
		strParts = append(strParts, strings.Join(valuesStrParts, ", "))
	}
	return strings.Join(strParts, ", ")
}
