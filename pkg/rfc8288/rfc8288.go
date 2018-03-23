package rfc8288

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

var ErrExtensionKeyIsReserved = errors.New("rfc8288: the given extension key name is reserved please choose another name")

var reservedKeys = map[string]struct{}{
	"rel":      {},
	"hreflang": {},
	"media":    {},
	"title":    {},
	"title*":   {},
	"type":     {},
}

type Link struct {
	HREF       url.URL
	Rel        string
	HREFLang   string
	Media      string
	Title      string
	TitleStar  string
	Type       string
	extensions map[string]interface{}
}

func (w Link) String() string {

	var result []string

	result = append(result, fmt.Sprintf(`<%s>`, w.HREF.String()))

	if w.Rel != "" {
		result = append(result, fmt.Sprintf(`rel="%s"`, w.Rel))
	}

	if w.HREFLang != "" {
		result = append(result, fmt.Sprintf(`hreflang="%s"`, w.HREFLang))
	}

	if w.Media != "" {
		result = append(result, fmt.Sprintf(`media="%s"`, w.Media))
	}

	if w.Title != "" {
		result = append(result, fmt.Sprintf(`title="%s"`, w.Title))
	}

	if w.TitleStar != "" {
		result = append(result, fmt.Sprintf(`title*="%s"`, w.TitleStar))
	}

	if w.Type != "" {
		result = append(result, fmt.Sprintf(`type="%s"`, w.Type))
	}

	for key, value := range w.extensions {
		result = append(result, fmt.Sprintf(`%s="%s"`, key, value))
	}

	return strings.Join(result, "; ")

}

func (w Link) Extensions() []string {

	extensions := make([]string, len(w.extensions))
	x := 0

	for extension := range w.extensions {
		extensions[x] = extension
		x++
	}

	return extensions

}

func (w Link) Extension(key string) (interface{}, bool) {
	val, ok := w.extensions[key]
	return val, ok
}

func (w *Link) Extend(key string, value interface{}) error {

	if _, reserved := reservedKeys[strings.ToLower(key)]; reserved {
		return ErrExtensionKeyIsReserved
	}

	if w.extensions == nil {
		w.extensions = make(map[string]interface{})
	}

	w.extensions[key] = value

	return nil

}

func (w Link) MarshalJSON() ([]byte, error) {

	out := map[string]interface{}{}

	var zero url.URL
	if w.HREF != zero {
		out["href"] = w.HREF.String()
	}

	if w.Rel != "" {
		out["rel"] = w.Rel
	}

	if w.HREFLang != "" {
		out["hreflang"] = w.HREFLang
	}

	if w.Media != "" {
		out["media"] = w.Media
	}

	if w.Title != "" {
		out["title"] = w.Title
	}

	if w.TitleStar != "" {
		out["title*"] = w.TitleStar
	}

	if w.Type != "" {
		out["type"] = w.Type
	}

	for k, v := range w.extensions {
		out[k] = v
	}

	return json.Marshal(out)

}

func (w *Link) UnmarshalJSON(data []byte) error {

	in := map[string]interface{}{}
	json.Unmarshal(data, &in)

	for k, v := range in {

		switch strings.ToLower(k) {
		case "href":

			if str, ok := v.(string); ok {

				if uri, err := url.Parse(str); err == nil {
					w.HREF = *uri
				} else {

					return &json.UnmarshalTypeError{
						Value:  "uri",
						Type:   reflect.TypeOf(w.HREF),
						Field:  "href",
						Struct: "Link",
					}

				}

			} else {

				return &json.UnmarshalTypeError{
					Value:  "uri",
					Type:   reflect.TypeOf(""),
					Field:  "instance",
					Struct: "Link",
				}

			}

		case "rel":

			if str, ok := v.(string); ok {
				w.Rel = str
			} else {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "rel",
					Struct: "Link",
				}
			}

		case "hreflang":

			if str, ok := v.(string); ok {
				w.HREFLang = str
			} else {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(int64(0)),
					Field:  "hreflang",
					Struct: "Link",
				}
			}

		case "media":

			if str, ok := v.(string); ok {
				w.Media = str
			} else {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "media",
					Struct: "Link",
				}
			}

		case "title":

			if str, ok := v.(string); ok {
				w.Title = str
			} else {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "title",
					Struct: "Link",
				}
			}

		case "title*":

			if str, ok := v.(string); ok {
				w.TitleStar = str
			} else {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "title*",
					Struct: "Link",
				}
			}

		case "type":

			if str, ok := v.(string); ok {
				w.Type = str
			} else {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "type",
					Struct: "Link",
				}
			}

		default:

			if err := w.Extend(k, v); err != nil {

				t := reflect.TypeOf(v)

				return &json.UnmarshalTypeError{
					Value:  t.Name(),
					Type:   t,
					Field:  k,
					Struct: "Link",
				}

			}

		}

	}

	return nil

}
