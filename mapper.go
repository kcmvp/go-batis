package batis

import (
	"bytes"
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

	cacheName := strings.TrimSpace(node.SelectAttr(CACHE_NAME_ATTR))
	cacheKeyExp := strings.TrimSpace(node.SelectAttr(CACHE_KEY_ATTR))
	if len(cacheName) > 0 && len(cacheKeyExp) < 1 || len(cacheName) < 1 && len(cacheKeyExp) > 0 {
		return fmt.Errorf("mapper#%v: empty cache name or key", id)
	} else if len(cacheKeyExp) > 0 && !pattern.MatchString(cacheKeyExp) {
		return fmt.Errorf("mapper#%v: cache key must be an expression %v", id, cacheKeyExp)
	} else if len(cacheName) > 0 && len(cacheKeyExp) > 0 {
		if v, err := clause.eval(cacheKeyExp); err != nil || v == nil {
			return fmt.Errorf("mapper#%v: invalid cache key %v", id, cacheKeyExp)
		} else {
			clause.cacheName = cacheName
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
	clause.statement = strings.TrimSpace(buff.String())
	return nil
}

//@FixMe need to remove redundant \n\r
func (clause *Clause) buildXmlNode(n *xmlquery.Node, buff *bytes.Buffer) (err error) {
	st := n.Data
	switch n.Type {
	case xmlquery.TextNode, xmlquery.CharDataNode:
		if st, err = clause.processHolder(n.Data); err == nil {
			prettySql(buff, st)
		}
		return
	// for xmlquery.ElementNode
	case xmlquery.ElementNode:
		if strings.EqualFold("where", st) || strings.EqualFold("set", st) {
			//buff.WriteString(st)
			prettySql(buff, st)
		} else if strings.EqualFold("include", st) {
			if node := clause.findChildById(n.SelectAttr("refid")); node != nil {
				//xml.EscapeText(buff, []byte(b.InnerText()))
				//prettySql(buff, b.InnerText())
				for node = node.FirstChild; node != nil; node = node.NextSibling {
					if err = clause.buildXmlNode(node, buff); err != nil {
						return err
					}
					for child := node.FirstChild; child != nil; child = child.NextSibling {
						if err = clause.buildXmlNode(child, buff); err != nil {
							return err
						}
					}
				}

			} else {
				err = fmt.Errorf("mapper#%v:failed to find the include %v", clause.id, n.SelectAttr("refid"))
			}
		} else if strings.EqualFold("if", st) && clause.args != nil {
			el := n.SelectAttr("test")
			if value, ignoreErr := expr.Eval(el, clause.args); ignoreErr == nil && value.(bool) {
				if st, err = clause.processHolder(n.InnerText()); err == nil {
					prettySql(buff, st)
				} else {
					//fmt.Println(fmt.Sprintf("mapper#%v:expression evaluate false : %v. %v", clause.id, el, err))
					return
				}
			} else {
				fmt.Println(fmt.Sprintf("mapper#%v:expression evaluate false : %v. %v", clause.id, el, ignoreErr))
			}
		}
		return
	default:
		fmt.Println(fmt.Sprintf("xml node value is: %v", n.Data))
		return
	}
}

func (clause Clause) eval(exp string) (interface{}, error) {
	par := strings.TrimSuffix(strings.TrimPrefix(exp, "#{"), "}")
	return expr.Eval(par, clause.args)
}

func (clause *Clause) processHolder(str string) (string, error) {
	var buff bytes.Buffer
	for _, par := range pattern.FindAllString(str, -1) {
		if value, err := clause.eval(par); err == nil && value != nil {
			clause.sqlParams = append(clause.sqlParams, value)
		} else {
			return "", fmt.Errorf("mapper#%v: failed to resolve the expression: %v. %w", clause.id, par, err)
		}
	}
	buff.WriteString(pattern.ReplaceAllString(str, "?"))

	return buff.String(), nil
}

var newLineSpacePattern = regexp.MustCompile(`\s+`)
var sqlWherePattern = regexp.MustCompile(`\s+where\s+$`)
var sqlAnd = regexp.MustCompile(`^(?i)and\s+`)

func prettySql(buff *bytes.Buffer, str string) {
	str = newLineSpacePattern.ReplaceAllString(str, " ")
	str = strings.TrimSpace(str)
	if len(str) > 0 {
		if sqlWherePattern.MatchString(buff.String()) && sqlAnd.MatchString(str) {
			str = sqlAnd.ReplaceAllString(str, "")
		}
		buff.WriteString(str + " ")
		//xml.Escape(buff, []byte(str+" "))
	}
}
