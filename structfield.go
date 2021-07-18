package typemeta

import (
	"github.com/fatih/structtag"
)

// StructField is type meta for a struct
type StructField struct {
	Name          string          // Struct field name
	Index         int             // Index of the field in its parent struct
	Anonymous     bool            // Whether the field is anonymous/embedded
	Private       bool            // Whether the field is private to the package it is defined in, i.e. starting with a lowercase letter
	JSONName      string          // JSON name of the field
	JSONExcluded  bool            // Whether the field is excluded when the struct is marshaled to JSON
	JSONOmitEmpty bool            // Whether the value of the field should be set to null when zero
	Description   string          // Description (for API schemas etc.)
	DefaultValue  interface{}     // Default value (for API schemas etc.)
	Tags          *structtag.Tags // Parsed struct field tags
	TypeMeta                      // Type meta of the field value
}

// JSONNonNull returns whether the value will never be defined as null in JSON
func (sf StructField) JSONNonNull() bool {
	return !(sf.JSONOmitEmpty || sf.TypeMeta.JSONNonNull())
}

func (sf StructField) String() string {
	return sf.Name + ": " + sf.TypeMeta.String()
}

// Tag returns a parsed struct field tag, or nil if it has not been defined
func (sf StructField) Tag(key string) *structtag.Tag {
	if sf.Tags == nil {
		return nil
	}
	if tag, _ := sf.Tags.Get(key); tag != nil {
		return tag
	}
	return nil
}
