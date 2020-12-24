package batis

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/antonmedv/expr"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

type SqlMapper string

type SqlType string

var sqlTypes = []SqlType{"insert", "select", "delete", "update", "sql"}

type StatementType string

var statementTypes = []StatementType{"where", "set", "if", "foreach"}

type Clause struct {
	xmlRoot   *xmlquery.Node
	args      interface{}
	id        string
	sqlType   SqlType
	statement string
	cacheName string
	cacheKey  string
	sqlParams []interface{}
}

const cacheKeyAttr, cacheNameAttr = "cacheKey", "cacheName"

var paramPattern = regexp.MustCompile(`#\{\w*\.?\w*\}`)

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
		clause.xmlRoot = root
		return clause, clause.buildMapperNode(entries[1])
	} else {
		return nil, errors.New(fmt.Sprintf("invalid naming standard %v", mapperName))
	}
}

func (clause *Clause) findChildById(id string) *xmlquery.Node {
	var node *xmlquery.Node
	for _, t := range sqlTypes {
		node = xmlquery.FindOne(clause.xmlRoot, fmt.Sprintf("//%v[@id='%v']", t, id))
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
	root := clause.findChildById(id)
	if root == nil {
		return errors.New(fmt.Sprintf("failed to find the node %v", id))
	}
	clause.id = id
	clause.sqlType = SqlType(root.Data)

	cacheName := strings.TrimSpace(root.SelectAttr(cacheNameAttr))
	cacheKeyExp := strings.TrimSpace(root.SelectAttr(cacheKeyAttr))
	if len(cacheName) > 0 && len(cacheKeyExp) < 1 || len(cacheName) < 1 && len(cacheKeyExp) > 0 {
		return fmt.Errorf("mapper#%v: empty cache name or key", id)
	} else if len(cacheKeyExp) > 0 && !paramPattern.MatchString(cacheKeyExp) {
		return fmt.Errorf("mapper#%v: cache key must be an expression %v", id, cacheKeyExp)
	} else if len(cacheName) > 0 && len(cacheKeyExp) > 0 {
		if v, err := clause.eval([]string{cacheKeyExp}); err != nil || len(v) == 0 {
			return fmt.Errorf("mapper#%v: invalid cache key %v", id, cacheKeyExp)
		} else {
			clause.cacheName = cacheName
			clause.cacheKey = fmt.Sprintf("%v", v[0])
		}
	}

	var buff bytes.Buffer
	for node := root.FirstChild; node != nil; node = node.NextSibling {
		if err := clause.buildXmlNode(root, node, &buff); err != nil {
			return err
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if err := clause.buildXmlNode(node, child, &buff); err != nil {
				return err
			}
		}
	}
	clause.statement = strings.TrimSpace(buff.String())
	return nil
}

//@FixMe need to remove redundant \n\r
func (clause *Clause) buildXmlNode(parent, current *xmlquery.Node, buff *bytes.Buffer) (err error) {
	st := current.Data
	switch current.Type {
	case xmlquery.TextNode, xmlquery.CharDataNode:
		if st, err = clause.sqlParameter(parent, st); err == nil {
			prettySql(buff, st)
		}
		return err
	// for xmlquery.ElementNode
	case xmlquery.ElementNode:
		if strings.EqualFold("where", st) || strings.EqualFold("set", st) {
			//buff.WriteString(st)
			prettySql(buff, st)
		} else if strings.EqualFold("include", st) {
			if node := clause.findChildById(current.SelectAttr("refid")); node != nil {
				//xml.EscapeText(buff, []byte(b.InnerText()))
				//prettySql(buff, b.InnerText())
				for node = node.FirstChild; node != nil; node = node.NextSibling {
					if err = clause.buildXmlNode(clause.xmlRoot, node, buff); err != nil {
						return err
					}
					for child := node.FirstChild; child != nil; child = child.NextSibling {
						if err = clause.buildXmlNode(node, child, buff); err != nil {
							return err
						}
					}
				}

			} else {
				return fmt.Errorf("mapper#%v:failed to find the include %v", clause.id, current.SelectAttr("refid"))
			}
		} else if strings.EqualFold("if", st) && clause.args != nil {
			el := current.SelectAttr("test")
			if value, ignoreErr := expr.Eval(el, clause.args); ignoreErr == nil && value.(bool) {
				if st, err = clause.sqlParameter(parent, current.InnerText()); err == nil {
					prettySql(buff, st)
				} else {
					//fmt.Println(fmt.Sprintf("mapper#%v:expression evaluate false : %v. %v", clause.id, el, err))
					return err
				}
			} else {
				fmt.Println(fmt.Sprintf("mapper#%v:expression evaluate false : %v. %v", clause.id, spacePattern.ReplaceAllString(el, " "), ignoreErr))
			}
		} else if strings.EqualFold("foreach", st) {
			fmt.Println(fmt.Sprintf("mapper#%v:foreach node : %v.", clause.id, current.SelectAttr("collection")))
		}
		return
	default:
		fmt.Println(fmt.Sprintf("xml node value is: %v", current.Data))
		return
	}
}

var nestParamPattern = regexp.MustCompile(`#\{\s+\.\w*\}`)

func (clause Clause) eval(exps []string, envs ...interface{}) ([]interface{}, error) {
	var rt []interface{}
	for _, exp := range exps {
		par := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(exp, "#{"), "}"))
		env := clause.args
		if strings.HasPrefix(par, ".") {
			if len(envs) < 1 {
				return nil, fmt.Errorf("missed additional property for dot property: '%v'", exp)
			}
			env = envs[0]
			par = strings.ReplaceAll(par, ".", "")
		}
		if v, err := expr.Eval(par, env); err != nil || v == nil {
			return nil, fmt.Errorf("can not resolve: '%v'", par)
		} else {
			rt = append(rt, v)
		}
	}
	return rt, nil
}

