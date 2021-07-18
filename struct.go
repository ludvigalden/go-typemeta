package typemeta

import (
	"fmt"
	"reflect"
	"strconv"
)

// Struct is type meta for a struct
type Struct struct {
	name         string
	Description  string
	StringParser func(string) (interface{}, error)
	Fields       map[int]StructField

	typ reflect.Type
}

// Field returns the field at the specified index. If no matching field is found, nil is returned.
func (s *Struct) Field(fieldIndex int) *StructField {
	field, ok := s.Fields[fieldIndex]
	if !ok {
		return nil
	}
	return &field
}

// FieldByName returns the field with the specified name. The specified name can either be the struct field name or the JSON field name.
// If no matching field is found, nil is returned.
func (s *Struct) FieldByName(fieldName string) *StructField {
	if fieldName == "" {
		return nil
	}
	field := s.FindField(func(field StructField) bool {
		return field.Name == fieldName || field.JSONName == fieldName
	})
	return field
}

// FieldIndexByName returns the field index of a field with the specified name. The specified name can either
// be the struct field name or the JSON field name. If no matching field is found, -1 is returned.
func (s *Struct) FieldIndexByName(fieldName string) int {
	field := s.FieldByName(fieldName)
	if field == nil {
		return -1
	}
	return field.Index
}

// EnsureFieldByName returns the field with the specified name. The specified name can either be the struct field name or the JSON field name.
// If no matching field is found, EnsureField panics.
func (s *Struct) EnsureFieldByName(fieldName string) StructField {
	f := s.FieldByName(fieldName)
	if f == nil {
		panic("Field with name \"" + fieldName + "\" does not exist in " + s.String())
	}
	return *f
}

// EnsureField returns the field with the specified name. The specified name can either be the struct field name or the JSON field name.
// If no matching field is found, EnsureField panics.
func (s *Struct) EnsureField(fieldIndex int) StructField {
	field := s.Field(fieldIndex)
	if field == nil {
		panic("Field with index [" + strconv.Itoa(fieldIndex) + "] does not exist in " + s.String())
	}
	return *field
}

// FindField iterates through all fields and returns the first field true is returned for
func (s *Struct) FindField(iteratee func(StructField) bool) *StructField {
	for i := 0; i < s.typ.NumField(); i++ {
		if field := s.Fields[i]; iteratee(field) {
			return &field
		}
	}
	return nil
}

// IterateFields iterates through all of the fields of the struct
func (s *Struct) IterateFields(iteratee func(StructField)) {
	for i := 0; i < s.typ.NumField(); i++ {
		iteratee(s.Fields[i])
	}
}

// SetName sets the name of the struct type meta
func (s *Struct) SetName(name string) *Struct {
	s.name = name
	return s
}

// SetDescription sets the description of the struct type meta
func (s *Struct) SetDescription(description string) *Struct {
	s.Description = description
	return s
}

// SetStringParser sets the string parser of the struct type meta
func (s *Struct) SetStringParser(stringParser func(string) (interface{}, error)) *Struct {
	s.StringParser = stringParser
	return s
}

// SetField sets the type meta of the specifield field
func (s *Struct) SetField(fieldName string, fieldType TypeMeta) *Struct {
	field := s.EnsureFieldByName(fieldName)
	field.TypeMeta = fieldType
	s.Fields[field.Index] = field
	return s
}

// SetFieldElem sets the elem type meta of the specifield field
func (s *Struct) SetFieldElem(fieldName string, fieldElemType TypeMeta) *Struct {
	field := s.EnsureFieldByName(fieldName)
	field.TypeMeta = setElem(field.TypeMeta, fieldElemType)
	s.Fields[field.Index] = field
	return s
}

func setElem(typeMeta TypeMeta, elemTypeMeta TypeMeta) TypeMeta {
	switch fieldType := typeMeta.(type) {
	case *Ptr:
		fieldType = fieldType.Copy()
		return setElem(fieldType, elemTypeMeta)
	case *Slice:
		fieldType = fieldType.Copy()
		fieldType.Elem = elemTypeMeta
		return fieldType
	default:
		panic("cannot to set elem" + fmt.Sprint(elemTypeMeta) + " to " + fmt.Sprint(typeMeta))
	}
}

// Copy returns a copy of the struct type meta
func (s *Struct) Copy() *Struct {
	cs := *s
	if s.Fields != nil {
		cs.Fields = map[int]StructField{}
		for index, field := range s.Fields {
			cs.Fields[index] = field
		}
	}
	return &cs
}

// Type returns
func (s *Struct) Type() reflect.Type {
	return s.typ
}

// Kind returns
func (s *Struct) Kind() reflect.Kind {
	return reflect.Struct
}

// JSONNonNull returns
func (s *Struct) JSONNonNull() bool {
	return true
}

// Primitive returns false for structs
func (s *Struct) Primitive() bool {
	return false
}

// Name returns the type meta's explicitly set name or the type's name within its package for a defined type. For other (non-defined) types it returns the empty string.
func (s *Struct) Name() string {
	if s.name != "" {
		return s.name
	}
	return s.typ.Name()
}

// String returns
func (s *Struct) String() string {
	return s.typ.String()
}
