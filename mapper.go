package sql

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/antonmedv/expr"
	"github.com/kcmvp/go-batis/session"
	"reflect"
	"regexp"
	"strings"
)

var nodeTypes = append(session.SqlTypes, "sql")

type Mapper string

type Clause struct {
	doc       *xmlquery.Node
	arg       interface{}
	ctx       context.Context
	id        string
	sqlType   session.SqlType
	statement bytes.Buffer
	cacheKey  string
	sqlParams []interface{}
}

const scopePrefix = "mapper_"
const cacheKeyAttr = "cacheKey"

var paramPattern = regexp.MustCompile(`#{\w*\.?\w*}`)
var atLeastOnCharacter = regexp.MustCompile(".*\\w.*")
var spacePattern = regexp.MustCompile(`\s+`)
var sqlWherePattern = regexp.MustCompile(`\s+where\s+$`)
var sqlAnd = regexp.MustCompile(`^(?i)and\s+`)

// https://www.freecodecamp.org/news/generics-in-golang/
func (m Mapper) Id() string {
	return string(m)
}

// arg must be one of below types
// map, struct, slice of map, slice of struct
func (m Mapper) Query(arg interface{}) ([]interface{}, error) {
	return m.QueryContext(context.Background(), arg)
}

func (m Mapper) QueryContext(ctx context.Context, arg interface{}) ([]interface{}, error) {
	mCtx := session.MapperContext(m.Id())
	c := &Clause{
		id:  m.Id(),
		doc: mCtx.Mapping(),
		//arg:     arg,
		ctx:     context.WithValue(ctx, fmt.Sprintf("%v_#", scopePrefix), arg),
		sqlType: "select",
	}
	if err := c.buildNode(c.root()); err != nil {
		return nil, errors.Unwrap(fmt.Errorf("failed to buildNode the mapper: %w ", err))
	}
	dest := []interface{}{}
	return dest, mCtx.Session().QueryCacheable(ctx, c.Statement(), c.cacheKey, arg, &dest)
}

func (m Mapper) Exec(arg interface{}) (sql.Result, error) {
	return m.ExecContext(context.Background(), arg)
}

func (m Mapper) ExecContext(ctx context.Context, arg interface{}) (sql.Result, error) {
	//mCtx := session.MapperContext(m.Id())
	//m.buildNode(mCtx.Mapping(), arg)
	//panic(fmt.Sprintf("session is %v", mCtx))
	panic("")
}

func (clause *Clause) root() *xmlquery.Node {
	return clause.findChildById(clause.id)
}

func (clause *Clause) findChildById(id string) *xmlquery.Node {
	var node *xmlquery.Node
	//@fixme no need nodeType
	for _, t := range nodeTypes {
		node = xmlquery.FindOne(clause.doc, fmt.Sprintf("//%v[@id='%v']", t, id))
		if node != nil && node.Data == string(t) {
			break
		}
	}
	return node
}

func (clause *Clause) SqlType() session.SqlType {
	return clause.sqlType
}

func (clause *Clause) Statement() string {
	return clause.statement.String()
}

func (clause *Clause) build() error {
	return clause.buildNode(clause.root())
}

