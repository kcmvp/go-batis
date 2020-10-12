package sqlx

import (
	"database/sql"
	"fmt"
	"reflect"
)

type MapperType uint

const (
	CREATE  MapperType = iota
	UPDATE
	DELETE
	FIND
	SEARCH
)

// return the type of doCreate's parameter
type Mapper func() (MapperType, string, interface{})

func (mapper Mapper) Exec(arg interface{}) (sql.Result, interface{})  {
	mapperType, mapperName, parmType := mapper()
	fmt.Print("mapper name is v%", mapperName)
	ta := reflect.TypeOf(parmType)
	tb := reflect.TypeOf(arg)
	switch mapperType {
	case CREATE:
		fmt.Print("it's a create mapper")
	case UPDATE:
		fmt.Print("it's a create mapper")
	case DELETE:
		fmt.Print("it's a create mapper")
	case FIND:
		fmt.Print("it's a create mapper")
	case SEARCH:
		fmt.Print("it's a create mapper")
	default:
		panic("it's not a valid mapper type")
	}
	if ta == tb {
		// todo
	}
	return nil,nil
}
