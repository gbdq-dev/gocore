// Package transformer provides helpers to serialize and deserialize values
// using supported wire formats.
package transformer

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported format")
	ErrTargetNil         = errors.New("target cannot be nil")
	ErrTargetNotPointer  = errors.New("target must be a pointer")
	ErrTargetPointerNil  = errors.New("target pointer cannot be nil")
)

// Format represents the wire format to use for serialization and deserialization.
type Format string

const (
	// JSON encodes and decodes data using encoding/json.
	JSON Format = "json"
	// XML encodes and decodes data using encoding/xml.
	XML Format = "xml"
)

// Transform serializes data into the requested format.
//
// Supported formats are JSON and XML. It returns an error when format is not
// supported.
func Transform(data any, format Format) ([]byte, error) {
	switch format {
	case JSON:
		return json.Marshal(data)
	case XML:
		return xml.Marshal(data)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
	}
}

// Parse deserializes data in the requested format into target.
//
// Target must be a non-nil pointer to the destination value. Supported formats
// are JSON and XML. It returns an error when format is not supported.
func Parse(data []byte, format Format, target any) error {
	if target == nil {
		return ErrTargetNil
	}

	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr {
		return ErrTargetNotPointer
	}

	if rv.IsNil() {
		return ErrTargetPointerNil
	}

	switch format {
	case JSON:
		return json.Unmarshal(data, target)
	case XML:
		return xml.Unmarshal(data, target)
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
	}
}