func (clause *Clause) buildNode(node *xmlquery.Node) error {
	var err error

	if node == nil || strings.TrimSpace(node.Data) == "" {
		return nil
	}
	if strings.EqualFold("mapper", node.Parent.Data) {
		if !strings.EqualFold(string(clause.sqlType), node.Data) {
			return fmt.Errorf("mapper#%v: expect clause type %v, but the type is %v", clause.id, clause.sqlType, node.Data)
		}
		cacheKeyExp := strings.TrimSpace(node.SelectAttr(cacheKeyAttr))
		if len(cacheKeyExp) > 0 {
			if v, err := clause.evaluate([]string{cacheKeyExp}); err != nil || len(v) == 0 {
				return errors.Unwrap(fmt.Errorf("mapper#%v: invalid cache key %v: %w", clause.id, cacheKeyExp, err))
			} else {
				clause.cacheKey = fmt.Sprintf("%v", v[0])
			}
		}
	}

	st := node.Data
	fmt.Printf("mapper is %v,node data is %v\n", clause.id, st)
	switch node.Type {
	case xmlquery.TextNode, xmlquery.CharDataNode:
		//if err = clause.eval(st); err != nil {
		//	//normalizeSql(buff, st)
		//}
		//return err
		//return clause.eval(st)

		clause.normalizeSql(st)
	case xmlquery.ElementNode:
		if strings.EqualFold("foreach", st) {
			/**
			1: if there is no collection then the parameter must be slice
			2: if there is collection then must be evaluate to slice.
			3: item is mandatory
			 */
			list := node.SelectAttr("collection")
			item := node.SelectAttr("item")
			fmt.Println(fmt.Sprintf("mapper#%v:foreach node : collection:%v, item:%v", clause.id, list, item))
			clause.ctx = context.WithValue(clause.ctx, fmt.Sprintf("%v_%v", scopePrefix, item), list)
		} else if strings.EqualFold("where", st) || strings.EqualFold("set", st) {
			// @do nothing
			//clause.normalizeSql(st)
		} else if strings.EqualFold("include", st) {
			if refNode := clause.findChildById(node.SelectAttr("refid")); refNode != nil {
				clause.buildNode(refNode)
			} else {
				return fmt.Errorf("mapper#%v:failed to find the include %v", clause.id, node.SelectAttr("refid"))
			}
			//@todo fix the case 1==1
		} else if strings.EqualFold("if", st) && clause.arg != nil {
			el := node.SelectAttr("test")
			if value, ignoreErr := expr.Eval(el, clause.arg); ignoreErr != nil || !value.(bool) {
				fmt.Println(fmt.Sprintf("mapper#%v:expression evaluate false : %v. %v", clause.id, spacePattern.ReplaceAllString(el, " "), ignoreErr))
				// skip the node and its children
				return nil
			}
		}
	default:
		fmt.Println(fmt.Sprintf("xml node value is: %v", st))
	}
	for node = node.FirstChild; node != nil; node = node.NextSibling {
		if err = clause.buildNode(node); err != nil {
			return err
		}
	}
	return err
}

/*

func (clause Clause) eval(exp string, test bool) ([]interface{}, error) {
	var err error
	exp = strings.TrimSpace(spacePattern.ReplaceAllString(exp, " "))
	var rt []interface{}
	for _, par := range paramPattern.FindAllString(exp, -1) {
		par = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(par, "#{"), "}"))
		var v interface{}
		items := strings.Split(par, ".");
		ctx := clause.arg
		if len(items) > 1 {
			if ctx = clause.argMap[items[0]]; ctx == nil {
				return nil, fmt.Errorf("can not find parameter %v", items[0])
			}
		}
		if v, err = expr.Eval(items[0], ctx); err == nil && v != nil {
			rt = append(rt, v)
			if !test {
				// @todo 1: add sqlMap 2: sqlstatement
			}
		} else {
			return nil, errors.Unwrap(fmt.Errorf("can not resolve: '%v', %w", par, err))
		}

		return rt, err
	}

	return nil, err
	if params := paramPattern.FindAllString(exp, -1); len(params) > 0 {
		if parent.Data == "foreach" {
			slicedArg := clause.args
			collection := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(parent.SelectAttr("collection"), "#{"), "}"))
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
				if values, err := clause.evaluate(params, s.Index(i).Interface()); err == nil && len(values) > 0 {
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
			if values, err := clause.evaluate(params); err == nil && len(values) > 0 {
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
*/

//@FixMe need to remove redundant \n\r
/*
func (clause *Clause) buildMapperNode(parent, current *xmlquery.Node, buff *bytes.Buffer) (err error) {
	st := current.Data
	switch current.Type {
	case xmlquery.TextNode, xmlquery.CharDataNode:
		if st, err = clause.sqlParameter(parent, st); err == nil {
			normalizeSql(buff, st)
		}
		return err
	// for xmlquery.ElementNode
	case xmlquery.ElementNode:
		if strings.EqualFold("where", st) || strings.EqualFold("set", st) {
			//buff.WriteString(st)
			normalizeSql(buff, st)
		} else if strings.EqualFold("include", st) {
			if node := clause.findChildById(current.SelectAttr("refid")); node != nil {
				//xml.EscapeText(buff, []byte(b.InnerText()))
				//normalizeSql(buff, b.InnerText())
				for node = node.FirstChild; node != nil; node = node.NextSibling {
					if err = clause.buildMapperNode(clause.doc, node, buff); err != nil {
						return err
					}
					for child := node.FirstChild; child != nil; child = child.NextSibling {
						if err = clause.buildMapperNode(node, child, buff); err != nil {
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
					normalizeSql(buff, st)
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
*/

