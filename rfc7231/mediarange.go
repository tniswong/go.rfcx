package rfc7231

import (
	"fmt"
	"math"
	"mime"
	"strings"
)

// mediaRange represents a mediaRange as defined for use in an HTTP Accept Header as defined in RFC 7231
type mediaRange struct {
	TypeName    string
	SubtypeName string
	Params      map[string]string
	Q           float64
}

// String returns the string representation of the mediaRange
func (m mediaRange) String() string {

	result := fmt.Sprintf("%s/%s", m.TypeName, m.SubtypeName)

	// quality is a numeric value between 0.000 and 1.000
	q := math.Min(math.Max(m.Q, 0.0), 1.0)

	// q=1.0 is semantically equivalent to q being omitted
	if q > 0.0 && q < 1.0 {
		result = fmt.Sprintf("%s; q=%.1f", result, q)
	}

	for k, v := range m.Params {

		if !strings.EqualFold(k, "q") {
			result = fmt.Sprintf("%s; %s=%s", result, k, v)
		}

	}

	return result

}

// Supports returns whether or not given mediaType is supported by the mediaRange
func (m mediaRange) Supports(t string) bool {

	mediaType, params, err := mime.ParseMediaType(t)

	if err != nil {
		return false
	}

	if m.TypeName == "*" {
		return true
	}

	if m.SubtypeName == "*" {
		return strings.HasPrefix(mediaType, fmt.Sprintf("%s/", m.TypeName))
	}

	return mediaType == fmt.Sprintf("%s/%s", m.TypeName, m.SubtypeName) && equalMaps(params, m.Params)

}

// equalMaps returns whether the two maps are equivalent
func equalMaps(x map[string]string, y map[string]string) bool {

	if x == nil || y == nil {
		return x == nil && y == nil
	}

	if len(x) != len(y) {
		return false
	}

	for k, vx := range x {
		if vy, ok := y[k]; !ok || vy != vx {
			return false
		}
	}

	return true

}

// weight is used to determine the current weight of the mediaRange as defined by the rfc
func (m mediaRange) weight() float64 {

	var (
		weight float64

		// quality is a numeric value between 0.000 and 1.000
		q = math.Min(math.Max(m.Q, 0.0), 1.0)
	)

	// the minimum non-zero value for q is 0.001 as defined by section 5.3.1
	// if q is 0, then its unset so use the default weight
	if q == 0 {
		// quality default weight is 1.000
		q = 1
	}

	if m.TypeName == "*" {
		// its */*, least weighty
		weight = 0.00000
	} else if m.SubtypeName == "*" {
		// its type/*, slightly more weighty than */*
		weight = 0.00001
	} else {
		// its type/subtype, slightly more weighty than type/*
		weight = 0.00002
	}

	// .0001 for each param. the more specific, the more weight.
	return weight + q + (0.00010 * float64(len(m.Params)))

}

type byWeight []mediaRange

// sort.Interface
func (b byWeight) Len() int {
	return len(b)
}

// sort.Interface
func (b byWeight) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// sort.Interface
func (b byWeight) Less(i, j int) bool {
	return b[i].weight() > b[j].weight()
}
