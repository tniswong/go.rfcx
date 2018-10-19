package rfc7231

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("mediaRange", func() {

	type stringExample struct {
		in  mediaRange
		out string
	}

	DescribeTable(
		"String()",
		func(e stringExample) {

			// when
			result := e.in.String()

			// then
			Expect(result).To(Equal(e.out))

		},
		Entry("Type/Subtype", stringExample{
			in: mediaRange{
				TypeName:    "type",
				SubtypeName: "subtype",
			},
			out: "type/subtype",
		}),

		Entry("Type/Subtype; q=x", stringExample{
			in: mediaRange{
				TypeName:    "type",
				SubtypeName: "subtype",
				Q:           0.5,
			},
			out: "type/subtype; q=0.5",
		}),

		Entry("Type/Subtype; q=x; key=value", stringExample{
			in: mediaRange{
				TypeName:    "type",
				SubtypeName: "subtype",
				Q:           0.5,
				Params:      map[string]string{"key": "value"},
			},
			out: "type/subtype; q=0.5; key=value",
		}),

		Entry("Type/Subtype; key=value", stringExample{
			in: mediaRange{
				TypeName:    "type",
				SubtypeName: "subtype",
				Params:      map[string]string{"key": "value"},
			},
			out: "type/subtype; key=value",
		}),

		Entry("Type/Subtype; q < 0", stringExample{
			in: mediaRange{
				TypeName:    "type",
				SubtypeName: "subtype",
				Q:           -1.0,
			},
			out: "type/subtype",
		}),

		Entry("Type/Subtype; q > 0", stringExample{
			in: mediaRange{
				TypeName:    "type",
				SubtypeName: "subtype",
				Q:           2.0,
			},
			out: "type/subtype",
		}),
	)

	type supportsExample struct {
		mediaRange mediaRange
		mediaType  string

		out bool
	}

	DescribeTable(
		"Supports(mediaType)",
		func(e supportsExample) {

			// given
			m := e.mediaRange

			// when
			result := m.Supports(e.mediaType)

			// then
			Expect(result).To(Equal(e.out))

		},
		Entry("Invalid Type -> */*", supportsExample{
			mediaRange: mediaRange{
				TypeName:    "*",
				SubtypeName: "*",
			},
			mediaType: "!@#$!@#$!@#%!@#%!@#$;;;",
			out:       false,
		}),
		Entry("any/type -> */*", supportsExample{
			mediaRange: mediaRange{
				TypeName:    "*",
				SubtypeName: "*",
			},
			mediaType: "any/type",
			out:       true,
		}),
		Entry("text/type -> text/*", supportsExample{
			mediaRange: mediaRange{
				TypeName:    "text",
				SubtypeName: "*",
			},
			mediaType: "text/type",
			out:       true,
		}),
		Entry("application/type -> text/*", supportsExample{
			mediaRange: mediaRange{
				TypeName:    "text",
				SubtypeName: "*",
			},
			mediaType: "application/type",
			out:       false,
		}),
		Entry("application/type -> text/plain", supportsExample{
			mediaRange: mediaRange{
				TypeName:    "text",
				SubtypeName: "plain",
			},
			mediaType: "application/type",
			out:       false,
		}),
		Entry("text/plain; key=value -> text/plain; key=value", supportsExample{
			mediaRange: mediaRange{
				TypeName:    "text",
				SubtypeName: "plain",
				Params: map[string]string{
					"key": "value",
				},
			},
			mediaType: "text/plain;key=value",
			out:       true,
		}),
	)

})