func (clause *Clause) sqlParameter(parent *xmlquery.Node, exp string) (string, error) {
	var buff bytes.Buffer
	var err error
	exp = strings.TrimSpace(spacePattern.ReplaceAllString(exp, " "))
	if params := paramPattern.FindAllString(exp, -1); len(params) > 0 {
		if parent.Data == "foreach" {
			slicedArg := clause.args
			collection := strings.TrimSpace(parent.SelectAttr("collection"))
			//item := strings.TrimSpace(parent.SelectAttr("item"))
			separator := strings.TrimSpace(parent.SelectAttr("separator"))
			if len(separator) == 0 {
				separator = ","
			}
			if len(collection) > 0 {
				if slicedArg, err = expr.Eval(collection, clause.args); err != nil ||
					reflect.ValueOf(slicedArg).Kind() != reflect.Slice {
					return "", fmt.Errorf("mapper#%v: foreach statement, collection property must be a slice: `%v`, %w", clause.id, collection, err)
				}
			} else if reflect.ValueOf(clause.args).Kind() != reflect.Slice {
				return "", fmt.Errorf("mapper#%v: foreach statement, parameter %+v is not a slice", clause.id, clause.args)
			}
			s := reflect.ValueOf(slicedArg)
			for i := 0; i < s.Len(); i++ {
				if values, err := clause.eval(params, s.Index(i)); err == nil && len(values) > 0 {
					clause.sqlParams = append(clause.sqlParams, values...)
				} else {
					return "", fmt.Errorf("mapper#%v: %v, %w", clause.id, exp, err)
				}
				buff.WriteString(paramPattern.ReplaceAllString(exp, "?"))
				if i != s.Len()-1 {
					buff.WriteString(separator)
				}
			}
		} else {
			if values, err := clause.eval(params); err == nil && len(values) > 0 {
				clause.sqlParams = append(clause.sqlParams, values...)
			} else {
				return "", fmt.Errorf("mapper#%v: %v, %w", clause.id, exp, err)
			}
			buff.WriteString(paramPattern.ReplaceAllString(exp, "?"))
		}
		return buff.String(), nil
	} else {
		return exp, nil
	}
}

var spacePattern = regexp.MustCompile(`\s+`)
var sqlWherePattern = regexp.MustCompile(`\s+where\s+$`)
var sqlAnd = regexp.MustCompile(`^(?i)and\s+`)

func prettySql(buff *bytes.Buffer, str string) {
	str = spacePattern.ReplaceAllString(str, " ")
	str = strings.TrimSpace(str)
	if len(str) > 0 {
		if sqlWherePattern.MatchString(buff.String()) && sqlAnd.MatchString(str) {
			str = sqlAnd.ReplaceAllString(str, "")
		}
		buff.WriteString(str + " ")
		//xml.Escape(buff, []byte(str+" "))
	}
}
