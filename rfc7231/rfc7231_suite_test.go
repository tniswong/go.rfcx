package rfc7231

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRfc7231(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rfc7231 Suite")
}
