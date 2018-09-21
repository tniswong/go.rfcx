package rfc7231

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("scanner", func() {

	type scanExample struct {
	    in string
	    out []struct {
            Token   token
            Literal string
        }
    }

	DescribeTable(
		"scan()",
		func(example scanExample) {

			// given
			rs := strings.NewReader(example.in)
			s := scanner{runeScanner: rs}

			x := 0
			for {

				// assert that we haven't scanned more than we expect to
				Expect(x < len(example.out)).To(BeTrue())

				// when
				token, literal, err := s.scan()

				// then
				Expect(err).To(BeNil())
				Expect(token).To(Equal(example.out[x].Token))
				Expect(literal).To(Equal(example.out[x].Literal))

				x++

				if token == EOF {
					break
				}

			}

			Expect(x).To(Equal(len(example.out)))

		},
		Entry(
			"example 1",
			scanExample{
                `text/plain; Q=0.5, text/html, text/x-dvi; Q=0.8, text/x-c`,
                []struct {
                    Token   token
                    Literal string
                }{
                    {Token: WORD, Literal: "text"},
                    {Token: SLASH, Literal: "/"},
                    {Token: WORD, Literal: "plain"},
                    {Token: SEMICOLON, Literal: ";"},
                    {Token: WS, Literal: " "},
                    {Token: WORD, Literal: "Q"},
                    {Token: EQ, Literal: "="},
                    {Token: WORD, Literal: "0.5"},
                    {Token: COMMA, Literal: ","},
                    {Token: WS, Literal: " "},
                    {Token: WORD, Literal: "text"},
                    {Token: SLASH, Literal: "/"},
                    {Token: WORD, Literal: "html"},
                    {Token: COMMA, Literal: ","},
                    {Token: WS, Literal: " "},
                    {Token: WORD, Literal: "text"},
                    {Token: SLASH, Literal: "/"},
                    {Token: WORD, Literal: "x-dvi"},
                    {Token: SEMICOLON, Literal: ";"},
                    {Token: WS, Literal: " "},
                    {Token: WORD, Literal: "Q"},
                    {Token: EQ, Literal: "="},
                    {Token: WORD, Literal: "0.8"},
                    {Token: COMMA, Literal: ","},
                    {Token: WS, Literal: " "},
                    {Token: WORD, Literal: "text"},
                    {Token: SLASH, Literal: "/"},
                    {Token: WORD, Literal: "x-c"},
                    {Token: EOF, Literal: ""},
                },
            },
		),
		Entry(
			"example 2",
			scanExample{
                `text/*, text/plain, text/plain;format=flowed, */*`,
                []struct {
                    Token   token
                    Literal string
                }{
                    {Token: WORD, Literal: "text"},
                    {Token: SLASH, Literal: "/"},
                    {Token: WORD, Literal: "*"},
                    {Token: COMMA, Literal: ","},
                    {Token: WS, Literal: " "},
                    {Token: WORD, Literal: "text"},
                    {Token: SLASH, Literal: "/"},
                    {Token: WORD, Literal: "plain"},
                    {Token: COMMA, Literal: ","},
                    {Token: WS, Literal: " "},
                    {Token: WORD, Literal: "text"},
                    {Token: SLASH, Literal: "/"},
                    {Token: WORD, Literal: "plain"},
                    {Token: SEMICOLON, Literal: ";"},
                    {Token: WORD, Literal: "format"},
                    {Token: EQ, Literal: "="},
                    {Token: WORD, Literal: "flowed"},
                    {Token: COMMA, Literal: ","},
                    {Token: WS, Literal: " "},
                    {Token: WORD, Literal: "*"},
                    {Token: SLASH, Literal: "/"},
                    {Token: WORD, Literal: "*"},
                    {Token: EOF, Literal: ""},
                },
            },
		),
	)

})
