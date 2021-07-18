package typemeta

import (
	"encoding/json"
	"reflect"
	"sync"
	"unicode"

	"github.com/fatih/structtag"
)

// Schema is a schema of type meta
type Schema struct {
	mu        sync.Mutex
	types     map[reflect.Type]TypeMeta
	enumTypes map[*Enum]TypeMeta
}

// Get returns type meta for the specified type. If type meta is passed, that is returned. If a reflect type is passed,
// the type meta for that type is returned. If a reflect kind is passed, the type meta for the default type for that kind is returned.
// In all other cases, the reflect type is retrieved of the passed argument using `reflect.TypeOf` and the type meta is retrieved from that type.
// If no type can be inferred, i.e. `reflect.TypeOf` returns nil, a *typemeta.Interface with a nil type is returned.
// Each inferred reflect type points to a single type meta object for a single schema.
func (s *Schema) Get(typ interface{}) TypeMeta {
	switch typ := typ.(type) {
	case TypeMeta:
		if typ == nil {
			return s.get(nil, false)
		}
		return typ
	case reflect.Type:
		return s.get(typ, false)
	case reflect.Value:
		return s.get(typ.Type(), false)
	case reflect.Kind:
		if kindType, ok := kindTypes[typ]; ok {
			return s.get(kindType, false)
		}
		panic("No default type has been defined for kind " + typ.String())
	case *Enum:
		if s.enumTypes[typ] != nil {
			return s.enumTypes[typ]
		}
		r, ok := s.get(typ.typ, false).(*Primitive)
		if !ok {
			r = &Primitive{typ: typ.typ, enum: typ}
		} else {
			r = r.Copy()
			r.enum = typ
		}
		s.enumTypes[typ] = r
		return r
	default:
		return s.get(reflect.TypeOf(typ), false)
	}
}

// SliceOf returns slie type meta for the specified type
func (s *Schema) SliceOf(typ interface{}) *Slice {
	return s.Get(reflect.SliceOf(s.Get(typ).Type())).(*Slice)
}

// PtrTo returns a pointer to the specified type
func (s *Schema) PtrTo(typ interface{}) *Ptr {
	elem := s.Get(typ)
	ptr := s.Get(reflect.PtrTo(elem.Type())).(*Ptr)
	if ptr.Elem != elem {
		ptr = ptr.Copy()
		ptr.Elem = elem
	}
	return ptr
}

func (s *Schema) get(rtyp reflect.Type, locked bool) TypeMeta {
	s.lock(locked)
	meta := s.types[rtyp]
	if meta != nil {
		s.unlock(locked)
		return meta
	}
	if rtyp == nil {
		i := &Interface{}
		s.types[rtyp] = i
		s.unlock(locked)
		return i
	}
	switch rtyp.Kind() {
	case reflect.Ptr:
		ptr := &Ptr{typ: rtyp}
		s.types[rtyp] = ptr
		ptr.Elem = s.get(rtyp.Elem(), true)
		s.unlock(locked)
		return ptr
	case reflect.Struct:
		if rtyp.Implements(primitiveType) {
			// primitive struct
			p := &Primitive{rtyp, nil, ""}
			s.types[rtyp] = p
			s.unlock(locked)
			return p
		}
		strct := &Struct{Fields: map[int]StructField{}, typ: rtyp}
		s.types[rtyp] = strct
		for fieldIndex := 0; fieldIndex < rtyp.NumField(); fieldIndex++ {
			rsf := rtyp.Field(fieldIndex)
			nameFirstChar := []rune(rsf.Name)[0]
			private := unicode.IsLower(nameFirstChar) || nameFirstChar == underscoreChar
			field := StructField{
				Name:         rsf.Name,
				Index:        fieldIndex,
				Anonymous:    rsf.Anonymous,
				Private:      private,
				JSONExcluded: private,
				TypeMeta:     s.get(rsf.Type, true),
			}

			tags, err := structtag.Parse(string(rsf.Tag))
			if err != nil {
				panic("failed parsing tags <" + string(rsf.Tag) + "> of field \"" + field.String() + "\": " + err.Error())
			}
			field.Tags = tags
			if jsonTag, _ := tags.Get("json"); jsonTag != nil {
				if jsonTag.Name != "" {
					if jsonTag.Name != "-" {
						field.JSONName = jsonTag.Name
					} else {
						field.JSONExcluded = true
					}
				} else {
					field.JSONName = rsf.Name
				}
				if jsonTag.HasOption("omitempty") {
					field.JSONOmitEmpty = true
				}
			} else if !field.JSONExcluded {
				field.JSONName = rsf.Name
			}
			if defaultValueTag, err := tags.Get("default"); defaultValueTag != nil && err == nil {
				defaultValueStr := defaultValueTag.Name
				// defaultReflectValue, err := field.ParseReflectValue(reflect.ValueOf(defaultValueStr))
				// if err != nil {
				// 	panic("Unable to convert default value in tag of field \"" + field.String() + "\": " + err.Error())
				// }
				field.DefaultValue = defaultValueStr
			}
			if descriptionTag, err := tags.Get("description"); descriptionTag != nil && err == nil {
				field.Description = descriptionTag.Value()
			}

			strct.Fields[fieldIndex] = field
		}
		s.unlock(locked)
		return strct
	case reflect.Map:
		mp := &Map{typ: rtyp}
		s.types[rtyp] = mp
		mp.Key = s.get(rtyp.Key(), true)
		mp.Elem = s.get(rtyp.Elem(), true)
		s.unlock(locked)
		return mp
	case reflect.Slice:
		sl := &Slice{typ: rtyp}
		s.types[rtyp] = sl
		sl.Elem = s.get(rtyp.Elem(), true)
		s.unlock(locked)
		return sl
	case reflect.Array:
		arr := &Array{typ: rtyp}
		s.types[rtyp] = arr
		arr.Elem = s.get(rtyp.Elem(), true)
		s.unlock(locked)
		return arr
	case reflect.Interface:
		i := &Interface{rtyp}
		s.types[rtyp] = i
		s.unlock(locked)
		return i
	default:
		p := &Primitive{rtyp, nil, ""}
		s.types[rtyp] = p
		s.unlock(locked)
		return p
	}
}

