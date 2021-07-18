package typemeta

import (
	"reflect"
	"testing"
	"time"
)

func TestStruct(t *testing.T) {
	// test simple struct
	type StructA struct {
		String     string                 `json:"string" description:"A string field"`
		Bool       bool                   `json:"bool,omitempty"`
		Record     map[string]interface{} `json:"-"`
		privateKey string
		PublicKey  string
		Time       time.Time
		TimePtr    *time.Time

		_ string
	}
	var structA = Get(StructA{}).(*Struct)

	// check that the properties are set correctly
	if structA.Name() != "StructA" {
		t.Error("does not have name of struct")
	}
	if structA.Type() != reflect.TypeOf(StructA{}) {
		t.Error("does not have type of struct")
	}
	if structA.Primitive() {
		t.Error("should not be defined as primitive")
	}

	// check string field
	stringField := structA.EnsureField(0)
	if stringField.Kind() != reflect.String {
		t.Error("String field should have string kind")
	}
	if stringField.Name != "String" {
		t.Error("String field does not have correctly defined name")
	}
	if stringField.JSONName != "string" {
		t.Error("String field does not have correctly defined JSON name")
	}
	if stringField.JSONOmitEmpty {
		t.Error("String field should not have JSONOmitEmpty option")
	}
	if stringField.Private {
		t.Error("String field should not have JSONOmitEmpty option")
	}
	if stringField.Description == "" {
		t.Error("String field should have a defined description")
	}

	// check bool field
	boolField := structA.EnsureField(1)
	if boolField.Kind() != reflect.Bool {
		t.Error("Bool field should have bool kind")
	}
	if boolField.Name != "Bool" {
		t.Error("Bool field does not have correctly defined name")
	}
	if boolField.JSONName != "bool" {
		t.Error("Bool field does not have correctly defined JSON name")
	}
	if !boolField.JSONOmitEmpty {
		t.Error("Bool field should have JSONOmitEmpty option")
	}
	if boolField.Description != "" {
		t.Error("Bool field should not have a defined description")
	}

	// check bool field
	recordField := structA.EnsureField(2)
	if recordField.Kind() != reflect.Map {
		t.Error("Record field should have map kind")
	}
	if recordField.JSONName != "" {
		t.Error("Record field has a defined JSON name")
	}
	if !recordField.JSONExcluded {
		t.Error("Record field should be JSON excluded")
	}
	if recordTypeMeta, ok := recordField.TypeMeta.(*Map); !ok {
		t.Error("Record field type meta should be a *typemeta.Map")
	} else {
		if !recordTypeMeta.Primitive() {
			t.Error("Record should be primitive")
		}
		if recordTypeMeta.Key.Kind() != reflect.String {
			t.Error("Record field key should have string kind")
		}
		if recordTypeMeta.Elem.Kind() != reflect.Interface {
			t.Error("Record field elem should have interface kind")
		}
		if _, ok := recordTypeMeta.Elem.(*Interface); !ok {
			t.Error("Record field type should be a *typemeta.Interface")
		}
	}

	// check private field
	privateKeyField := structA.EnsureField(3)
	if !privateKeyField.Private {
		t.Error("PrivateKey should be private")
	}
	if !privateKeyField.JSONExcluded {
		t.Error("PrivateKey should be JSON excluded")
	}
	if privateKeyField.JSONName != "" {
		t.Error("PrivateKey should have zero JSON name")
	}

	// check public key field
	publicKeyField := structA.EnsureField(4)
	if publicKeyField.JSONName != "PublicKey" {
		t.Error("PublicKey JSON name should be PublicKey")
	}
	if publicKeyField.JSONExcluded {
		t.Error("PublicKey should not be JSON excluded")
	}

	// check private key field
	timeField := structA.EnsureField(5)
	if timeField.JSONName != "Time" {
		t.Error("Time JSON name should be Time")
	}
	if !timeField.Primitive() {
		t.Error("Time should be primitive")
	}
	if timeFieldTypeMeta, ok := timeField.TypeMeta.(*Primitive); !ok {
		t.Error("Time field type should be a *typemeta.Primitive")
	} else if timeFieldTypeMeta.Kind() != reflect.Struct {
		t.Error("Time field kind be reflect.Struct")
	}

	privateUnderscoreField := structA.EnsureField(7)
	if !privateUnderscoreField.Private {
		t.Error("_ should be private")
	}
}
