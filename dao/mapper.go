package dao

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/dgraph-io/ristretto"
	"github.com/kcmvp/go-batis"
	. "github.com/kcmvp/go-batis/dao/internal/syntax"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var maxCacheMaxEntries = 100

var (
	once        sync.Once
	clauseCache *ristretto.Cache
)

type MapperType uint

const (
	Insert MapperType = iota
	Update
	Delete
	Find
	Search
)

func (mapperType MapperType) name() string {
	return [...]string{"insert", "update", "delete", "select", "select"}[mapperType]
}

type Mapper interface {
	Exec(dest interface{}, arg interface{}) error
	ExecWithTx(tx sql.Tx, dest interface{}, arg interface{}) error
}

func NewMapper(mType MapperType, name string /*, argType interface{}*/) Mapper {
	return &mapper{
		MapperType: mType,
		mapperName: name,
		//argType:    argType,
	}
}

type mapper struct {
	MapperType
	mapperName string
}


func (m mapper) ExecWithTx(tx sql.Tx, dest interface{}, arg interface{}) error {
	if clause, err := m.build(); err != nil {
		if err := clause.Build(arg); err != nil {
			//@todo panic
			return err
		}
	}
	return nil
}

func (m mapper) Exec(dest interface{}, arg interface{}) error {
	if clause, err := m.build(); err != nil {
		if err := clause.Build(arg); err != nil {
			return err
			panic(fmt.Sprintf("@todo %v", err))
		}
	}
	return nil
}

// mapper file naming pattern is ${struct}Mapper.xml
func (mapper mapper) build() (*Clause, error) {

	once.Do(func() {
		var err error
		clauseCache, err = ristretto.NewCache(&ristretto.Config{
			NumCounters: 10000,     //10K
			MaxCost:     100000000, //100MB
			BufferItems: 64,
		})
		if err != nil {
			panic(err)
		}
	})

	if c, ok := clauseCache.Get(mapper.mapperName); ok {
		return c.(*Clause), nil
	}

	if entries := strings.Split(mapper.mapperName, "."); len(entries) == 2 {
		path, err := filepath.Abs(fmt.Sprintf("%v/%vMapper.xml", batis.Config.MapperDir(), entries[0]))
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
			if batis.TestEnv() {
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
