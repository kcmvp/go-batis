package syntax

import (
	"database/sql"
	"encoding/xml"
	"github.com/Knetic/govaluate"
	. "github.com/go-batis/dao"
	"github.com/go-batis/dao/internal/cache"
	"net/http"
)

type Condition string

func (c Condition) Name() string {
	http.Handle("/foo", nil)
	return string(c)
}

func (c Condition) Value() (interface{}, error) {
	panic("todo")
	exp, _ := govaluate.NewEvaluableExpression(c.Name())
	//@todo fix return type
	return exp.Evaluate(nil)
}

type Foreach struct {
	XMLName    xml.Name `xml:"foreach"`
	Collection string   `xml:"collection,attr"`
	Item       string   `xml:"item,attr"`
	Separator  string   `xml:"separator,attr"`
	Value      string   `xml:",chardata"`
}

type If struct {
	XMLName xml.Name  `xml:"if"`
	Test    Condition `xml:"test,attr"`
	Value   string    `xml:",chardata"`
}

// set_value_if only is applicable for update
type SetIf struct {
	XMLName xml.Name `xml:"set"`
	Ifs     []If
	Value   string `xml:",chardata"`
}

// where_condition_if
type WhereIf struct {
	XMLName xml.Name `xml:"where"`
	Ifs     []If
}

type Clause struct {
	XMLName       xml.Name `xml:"insert"`
	Id            string   `xml:"id,attr"`
	CacheName     string   `xml:"cacheName,attr"`
	CacheKey      string   `xml:"cacheKey,attr"`
	ParameterType string   `xml:"parameterType,attr"`
	ResultType    string   `xml:"resultType,attr"`
	HardText1     string   `xml:",chardata"`
	Each          Foreach  //
	Sets          SetIf
	Wheres        []WhereIf
	HardText2     string `xml:",chardata"`
}

// Build returns sql clause and cache key if exists.
func (receiver *Clause) Build(arg interface{}) (string, error) {
	panic("todo")
}



func (c Clause) Exec(arg interface{}) (CacheFunc, error) {
	var rs sql.Result
	switch c.XMLName.Local {
	case "select":
		return func(s string, i interface{}) (interface{}, error) {
			return cache.Put(s, i)
		}, nil
	case "delete", "update":
		panic("")
	}
}
