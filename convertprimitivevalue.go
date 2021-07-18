package typemeta

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

func convertPrimitiveValue(value reflect.Value, valueTypeMeta *Primitive, toTypeMeta *Primitive) (reflect.Value, error) {
	valueKind := valueTypeMeta.Kind()
	toType := toTypeMeta.Type()
	toKind := toType.Kind()
	if toKind == reflect.String {
		return convertValueToString(value, valueTypeMeta, toTypeMeta)
	}
	if value.IsZero() {
		// a zero primitive will always correspond with the zero value of the target type
		return reflect.New(toTypeMeta.typ), nil
	}
	// !TODO: convert all primitive values not convertible by reflect package and use custom parsers
	if valueKind == toKind {
		return value.Convert(toTypeMeta.Type()), nil
	}
	switch toKind {
	case reflect.Bool:
		switch valueKind {
		case reflect.String:
			bool, err := strconv.ParseBool(value.String())
			if err == nil {
				return reflect.ValueOf(bool).Convert(toType), nil
			}
			floatValue, err := convertPrimitiveValue(value, valueTypeMeta, Get(reflect.Float64).(*Primitive))
			if err != nil {
				return value, valueNotAssignibleError(value, toTypeMeta)
			}
			return reflect.ValueOf(floatValue.Float() > 0).Convert(toType), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return reflect.ValueOf(value.Int() > 0).Convert(toType), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return reflect.ValueOf(value.Uint() > 0).Convert(toType), nil
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(value.Float() > 0).Convert(toType), nil
		}
	case reflect.Struct:
		switch valueKind {
		case reflect.String:
			return convertString(value.String(), toTypeMeta)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch valueKind {
		case reflect.String:
			int, err := strconv.ParseInt(value.String(), 10, 64)
			if err != nil {
				return value, valueNotAssignibleError(value, toTypeMeta)
			}
			return reflect.ValueOf(int).Convert(toType), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64:
			return value.Convert(toType), nil
		case reflect.Bool:
			if value.Bool() {
				return reflect.ValueOf(1).Convert(toType), nil
			}
			return reflect.ValueOf(0).Convert(toType), nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch valueKind {
		case reflect.String:
			uint, err := strconv.ParseUint(value.String(), 10, 64)
			if err != nil {
				return value, valueNotAssignibleError(value, toTypeMeta)
			}
			return reflect.ValueOf(uint).Convert(toType), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64:
			return value.Convert(toType), nil
		case reflect.Bool:
			if value.Bool() {
				return reflect.ValueOf(1).Convert(toType), nil
			}
			return reflect.ValueOf(0).Convert(toType), nil
		}
	case reflect.Float32, reflect.Float64:
		switch valueKind {
		case reflect.String:
			float, err := strconv.ParseFloat(value.String(), 64)
			if err != nil {
				return value, valueNotAssignibleError(value, toTypeMeta)
			}
			return reflect.ValueOf(float).Convert(toType), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return value.Convert(toType), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return value.Convert(toType), nil
		case reflect.Float32, reflect.Float64:
			return value.Convert(toType), nil
		case reflect.Bool:
			if value.Bool() {
				return reflect.ValueOf(1).Convert(toType), nil
			}
			return reflect.ValueOf(0).Convert(toType), nil
		}
	}
	return value, notAssignibleError(valueTypeMeta, toTypeMeta)
}

func convertValueToString(value reflect.Value, valueTypeMeta TypeMeta, toTypeMeta TypeMeta) (reflect.Value, error) {
	marshaler := marshalerOf(value)
	if marshaler != nil {
		b, err := marshaler()
		if err != nil {
			return value, err
		}
		return reflect.ValueOf(bytesToString(b)).Convert(toTypeMeta.Type()), nil
	}
	return reflect.ValueOf(fmt.Sprint(value.Interface())).Convert(toTypeMeta.Type()), nil
}

func convertString(str string, toTypeMeta TypeMeta) (reflect.Value, error) {
	value := reflect.New(toTypeMeta.Type())
	unmarshaler := unmarshalerOf(value)
	if unmarshaler != nil {
		err := unmarshaler(stringToBytes(str))
		if err != nil {
			return value, err
		}
		return value, nil
	}
	return value, nil
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func marshalerOf(value reflect.Value) func() ([]byte, error) {
	marshalerValue := value.MethodByName("MarshalText")
	if !marshalerValue.IsValid() {
		marshalerValue = value.MethodByName("MarshalJSON")
	}
	if !marshalerValue.IsValid() {
		marshalerValue = value.MethodByName("MarshalBinary")
	}
	if marshalerValue.IsValid() {
		unmarshaler, ok := marshalerValue.Interface().(func() ([]byte, error))
		if ok {
			return unmarshaler
		}
	}
	return nil
}
func unmarshalerOf(value reflect.Value) func([]byte) error {
	unmarshalerValue := value.MethodByName("UnmarshalText")
	if !unmarshalerValue.IsValid() {
		unmarshalerValue = value.MethodByName("UnmarshalJSON")
	}
	if !unmarshalerValue.IsValid() {
		unmarshalerValue = value.MethodByName("UnmarshalBinary")
	}
	if unmarshalerValue.IsValid() {
		unmarshaler, ok := unmarshalerValue.Interface().(func([]byte) error)
		if ok {
			return unmarshaler
		}
	}
	return nil
}
