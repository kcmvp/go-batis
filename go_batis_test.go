package batis_test

import (
	"github.com/kcmvp/go-batis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)


var _ = Describe("GoBatis", func() {
	os.Setenv("env","test")
	db := batis.DB()
	Describe("Test system initialization", func() {
		Context("configuration merge", func() {
			When("there is a value in the env", func() {
				It("should merge default and test environment", func() {
					Expect(batis.Config.Driver()).Should(Equal("sqlite3"))
					Expect(batis.Config.DBUrl()).Should(Equal("./testdb"))
					Expect(batis.Config.MaxOpen()).Should(Equal(20))
					Expect(db).ShouldNot(BeNil())
				})
			})
		})
	})
})