package rfc7231

import (
	"fmt"
	"math"
	"strings"
)

// MediaRange represents a MediaRange as defined for use in an HTTP Accept Header as defined in RFC 7231
type MediaRange struct {
	TypeName    string
	SubtypeName string
	Params      map[string]string
	Q           float64
}

// String returns the string representation of the MediaRange
func (m MediaRange) String() string {

	result := fmt.Sprintf("%s/%s", m.TypeName, m.SubtypeName)

	if m.Q > 0.0 {
		result = fmt.Sprintf("%s; q=%.1f", result, m.Q)
	}

	for k, v := range m.Params {

		if !strings.EqualFold(k, "q") {
			result = fmt.Sprintf("%s; %s=%s", result, k, v)
		}

	}

	return result

}

func (m MediaRange) weight() float64 {

	var (
		weight float64

		// quality default is 1, unless overridden
		q = math.Min(math.Max(m.Q, 0.0), 1.0)
	)

	// minimum non-zero value for q is 0.001 as defined by section 5.3.1

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

type byWeight []MediaRange

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
