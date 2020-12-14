package batis

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/antonmedv/expr"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SqlMapper string

type SqlType string

var sqlTypes = []SqlType{"insert", "select", "delete", "update", "sql"}

type StatementType string

var statementTypes = []StatementType{"where", "set", "if", "foreach"}

type Clause struct {
	root      *xmlquery.Node
	args      interface{}
	id        string
	sqlType   SqlType
	statement string
	cacheName string
	cacheKey  string
	sqlParams []interface{}
}

const CACHE_KEY_ATTR, CACHE_NAME_ATTR = "cacheKey", "cacheName"

var pattern = regexp.MustCompile(`#\{\w*\.?\w*\}`)

// mapperId file naming pattern is ${struct}Mapper.xml
// naming standard of mapperId is ${file name}.${mapperId}
// ex: `dog.findByName` means its definition in the `dog.xml` and the `id' attribute is `findByName`
func (mapper SqlMapper) build(mapperDir string, args interface{}) (*Clause, error) {
	mapperName := string(mapper)
	clause := &Clause{
		args: args,
	}
	var fName string
	if entries := strings.Split(mapperName, "."); len(entries) == 2 {
		fName = fmt.Sprintf("%v/%v.xml", mapperDir, entries[0])
		path, err := filepath.Abs(fName)
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
		clause.root = root
		clause.buildMapperNode(entries[1])
		return clause, nil
	} else {
		return nil, errors.New(fmt.Sprintf("invalid naming standard %v", mapperName))
	}
}

func (clause *Clause) findChildById(id string) *xmlquery.Node {
	var node *xmlquery.Node
	for _, t := range sqlTypes {
		node = xmlquery.FindOne(clause.root, fmt.Sprintf("//%v[@id='%v']", t, id))
		if node != nil && node.Data == string(t) {
			break
		}
	}
	return node
}

func (clause *Clause) CacheKey() (string, error) {
	if len(clause.cacheName) > 0 && len(clause.cacheKey) > 0 {
		return fmt.Sprintf("%v::%v", clause.cacheName, clause.cacheKey), nil
	} else {
		return "", fmt.Errorf("invalid cache key. cachePrefgix: %v, cacheId : %v", clause.cacheName, clause.cacheKey)
	}
}
func (clause *Clause) SqlType() SqlType {
	return clause.sqlType
}

func (clause *Clause) Statement() string {
	return clause.statement
}

func (clause *Clause) buildMapperNode(id string) error {
	node := clause.findChildById(id)
	if node == nil {
		return errors.New(fmt.Sprintf("failed to find the node %v", id))
	}
	clause.id = id
	clause.sqlType = SqlType(node.Data)
	clause.cacheName = strings.TrimSpace(node.SelectAttr(CACHE_NAME_ATTR))
	key := strings.TrimSpace(node.SelectAttr(CACHE_KEY_ATTR))
	if len(clause.cacheName) > 0 && len(key) > 0 {
		if v, err := expr.Eval(key, clause.args); err != nil {
			fmt.Println(fmt.Sprintf("invalid cache key %s, ignore cache for sql %v", key, id))
		} else {
			clause.cacheKey = fmt.Sprintf("%v", v)
		}
	}
	var buff bytes.Buffer
	for node = node.FirstChild; node != nil; node = node.NextSibling {
		clause.buildXmlNode(node, &buff)
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			clause.buildXmlNode(child, &buff)
		}
	}
	clause.statement = buff.String()

	return nil
}

//@FixMe need to remove redundant \n\t
func (clause *Clause) buildXmlNode(n *xmlquery.Node, buff *bytes.Buffer) (err error) {
	//buf := bytes.NewBufferString(clause.statement)
	switch n.Type {
	case xmlquery.TextNode, xmlquery.CharDataNode:
		//@FixMe need to check #{}, in some case there #{} in the statement
		buff.WriteString(n.Data)
		return
	// for xmlquery.ElementNode
	case xmlquery.ElementNode:
		//@todo
		t := n.Data
		if strings.EqualFold("where", t) || strings.EqualFold("set", t) {
			//xml.EscapeText(buff, []byte(n.Data))
			buff.WriteString(t)
		} else if strings.EqualFold("include", t) {
			if b := clause.findChildById(n.SelectAttr("refid")); b != nil {
				xml.EscapeText(buff, []byte(b.InnerText()))
			} else {
				err = fmt.Errorf("failed to find the include %v", n.SelectAttr("refid"))
			}
		} else if strings.EqualFold("if", t) && clause.args != nil {
			el := n.SelectAttr("test")
			var value interface{}
			if value, err = expr.Eval(el, clause.args); err == nil {
				if value.(bool) {
					for _, par := range pattern.FindAllString(n.InnerText(), -1) {
						par = strings.TrimSuffix(strings.TrimPrefix(par, "#{"), "}")
						if value, err = expr.Eval(par, clause.args); err == nil {
							clause.sqlParams = append(clause.sqlParams, value)
						} else {
							fmt.Println(fmt.Sprintf("failed to evaluate the expression: %v, %v", par, err))
							return
						}
					}
					buff.WriteString(pattern.ReplaceAllString(n.InnerText(), " ? "))
				}
			} else {
				fmt.Println(fmt.Sprintf("failed to evaluate the expression: %v, %v", el, err))
			}
		}
		return
	default:
		fmt.Println(fmt.Sprintf("xml node value is: %v", n.Data))
		return
	}
}
