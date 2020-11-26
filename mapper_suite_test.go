package batis_test

import (
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestDal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dal Suite")
}
