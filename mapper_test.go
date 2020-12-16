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
