package rfc8288

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"
)

var (
	// ErrExtensionKeyIsReserved describes an attempt to call Link.Extend(k,v) with a reserved key name
	ErrExtensionKeyIsReserved = errors.New("rfc8288: the given extension key name is reserved please choose another name")

	// ReservedKeys holds the names of all the reserved key names that are not allowed to be used as extensions
	ReservedKeys = map[string]struct{}{
		"href":     {},
		"rel":      {},
		"hreflang": {},
		"media":    {},
		"title":    {},
		"title*":   {},
		"type":     {},
	}
)

// ParseLink attempts to parse a link string
func ParseLink(link string) (Link, error) {

	var (
		rs io.RuneScanner = strings.NewReader(link)
		s                 = scanner{runeScanner: rs}
		p                 = parser{scanner: s}
	)

	return p.parse()

}

// Link is an implementation of the structure defined by RFC8288 Web Linking
type Link struct {
	HREF      url.URL
	Rel       string
	HREFLang  string
	Media     string
	Title     string
	TitleStar string
	Type      string

	extensionKeys []string
	extensions    map[string]interface{}
}

// String returns the Link in a format usable for HTTP Headers as defined by RFC8288
func (l Link) String() string {

	var result []string

	result = append(result, fmt.Sprintf(`<%s>`, l.HREF.String()))

	if l.Rel != "" {
		result = append(result, fmt.Sprintf(`rel="%s"`, l.Rel))
	}

	if l.HREFLang != "" {
		result = append(result, fmt.Sprintf(`hreflang="%s"`, l.HREFLang))
	}

	if l.Media != "" {
		result = append(result, fmt.Sprintf(`media="%s"`, l.Media))
	}

	if l.Title != "" {
		result = append(result, fmt.Sprintf(`title="%s"`, l.Title))
	}

	if l.TitleStar != "" {
		result = append(result, fmt.Sprintf(`title*="%s"`, l.TitleStar))
	}

	if l.Type != "" {
		result = append(result, fmt.Sprintf(`type="%s"`, l.Type))
	}

	for key, value := range l.extensions {
		result = append(result, fmt.Sprintf(`%s="%s"`, key, value))
	}

	return strings.Join(result, "; ")

}

// ExtensionKeys returns a slice of strings representing the names of extension keys for this Link struct in the order
// they were added
func (l Link) ExtensionKeys() []string {
	return l.extensionKeys
}

// Extension retrieves the value for an extension if present. A bool is also returned to signify whether the value was
// present upon retrieval
func (l *Link) Extension(key string) (interface{}, bool) {

	if l.extensions == nil {
		l.extensions = make(map[string]interface{})
	}

	val, ok := l.extensions[key]
	return val, ok

}

// Extend adds an extension to the Link. Only non-reserved extension keys are allowed.
// Setting the value to nil will remove the extension.
func (l *Link) Extend(key string, value interface{}) error {

	if _, reserved := ReservedKeys[strings.ToLower(key)]; reserved {
		return ErrExtensionKeyIsReserved
	}

	_, keyFound := l.Extension(key)
	if !keyFound {
		l.extensionKeys = append(l.extensionKeys, key)
	}

	if value != nil {
		l.extensions[key] = value
	} else {

		delete(l.extensions, key)

		for x := 0; x < len(l.extensionKeys); {

			if strings.EqualFold(key, l.extensionKeys[x]) {
				l.extensionKeys = append(l.extensionKeys[:x], l.extensionKeys[x+1:]...)
				break
			}

			x++

		}

	}

	return nil

}

// MarshalJSON Marshals JSON
func (l Link) MarshalJSON() ([]byte, error) {

	out := map[string]interface{}{}

	var zero url.URL
	if l.HREF != zero {
		out["href"] = l.HREF.String()
	}

	if l.Rel != "" {
		out["rel"] = l.Rel
	}

	if l.HREFLang != "" {
		out["hreflang"] = l.HREFLang
	}

	if l.Media != "" {
		out["media"] = l.Media
	}

	if l.Title != "" {
		out["title"] = l.Title
	}

	if l.TitleStar != "" {
		out["title*"] = l.TitleStar
	}

	if l.Type != "" {
		out["type"] = l.Type
	}

	for _, extensionKey := range l.extensionKeys {
		out[extensionKey] = l.extensions[extensionKey]
	}

	return json.Marshal(out)

}

// UnmarshalJSON unmarshal JSON
func (l *Link) UnmarshalJSON(data []byte) error {

	in := map[string]interface{}{}
	json.Unmarshal(data, &in)

	for k, v := range in {

		switch strings.ToLower(k) {
		case "href":

			if str, ok := v.(string); ok {

				if uri, err := url.Parse(str); err == nil {
					l.HREF = *uri
				} else {

					return &json.UnmarshalTypeError{
						Value:  "uri",
						Type:   reflect.TypeOf(l.HREF),
						Field:  "href",
						Struct: "Link",
					}

				}

			} else {

				return &json.UnmarshalTypeError{
					Value:  "uri",
					Type:   reflect.TypeOf(""),
					Field:  "href",
					Struct: "Link",
				}

			}

		case "rel":

			if str, ok := v.(string); ok {
				l.Rel = str
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
				l.HREFLang = str
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
				l.Media = str
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
				l.Title = str
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
				l.TitleStar = str
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
				l.Type = str
			} else {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "type",
					Struct: "Link",
				}
			}

		default:

			if err := l.Extend(k, v); err != nil {

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
