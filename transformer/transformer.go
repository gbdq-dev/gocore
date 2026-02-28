package transformer

import (
	"encoding/json"
	"encoding/xml"
	"errors"
)

type Format string

const (
	JSON Format = "json"
	XML  Format = "xml"
)

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
