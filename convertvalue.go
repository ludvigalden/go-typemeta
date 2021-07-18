package typemeta

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ConvertValue converts a value to a specified type and returns a detailed error if it fails
func ConvertValue(value reflect.Value, toType interface{}) (reflect.Value, error) {
	return DefaultSchema.ConvertValue(value, toType)
}

// ConvertInterfaceValue converts a value to a specified type and returns a detailed error if it fails
func ConvertInterfaceValue(value interface{}, toType interface{}) (interface{}, error) {
	return DefaultSchema.ConvertInterfaceValue(value, toType)
}

// ConvertValue converts a value to a specified type and returns an error if it fails
func (s *Schema) ConvertValue(value reflect.Value, toType interface{}) (reflect.Value, error) {
	if !value.IsValid() {
		return value, errors.New("received invalid value")
	} else if value.Type() == nil {
		return value, errors.New("received value with nil type")
	}
	return convertValue(s, value, s.Get(value.Type()), s.Get(toType))
}

// ConvertInterfaceValue converts a value to a specified type and returns a detailed error if it fails
func (s *Schema) ConvertInterfaceValue(value interface{}, toType interface{}) (interface{}, error) {
	var rv reflect.Value
	if vrv, ok := value.(reflect.Value); ok {
		rv = vrv
	} else {
		rv = reflect.ValueOf(value)
	}
	cv, err := s.ConvertValue(rv, toType)
	if err != nil {
		return nil, err
	}
	return cv.Interface(), nil
}

func convertValue(s *Schema, value reflect.Value, valueTypeMeta TypeMeta, toTypeMeta TypeMeta) (reflect.Value, error) {
	if !value.IsValid() {
		return value, errors.New("received invalid value")
	}
	if value.Kind() == reflect.Interface {
		if reflect.TypeOf(value.Interface()) == nil {
			if value.IsZero() {
				return reflect.New(toTypeMeta.Type()).Elem(), nil
			}
			// return zero value
			return value, errors.New("received interface value with nil type")
		}
		value = reflect.ValueOf(value.Interface())
		valueTypeMeta = s.get(value.Type(), false)
	}
	toType := toTypeMeta.Type()
	if value.Type() == toType {
		return value, nil
	}
	if toType.Kind() == reflect.Interface {
		if !value.Type().Implements(toType) {
			return value, errors.New(value.Type().String() + " does not implement " + toType.String())
		}
		return value, nil
	}
	if toType.Kind() == reflect.Ptr {
		if value.Kind() != reflect.Ptr {
			// parse a pointer to the value (if not equal, keep parsing)
			newValue := reflect.New(value.Type())
			newValue.Elem().Set(value)
			return convertValue(s, newValue, s.Get(newValue.Type()), toTypeMeta)
		}
		// the value is a pointer, too, but on different levels (e.g. ***string vs. **string or *string vs. **string)
		// here, we check if the inner type is equal, and in that case returns it
		valueNonPtrType := value.Type()
		valuePtrLevel := 0
		for valueNonPtrType.Kind() == reflect.Ptr {
			valueNonPtrType = valueNonPtrType.Elem()
			valuePtrLevel++
		}
		nonPtrToType := toType
		toPtrLevel := 0
		for nonPtrToType.Kind() == reflect.Ptr {
			nonPtrToType = nonPtrToType.Elem()
			toPtrLevel++
		}
		if nonPtrToType == valueNonPtrType {
			newValue := value
			if valuePtrLevel > toPtrLevel {
				// return element of value
				for i := 0; i < valuePtrLevel-toPtrLevel; i++ {
					if newValue.IsNil() {
						// return zero value
						return reflect.New(toType).Elem(), nil
					}
					newValue = newValue.Elem()
				}
			} else {
				// return pointer to value
				for i := 0; i < toPtrLevel-valuePtrLevel; i++ {
					nextNewValue := reflect.New(newValue.Type())
					nextNewValue.Elem().Set(newValue)
					newValue = nextNewValue
				}
			}
			return newValue, nil
		}
		// unequal inner type, convert inner value
	}
	nonPtrValue := value
	for nonPtrValue.Kind() == reflect.Ptr {
		if nonPtrValue.IsNil() {
			// return zero value
			return reflect.New(toType).Elem(), nil
		}
		nonPtrValue = nonPtrValue.Elem()
		if nonPtrValue.Type() == toType {
			return nonPtrValue, nil
		}
	}
	if nonPtrValue.Type() == toType {
		return nonPtrValue, nil
	}
	var err error
	nonPtrValue, err = convertNonPtrValue(s, nonPtrValue, NonPtr(valueTypeMeta), NonPtr(toTypeMeta))
	if err == nil {
		if nonPtrValue.Type() == toTypeMeta.Type() || NonPtr(toTypeMeta).Type().Kind() == reflect.Interface {
			return nonPtrValue, nil
		}
		return convertValue(s, nonPtrValue, s.Get(nonPtrValue.Type()), toTypeMeta)
	}
	return value, err
}

