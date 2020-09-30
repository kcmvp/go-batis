package sqlx_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoBatis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoBatis Suite")
}
