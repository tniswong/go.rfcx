package rfc8288

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"net/url"
	"reflect"
)

var _ = Describe("Rfc8288", func() {

	Describe("Link", func() {

		DescribeTable(
			"String()",
			func(in Link, out string) {

				// given
				result := in.String()

				// expect
				Expect(result).To(Equal(out))

			},
			Entry(
				"with href",
				Link{
					HREF: URL("https://www.google.com"),
				},
				"<https://www.google.com>",
			),
			Entry(
				"with href, hreflang",
				Link{
					HREF:     URL("https://www.google.com"),
					HREFLang: "en",
				},
				`<https://www.google.com>; hreflang="en"`,
			),
			Entry(
				"with href, rel",
				Link{
					HREF: URL("https://www.google.com"),
					Rel:  "next",
				},
				`<https://www.google.com>; rel="next"`,
			),
			Entry(
				"with href, media",
				Link{
					HREF:  URL("https://www.google.com"),
					Media: "media",
				},
				`<https://www.google.com>; media="media"`,
			),
			Entry(
				"with href, title",
				Link{
					HREF:  URL("https://www.google.com"),
					Title: "title",
				},
				`<https://www.google.com>; title="title"`,
			),
			Entry(
				"with href, title*",
				Link{
					HREF:      URL("https://www.google.com"),
					TitleStar: "title*",
				},
				`<https://www.google.com>; title*="title*"`,
			),
			Entry(
				"with href, type",
				Link{
					HREF: URL("https://www.google.com"),
					Type: "type",
				},
				`<https://www.google.com>; type="type"`,
			),
			Entry(
				"with href, extension",
				func() Link {

					w := Link{
						HREF: URL("https://www.google.com"),
					}

					w.Extend("extension", "value")

					return w
				}(),
				`<https://www.google.com>; extension="value"`,
			),
		)

		DescribeTable(
			"MarshalJSON()",
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

				for _, key := range in.ExtensionKeys() {
					Expect(result[key]).To(Equal(out[key]))
				}

			},
			Entry(
				"should marshal all fields and extensions",
				func() Link {

					l := Link{
						HREF:      URL("https://www.google.com"),
						Rel:       "rel",
						HREFLang:  "hreflang",
						Media:     "media",
						Title:     "title",
						TitleStar: "title*",
						Type:      "type",
					}

					l.Extend("extension", "value")

					return l

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
			"UnmarshalJSON()",
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
				"should unmarshal all fields and extensions",
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

					l := Link{
						HREF:      URL("https://www.google.com"),
						Rel:       "rel",
						HREFLang:  "hreflang",
						Media:     "media",
						Title:     "title",
						TitleStar: "title*",
						Type:      "type",
					}

					l.Extend("extension", "value")

					return l

				}(),
			),
		)

		DescribeTable(
			"UnmarshalJSON() error cases",
			func(in string, out error) {

				// given
				l := Link{}

				// when
				err := json.Unmarshal([]byte(in), &l)

				// then
				Expect(err).To(Equal(out))

			},
			Entry(
				"should return json.UnmarshalTypeError describing href field",
				`{
                    "href": "Not a valid url !@#$%^&*()_+"
                }`,
				&json.UnmarshalTypeError{
					Value:  "uri",
					Type:   reflect.TypeOf(url.URL{}),
					Field:  "href",
					Struct: "Link",
				},
			),
			Entry(
				"should return json.UnmarshalTypeError describing href field",
				`{
                    "href": false
                }`,
				&json.UnmarshalTypeError{
					Value:  "uri",
					Type:   reflect.TypeOf(""),
					Field:  "href",
					Struct: "Link",
				},
			),
			Entry(
				"should return json.UnmarshalTypeError describing rel field",
				`{
                    "rel": false
                }`,
				&json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "rel",
					Struct: "Link",
				},
			),
			Entry(
				"should return json.UnmarshalTypeError describing hreflang field",
				`{
                    "hreflang": false
                }`,
				&json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(int64(0)),
					Field:  "hreflang",
					Struct: "Link",
				},
			),
			Entry(
				"should return json.UnmarshalTypeError describing media field",
				`{
                    "media": false
                }`,
				&json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "media",
					Struct: "Link",
				},
			),
			Entry(
				"should return json.UnmarshalTypeError describing title field",
				`{
                    "title": false
                }`,
				&json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "title",
					Struct: "Link",
				},
			),
			Entry(
				"should return json.UnmarshalTypeError describing title* field",
				`{
                    "title*": false
                }`,
				&json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "title*",
					Struct: "Link",
				},
			),
			Entry(
				"should return json.UnmarshalTypeError describing type field",
				`{
                    "type": false
                }`,
				&json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(""),
					Field:  "type",
					Struct: "Link",
				},
			),
		)

		Describe("Extend(key, value)", func() {

			It("should make the value accessible via Extension(key)", func() {

				// given
				l := Link{}
				key := "extension"
				value := "value"

				// when
				l.Extend(key, value)
				result, ok := l.Extension(key)

				// then
				Expect(ok).To(BeTrue())
				Expect(result).To(Equal(value))

			})

			It("should delete the extension if assigned a nil value", func() {

				// given
				l := Link{}
				key := "extension"
				value := "value"

				// when
				l.Extend(key, value)
				l.Extend(key, nil)
				result, ok := l.Extension(key)

				// then
				Expect(ok).To(BeFalse())
				Expect(result).To(BeNil())

			})

			It("should return ErrExtensionKeyIsReserved if attempting to Extend with reserved key", func() {

				// given
				l := Link{}

				// expect
				for reservedKey := range ReservedKeys {

					err := l.Extend(reservedKey, "any")

					Expect(err).To(Equal(ErrExtensionKeyIsReserved))

				}

			})

		})

	})

})
