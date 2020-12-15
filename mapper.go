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
		return clause, clause.buildMapperNode(entries[1])
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
		if v, err := expr.Eval(key, clause.args); err != nil || v == nil {
			return fmt.Errorf("invalid cache key %s, ignore cache for sql %v", key, id)
		} else {
			clause.cacheKey = fmt.Sprintf("%v", v)
		}
	}
	var buff bytes.Buffer
	for node = node.FirstChild; node != nil; node = node.NextSibling {
		if err := clause.buildXmlNode(node, &buff); err != nil {
			return err
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if err := clause.buildXmlNode(child, &buff); err != nil {
				return err
			}
		}
	}
	clause.statement = buff.String()
	return nil
}

//@FixMe need to remove redundant \n\t
func (clause *Clause) buildXmlNode(n *xmlquery.Node, buff *bytes.Buffer) (err error) {
	st := strings.TrimSpace(n.Data)
	switch n.Type {
	case xmlquery.TextNode, xmlquery.CharDataNode:
		//@FixMe need to check #{}, in some case there #{} in the statement
		if st, err = clause.placeHolder(n.Data); err == nil {
			buff.WriteString(st)
		}
		return
	// for xmlquery.ElementNode
	case xmlquery.ElementNode:
		if strings.EqualFold("where", st) || strings.EqualFold("set", st) {
			buff.WriteString(st)
		} else if strings.EqualFold("include", st) {
			if b := clause.findChildById(n.SelectAttr("refid")); b != nil {
				xml.EscapeText(buff, []byte(b.InnerText()))
			} else {
				err = fmt.Errorf("failed to find the include %v", n.SelectAttr("refid"))
			}
		} else if strings.EqualFold("if", st) && clause.args != nil {
			el := n.SelectAttr("test")
			var value interface{}
			if value, err = expr.Eval(el, clause.args); err == nil && value.(bool) {
				if st, err = clause.placeHolder(n.InnerText()); err == nil {
					buff.WriteString(st)
				}
			} else if err != nil {
				return fmt.Errorf("failed to resolve the expression: %v%w", el, err)
			}
		}
		return
	default:
		fmt.Println(fmt.Sprintf("xml node value is: %v", n.Data))
		return
	}
}

func (clause *Clause) placeHolder(str string) (string, error) {
	var buff bytes.Buffer
	for _, par := range pattern.FindAllString(str, -1) {
		par = strings.TrimSuffix(strings.TrimPrefix(par, "#{"), "}")
		if value, err := expr.Eval(par, clause.args); err == nil && value != nil {
			clause.sqlParams = append(clause.sqlParams, value)
		} else {
			return "", fmt.Errorf("failed to resolve the expression: #{%v} for mapper:%v. %w", par, clause.id, err)
		}
	}
	buff.WriteString(pattern.ReplaceAllString(str, "?"))

	return buff.String(), nil
}

//func prettyPrint(str string) string {
//	str = strings.ReplaceAll(str, "\n"," ")
//
//}
