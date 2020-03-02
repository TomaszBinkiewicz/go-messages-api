package project_package

import (
	"reflect"
	"testing"
)

func TestGetTime(t *testing.T) {
	var intTime interface{} = GetTime()
	_, ok := intTime.(int)
	if !ok {
		t.Errorf("GetTime return value type is %v, want int", reflect.TypeOf(intTime))
	}
}

func TestValidateEmail(t *testing.T) {
	// proper email address
	validation := ValidateEmail("some.email@example.com")
	if validation == false {
		t.Errorf("Function ValidateEmail should return true, but returned false.")
	}
	// proper email address
	validation = ValidateEmail("some12.!#$%&'*+/=?^_`{|}~-EMAIL@example.com")
	if validation == false {
		t.Errorf("Function ValidateEmail should return true, but returned false.")
	}
	// incorrect email address
	validation = ValidateEmail("some.email.example.com")
	if validation == true {
		t.Errorf("Function ValidateEmail should return false, but returned true.")
	}
	// incorrect email address
	validation = ValidateEmail("some.email@")
	if validation == true {
		t.Errorf("Function ValidateEmail should return false, but returned true.")
	}
	// proper email address
	validation = ValidateEmail("@example.com")
	if validation == true {
		t.Errorf("Function ValidateEmail should return false, but returned true.")
	}
}