// locks to schema mutex if locked is false
func (s *Schema) lock(locked bool) {
	if !locked {
		s.mu.Lock()
	}
}

// unlocks to schema mutex if locked is false
func (s *Schema) unlock(locked bool) {
	if !locked {
		s.mu.Unlock()
	}
}

// GetPrimitive is `Get` but asserts the returned type meta to `*Primitive`, meaning it panics if the specified type is not primitive.
func (s *Schema) GetPrimitive(typ interface{}) *Primitive {
	t, ok := s.Get(typ).(*Primitive)
	if !ok {
		panic(s.Get(typ).String() + " is not a primitive type")
	}
	return t
}

// GetStruct is `Get` but asserts the returned type meta to `*Struct`, meaning it panics if the specified type is not a struct.
func (s *Schema) GetStruct(typ interface{}) *Struct {
	t, ok := s.Get(typ).(*Struct)
	if !ok {
		panic(s.Get(typ).String() + " is not a struct type")
	}
	return t
}

// NewSchema returns a new type meta schema
func NewSchema() *Schema {
	return &Schema{sync.Mutex{}, make(map[reflect.Type]TypeMeta), make(map[*Enum]TypeMeta)}
}

// DefaultSchema is the default type meta schema
var DefaultSchema = NewSchema()

// Get returns type meta for the specified type. If type meta is passed, that is returned. If a reflect type is passed,
// the type meta for that type is returned. If a reflect kind is passed, the type meta for the default type for that kind is returned.
// In all other cases, the reflect type is retrieved of the passed argument using `reflect.TypeOf` and the type meta is retrieved from that type.
// If no type can be inferred, i.e. `reflect.TypeOf` returns nil, a *typemeta.Interface with a nil type is returned.
// Each inferred reflect type points to a single type meta object.
func Get(typ interface{}) TypeMeta {
	return DefaultSchema.Get(typ)
}

// GetPrimitive is `Get` but asserts the returned type meta to `*Primitive`, meaning it panics if the specified type is not primitive.
func GetPrimitive(typ interface{}) *Primitive {
	return DefaultSchema.GetPrimitive(typ)
}

// GetStruct is `Get` but asserts the returned type meta to `*Struct`, meaning it panics if the specified type is not a struct.
func GetStruct(typ interface{}) *Struct {
	return DefaultSchema.GetStruct(typ)
}

// SliceOf returns slice type meta for the specified type
func SliceOf(typ interface{}) *Slice {
	return DefaultSchema.SliceOf(typ)
}

// PtrTo returns type meta for a pointer to the specified type
func PtrTo(typ interface{}) *Ptr {
	return DefaultSchema.PtrTo(typ)
}

var kindTypes = map[reflect.Kind]reflect.Type{
	reflect.Bool:       reflect.TypeOf(false),
	reflect.Int:        reflect.TypeOf(int(0)),
	reflect.Int8:       reflect.TypeOf(int8(0)),
	reflect.Int16:      reflect.TypeOf(int16(0)),
	reflect.Int32:      reflect.TypeOf(int32(0)),
	reflect.Int64:      reflect.TypeOf(int64(0)),
	reflect.Uint:       reflect.TypeOf(uint(0)),
	reflect.Uint8:      reflect.TypeOf(uint8(0)),
	reflect.Uint16:     reflect.TypeOf(uint16(0)),
	reflect.Uint32:     reflect.TypeOf(uint32(0)),
	reflect.Uint64:     reflect.TypeOf(uint64(0)),
	reflect.Uintptr:    reflect.TypeOf(uintptr(0)),
	reflect.Float32:    reflect.TypeOf(float32(0)),
	reflect.Float64:    reflect.TypeOf(float64(0)),
	reflect.Complex64:  reflect.TypeOf(complex64(0)),
	reflect.Complex128: reflect.TypeOf(complex128(0)),
	reflect.Array:      reflect.TypeOf([]interface{}{}),
	reflect.Map:        reflect.TypeOf(map[interface{}]interface{}{}),
	reflect.Slice:      reflect.TypeOf([]interface{}{}),
	reflect.String:     reflect.TypeOf(string("")),
	reflect.Interface:  reflect.TypeOf(nil),
}

var primitiveType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()

const underscoreChar rune = '_'
