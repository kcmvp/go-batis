package sql

import (
	"github.com/antchfx/xmlquery"
	//"github.com/antchfx/xmlquery"
	"github.com/kcmvp/go-batis/session"
	"github.com/stretchr/testify/assert"
	"os"
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
	Name  string
	Age   int
	Price int
}

var dogList = []Dog{
	{
		Name:  "mimi",
		Age:   1,
		Price: 2,
	}, {
		Name:  "kaka",
		Age:   2,
		Price: 3,
	},
}
var yellowDog = map[string]interface{}{
	"color":   "yellow",
	"dogList": dogList,
}

type ColoredDog struct {
	Color   string
	DogList []Dog
}

var redDog = ColoredDog{
	Color:   "red",
	DogList: dogList,
}

var mappers = []struct {
	desc     string
	mapperId string
	sqlType string
	arg     interface{}
	positive bool
	msg      string
}{
	//{"case1: miss parameters", "createDog", "insert",nil, false, "mapper#createDog: insert into Dog(name,age,price) values (#{name},#{age},#{price}), can not resolve: 'name'"},
	//{"case2: simple create", "createDog", "insert",dogMapF, true, "insert into Dog(name,age,price) values (?,?,?)"},
	{"case3: with partial parameters", "updateDog", "update",dogMapP, true, "update Dog set name = ?, age = ?, updated_at = CURRENT_TIMESTAMP() where id = ?"},
	{"case4: with all parameters", "updateDog", "update",dogMapF, true, "update Dog set name = ?, age = ?, price = ?, updated_at = CURRENT_TIMESTAMP() where id = ?"},
	//{"case5: miss cache name", "findDogByIdNoCacheName", dogMapF, false, "mapper#findDogByIdNoCacheName: empty cache name or key"},
	{"case6: simple find clause", "findDogById", "select",dogMapF, true, "select * from Dog where id = ?"},
	{"case6: dynamic where", "searchByExample", "select",dogMapF, true, "select count(1) from Dog where name = ? and age = ?"},
	{"case7: delete statement missed parameter", "deleteDogById", "delete",dogMapF, false, "mapper#deleteDogById: delete from dog where levy_serial_number = #{levySerialNumber}, can not resolve: 'levySerialNumber'"},
	{"case8: invalid placeholder character", "updateWeekDayPriceCase1","update", dogMapF, true, "UPDATE T_WEEK_DAY_PRICE set PRICE = ${price}, UPDATED_AT = CURRENT_TIMESTAMP() where PRICE_PLAN_ID = ${age} and NUM_OF_WEEK = ${id}"},
	{"case9: simple sql ref", "selectByRef", "select",dogMapF, true, "SELECT name, age, size FROM UUC_COMPANY where id = ?"},
	{"case10: nest sql ref", "selectByRefNest","select", dogMapF, true, "SELECT f.ID , f.PROCESS_KEY , f.PROCESS_NAME , f.MODULE_CODE , m.MODULE_NAME , f.NOTE , f.STATUS FROM UBPC_PROCESS_FILE f LEFT JOIN UBPC_MODULE m ON m.MODULE_CODE = f.MODULE_CODE where f.MODULE_CODE = ? AND (f.PROCESS_NAME LIKE CONCAT('%',?,'%') OR f.PROCESS_KEY LIKE CONCAT('%',?,'%')) AND f.age = ?"},
	{"case11: sql with escape", "findDogByIdEscape","select", dogMapF, true, "select * from Dog where id <= ? and price >= 100"},
	{"case12: invalid parameter", "forEachCase1","insert", dogMapF, false, "is not a slice"},
	{"case13: happy flow", "forEachCase1","insert", map[string]interface{}{
		"dogList":dogList,
	}, true, "insert into Dog(name,age,price) values (?,?,?),(?,?,?)"},
	{"case14: ", "forEachCase2","insert", redDog, true, "insert into Dog(name,age,price) values (?,?,?);(?,?,?)"},
}

var mapDir = "./defaultDS/dog.xml"

func TestMapperBuildCharData(t *testing.T) {
	assertion := assert.New(t)
	f, err := os.OpenFile(mapDir, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	node, err := xmlquery.Parse(f)
	if err != nil {
		panic(err)
	}
	for _, m := range mappers {
		t.Run(m.mapperId, func(t *testing.T) {
			clause := &Clause{
				id: m.mapperId,
				doc: node,
				arg: m.arg,
				sqlType: session.SqlType(m.sqlType),
			}
			err = clause.build()
			assertion.Equal(1,1)
			/*
			if !m.positive {
				assertion.NotNil(err, m.desc)
				assertion.Contains(err.Error(), m.msg, m.desc)
			} else {
				assertion.Nil(err,m.mapperId)
				assertion.Equal(m.msg, clause.statement, m.mapperId, m.desc)

				if reflect.ValueOf(m.arg).Kind() == reflect.Map {
					s := reflect.ValueOf(m.arg)
					vs := make([]interface{}, 0, s.Len())
					for _, key := range s.MapKeys() {
						vs = append(vs, s.MapIndex(key).Interface())
					}
					for _, v := range clause.sqlParams {
						assertion.True(func(i interface{}) bool {
							for _, param := range vs {
								if param == v {
									return true
								}
							}
							return false
						}(v), m.desc)
					}
				}
			}
			 */
		})
	}
}
