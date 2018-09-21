package rfc8288

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"testing"
)

func URL(u string) url.URL {

	defer GinkgoRecover()
	uri, err := url.Parse(u)

	if err != nil {
		Fail(err.Error())
	}

	return *uri

}

func TestRfc8288(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rfc8288 Suite")
}
