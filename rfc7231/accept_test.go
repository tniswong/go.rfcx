package rfc7231

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Accept", func() {

	type parseAcceptExample struct {
		in  string
		out Accept
	}

	DescribeTable("ParseAccept()",
		func(example parseAcceptExample) {

			// when
			result, err := ParseAccept(example.in)

			// then
			Expect(err).To(BeNil())
			Expect(result).To(Equal(example.out))

		},
		Entry("example 1", parseAcceptExample{
			in: "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
			out: Accept{
				mediaRanges: []mediaRange{
					{TypeName: "text", SubtypeName: "plain", Params: map[string]string{}, Q: 0.5},
					{TypeName: "text", SubtypeName: "html", Params: map[string]string{}, Q: 0},
					{TypeName: "text", SubtypeName: "x-dvi", Params: map[string]string{}, Q: 0.8},
					{TypeName: "text", SubtypeName: "x-c", Params: map[string]string{}, Q: 0},
				},
			}},
		),
		Entry("example 2", parseAcceptExample{
			in: "text/*, text/html, text/html; level=1, */*",
			out: Accept{
				mediaRanges: []mediaRange{
					{TypeName: "text", SubtypeName: "*", Params: map[string]string{}, Q: 0},
					{TypeName: "text", SubtypeName: "html", Params: map[string]string{}, Q: 0},
					{TypeName: "text", SubtypeName: "html", Params: map[string]string{"level": "1"}, Q: 0},
					{TypeName: "*", SubtypeName: "*", Params: map[string]string{}, Q: 0},
				},
			},
		}),
	)

	type mostAcceptableExample struct {
		header     string
		mediaTypes []string

		result string
		ok     bool
	}

	DescribeTable("MostAcceptable(mediaTypes)",
		func(e mostAcceptableExample) {

			// given
			accept, err := ParseAccept(e.header)
			Expect(err).To(BeNil())

			// when
			result, ok := accept.MostAcceptable(e.mediaTypes)

			// then
			Expect(ok).To(Equal(e.ok))
			Expect(result).To(Equal(e.result))

		},
		Entry("Example 1", mostAcceptableExample{
			header:     "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
			mediaTypes: []string{"text/plain", "text/x-dvi", "text/html"},
			result:     "text/html",
			ok:         true,
		}),
		Entry("Example 2", mostAcceptableExample{
			header:     "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
			mediaTypes: []string{"text/plain", "text/x-dvi"},
			result:     "text/x-dvi",
			ok:         true,
		}),
		Entry("Example 3", mostAcceptableExample{
			header:     "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
			mediaTypes: []string{"application/json"},
			result:     "",
			ok:         false,
		}),
	)

	type acceptableExample struct {
		header    string
		mediaType string
		ok        bool
	}

	DescribeTable("Acceptable(mediaType)",
		func(e acceptableExample) {

			// given
			accept, err := ParseAccept(e.header)
			Expect(err).To(BeNil())

			// when
			ok := accept.Acceptable(e.mediaType)

			// then
			Expect(ok).To(Equal(e.ok))

		},
		Entry("Example 1", acceptableExample{
			header:    "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
			mediaType: "text/plain",
			ok:        true,
		}),
		Entry("Example 2", acceptableExample{
			header:    "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
			mediaType: "text/plain",
			ok:        true,
		}),
		Entry("Example 3", acceptableExample{
			header:    "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
			mediaType: "application/json",
			ok:        false,
		}),
	)

	Describe("Acceptable() given a zero-value Accept{}", func() {

		It("should return true", func() {

			// given
			accept := Accept{}

			// when
			result := accept.Acceptable("any/type")

			// then
			Expect(result).To(BeTrue())

		})

	})

	type stringExample struct {
		in  Accept
		out string
	}

	DescribeTable("String()",
		func(example stringExample) {

			// when
			result := example.in.String()

			// then
			Expect(result).To(Equal(example.out))

		},
		Entry("example 1", stringExample{
			in: Accept{
				mediaRanges: []mediaRange{
					{TypeName: "text", SubtypeName: "plain", Params: map[string]string{}, Q: 0.5},
					{TypeName: "text", SubtypeName: "html", Params: map[string]string{}, Q: 0},
					{TypeName: "text", SubtypeName: "x-dvi", Params: map[string]string{}, Q: 0.8},
					{TypeName: "text", SubtypeName: "x-c", Params: map[string]string{}, Q: 0},
				},
			},
			out: "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
		}),
		Entry("example 2", stringExample{
			in: Accept{
				mediaRanges: []mediaRange{
					{TypeName: "text", SubtypeName: "*", Params: map[string]string{}, Q: 0},
					{TypeName: "text", SubtypeName: "html", Params: map[string]string{}, Q: 0},
					{TypeName: "text", SubtypeName: "html", Params: map[string]string{"level": "1"}, Q: 0},
					{TypeName: "*", SubtypeName: "*", Params: map[string]string{}, Q: 0},
				},
			},
			out: "text/*, text/html, text/html; level=1, */*",
		}),
		Entry("Zero-value Accept == */*", stringExample{
			in:  Accept{},
			out: "*/*",
		}),
	)

})
