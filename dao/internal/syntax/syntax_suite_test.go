package syntax_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSyntax(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Syntax Suite")
}
