package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type MapperType uint

const (
	InsertMapper MapperType = iota
	UpdateMapper
	DeleteMapper
	FindMapper
	SearchMapper
)

func (mapperType MapperType) name() string {
	return [...]string{"insert", "update", "delete", "select", "select"}[mapperType]
}

func (mapperType MapperType) NewClause() interface{} {
	switch mapperType {
	case InsertMapper:
		return &InsertClause{}
	case UpdateMapper:
		return &UpdateClause{}
	case DeleteMapper:
		return &DeleteClause{}
	case FindMapper, SearchMapper:
		return SelectClause{}
	default:
		return nil
	}

}

// return the type of doCreate's parameter
// InsertMapper, UpdateMapper: interface{} is parameterType
// FindMapper, SearchMapper: interface{} is the resultType
type Mapper func() (MapperType, string, interface{})

func (mapper Mapper) Exec(arg interface{}) (sql.Result, interface{}) {
	mapperType, mapperName, parmType := mapper()
	fmt.Print("dao name is v%", mapperName)
	ta := reflect.TypeOf(parmType)
	tb := reflect.TypeOf(arg)
	switch mapperType {
	case InsertMapper:
		// paramType is input
		fmt.Print("it's a create dao")
	case UpdateMapper:
		// paramType is input
		fmt.Print("it's a create dao")
	case DeleteMapper:
		// paramType is input
		fmt.Print("it's a create dao")
	case FindMapper:
		// paramType is output Type
		fmt.Print("it's a create dao")
	case SearchMapper:
		// paramType is output Type
		fmt.Print("it's a create dao")
	default:
		panic("it's not a valid dao type")
	}
	if ta == tb {
		// todo
	}
	return nil, nil
}

// mapper file naming pattern is ${struct}Mapper.xml
func (mapper Mapper) clause() (*xmlquery.Node, error) {
	mapperType, mapperName, _ := mapper()
	if entries := strings.Split(mapperName, "."); len(entries) == 2 {
		path, err := filepath.Abs(fmt.Sprintf("../mapper/%vMapper.xml", entries[0]))
		if err != nil {
			return nil, err
		}
		f, err := os.OpenFile(path, os.O_RDONLY, 0666)
		if err != nil {
			return nil, err
		}
		root, err := xmlquery.Parse(f)
		if err != nil {
			return nil, err
		}
		node := xmlquery.FindOne(root, fmt.Sprintf("//%v[@id='%v']", mapperType.name(), entries[1]))
		if node == nil {
			return node, errors.New("can't find the node")
		}
		return node, nil
		/*
			dataStr := node.OutputXML(false)
			fmt.Print(dataStr)
			return strings.TrimFunc(node.InnerText(), func(r rune) bool {
				return !(unicode.IsGraphic(r)) || unicode.IsSpace(r)
			}), nil
		*/
	} else {
		return nil, errors.New("can't find the dao")
	}

}
