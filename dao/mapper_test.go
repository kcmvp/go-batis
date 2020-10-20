package dao

import (
	. "github.com/onsi/ginkgo"
)

type Dog struct {
	Name string
}

var _ = Describe("Mapper", func() {

	//var (
	//	createDogMapper Mapper = NewMapper(MapperType, "") {
	//		return Insert, "dog.createDog", Dog{}
	//	}
	//	batchInsert Mapper = func() (MapperType, string, interface{}) {
	//		return Insert, "dog.batchInsert", Dog{}
	//	}
	//
	//)
	//
	//Describe("Test create mapper nodes", func() {
	//	Context("simple insert", func() {
	//		It("should be a simple insert sql getMapper", func() {
	//			node, _ := createDogMapper.getMapper()
	//			Expect(len(node.Attr)).To(Equal(1))
	//			Expect(node.Attr[0]).To(Equal("createDog"))
	//		})
	//
	//})


})