var nestParamPattern = regexp.MustCompile(`#{\s+\.\w*}`)

func (clause Clause) evaluate(exps []string) ([]interface{}, error) {
	var rt []interface{}
	for _, exp := range exps {
		par := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(exp, "#{"), "}"))
		var v interface{}
		var err error
		if strings.Contains(par, ".") {
			// @todo should evaluate from map
		} else {
			if v, err = expr.Eval(par, clause.arg); err == nil && v != nil {
				rt = append(rt, v)
			} else {
				return nil, errors.Unwrap(fmt.Errorf("can not resolve: '%v', %w", par, err))
			}
		}
	}
	return rt, nil
}

func (clause *Clause) sqlParameter(node *xmlquery.Node) (string, error) {
	var buff bytes.Buffer
	var err error
	exp := node.Data
	parent := node.Parent
	exp = strings.TrimSpace(spacePattern.ReplaceAllString(exp, " "))
	if params := paramPattern.FindAllString(exp, -1); len(params) > 0 {
		if strings.EqualFold(parent.Data, "foreach") {
			slicedArg := clause.arg
			collection := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(node.SelectAttr("collection"), "#{"), "}"))
			//item := strings.TrimSpace(node.SelectAttr("item"))
			separator := strings.TrimSpace(node.SelectAttr("separator"))
			if len(separator) == 0 {
				separator = ","
			}
			if len(collection) > 0 {
				if slicedArg, err = expr.Eval(collection, slicedArg); err != nil ||
					reflect.ValueOf(slicedArg).Kind() != reflect.Slice {
					return "", fmt.Errorf("mapper#%v: foreach statement, collection property must be a slice: `%v`, %w", clause.id, collection, err)
				}
			} else if reflect.ValueOf(slicedArg).Kind() != reflect.Slice {
				return "", fmt.Errorf("mapper#%v: foreach statement, parameter %+v is not a slice", clause.id, slicedArg)
			}
			s := reflect.ValueOf(slicedArg)
			for i := 0; i < s.Len(); i++ {
				if values, err := clause.evaluate(params, s.Index(i).Interface()); err == nil && len(values) > 0 {
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
			if values, err := clause.evaluate(params); err == nil && len(values) > 0 {
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

/*
func (clause *Clause) normalizeSql(node *xmlquery.Node) {
	str := strings.TrimSpace(node.Data)

	//values := clause.ctx.Value(fmt.Sprintf("%v_", scopePrefix))
	parent := node.Parent
	if strings.EqualFold("foreach", parent.Data) {
		item := parent.SelectAttr("item")
		values = clause.ctx.Value(fmt.Sprintf("%v_%v", scopePrefix, item))
	}


	for _, par := range paramPattern.FindAllString(str, -1) {
		par = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(par, "#{"), "}"))
		var v interface{}
		items := strings.Split(par, ".")
		ctx := clause.arg
		if len(items) > 1 {
			if ctx = clause.argMap[items[0]]; ctx == nil {
				return nil, fmt.Errorf("can not find parameter %v", items[0])
			}
		}
		if v, err = expr.Eval(items[0], ctx); err == nil && v != nil {
			rt = append(rt, v)
			if !test {
				// @todo 1: add sqlMap 2: sqlstatement
			}
		} else {
			return nil, errors.Unwrap(fmt.Errorf("can not resolve: '%v', %w", par, err))
		}

		return rt, err
	}
	if len(str) > 0 {
		if sqlWherePattern.MatchString(clause.statement.String()) && sqlAnd.MatchString(str) {
			str = sqlAnd.ReplaceAllString(str, "")
		}
		clause.statement.WriteString(str + " ")
	}
}
*/
