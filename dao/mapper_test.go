package dao

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

type Dog struct {
	Name string
}

var _ = Describe("Mapper", func() {

	var (
		createDogMapper Mapper = func() (MapperType, string, interface{}) {
			return InsertMapper, "dog.createDog", Dog{}
		}
		batchInsert Mapper = func() (MapperType, string, interface{}) {
			return InsertMapper, "dog.batchInsert", Dog{}
		}

	)

	Describe("Test create mapper nodes", func() {
		Context("simple insert", func() {
			It("should be a simple insert sql clause", func() {
				node, _ := createDogMapper.clause()
				Expect(len(node.Attr)).To(Equal(1))
				Expect(node.Attr[0]).To(Equal("createDog"))
			})
		})

		Context("batch insert", func() {
			It("should be a batch insert sql clause", func() {
				node, _ := batchInsert.clause()
				Expect(len(node.Attr)).To(Equal(1))
				Expect(node.Attr[0]).To(Equal("batchInsert"))
				Expect(strings.Contains(node.OutputXML(true),"foreach collection")).To(BeTrue())
			})
		})
	})


})
