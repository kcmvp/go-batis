package dao

import "encoding/xml"
import "github.com/Knetic/govaluate"

type Condition string

func (c Condition) Name() string {
	return string(c)
}

func (c Condition) Value() (interface{}, error) {
	panic("todo")
	exp, _ := govaluate.NewEvaluableExpression(c.Name())
	//@todo fix return type
	return exp.Evaluate(nil)
}

type CommonAttr struct {
	Id string `xml:"id,attr"`
}

type Foreach struct {
	XMLName    xml.Name `xml:"foreach"`
	Collection string   `xml:"collection,attr"`
	Item       string   `xml:"item,attr"`
	Separator  string   `xml:"separator,attr"`
	Value       string   `xml:",chardata"`
}


type If struct {
	XMLName xml.Name  `xml:"if"`
	Test    Condition `xml:"test,attr"`
	Value    string    `xml:",chardata"`
}

// set_value_if only is applicable for update
type SetIf struct {
	XMLName xml.Name `xml:"set"`
	Ifs     []If
	Value    string `xml:",chardata"`
}
// where_condition_if
type WhereIf struct {
	XMLName xml.Name `xml:"where"`
	Ifs     []If
}

// sql clause

type Clause interface {
	Build(ctx interface{}) (string, error)
}

type InsertClause struct {
	XMLName xml.Name `xml:"insert"`
	CommonAttr
	ParameterType string `xml:"parameterType,attr"`
	Text string `xml:",chardata"`
	Each Foreach
}

func (i InsertClause) Build(ctx interface{}) (string, error) {
	panic("implement me")
}

type UpdateClause struct {
	XMLName xml.Name `xml:"update"`
	CommonAttr
	ParameterType string `xml:"parameterType,attr"`
	CacheEvict string `xml:"cacheEvict,attr"`
	BeforeValue string `xml:",chardata"`
	Sets        SetIf
	AfterValue  string `xml:",chardata"`
}

func (u UpdateClause) Build(ctx interface{}) (string, error) {
	panic("implement me")
}

type SelectClause struct {
	XMLName xml.Name `xml:"select"`
	CommonAttr
	ResultType string `xml:"resultType,attr"`
	CacheKey string `xml:"cacheKey,attr"`
	BeforeValue string `xml:",chardata"`
	Wheres []WhereIf
	AfterValue string `xml:",chardata"`
}

func (s SelectClause) Build(ctx interface{}) (string, error) {
	panic("implement me")
}

type DeleteClause struct {
	XMLName xml.Name `xml:"delete"`
	CommonAttr
	Value string `xml:",chardata"`
}

func (d DeleteClause) Build(ctx interface{}) (string, error) {
	panic("implement me")
}

