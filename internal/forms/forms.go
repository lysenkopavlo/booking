package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
)

// Form is a custom form struct
type Form struct {
	url.Values
	Errors errors
}

// New initilaizes a Form struct with data field
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Valid returns true if there are no errors,
// otherwise returns false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")

		}
	}
}

// MinLength check for minimum length
func (f *Form) MinLength(field string, lenght int) bool {
	x := f.Get(field)
	if utf8.RuneCountInString(x)-1 < lenght {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", lenght))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) bool {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
		return false
	}
	return true
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}
