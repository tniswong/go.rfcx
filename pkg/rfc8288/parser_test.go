package rfc8288_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/tniswong/go.rfcx/pkg/rfc8288"
)

var _ = Describe("parser", func() {

	DescribeTable(
		"Valid Cases",
		func(in string, out Link) {

			// given
			result, err := ParseLink(in)

			// expect
			Expect(err).To(BeNil())
			Expect(out.HREF).To(Equal(result.HREF))

		},
		Entry(
			"href",
			`<https://www.google.com>`,
			Link{
				HREF: URL("https://www.google.com"),
			},
		),
		Entry(
			"href, rel",
			`<https://www.google.com>; rel="next"`,
			Link{
				HREF: URL("https://www.google.com"),
				Rel:  "next",
			},
		),
		Entry(
			"href, rel, hreflang",
			`<https://www.google.com>; rel="next"; hreflang="en"`,
			Link{
				HREF:     URL("https://www.google.com"),
				Rel:      "next",
				HREFLang: "en",
			},
		),
		Entry(
			"href, rel, hreflang, media",
			`<https://www.google.com>; rel="next"; hreflang="en"; media="media"`,
			Link{
				HREF:     URL("https://www.google.com"),
				Rel:      "next",
				HREFLang: "en",
				Media:    "media",
			},
		),
		Entry(
			"href, rel, hreflang, title",
			`<https://www.google.com>; rel="next"; hreflang="en"; title="title"`,
			Link{
				HREF:     URL("https://www.google.com"),
				Rel:      "next",
				HREFLang: "en",
				Title:    "title",
			},
		),
		Entry(
			"href, rel, hreflang, title, title*",
			`<https://www.google.com>; rel="next"; hreflang="en"; title="title"; title*="title*"`,
			Link{
				HREF:      URL("https://www.google.com"),
				Rel:       "next",
				HREFLang:  "en",
				Title:     "title",
				TitleStar: "title*",
			},
		),
		Entry(
			"href, rel, hreflang, title, title*, type",
			`<https://www.google.com>; rel="next"; hreflang="en"; title="title"; title*="title*"; type="type"`,
			Link{
				HREF:      URL("https://www.google.com"),
				Rel:       "next",
				HREFLang:  "en",
				Title:     "title",
				TitleStar: "title*",
				Type:      "type",
			},
		),
	)

	Describe("Parse", func() {

		It("should parse extensions", func() {

			// given
			l := `<https://www.google.com>; extension="value"`

			// when
			w, err := ParseLink(l)

			// then
			Expect(err).To(BeNil())

			value, present := w.Extension("extension")
			Expect(len(w.Extensions())).To(Equal(1))
			Expect(present).To(BeTrue())
			Expect(value).To(Equal("value"))

		})

	})

})
