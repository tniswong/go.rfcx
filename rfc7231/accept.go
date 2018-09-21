package rfc7231

import (
	"io"
	"strings"
)

// ParseAccept parses the value of an HTTP Accept Header as defined in RFC 7231 Sec. 5.3.2
func ParseAccept(accept string) (Accept, error) {

	var (
		rs io.RuneScanner = strings.NewReader(accept)
		s                 = scanner{runeScanner: rs}
		p                 = parser{scanner: s}
	)

	return p.parse()

}

// Accept represents the value of an HTTP Accept Header as defined in RFC 7231 Sec. 5.3.2
type Accept struct {
	MediaRanges []MediaRange
}

// Type returns the string representation of the HTTP Accept Header
func (a Accept) String() string {

	var mrStrings []string

	for _, mr := range a.MediaRanges {
		mrStrings = append(mrStrings, mr.String())
	}

	return strings.Join(mrStrings, ", ")

}
