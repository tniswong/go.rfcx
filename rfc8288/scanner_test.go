package rfc8288

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("scanner", func() {

	type TokenLiteral struct {
		Token   token
		Literal string
	}

	DescribeTable("scan",
		func(in string, out []TokenLiteral) {

			// given
			r := strings.NewReader(in)
			s := scanner{runeScanner: r}

			x := 0
			for {

				// assert that we haven't scanned more than we expect to
				Expect(x < len(out)).To(BeTrue())

				// when
				token, literal, err := s.Scan()

				// then
				Expect(err).To(BeNil())
				Expect(token).To(Equal(out[x].Token))
				Expect(literal).To(Equal(out[x].Literal))

				x++

				if token == EOF {
					break
				}

			}

			Expect(x).To(Equal(len(out)))

		},
		Entry(
			"example 1",
			`<about:blank>; rel="prev"; title*="title"; media="media"; custom*="custom"`,
			[]TokenLiteral{
				{Token: LT, Literal: "<"},
				{Token: WORD, Literal: "about:blank"},
				{Token: GT, Literal: ">"},
				{Token: SEMICOLON, Literal: ";"},
				{Token: WS, Literal: " "},
				{Token: REL, Literal: "rel"},
				{Token: EQ, Literal: "="},
				{Token: QUOTE, Literal: `"`},
				{Token: WORD, Literal: "prev"},
				{Token: QUOTE, Literal: `"`},
				{Token: SEMICOLON, Literal: ";"},
				{Token: WS, Literal: " "},
				{Token: TITLE, Literal: "title"},
				{Token: STAR, Literal: "*"},
				{Token: EQ, Literal: "="},
				{Token: QUOTE, Literal: `"`},
				{Token: WORD, Literal: "title"},
				{Token: QUOTE, Literal: `"`},
				{Token: SEMICOLON, Literal: ";"},
				{Token: WS, Literal: " "},
				{Token: MEDIA, Literal: "media"},
				{Token: EQ, Literal: "="},
				{Token: QUOTE, Literal: `"`},
				{Token: WORD, Literal: "media"},
				{Token: QUOTE, Literal: `"`},
				{Token: SEMICOLON, Literal: ";"},
				{Token: WS, Literal: " "},
				{Token: WORD, Literal: "custom"},
				{Token: STAR, Literal: "*"},
				{Token: EQ, Literal: "="},
				{Token: QUOTE, Literal: `"`},
				{Token: WORD, Literal: "custom"},
				{Token: QUOTE, Literal: `"`},
				{Token: EOF, Literal: ""},
			},
		),
		Entry(
			"example two. it's ok that this is an invalid link. lexer don't care",
			"<https://www.google.com> asdf",
			[]TokenLiteral{
				{Token: LT, Literal: "<"},
				{Token: WORD, Literal: "https://www.google.com"},
				{Token: GT, Literal: ">"},
				{Token: WS, Literal: " "},
				{Token: WORD, Literal: "asdf"},
				{Token: EOF, Literal: ""},
			},
		),
		Entry(
			"Edge Case: Ends with whitespace",
			" ",
			[]TokenLiteral{
				{Token: WS, Literal: " "},
				{Token: EOF, Literal: ""},
			},
		),
	)

})
