package syntax

import (
	"encoding/xml"
	"github.com/Knetic/govaluate"
	"github.com/kcmvp/go-batis/dao/internal/cache"
)

type Test string

func (c Test) Name() string {
	return string(c)
}

func (c Test) Value() (interface{}, error) {
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
	XMLName xml.Name `xml:"if"`
	Test    Test     `xml:"test,attr"`
	Value   string   `xml:",chardata"`
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
	XMLName       xml.Name
	Id            string   `xml:"id,attr"`
	CacheName     string   `xml:"cacheName,attr"`
	CacheKey      string   `xml:"cacheKey,attr"`
	//ParameterType string   `xml:"parameterType,attr"`
	ResultType    string   `xml:"resultType,attr"`
	CharData1     string   `xml:",chardata"`
	Each          Foreach  //
	Sets          SetIf
	Wheres        WhereIf
	CharData2     string `xml:",chardata"`
	before        cacheHook
	after         cacheHook
}

type cacheHook func(string, ...interface{}) (interface{}, error)

var get cacheHook = func(s string, i ...interface{}) (interface{}, error) {
	return cache.Get(s)
}

var put cacheHook = func(s string, i ...interface{}) (interface{}, error) {
	return cache.Put(s, i)
}

var evict cacheHook = func(s string, i ...interface{}) (interface{}, error) {
	return cache.Evict(s)
}

// Build returns sql clause and cache key if exists.
func (receiver *Clause) Build(arg interface{}) (string, error) {
	panic("todo")
}

func (c Clause) Exec(arg interface{}) (interface{}, error) {
	if c.before != nil {
		if v, err := c.before(c.CacheKey, arg); v != nil && err == nil {
			return v, nil
		} else {
			c.after = put
		}
	}
	// exec()

	if c.after != nil {
		defer c.after(c.CacheKey, arg)
	}

	return nil, nil

}
