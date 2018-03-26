package rfc7807

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strings"
)

const (
	// JSONMediaType is the MIME Media type for the Problem struct
	JSONMediaType = "application/problem+json"
)

// ErrExtensionKeyIsReserved is thrown when attempting to call .Extend(k,v) on a problem with a reserved key name
var ErrExtensionKeyIsReserved = errors.New("rfc7807: the given extension key name is reserved please choose another name")

var reservedKeys = map[string]struct{}{
	"type":     {},
	"title":    {},
	"status":   {},
	"detail":   {},
	"instance": {},
}

// Problem is a struct representing Problem Details as descrbed in rfc7807
type Problem struct {
	Type       string
	Title      string
	Status     int64
	Detail     string
	Instance   url.URL
	extensions map[string]interface{}
}

// Extensions returns a slice of strings representing the names of extension keys for this Problem struct
func (p Problem) Extensions() []string {

	extensions := make([]string, len(p.extensions))
	x := 0

	for extension := range p.extensions {
		extensions[x] = extension
		x++
	}

	return extensions

}

// Extension retrieves the value for an extension if present. A bool is also returned to signify whether the value was
// present upon retrieval
func (p Problem) Extension(key string) (interface{}, bool) {
	val, ok := p.extensions[key]
	return val, ok
}

// Extend adds an extension to the Problem. Only non-reserved extension keys are allowed
func (p *Problem) Extend(key string, value interface{}) error {

	if _, reserved := reservedKeys[strings.ToLower(key)]; reserved {
		return ErrExtensionKeyIsReserved
	}

	if p.extensions == nil {
		p.extensions = make(map[string]interface{})
	}

	p.extensions[key] = value

	return nil

}

// MarshalJSON Marshals JSON
func (p Problem) MarshalJSON() ([]byte, error) {

	out := map[string]interface{}{}

	if p.Type != "" {
		out["type"] = p.Type
	}

	if p.Title != "" {
		out["title"] = p.Title
	}

	if p.Status != 0 {
		out["status"] = p.Status
	}

	if p.Detail != "" {
		out["detail"] = p.Detail
	}

	var zero url.URL
	if p.Instance != zero {
		out["instance"] = p.Instance.String()
	}

	for k, v := range p.extensions {
		out[k] = v
	}

	return json.Marshal(out)

}

// UnmarshalJSON unmarshalls JSON
func (p *Problem) UnmarshalJSON(data []byte) error {

	in := map[string]interface{}{}
	json.Unmarshal(data, &in)

	for k, v := range in {

		switch strings.ToLower(k) {
		case "type":

			if str, ok := v.(string); ok {
				p.Type = str
			} else {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "type",
					Struct: "Problem",
				}
			}

		case "title":

			if str, ok := v.(string); ok {
				p.Title = str
			} else {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "title",
					Struct: "Problem",
				}
			}

		case "status":

			if num, ok := v.(float64); ok {
				p.Status = int64(num)
			} else {
				return &json.UnmarshalTypeError{
					Value:  "number",
					Type:   reflect.TypeOf(int64(0)),
					Field:  "status",
					Struct: "Problem",
				}
			}

		case "detail":

			if str, ok := v.(string); ok {
				p.Detail = str
			} else {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "detail",
					Struct: "Problem",
				}
			}

		case "instance":

			if str, ok := v.(string); ok {

				if uri, err := url.Parse(str); err == nil {
					p.Instance = *uri
				} else {

					return &json.UnmarshalTypeError{
						Value:  "uri",
						Type:   reflect.TypeOf(p.Instance),
						Field:  "instance",
						Struct: "Problem",
					}

				}

			} else {

				return &json.UnmarshalTypeError{
					Value:  "uri",
					Type:   reflect.TypeOf(""),
					Field:  "instance",
					Struct: "Problem",
				}

			}

		default:

			if err := p.Extend(k, v); err != nil {

				t := reflect.TypeOf(v)

				return &json.UnmarshalTypeError{
					Value:  t.Name(),
					Type:   t,
					Field:  k,
					Struct: "Problem",
				}

			}

		}

	}

	return nil

}
