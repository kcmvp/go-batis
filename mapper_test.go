package batis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var dogMapP = map[string]interface{}{
	"id":   1234567,
	"name": "hello",
	"age":  12,
}
var dogMapF = map[string]interface{}{
	"id":    1234567,
	"name":  "hello",
	"age":   12,
	"price": 136,
}


// test data for foreach
type Dog struct {
	Name string
	Age int
	Price int
}
var dogList = []Dog{
	{
		Name: "mimi",
		Age: 1,
		Price: 2,
	},{
		Name: "kaka",
		Age: 2,
		Price: 3,
	},
}
var yellowDog = map[string]interface{}{
	"color": "yellow",
	"dogList":dogList,
}
type ColoredDog struct {
	Color string
	DogList []Dog
}
var redDog = ColoredDog{
	Color: "red",
	DogList: dogList,
}




var mappers = []struct {
	desc      string
	mapperId  string
	argMap    map[string]interface{}
	positive  bool
	msg    string
}{
	{"case1: miss parameters", "dog.createDog", nil,false, "mapper#createDog: failed to resolve the expression: #{name}."},
	{"case2: simple create", "dog.createDog", dogMapF,true, "insert into Dog(name,age,price) values (?,?,?)"},
	{"case3: with partial parameters", "dog.updateDog", dogMapP,true ,"update Dog set name = ?, age = ?, updated_at = CURRENT_TIMESTAMP() where id = ?"},
	{"case4: with all parameters", "dog.updateDog", dogMapF,true, "update Dog set name = ?, age = ?, price = ?, updated_at = CURRENT_TIMESTAMP() where id = ?"},
	{"case5: miss cache name", "dog.findDogByIdNoCacheName", dogMapF,false, "mapper#findDogByIdNoCacheName: empty cache name or key"},
	{"case6: simple find clause", "dog.findDogById", dogMapF, true, "select * from Dog where id = ?"},
	{"case6: dynamic where", "dog.searchByExample", dogMapF, true, "select count(1) from Dog where name = ? and age = ?"},
	{"case7: delete statement missed parameter", "dog.deleteDogById", dogMapF,false, "mapper#deleteDogById: failed to resolve the expression: #{levySerialNumber}"},
	{"case8: invalid placeholder character", "dog.updateWeekDayPriceCase1", dogMapF,true, "UPDATE T_WEEK_DAY_PRICE set PRICE = ${price}, UPDATED_AT = CURRENT_TIMESTAMP() where PRICE_PLAN_ID = ${age} and NUM_OF_WEEK = ${id}"},
	{"case9: simple sql ref", "dog.selectByRef", dogMapF,true, "SELECT name, age, size FROM UUC_COMPANY where id = ?"},
	{"case10: nest sql ref", "dog.selectByRefNest", dogMapF,true, "SELECT f.ID , f.PROCESS_KEY , f.PROCESS_NAME , f.MODULE_CODE , m.MODULE_NAME , f.NOTE , f.STATUS FROM UBPC_PROCESS_FILE f LEFT JOIN UBPC_MODULE m ON m.MODULE_CODE = f.MODULE_CODE where f.MODULE_CODE = ? AND (f.PROCESS_NAME LIKE CONCAT('%',?,'%') OR f.PROCESS_KEY LIKE CONCAT('%',?,'%')) AND f.age = ?"},
	{"case11: sql with escape", "dog.findDogByIdEscape", dogMapF,true, "select * from Dog where id <= ? and price >= 100"},
}

var mapDir = "./mapper"

func TestMapperBuildCharData(t *testing.T) {
	assert := assert.New(t)
	for _, m := range mappers {
		t.Run(m.mapperId, func(t *testing.T) {
			mapper := SqlMapper(m.mapperId)
			clause, err := mapper.build(mapDir, m.argMap)
			if !m.positive {
				assert.NotNil(err, m.desc)
				assert.Contains(err.Error(), m.msg, m.desc)
			} else {
				assert.Nil(err)
				vs := make([]interface{}, 0, len(m.argMap))
				for _, value := range m.argMap {
					vs = append(vs, value)
				}
				assert.Equal(m.msg, clause.statement, m.mapperId, m.desc)
				for _, v := range clause.sqlParams {
					assert.True(func(i interface{}) bool {
						for _, param := range m.argMap {
							if param == v {
								return true
							}
						}
						return false
					}(v), m.desc)
				}
			}
		})
	}
}
