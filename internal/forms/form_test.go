package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	var tF *Form
	r := httptest.NewRequest("POST", "/sumeURL", nil)

	tF = New(r.PostForm)

	if !tF.Valid() {
		t.Error("empty Errors field is not empty, but it should.")
	}

}

func TestForm_Required(t *testing.T) {
	var tF *Form
	r := httptest.NewRequest("POST", "/someURL", nil)

	tF = New(r.PostForm)

	tF.Required("a", "b", "c")
	if tF.Valid() {
		t.Error("empty Errors field is not empty, but it should.")
	}

	// Fill in some data into Form
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r.PostForm = postedData
	tF = New(r.PostForm)

	tF.Required("a", "b", "c")

	if !tF.Valid() {
		t.Error("Errors field is empty, but it must be filled!")
	}

}

func TestForm_MinLenght(t *testing.T) {
	postedData := url.Values{}
	tF := New(postedData)
	if tF.MinLength("x", 10) {
		t.Error("it should not have a field, but it does.")
	}
	// Testing Get function from errors.go here
	if tF.Errors.Get("x") == "" {
		t.Error("failed error detecton.")
	}

	postedData = url.Values{}
	postedData.Add("anyField", "anyValue")
	tF = New(postedData)
	if tF.MinLength("anyField", 100) {
		t.Error("incorrect length messurement.")
	}

	postedData = url.Values{}
	postedData.Add("someField", "someValue")
	tF = New(postedData)
	if !tF.MinLength("someField", 1) {
		t.Error("it should have a field, but it doesn't.")
	}
	if tF.Errors.Get("someField") != "" {
		t.Error("Fake error detection.")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	tF := New(postedData)

	if tF.IsEmail("x") {
		t.Error("it should not have a field, but it does.")
	}

	postedData = url.Values{}
	postedData.Add("email", "me@here.com")

	tF = New(postedData)

	if !tF.IsEmail("email") {
		t.Error("it should have a field, but it doesn't.")
	}

	postedData = url.Values{}
	postedData.Add("email", "invalidAddr")

	tF = New(postedData)

	if tF.IsEmail("email") {
		t.Error("it shows that email is valid.")
	}
}
func TestForm_Has(t *testing.T) {
	postedData := url.Values{}

	tF := New(postedData)

	if tF.Has("non-existant-field") {
		t.Error("it should not have a field, but it does.")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")

	tF = New(postedData)

	if !tF.Has("a") {
		t.Error("It should has a field, but it doesn't.")
	}

}
