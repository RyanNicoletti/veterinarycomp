package forms

import (
	"io"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

func NewForm(data url.Values) *Form {
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
	fieldFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		f.Errors.Add(field, "Enter a number")
		return 0, nil
	}
	return fieldFloat, nil
}

func (f *Form) StringToInt(field string) (int, error) {
	value := f.Get(field)
	if value == "" {
		return 0, nil
	}
	fieldInt, err := strconv.Atoi(value)
	if err != nil {
		f.Errors.Add(field, "Enter a number")
		return 0, nil
	}
	return fieldInt, nil
}

func (f *Form) ProcessFileUpload(fieldName string, fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		f.Errors.Add(fieldName, "Error opening file")
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		f.Errors.Add(fieldName, "Error reading file")
		return nil, err
	}

	return data, nil
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
