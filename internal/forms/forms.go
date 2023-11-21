package forms

import (
	"net/url"
	"strconv"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{data, errors(map[string][]string{})}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) StringToFloat(field string) (float64, error) {
	value := f.Get(field)
	if value == "" {
		return 0, nil
	}
	return strconv.ParseFloat(value, 64)
}

func (f *Form) StringToInt(field string) (int, error) {
	value := f.Get(field)
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
