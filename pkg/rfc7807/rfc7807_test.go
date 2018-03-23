package rfc7807_test

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/tniswong/go.rfcx/pkg/rfc7807"
    "net/url"
)

var _ = Describe("Rfc7807", func() {

	DescribeTable(
		"MarshalJSON",
		func(in Problem, out map[string]interface{}) {

			// given
			jsonBytes, err := json.Marshal(in)
			fmt.Println(string(jsonBytes))
			result := make(map[string]interface{})

			Expect(err).To(BeNil())

			// when
			json.Unmarshal(jsonBytes, &result)

			// then
			if in.Type != "" {
				Expect(result["type"]).To(Equal(out["type"]))
			}

			if in.Title != "" {
				Expect(result["title"]).To(Equal(out["title"]))
			}

			if in.Status != 0 {
				Expect(result["status"]).To(BeEquivalentTo(out["status"]))
			}

			if in.Detail != "" {
				Expect(result["detail"]).To(Equal(out["detail"]))
			}

			var zero url.URL
			if in.Instance != zero {
				Expect(result["instance"]).To(Equal(out["instance"]))
			}

			for _, key := range in.Extensions() {
				Expect(result[key]).To(Equal(out[key]))
			}

		},
		Entry(
			"type",
			Problem{
				Type: "type",
			},
			map[string]interface{}{
				"type": "type",
			},
		),
		Entry(
			"title",
			Problem{
				Title: "title",
			},
			map[string]interface{}{
				"title": "title",
			},
		),
		Entry(
			"status",
			Problem{
				Status: 500,
			},
			map[string]interface{}{
				"status": 500,
			},
		),
		Entry(
			"detail",
			Problem{
				Detail: "detail",
			},
			map[string]interface{}{
				"detail": "detail",
			},
		),
		Entry(
			"instance",
			Problem{
				Instance: URL("about:blank"),
			},
			map[string]interface{}{
				"instance": "about:blank",
			},
		),
		Entry(
			"all basic, no extensions",
			Problem{
				Type:     "type",
				Title:    "title",
				Status:   500,
				Detail:   "detail",
				Instance: URL("about:blank"),
			},
			map[string]interface{}{
				"type":     "type",
				"title":    "title",
				"status":   500,
				"detail":   "detail",
				"instance": "about:blank",
			},
		),
		Entry(
			"extension",
			func() Problem {
				p := Problem{}
				p.Extend("extension", "extension")
				return p
			}(),
			map[string]interface{}{
				"extension": "extension",
			},
		),
		Entry(
			"all basic, with extensions",
			func() Problem {
				p := Problem{
					Type:     "type",
					Title:    "title",
					Status:   500,
					Detail:   "detail",
					Instance: URL("about:blank"),
				}
				p.Extend("extension", "extension")
				return p
			}(),
			map[string]interface{}{
				"type":      "type",
				"title":     "title",
				"status":    500,
				"detail":    "detail",
				"instance":  "about:blank",
				"extension": "extension",
			},
		),
	)

	DescribeTable(
		"UnmarshalJSON",
		func(in map[string]interface{}, out Problem) {

			// given
			jsonBytes, err := json.Marshal(in)
			result := Problem{}
			Expect(err).To(BeNil())

			// when
			json.Unmarshal(jsonBytes, &result)

			// then
			if _, ok := in["type"]; ok {
				Expect(result.Type).To(Equal(out.Type))
			}

			if _, ok := in["title"]; ok {
				Expect(result.Title).To(Equal(out.Title))
			}

			if _, ok := in["status"]; ok {
				Expect(result.Status).To(Equal(out.Status))
			}

			if _, ok := in["detail"]; ok {
				Expect(result.Detail).To(Equal(out.Detail))
			}

			if _, ok := in["instance"]; ok {
				Expect(result.Instance.String()).To(Equal(out.Instance.String()))
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
				"type": "type",
			},
			Problem{
				Type: "type",
			},
		),
		Entry(
			"title",
			map[string]interface{}{
				"title": "title",
			},
			Problem{
				Title: "title",
			},
		),
		Entry(
			"status",
			map[string]interface{}{
				"status": 500,
			},
			Problem{
				Status: 500,
			},
		),
		Entry(
			"detail",
			map[string]interface{}{
				"detail": "detail",
			},
			Problem{
				Detail: "detail",
			},
		),
		Entry(
			"instance",
			map[string]interface{}{
				"instance": "about:blank",
			},
			Problem{
				Instance: URL("about:blank"),
			},
		),
		Entry(
			"all basic, no extensions",
			map[string]interface{}{
				"type":     "type",
				"title":    "title",
				"status":   500,
				"detail":   "detail",
				"instance": "about:blank",
			},
			Problem{
				Type:     "type",
				Title:    "title",
				Status:   500,
				Detail:   "detail",
				Instance: URL("about:blank"),
			},
		),
		Entry(
			"extension",
			map[string]interface{}{
				"extension": "extension",
			},
			func() Problem {
				p := Problem{}
				p.Extend("extension", "extension")
				return p
			}(),
		),
		Entry(
			"all basic, with extensions",
			map[string]interface{}{
				"type":      "type",
				"title":     "title",
				"status":    500,
				"detail":    "detail",
				"instance":  "about:blank",
				"extension": "extension",
			},
			func() Problem {
				p := Problem{
					Type:     "type",
					Title:    "title",
					Status:   500,
					Detail:   "detail",
					Instance: URL("about:blank"),
				}
				p.Extend("extension", "extension")
				return p
			}(),
		),
	)

})
