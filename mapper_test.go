package sqlx

import (
	"testing"
)

type Dog struct {
	Name string
}

var createDogMapper Mapper = func() (MapperType, string, interface{}) {
	return CREATE, "createDog", Dog{}
}

var findDogMapper Mapper = func() (MapperType, string, interface{}) {
	return FIND, "findDogById", Dog{}
}

var updateDogMapper Mapper = func() (MapperType, string, interface{}) {
	return UPDATE, "updateDogAge", Dog{}
}

var searchDogMapper Mapper = func() (MapperType, string, interface{}) {
	return SEARCH, "searchDogBySample", Dog{}
}




func TestName(t *testing.T) {
	createDogMapper.Exec(Dog{})
}

