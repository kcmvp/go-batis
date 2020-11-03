package dao

import (
	"github.com/kcmvp/go-batis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strings"
)


var _ = Describe("Mapper", func() {
	os.Setenv("env", "test")
	batis.DB()
	var (
		createDogMapper Mapper = NewMapper(Insert, "dog.createDog")
	)

	Describe("Test create mapper nodes", func() {
		Context("simple insert", func() {
			It("should be a simple insert sql build", func() {
				m, ok := createDogMapper.(*mapper)
				Expect(ok).To(BeTrue())
				c, err := m.build()
				Expect(err).To(BeNil())
				Expect(strings.ToLower(c.XMLName.Local)).To(Equal("insert"))
				Expect(c.Id).To(Equal("createDog"))
				Expect(c.ResultType).To(Equal("int64"))
				Expect(c.CharData1).To(Equal("insert into Dog(name,age,price) values (#{name},#{age},#{price})"))
				cv,_ := clauseCache.Get(m.mapperName)
				Expect(cv).To(Equal(c))
			})
		})
	})
})
