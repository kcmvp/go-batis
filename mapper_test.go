package batis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var mappers = []struct {
	mapper    string
	charData1 string
}{
	{"dog.createDog", "insert into Dog(name,age,price) values (#{name},#{age},#{price})"},
	{"dog.batchInsert", "insert into Dog(name,age,price) values"},
	//{"dog.findMyDog", "update Dog"},
}

var mapDir = "./mapper"

func TestMapperBuildCharData(t *testing.T) {
	assert := assert.New(t)
	for _, m := range mappers {
		t.Run(m.mapper, func(t *testing.T) {
			mp := SqlMapper(m.mapper)
			c, err := mp.build(mapDir)
			assert.Nil(err, "error should be nil")
			assert.Equal(m.charData1,c.statement , m.mapper, "charData1 should be equal")
		})
	}
}
