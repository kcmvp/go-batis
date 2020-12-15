package batis

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var dogMap = map[string]interface{}{
	"id":    1234567,
	"name":  "hello",
	"age":   12,
	"price": 123.5,
}

var mappers = []struct {
	mapperId  string
	argMap    map[string]interface{}
	errMsg    string
	statement string
}{
	//{"dog.createDog", nil, "failed to resolve the expression: #{name} for mapper:createDog.", ""},
	//{"dog.createDog", dogMap, "", "insert into Dog(name,age,price) values (?,?,?)"},
	{"dog.updateDog", dogMap, "", ""},
}

var mapDir = "./mapper"

func TestMapperBuildCharData(t *testing.T) {
	assert := assert.New(t)
	for _, m := range mappers {
		t.Run(m.mapperId, func(t *testing.T) {
			mapper := SqlMapper(m.mapperId)
			clause, err := mapper.build(mapDir, m.argMap)
			if len(m.errMsg) > 0 {
				assert.NotNil(err)
				assert.Contains(err.Error(), m.errMsg)
			} else {
				vs := make([]interface{}, 0, len(m.argMap))
				for _, value := range m.argMap {
					vs = append(vs, value)
				}
				assert.Equal(m.statement, strings.TrimSpace(strings.Trim(clause.statement, "\n")), m.mapperId, "charData1 should be equal")
				for _, v := range vs {
					assert.True(func(i interface{}) bool {
						for _, param := range clause.sqlParams {
							if param == v {
								return true
							}
						}
						return false
					}(v))
				}
			}
		})
	}
}
