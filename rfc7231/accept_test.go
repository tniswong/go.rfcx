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

    DescribeTable(
        "ParseAccept()",
        func(example parseAcceptExample) {

            // when
            result, err := ParseAccept(example.in)

            // then
            Expect(err).To(BeNil())
            Expect(result).To(Equal(example.out))

        },
        Entry(
            "example 1",
            parseAcceptExample{
                "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
                Accept{
                    []MediaRange{
                        {TypeName: "text", SubtypeName: "plain", Params: map[string]string{}, Q: 0.5},
                        {TypeName: "text", SubtypeName: "html", Params: map[string]string{}, Q: 0},
                        {TypeName: "text", SubtypeName: "x-dvi", Params: map[string]string{}, Q: 0.8},
                        {TypeName: "text", SubtypeName: "x-c", Params: map[string]string{}, Q: 0},
                    },
                },
            },
        ),
        Entry(
            "example 2",
            parseAcceptExample{
                "text/*, text/html, text/html; level=1, */*",
                Accept{
                    []MediaRange{
                        {TypeName: "text", SubtypeName: "*", Params: map[string]string{}, Q: 0},
                        {TypeName: "text", SubtypeName: "html", Params: map[string]string{}, Q: 0},
                        {TypeName: "text", SubtypeName: "html", Params: map[string]string{"level": "1"}, Q: 0},
                        {TypeName: "*", SubtypeName: "*", Params: map[string]string{}, Q: 0},
                    },
                },
            },
        ),
    )

	type stringExample struct {
        in Accept
        out  string
    }

    DescribeTable(
        "String()",
        func(example stringExample) {

            // when
            result := example.in.String()

            // then
            Expect(result).To(Equal(example.out))

        },
        Entry(
            "example 1",
            stringExample{
                Accept{
                    []MediaRange{
                        {TypeName: "text", SubtypeName: "plain", Params: map[string]string{}, Q: 0.5},
                        {TypeName: "text", SubtypeName: "html", Params: map[string]string{}, Q: 0},
                        {TypeName: "text", SubtypeName: "x-dvi", Params: map[string]string{}, Q: 0.8},
                        {TypeName: "text", SubtypeName: "x-c", Params: map[string]string{}, Q: 0},
                    },
                },
                "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
            },
        ),
        Entry(
            "example 2",
            stringExample{
                Accept{
                    []MediaRange{
                        {TypeName: "text", SubtypeName: "*", Params: map[string]string{}, Q: 0},
                        {TypeName: "text", SubtypeName: "html", Params: map[string]string{}, Q: 0},
                        {TypeName: "text", SubtypeName: "html", Params: map[string]string{"level": "1"}, Q: 0},
                        {TypeName: "*", SubtypeName: "*", Params: map[string]string{}, Q: 0},
                    },
                },
                "text/*, text/html, text/html; level=1, */*",
            },
        ),
    )

})
