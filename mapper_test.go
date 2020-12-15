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
	errMsg    string
	statement string
}{
	{"case1: miss parameters", "dog.createDog", nil, "mapper#createDog: failed to resolve the expression: #{name}.", ""},
	{"case2: simple create", "dog.createDog", dogMapF, "", "insert into Dog(name,age,price) values (?,?,?)"},
	{"case3: with partial parameters", "dog.updateDog", dogMapP, "", "update Dog set name = ?, age = ?, updated_at = CURRENT_TIMESTAMP() where id = ?"},
	{"case4: with all parameters", "dog.updateDog", dogMapF, "", "update Dog set name = ?, age = ?, price = ?, updated_at = CURRENT_TIMESTAMP() where id = ?"},
	{"case5: miss cache name", "dog.findDogById", dogMapF, "mapper#findDogById: empty cache name or key", ""},
}

var mapDir = "./mapper"

func TestMapperBuildCharData(t *testing.T) {
	assert := assert.New(t)
	for _, m := range mappers {
		t.Run(m.mapperId, func(t *testing.T) {
			mapper := SqlMapper(m.mapperId)
			clause, err := mapper.build(mapDir, m.argMap)
			if len(m.errMsg) > 0 {
				assert.NotNil(err, m.desc)
				assert.Contains(err.Error(), m.errMsg, m.desc)
			} else {
				vs := make([]interface{}, 0, len(m.argMap))
				for _, value := range m.argMap {
					vs = append(vs, value)
				}
				assert.Equal(m.statement, clause.statement, m.mapperId, m.desc)
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
