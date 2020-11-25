package syntax

import (
	"encoding/xml"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/kcmvp/go-batis/dao/internal/cache"
	"regexp"
	"strings"
)

type TestIf string

func (c TestIf) Name() string {
	return string(c)
}

func (c TestIf) Value() (interface{}, error) {
	panic("todo")
	exp, _ := govaluate.NewEvaluableExpression(c.Name())
	//@todo fix return type
	return exp.Evaluate(nil)
}

type DynamicSql interface {
	build(ctx interface{}) (string, []interface{}, error)
}

type Foreach struct {
	XMLName    xml.Name `xml:"foreach"`
	Collection string   `xml:"collection,attr"`
	Item       string   `xml:"item,attr"`
	Separator  string   `xml:"separator,attr"` // default ,
	Value      string   `xml:",chardata"`
}

func (f Foreach) build(ctx interface{}) (string, []interface{}, error) {
	panic("implement me")
}

type If struct {
	XMLName xml.Name `xml:"if"`
	Test    TestIf   `xml:"test,attr"`
	Value   string   `xml:",chardata"`
}

// set_value_if only is applicable for update
type SetIf struct {
	XMLName xml.Name `xml:"set"`
	Ifs     []If
	Value   string `xml:",chardata"`
}

func (sf SetIf) build(ctx interface{}) (string, []interface{}, error) {
	panic("implement me")
}

// where_condition_if
type WhereIf struct {
	XMLName xml.Name `xml:"where"`
	Ifs     []If
}

func (wf WhereIf) build(ctx interface{}) (string, []interface{}, error) {
	panic("implement me")
}

type CharData string

var reg = regexp.MustCompile(`#\{\w*\.?\w*\}`)

func (c CharData) build(ctx interface{}) (string, []interface{}, error) {
	//reg.FindAllString(string(c),-1)
	//for _, arg := range args {
	//	s := strings.ReplaceAll(arg,"#{","")
	//	_ := strings.ReplaceAll(s, "}","")
	//
	//}
	//reg.ReplaceAllString(string(c),"?")
	panic("")
}

type Clause struct {
	XMLName   xml.Name
	Id        string `xml:"id,attr"`
	CacheName string `xml:"cacheName,attr"`
	CacheKey  string `xml:"cacheKey,attr"`
	CharData1  CharData `xml:",chardata"`
	Foreach             //
	SetIf
	WhereIf
	CharData2 CharData `xml:",chardata"`
	sql       *string
	args      *[]interface{}
}



// Eval returns sql clause and cache key if exists.
func (clause *Clause) Eval(arg interface{}) error {
	ds := []DynamicSql{clause.CharData1, clause.Foreach, clause.Foreach, clause.SetIf, clause.WhereIf}

	for _, d := range ds {
		if s, _, err := d.build(arg); err != nil {
			panic("todo")
		} else {
			// s build sql
			fmt.Sprintf(s)
		}
	}
	panic("todo")

}


func (clause *Clause) Exec(arg interface{}) (interface{}, error) {
	if len(clause.CacheName) > 0 && len (clause.CacheKey) > 0 {
		if strings.EqualFold(clause.XMLName.Local, "select") {
			// hit cache
			if v, err := cache.Get(clause.CacheName+"::"+clause.CacheKey); err == nil {
				return v, nil
			} else {
				// execute sql @todo
				// cache the result @todo
				cache.Put(clause.CacheName+"::"+clause.CacheKey, "")
			}
		} else {
			// execute sql @todo
			cache.Evict(clause.CacheName+"::"+clause.CacheKey)
		}
	} else {
		// execute sql
		//@todo
	}

	return nil, nil

}
