package rfc8288_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/tniswong/go.rfcx/pkg/rfc8288"
	"net/url"
)

var _ = Describe("Rfc8288", func() {

	type StringExample struct {
		Input  Link
		Output string
	}

	DescribeTable(
		"String",
		func(example StringExample) {

			// given
			result := example.Input.String()

			// expect
			Expect(result).To(Equal(example.Output))

		},
		Entry(
			"href",
			StringExample{
				Link{
					HREF: URL("https://www.google.com"),
				},
				"<https://www.google.com>",
			},
		),
		Entry(
			"href, hreflang",
			StringExample{
				Link{
					HREF:     URL("https://www.google.com"),
					HREFLang: "en",
				},
				`<https://www.google.com>; hreflang="en"`,
			},
		),
		Entry(
			"href, rel",
			StringExample{
				Link{
					HREF: URL("https://www.google.com"),
					Rel:  "next",
				},
				`<https://www.google.com>; rel="next"`,
			},
		),
		Entry(
			"href, media",
			StringExample{
				Link{
					HREF:  URL("https://www.google.com"),
					Media: "media",
				},
				`<https://www.google.com>; media="media"`,
			},
		),
		Entry(
			"href, title",
			StringExample{
				Link{
					HREF:  URL("https://www.google.com"),
					Title: "title",
				},
				`<https://www.google.com>; title="title"`,
			},
		),
		Entry(
			"href, title*",
			StringExample{
				Link{
					HREF:      URL("https://www.google.com"),
					TitleStar: "title*",
				},
				`<https://www.google.com>; title*="title*"`,
			},
		),
		Entry(
			"href, type",
			StringExample{
				Link{
					HREF: URL("https://www.google.com"),
					Type: "type",
				},
				`<https://www.google.com>; type="type"`,
			},
		),
		Entry(
			"href, extension",
			StringExample{
				func() Link {

					w := Link{
						HREF: URL("https://www.google.com"),
					}

					w.Extend("extension", "value")

					return w
				}(),
				`<https://www.google.com>; extension="value"`,
			},
		),
	)

	DescribeTable(
		"MarshalJSON",
		func(in Link, out map[string]interface{}) {

			// given
			jsonBytes, err := json.Marshal(in)
			result := make(map[string]interface{})

			Expect(err).To(BeNil())

			// when
			json.Unmarshal(jsonBytes, &result)

			// then
			var zero url.URL
			if in.HREF != zero {
				Expect(result["href"]).To(Equal(out["href"]))
			}

			if in.Rel != "" {
				Expect(result["rel"]).To(Equal(out["rel"]))
			}

			if in.HREFLang != "" {
				Expect(result["hreflang"]).To(Equal(out["hreflang"]))
			}

			if in.Media != "" {
				Expect(result["media"]).To(Equal(out["media"]))
			}

			if in.Title != "" {
				Expect(result["title"]).To(Equal(out["title"]))
			}

			if in.TitleStar != "" {
				Expect(result["title*"]).To(Equal(out["title*"]))
			}

			if in.Type != "" {
				Expect(result["type"]).To(Equal(out["type"]))
			}

			for _, key := range in.Extensions() {
				Expect(result[key]).To(Equal(out[key]))
			}

		},
		Entry(
			"href",
			func() Link {

				w := Link{
					HREF:      URL("https://www.google.com"),
					Rel:       "rel",
					HREFLang:  "hreflang",
					Media:     "media",
					Title:     "title",
					TitleStar: "title*",
					Type:      "type",
				}

				w.Extend("extension", "value")

				return w

			}(),
			map[string]interface{}{
				"href":      "https://www.google.com",
				"rel":       "rel",
				"hreflang":  "hreflang",
				"media":     "media",
				"title":     "title",
				"title*":    "title*",
				"type":      "type",
				"extension": "value",
			},
		),
	)

	DescribeTable(
		"UnmarshalJSON",
		func(in map[string]interface{}, out Link) {

			// given
			jsonBytes, err := json.Marshal(in)
			result := Link{}
			Expect(err).To(BeNil())

			// when
			json.Unmarshal(jsonBytes, &result)

			// then
			if _, ok := in["href"]; ok {
				Expect(result.HREF.String()).To(Equal(out.HREF.String()))
			}

			if _, ok := in["rel"]; ok {
				Expect(result.Rel).To(Equal(out.Rel))
			}

			if _, ok := in["hreflang"]; ok {
				Expect(result.HREFLang).To(Equal(out.HREFLang))
			}

			if _, ok := in["media"]; ok {
				Expect(result.Media).To(Equal(out.Media))
			}

			if _, ok := in["title"]; ok {
				Expect(result.Title).To(Equal(out.Title))
			}

			if _, ok := in["title*"]; ok {
				Expect(result.TitleStar).To(Equal(out.TitleStar))
			}

			if _, ok := in["type"]; ok {
				Expect(result.Type).To(Equal(out.Type))
			}

			for key := range in {

				if _, isReserved := ReservedKeys[key]; isReserved {
					continue
				}

				resultValue, resultExists := result.Extension(key)
				Expect(resultExists).To(BeTrue())

				outValue, valueExists := out.Extension(key)
				Expect(valueExists).To(BeTrue())

				Expect(resultValue).To(Equal(outValue))

			}

		},
		Entry(
			"type",
			map[string]interface{}{
				"href":      "https://www.google.com",
				"rel":       "rel",
				"hreflang":  "hreflang",
				"media":     "media",
				"title":     "title",
				"title*":    "title*",
				"type":      "type",
				"extension": "value",
			},
			func() Link {

				w := Link{
					HREF:      URL("https://www.google.com"),
					Rel:       "rel",
					HREFLang:  "hreflang",
					Media:     "media",
					Title:     "title",
					TitleStar: "title*",
					Type:      "type",
				}

				w.Extend("extension", "value")

				return w

			}(),
		),
	)

})
