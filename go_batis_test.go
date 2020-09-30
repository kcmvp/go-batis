package sqlx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("GoBatis", func() {
	Expect(strings.Contains("Ginkgo is awesome", "is")).To(BeTrue())
})
