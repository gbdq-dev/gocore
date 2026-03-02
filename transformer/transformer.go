// Package transformer provides helpers to serialize and deserialize values
// using supported wire formats.
package transformer

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"reflect"
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
		return nil, errors.New("unsupported format: " + string(format))
	}
}

// Parse deserializes data in the requested format into target.
//
// Target must be a non-nil pointer to the destination value. Supported formats
// are JSON and XML. It returns an error when format is not supported.
func Parse(data []byte, format Format, target any) error {
	if target == nil {
		return errors.New("target cannot be nil")
	}

	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	if rv.IsNil() {
		return errors.New("target pointer cannot be nil")
	}

	switch format {
	case JSON:
		return json.Unmarshal(data, target)
	case XML:
		return xml.Unmarshal(data, target)
	default:
		return errors.New("unsupported format: " + string(format))
	}
}
