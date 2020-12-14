package batis

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var mappers = []struct {
	mapperId  string
	arg       interface{}
	statement string
}{
	//{"dog.createDog", nil, "insert into Dog(name,age,price) values (#{name},#{age},#{price})"},
	{"dog.updateDog", map[string]interface{}{
		"name": "hello",
	}, "update Dog"},
}

var mapDir = "./mapper"

func TestMapperBuildCharData(t *testing.T) {
	assert := assert.New(t)
	for _, m := range mappers {
		t.Run(m.mapperId, func(t *testing.T) {
			mp := SqlMapper(m.mapperId)
			c, err := mp.build(mapDir, m.arg)
			assert.Nil(err, "error should be nil")
			assert.Equal(m.statement, strings.TrimSpace(strings.Trim(c.statement, "\n")), m.mapperId, "charData1 should be equal")
		})
	}
}
