package typemeta

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// NonPtr returns the non-pointer type meta
func NonPtr(t TypeMeta) TypeMeta {
	switch t := t.(type) {
	case *Ptr:
		return NonPtr(t.Elem)
	default:
		return t
	}
}

// AssertSlice returns the non-pointer type meta or nil if not found
func AssertSlice(t TypeMeta) *Slice {
	t = NonPtr(t)
	switch t := t.(type) {
	case *Slice:
		return t
	default:
		return nil
	}
}

// AssertArray returns the non-pointer or nil if not found
func AssertArray(t TypeMeta) *Array {
	t = NonPtr(t)
	switch t := t.(type) {
	case *Array:
		return t
	default:
		return nil
	}
}

// MapOf returns the non-pointer or nil if not found
func MapOf(t TypeMeta) *Map {
	t = NonPtr(t)
	switch t := t.(type) {
	case *Map:
		return t
	default:
		return nil
	}
}

// StructOf attempts to get related struct type meta of the specified type meta.
// If the passed type meta is a `*typemeta.Struct`, that is returned.
// If the passed type meta is for a slice, array, map, or pointer, it attempts to get the struct of the element.
// If no struct type meta is found, nil is returned.
func StructOf(t TypeMeta) *Struct {
	switch t := t.(type) {
	case *Struct:
		return t
	case *Ptr:
		return StructOf(t.Elem)
	case *Slice:
		return StructOf(t.Elem)
	case *Array:
		return StructOf(t.Elem)
	case *Map:
		return StructOf(t.Elem)
	default:
		return nil
	}
}

// StructOrPrimitiveOf attempts to get struct or primitive type meta of the specified type meta.
// If the passed type meta is for a slice, array, map, or pointer, it attempts to get the struct or primtiive type meta of the element.
func StructOrPrimitiveOf(t TypeMeta) TypeMeta {
	switch t := t.(type) {
	case *Struct, *Primitive, *Interface:
		return t
	case *Ptr:
		if t.Elem == nil {
			panic("PTR WITHOUT ELEM " + fmt.Sprint(t.Elem) + " _ " + fmt.Sprint(t))
		}
		return StructOrPrimitiveOf(t.Elem)
	case *Slice:
		return StructOrPrimitiveOf(t.Elem)
	case *Array:
		return StructOrPrimitiveOf(t.Elem)
	case *Map:
		return StructOrPrimitiveOf(t.Elem)
	default:
		return t
	}
}

// InterfaceOf attempts to get related interface type meta of the specified type meta.
// If the passed type meta is for a slice, array, map, or pointer, it attempts to get the interface of the element.
// If no interface type meta is found, nil is returned.
func InterfaceOf(t TypeMeta) *Interface {
	switch t := t.(type) {
	case *Interface:
		return t
	case *Ptr:
		return InterfaceOf(t.Elem)
	case *Slice:
		return InterfaceOf(t.Elem)
	case *Array:
		return InterfaceOf(t.Elem)
	case *Map:
		return InterfaceOf(t.Elem)
	default:
		return nil
	}
}

// ElemOf returns the elem of the type meta.
// If the passed type meta is a `*typemeta.Ptr`, the elem of that is returned.
// If the passed type meat is a slice, array, or map, it returns the elem of those.
// In all other cases, nil is returned.
func ElemOf(t TypeMeta) TypeMeta {
	switch t := t.(type) {
	case *Ptr:
		return ElemOf(t.Elem)
	case *Slice:
		return t.Elem
	case *Array:
		return t.Elem
	case *Map:
		return t.Elem
	}
	return nil
}

// IsArray returns whether the passed type meta is for an array (may be a pointer to an array)
func IsArray(t TypeMeta) bool {
	return AssertArray(t) != nil
}

// IsSlice returns whether the passed type meta is for a slice (may be a pointer to a slice)
func IsSlice(t TypeMeta) bool {
	return AssertSlice(t) != nil
}

// IsArrayOrSlice returns whether the passed type meta is for a slice or an array (may be a pointer to a slice or an array)
func IsArrayOrSlice(t TypeMeta) bool {
	return IsSlice(t) || IsArray(t)
}

// FuncName returns the name of a function
func FuncName(fn interface{}) string {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		return ""
	}
	funcName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	funcNameParts := strings.Split(funcName, "/")
	funcName = strings.Join(funcNameParts[len(funcNameParts)-1:], "")
	return funcName
}

// EnsureSlice returns slice type meta for the specified type. It attempts to use `AssertSlice`
// and if nil uses `SliceOf` to ensure a non-nil returned slice type.
func EnsureSlice(t TypeMeta) *Slice {
	s := AssertSlice(t)
	if s == nil {
		s = SliceOf(t)
	}
	return s
}

// StructFieldOf returns a struct field at the specified path
func StructFieldAt(t TypeMeta, fieldPath []int) *StructField {
	fieldPathLen := len(fieldPath)
	if fieldPathLen == 0 {
		return nil
	}
	return structFieldAt(StructOf(t), fieldPath, fieldPathLen)
}

// EnsureStructFieldAt returns a struct field at the specified path
func EnsureStructFieldAt(t TypeMeta, fieldPath []int) StructField {
	fieldPathLen := len(fieldPath)
	if fieldPathLen == 0 {
		panic("No field path passed")
	}
	sf := structFieldAt(StructOf(t), fieldPath, fieldPathLen)
	if sf == nil {
		fieldPathStrs := []string{}
		for _, fieldIndex := range fieldPath {
			fieldPathStrs = append(fieldPathStrs, fmt.Sprint(fieldIndex))
		}
		panic("No struct field of " + t.String() + " found at path [" + strings.Join(fieldPathStrs, "][") + "]")
	}
	return *sf
}

func structFieldAt(st *Struct, fieldPath []int, fieldPathLen int) *StructField {
	if st == nil {
		return nil
	}
	field := st.Field(fieldPath[0])
	if fieldPathLen == 1 {
		return field
	}
	return structFieldAt(StructOf(field.TypeMeta), fieldPath[1:], fieldPathLen-1)
}