func convertNonPtrValue(s *Schema, value reflect.Value, valueTypeMeta TypeMeta, toTypeMeta TypeMeta) (reflect.Value, error) {
	toType := toTypeMeta.Type()
	newValue := reflect.New(toType).Elem()
	if toType.Kind() == reflect.String {
		return convertValueToString(value, valueTypeMeta, toTypeMeta)
	}
	if valueTypeMeta, ok := valueTypeMeta.(*Slice); ok {
		valueLen := value.Len()
		if value.IsNil() || valueLen == 0 { // empty slice
			return newValue, nil
		}
		switch toTypeMeta := toTypeMeta.(type) {
		case *Slice: // slice to slice
			newValue = reflect.MakeSlice(toType, valueLen, valueLen)
			for i := 0; i < valueLen; i++ {
				newElem, err := convertValue(s, value.Index(i), valueTypeMeta.Elem, toTypeMeta.Elem)
				if err != nil {
					return value, err
				}
				newValue.Index(i).Set(newElem)
			}
			return newValue, nil
		default:
			if valueLen == 1 {
				// try to convert first value
				newValue, err := ConvertValue(value.Index(0), toTypeMeta)
				if err == nil {
					return newValue, nil
				}
			} else if valueTypeMeta.Elem.Kind() == reflect.String {
				// join slice of strings
				strs := []string{}
				for i := 0; i < valueLen; i++ {
					strs = append(strs, value.Index(i).String())
				}
				newValue, err := ConvertValue(reflect.ValueOf(strings.Join(strs, ", ")), toTypeMeta)
				if err == nil {
					return newValue, nil
				}
			}
		}
		return value, notAssignibleError(valueTypeMeta, toTypeMeta)
	}
	switch toTypeMeta := toTypeMeta.(type) {
	case *Struct:
		if toTypeMeta.Primitive() {
			return value, errors.New("converting primitive structs not implemented")
			// if valueTypeMeta.Primitive() {
			// 	// primitive struct to primitive struct
			// } else {
			// 	// struct to primitive struct
			// 	switch valueTypeMeta := valueTypeMeta.(type) {

			// 	}
			// }
		}
		switch valueTypeMeta := valueTypeMeta.(type) {
		case *Map: // map to struct
			mapIter := value.MapRange()
			for mapIter.Next() {
				key := mapIter.Key()
				structField := toTypeMeta.FieldByName(key.String())
				if structField == nil {
					return value, errors.New("unrecognized key \"" + fmt.Sprint(key.Interface()) + "\" does not exist in struct \"" + toTypeMeta.String() + "\"")
				}
				keyValue := mapIter.Value()
				convertedValue, err := convertValue(s, keyValue, valueTypeMeta.Elem, structField.TypeMeta)
				if err != nil {
					return value, errors.New("could not convert value of key \"" + fmt.Sprint(key.Interface()) + "\" to field \"" + structField.String() + "\". " + err.Error())
				}
				fieldValue := newValue.Field(structField.Index)
				fieldValue.Set(convertedValue)
			}
			return newValue, nil
		case *Struct: // struct to struct
			if valueTypeMeta.Primitive() {
				return value, errors.New("converting primitive structs not implemented")
			}
		default:
			return value, notAssignibleError(valueTypeMeta, toTypeMeta)
		}
	case *Interface: // any type of value
		return value, nil
	case *Primitive:
		switch valueTypeMeta := valueTypeMeta.(type) {
		case *Primitive: // primitive to primitive
			return convertPrimitiveValue(value, valueTypeMeta, toTypeMeta)
		default:
			return value, notAssignibleError(valueTypeMeta, toTypeMeta)
		}
	case *Slice:
		switch valueTypeMeta := valueTypeMeta.(type) {
		case *Array: // !TODO array to slice
		case *Slice:
			newValue = reflect.MakeSlice(toType, value.Len(), value.Cap())
			for i := 0; i < value.Len(); i++ {
				convertedElem, err := convertValue(s, value.Index(i), valueTypeMeta.Elem, toTypeMeta.Elem)
				if err != nil {
					return value, err
				}
				newValue.Index(i).Set(convertedElem)
			}
			return newValue, nil
		case *Primitive:
			if valueTypeMeta.Kind() == reflect.String {
				// attempt to unmarshal the value
				newValuePtr := reflect.New(toType)
				err := json.Unmarshal(stringToBytes(value.String()), newValuePtr.Interface())
				if err == nil {
					return newValuePtr.Elem(), nil
				}
			}
		default: // any element to slice
			newValue = reflect.MakeSlice(toType, 1, 1)
			newElemValue, err := convertValue(s, value, valueTypeMeta, toTypeMeta.Elem)
			if err != nil {
				return value, err
			}
			newValue.Index(0).Set(newElemValue)
			return newValue, nil
		}
		return newValue, nil
	case *Map:
		switch valueTypeMeta := valueTypeMeta.(type) {
		case *Map: // map to struct
			convertKey := valueTypeMeta.Key.Type() != toTypeMeta.Key.Type()
			convertElem := valueTypeMeta.Elem.Type() != toTypeMeta.Elem.Type()
			if !convertKey {
				if !convertElem {
					newValue = value.Convert(toType)
				} else {
					newValue = reflect.MakeMap(toType)
					mapIter := value.MapRange()
					for mapIter.Next() {
						key := mapIter.Key()
						keyValue := mapIter.Value()
						convertedKeyValue, err := convertValue(s, keyValue, valueTypeMeta.Elem, toTypeMeta.Elem)
						if err != nil {
							return value, errors.New("could not convert value of key \"" + fmt.Sprint(key.Interface()) + "\" to \"" + toTypeMeta.Elem.String() + "\". " + err.Error())
						}
						newValue.SetMapIndex(key, convertedKeyValue)
					}
				}
			} else {
				newValue = reflect.MakeMap(toType)
				mapIter := value.MapRange()
				for mapIter.Next() {
					key := mapIter.Key()
					keyValue := mapIter.Value()
					convertedKey, err := convertValue(s, key, valueTypeMeta.Key, toTypeMeta.Key)
					if err != nil {
						return value, errors.New("could not convert key \"" + fmt.Sprint(key.Interface()) + "\" to \"" + toTypeMeta.Key.String() + "\". " + err.Error())
					}
					convertedKeyValue, err := convertValue(s, keyValue, valueTypeMeta.Elem, toTypeMeta.Elem)
					if err != nil {
						return value, errors.New("could not convert value of key \"" + fmt.Sprint(key.Interface()) + "\" to \"" + toTypeMeta.Elem.String() + "\". " + err.Error())
					}
					newValue.SetMapIndex(convertedKey, convertedKeyValue)
				}
			}
			return newValue, nil
		default:
			return value, notAssignibleError(valueTypeMeta, toTypeMeta)
		}
	case *Ptr:
		return value, errors.New("expected non-ptr")
	default:
		return value, notAssignibleError(valueTypeMeta, toTypeMeta)
	}
	return newValue, nil
}

func notAssignibleError(valueTypeMeta TypeMeta, toTypeMeta TypeMeta) error {
	return errors.New(valueTypeMeta.String() + " not assignable to " + toTypeMeta.String())
}

func valueNotAssignibleError(value reflect.Value, toTypeMeta TypeMeta) error {
	return errors.New(value.String() + " not assignable to " + toTypeMeta.String())
}
