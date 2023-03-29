package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	form := New(postedData)

	if !form.Has("b") {
		t.Error("form says it does not have a field when it does")
	}

	if form.Has("d") {
		t.Error("form says it has a field when it doesnt")
	}
}

func TestForm_Minlength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "aaaaa")
	form := New(postedData)

	// check for field that doesnt exist
	form.MinLength("b", 5)
	if form.Valid() {
		t.Error("field doesnt exist on form but was accepted")
	}

	isError := form.Errors.Get("b")
	if isError == "" {
		t.Error("should have an error but did not get one")
	}

	form = New(postedData)
	form.MinLength("a", 5)
	if !form.Valid() {
		t.Error("length of form entry was 5 but a min length of 5 was not accepted")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have an error but has one")
	}

	form = New(postedData)
	form.MinLength("a", 6)
	if form.Valid() {
		t.Error("length of form entry was 5 but a min length of 6 was accepted")
	}
}

func TestForm_IsEmail(t *testing.T) {
	testData := []struct {
		str   string
		legit bool
	}{
		{"abcde", false},
		{"aaaa@aaaa", false},
		{"@abcd.com", false},
		{"abdc.co.uk", false},
		{"abcds@", false},
		{"abcds@abcd.", false},
		{"abcds@abcd.co", true},
		{"abcds@abcd.co.uk", true},
		{"abcds@abcd.com", true},
	}

	for _, data := range testData {
		postedData := url.Values{}
		postedData.Add("email", data.str)
		form := New(postedData)

		form.IsEmail("email")

		valid := form.Valid()
		if !valid && data.legit {
			t.Error("valid email address seen as invalid")
		} else if valid && !data.legit {
			t.Error("invalid email address seen as valid")
		}
	}
}
