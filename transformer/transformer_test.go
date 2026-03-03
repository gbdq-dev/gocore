package transformer

import (
	"errors"
	"testing"
)

type person struct {
	Name string `json:"name" xml:"name"`
	Age  int    `json:"age" xml:"age"`
}

func TestTransform(t *testing.T) {
	tests := []struct {
		name    string
		data    any
		format  Format
		want    string
		wantErr error
		errText string
	}{
		{
			name:   "json",
			data:   person{Name: "Alice", Age: 30},
			format: JSON,
			want:   `{"name":"Alice","age":30}`,
		},
		{
			name:   "xml",
			data:   person{Name: "Alice", Age: 30},
			format: XML,
			want:   `<person><name>Alice</name><age>30</age></person>`,
		},
		{
			name:    "unsupported format",
			data:    person{Name: "Alice", Age: 30},
			format:  Format("yaml"),
			wantErr: ErrUnsupportedFormat,
			errText: "unsupported format: yaml",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := Transform(tc.data, tc.format)

			if tc.wantErr != nil {
				assertErrorIs(t, err, tc.wantErr)
				if tc.errText != "" {
					assertErrorString(t, err, tc.errText)
				}
				return
			}

			if err != nil {
				t.Fatalf("Transform() returned error: %v", err)
			}

			if string(got) != tc.want {
				t.Fatalf("Transform() = %q, want %q", string(got), tc.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	var nilPerson *person

	tests := []struct {
		name    string
		data    []byte
		format  Format
		target  any
		want    person
		wantErr error
		errText string
	}{
		{
			name:   "json",
			data:   []byte(`{"name":"Bob","age":25}`),
			format: JSON,
			target: &person{},
			want:   person{Name: "Bob", Age: 25},
		},
		{
			name:   "xml",
			data:   []byte(`<person><name>Bob</name><age>25</age></person>`),
			format: XML,
			target: &person{},
			want:   person{Name: "Bob", Age: 25},
		},
		{
			name:    "nil target",
			data:    []byte(`{"name":"Bob","age":25}`),
			format:  JSON,
			target:  nil,
			wantErr: ErrTargetNil,
			errText: "target cannot be nil",
		},
		{
			name:    "non-pointer target",
			data:    []byte(`{"name":"Bob","age":25}`),
			format:  JSON,
			target:  person{},
			wantErr: ErrTargetNotPointer,
			errText: "target must be a pointer",
		},
		{
			name:    "nil pointer target",
			data:    []byte(`{"name":"Bob","age":25}`),
			format:  JSON,
			target:  nilPerson,
			wantErr: ErrTargetPointerNil,
			errText: "target pointer cannot be nil",
		},
		{
			name:    "unsupported format",
			data:    []byte(`{"name":"Bob","age":25}`),
			format:  Format("yaml"),
			target:  &person{},
			wantErr: ErrUnsupportedFormat,
			errText: "unsupported format: yaml",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := Parse(tc.data, tc.format, tc.target)

			if tc.wantErr != nil {
				assertErrorIs(t, err, tc.wantErr)
				if tc.errText != "" {
					assertErrorString(t, err, tc.errText)
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse() returned error: %v", err)
			}

			got, ok := tc.target.(*person)
			if !ok {
				t.Fatalf("test setup error: expected *person target, got %T", tc.target)
			}

			if *got != tc.want {
				t.Fatalf("Parse() = %+v, want %+v", *got, tc.want)
			}
		})
	}
}

func assertErrorString(t *testing.T, err error, want string) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error %q, got nil", want)
	}

	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func assertErrorIs(t *testing.T, err error, want error) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error matching %v, got nil", want)
	}

	if !errors.Is(err, want) {
		t.Fatalf("error %q does not match %v", err, want)
	}
}
