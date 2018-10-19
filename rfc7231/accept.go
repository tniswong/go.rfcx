package rfc7231

import (
	"io"
	"sort"
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
	mediaRanges []mediaRange
}

// MostAcceptable returns the single most ok mediaType by weight. If no given mediaType is ok,
// will return "", false
func (a Accept) MostAcceptable(mediaTypes []string) (string, bool) {

	sort.Sort(byWeight(a.mediaRanges))

	for _, mediaRange := range a.mediaRanges {

		for _, mediaType := range mediaTypes {

			if mediaRange.Supports(mediaType) {
				return mediaType, true
			}

		}

	}

	return "", false

}

// Acceptable returns whether or not the mediaType is ok to the HTTP Accept Header
func (a Accept) Acceptable(mediaType string) bool {

	// treat zero-value Accept{} as "*/*"
	if len(a.mediaRanges) == 0 {
		return true
	}

	acceptable := false

	for _, mediaRange := range a.mediaRanges {

		acceptable = acceptable || mediaRange.Supports(mediaType)

		if acceptable {
			break
		}

	}

	return acceptable

}

// String the string representation of the HTTP Accept Header
func (a Accept) String() string {

	// treat zero-value Accept{} as "*/*"
	if len(a.mediaRanges) == 0 {
		return "*/*"
	}

	var mrStrings []string

	for _, mr := range a.mediaRanges {
		mrStrings = append(mrStrings, mr.String())
	}

	return strings.Join(mrStrings, ", ")

}
