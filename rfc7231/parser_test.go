package rfc7231

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("parser", func() {

	type parserExample struct {
		in  string
		out []MediaRange
	}

	DescribeTable(
		"parse()",
		func(example parserExample) {

			// given
			rs := strings.NewReader(example.in)
			s := scanner{runeScanner: rs}
			p := parser{scanner: s}

			// when
			result, err := p.parse()

			// then
			Expect(err).To(BeNil())
			Expect(len(result.MediaRanges)).To(Equal(len(example.out)))

			for i, mr := range result.MediaRanges {
				Expect(example.out[i].TypeName).To(Equal(mr.TypeName))
				Expect(example.out[i].SubtypeName).To(Equal(mr.SubtypeName))
				Expect(example.out[i].Params).To(Equal(mr.Params))
				Expect(example.out[i].Q).To(Equal(mr.Q))
			}

		},
		Entry("should parse all these MediaRanges",
			parserExample{
				"text/*, text/html, text/html;level=1, */*",
				[]MediaRange{
					{TypeName: "text", SubtypeName: "*", Params: map[string]string{}, Q: 0},
					{TypeName: "text", SubtypeName: "html", Params: map[string]string{}, Q: 0},
					{TypeName: "text", SubtypeName: "html", Params: map[string]string{"level": "1"}, Q: 0},
					{TypeName: "*", SubtypeName: "*", Params: map[string]string{}, Q: 0},
				},
			},
		),
	)

})
