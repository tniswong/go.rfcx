# go.rfcx

Go Library with implementations of RFC's that I find useful.

# RFC 7807 Problem Details

https://tools.ietf.org/html/rfc7807

```go
package rfc7807

const (
	// JSONMediaType is the MIME Media type for the Problem struct
	JSONMediaType = "application/problem+json"
)

// ErrExtensionKeyIsReserved is thrown when attempting to call .Extend(k,v) on a problem with a reserved key name
var ErrExtensionKeyIsReserved = errors.New("rfc7807: the given extension key name is reserved please choose another name")

// Problem is a struct representing Problem Details as descrbed in rfc7807
type Problem struct {
	Type       string
	Title      string
	Status     int64
	Detail     string
	Instance   url.URL
}

// Extensions returns a slice of strings representing the names of extension keys for this Problem struct
func (p Problem) Extensions() []string { ... }

// Extension retrieves the value for an extension if present. A bool is also returned to signify whether the value was
// present upon retrieval
func (p Problem) Extension(key string) (interface{}, bool) { ... }

// Extend adds an extension to the Problem. Only non-reserved extension keys are allowed
func (p *Problem) Extend(key string, value interface{}) { ... }
```

# RFC 8288 Web Linking

https://tools.ietf.org/html/rfc8288

```go
package rfc8288

// ErrExtensionKeyIsReserved is thrown when attempting to call .Extend(k,v) on a Link with a reserved key name
var ErrExtensionKeyIsReserved = errors.New("rfc8288: the given extension key name is reserved please choose another name")

// Parse Error Types
var (
	ErrInvalidLink         = errors.New("[ERR] invalid link")
	ErrMissingSemicolon    = errors.New("[ERR] invalid link: missing semicolon")
	ErrMissingClosingQuote = errors.New("[ERR] invalid link: missing closing quote")
	ErrMissingAttrValue    = errors.New("[ERR] invalid link: missing attribute value")
)

// ParseLink attempts to parse a link string
func ParseLink(link string) (Link, error) {
	return newParser(strings.NewReader(link)).Parse()
}

// Link is an implementation of the structure defined by RFC8288 Web Linking
type Link struct {
	HREF       url.URL
	Rel        string
	HREFLang   string
	Media      string
	Title      string
	TitleStar  string
	Type       string
}

// String returns the Link in a format usable for HTTP Headers as defined by RFC8288
func (w Link) String() string { ... }

// Extensions returns a slice of strings representing the names of extension keys for this Link struct
func (w Link) Extensions() []string { ... }

// Extension retrieves the value for an extension if present. A bool is also returned to signify whether the value was
// present upon retrieval
func (w Link) Extension(key string) (interface{}, bool) { ... }

// Extend adds an extension to the Link. Only non-reserved extension keys are allowed
func (w *Link) Extend(key string, value interface{}) error { ... }
```