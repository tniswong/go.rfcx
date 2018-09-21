package rfc8288

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("parser", func() {

	DescribeTable(
		"parse()",
		func(in string, out Link) {

			// when
			result, err := ParseLink(in)

			// expect
			Expect(err).To(BeNil())
			Expect(out).To(Equal(result))

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
		Entry(
			"href, rel, hreflang, title, title*, type, extensions",
			`<https://www.google.com>; rel="next"; hreflang="en"; title="title"; title*="title*"; type="type"; extension="value"`,
			Link{
				HREF:          URL("https://www.google.com"),
				Rel:           "next",
				HREFLang:      "en",
				Title:         "title",
				TitleStar:     "title*",
				Type:          "type",
				extensionKeys: []string{"extension"},
				extensions: map[string]interface{}{
					"extension": "value",
				},
			},
		),
	)

})
