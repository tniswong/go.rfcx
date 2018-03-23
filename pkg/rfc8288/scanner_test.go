package rfc8288_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/tniswong/go.rfcx/pkg/rfc8288"
	"io"
	"strings"
)

var _ = Describe("Scanner", func() {

	type TokenLiteral struct {
		Token   Token
		Literal string
		Err     error
	}

	DescribeTable("Scan",
		func(in string, out []TokenLiteral) {

			// given
			r := strings.NewReader(in)
			s := NewScanner(r)

			for _, t := range out {

				// when
				token, literal, err := s.Scan()

				// then

				if t.Err == nil {
					Expect(err).To(BeNil())
				} else {
					Expect(err).To(Equal(t.Err))
				}

				Expect(token).To(Equal(t.Token))
				Expect(literal).To(Equal(t.Literal))

			}

		},
		Entry(
			"example 1",
			`<about:blank>; rel="prev"; title*="title"; media="media"; custom*="custom"`,
			[]TokenLiteral{
				{Token: LT, Literal: "<", Err: nil},
				{Token: WORD, Literal: "about:blank", Err: nil},
				{Token: GT, Literal: ">", Err: nil},
				{Token: SEMICOLON, Literal: ";", Err: nil},
				{Token: WS, Literal: " ", Err: nil},
				{Token: REL, Literal: "rel", Err: nil},
				{Token: EQ, Literal: "=", Err: nil},
				{Token: QUOTE, Literal: `"`, Err: nil},
				{Token: WORD, Literal: `prev`, Err: nil},
				{Token: QUOTE, Literal: `"`, Err: nil},
				{Token: SEMICOLON, Literal: ";", Err: nil},
				{Token: WS, Literal: " ", Err: nil},
				{Token: TITLE, Literal: "title", Err: nil},
				{Token: STAR, Literal: "*", Err: nil},
				{Token: EQ, Literal: "=", Err: nil},
				{Token: QUOTE, Literal: `"`, Err: nil},
				{Token: WORD, Literal: `title`, Err: nil},
				{Token: QUOTE, Literal: `"`, Err: nil},
				{Token: SEMICOLON, Literal: ";", Err: nil},
				{Token: WS, Literal: " ", Err: nil},
				{Token: MEDIA, Literal: "media", Err: nil},
				{Token: EQ, Literal: "=", Err: nil},
				{Token: QUOTE, Literal: `"`, Err: nil},
				{Token: WORD, Literal: `media`, Err: nil},
				{Token: QUOTE, Literal: `"`, Err: nil},
				{Token: SEMICOLON, Literal: ";", Err: nil},
				{Token: WS, Literal: " ", Err: nil},
				{Token: WORD, Literal: `custom`, Err: nil},
				{Token: STAR, Literal: "*", Err: nil},
				{Token: EQ, Literal: "=", Err: nil},
				{Token: QUOTE, Literal: `"`, Err: nil},
				{Token: WORD, Literal: `custom`, Err: nil},
				{Token: QUOTE, Literal: `"`, Err: nil},
				{Token: EOF, Literal: ``, Err: io.EOF},
			},
		),
		Entry(
			"example two. it's ok that this is an invalid link. lexer don't care",
			`<https://www.google.com> asdf`,
			[]TokenLiteral{
				{Token: LT, Literal: "<", Err: nil},
				{Token: WORD, Literal: "https://www.google.com", Err: nil},
				{Token: GT, Literal: ">", Err: nil},
				{Token: WS, Literal: " ", Err: nil},
				{Token: WORD, Literal: "asdf", Err: nil},
				{Token: EOF, Literal: ``, Err: io.EOF},
			},
		),
	)

})
