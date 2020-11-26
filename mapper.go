package batis

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"os"
	"path/filepath"
	"strings"
)

type SqlMapper string

type SqlHookFunc func(string) string

func (mapper SqlMapper) buildWithHook(cfg Config, hook SqlHookFunc, args ...interface{})  {

}

// mapper file naming pattern is ${struct}Mapper.xml
// naming standard of mapper is ${file name}.${mapper id}
// ex: `dog.findByName` means its definition in the `dog.xml` and the `id' attribute is `findByName`
func (mapper SqlMapper) build(cfg Config, args ...interface{}) (*Clause, error) {


	if c, ok := clauseCache.Get(mapper.mapperName); ok {
		return c.(*Clause), nil
	}

	if entries := strings.Split(mapper.mapperName, "."); len(entries) == 2 {
		path, err := filepath.Abs(fmt.Sprintf("%v/%vMapper.xml", settings.MapperDir(), entries[0]))
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
		node := xmlquery.FindOne(root, fmt.Sprintf("//%v[@id='%v']", mapper.MapperType.name(), entries[1]))
		if node == nil {
			return nil, errors.New("can't find the node")
		}
		var c Clause
		xmlNode := node.OutputXML(true)
		if err = xml.Unmarshal([]byte(xmlNode), &c); err == nil {
			if TestEnv() {
				clauseCache.Set(mapper.mapperName, &c, 0)
			} else {
				defer func() {
					clauseCache.Set(mapper.mapperName, &c, 0)
				}()
			}
		} else {
			//@todo add log info
		}
		defer func() {
			f.Close()
		}()
		return &c, nil

	} else {
		return nil, errors.New("can't find the dao")
	}

}
