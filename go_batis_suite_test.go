package batis_test

import (
	"testing"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


func TestGoBatis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoBatis Suite")
}
