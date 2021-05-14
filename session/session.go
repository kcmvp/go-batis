package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

var once sync.Once
var sessionMap map[string]*Session
var mapperSession map[string]*Session

const MaxOpenConns = 100
const MaxIdleConns = 50

type SqlType string

var SqlTypes = []SqlType{"insert", "select", "delete", "update"}

type StatementHookFunc func(ctx context.Context, statement *string) error

type Configuration struct {
	Name                  string
	Url                   string
	UserName              string
	Password              string
	DriverName            string
	MaxOpenConns          int
	MaxIdleConns          int
	MaxTransactionRetries int
}

type Session struct {
	*sql.DB
	name     string
	cache    Cache
	sqlHooks []StatementHookFunc
}

func InitSessionDefault(cfg *Configuration) *Session {
	return InitSession(cfg, defaultCache())
}

func InitSession(cfg *Configuration, cache Cache) *Session {
	if err := validate(cfg); err != nil {
		log.Printf("Error: %v", err.Error())
		panic(err)
	}
	session := &Session{cache: cache, name: cfg.Name}
	if err := session.validate(); err != nil {
		panic(err)
	}
	if _, ok := sessionMap[cfg.Name]; !ok {
		if db, err := sql.Open(cfg.DriverName, cfg.Url); err != nil {
			panic(err)
		} else {
			db.SetMaxOpenConns(MaxOpenConns)
			db.SetMaxIdleConns(MaxIdleConns)
			if cfg.MaxIdleConns > 0 {
				db.SetMaxIdleConns(cfg.MaxIdleConns)
			} else if cfg.MaxOpenConns > 0 {
				v := math.Max(float64(cfg.MaxOpenConns), float64(cfg.MaxIdleConns+1))
				db.SetMaxOpenConns(int(v))
			}
			session.DB = db

			sessionMap[cfg.Name] = session
		}
	}
	return session
}

func validate(cfg *Configuration) error {
	_, err := os.Stat(cfg.Name)
	return err
}

func (session *Session) validate() error {
	var err error
	var mapperPath string
	if mapperPath, err = filepath.Abs(session.name); err != nil {
		return errors.Unwrap(fmt.Errorf("load mapper %w ", err))
	} else {
		log.Printf("loading mapper for: %v", mapperPath)
		return filepath.WalkDir(mapperPath, func(path string, d fs.DirEntry, err error) error {
			if err == nil {
				if d.IsDir() {
					return fs.SkipDir
				}
				mp := filepath.Join(mapperPath, path)
				session.name = mp
				if mf, e := os.OpenFile(mp, os.O_RDONLY, 0666); e != nil {
					if root, e := xmlquery.Parse(mf); e != nil {
						return errors.Unwrap(fmt.Errorf("failed to parse the file %v %w ", d.Name(), e))
					} else {
						for _, node := range xmlquery.Find(root, "//*/@id") {
							for _, sqlType := range SqlTypes {
								if node.Data == string(sqlType) {
									m := strings.TrimSpace(node.SelectAttr("id"))
									if _, ok := mapperSession[m]; !ok {
										mapperSession[m] = session
									} else {
										panic(fmt.Errorf("duplicated sql statement : %v", m))
									}
								}
							}
						}
					}
				} else {
					return errors.Unwrap(fmt.Errorf("failed to open file %w ", e))
				}
			}
			return err
		})
	}
}

func GetSession(name string) *Session {
	if rt, ok := sessionMap[name]; ok {
		return rt
	} else {
		panic(fmt.Sprintf("failed to get session with name %v", name))
	}
}

func MapperSession(mapperName string) *Session {
	if v, o := mapperSession[mapperName]; !o {
		panic(fmt.Errorf("can not find the session for the mapper %v", mapperName))
	} else {
		return v
	}
}

func AllSessions() []*Session {
	var rt []*Session
	for _, session := range sessionMap {
		rt = append(rt, session)
	}
	return rt
}

func (session Session) Name() string {
	return session.name
}

func (session *Session) Set(key string, value interface{}) error {
	return session.cache.Set(key, value)
}

func (session *Session) Get(key string) (interface{}, error) {
	return session.cache.Get(key)
}

func (session *Session) Del(key string) error {
	return session.cache.Del(key)
}

func (session *Session) Ttl(key string, value interface{}, duration time.Duration) error {
	return session.cache.Ttl(key, value, duration)
}

func (session *Session) SetSqlHook(h ...StatementHookFunc) {
	session.sqlHooks = append(session.sqlHooks, h...)
}

func (session *Session) SetCache(cache Cache) {
	if session.cache != nil {
		session.cache = cache
	} else {
		log.Printf("Warning: there is a cache has been attahced to this session already")
	}
}

func InitDefault() *viper.Viper {
	v := viper.New()
	v.SetConfigName("batis")
	v.SetConfigType("yml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %s \n", err)
	}
	Init(v)
	return v
}

func Init(v *viper.Viper) {
	dsViper := v.Sub("datasource")
	if dsViper == nil {
		log.Panic("Can not find the 'database' section in the configuration.")
	}
	cfgMap := map[string]*Configuration{}
	once.Do(func() {
		for _, k := range dsViper.AllKeys() {
			name := strings.Split(k, ".")[0]
			if nil == cfgMap[name] {
				c := &Configuration{}
				if e := dsViper.Sub(name).Unmarshal(c); e != nil {
					log.Panicf("failed to parse database configuration %v", e)
				}
				c.Name = name
				cfgMap[name] = c
				InitSessionDefault(c)
			}
		}
	})
	log.Printf("datasource %s initialized successfully", reflect.ValueOf(cfgMap).MapKeys())
}

func (s Session) statementHook(ctx context.Context, statement *string) error {
	for _, hook := range s.sqlHooks {
		if err := hook(ctx, statement); err != nil {
			log.Printf("error statement: %v; %w", statement, err)
			return err
		}
	}
	return nil
}

func (s Session) QueryCacheable(ctx context.Context, statement, key string, arg interface{}) (interface{}, error) {
	if len(key) > 0 {
		if v, err := s.Get(key); err != nil {
			goto withCache
		} else {
			return v, nil
		}
	} else {
		goto withCache
	}

withCache:
	s.statementHook(ctx, &statement)
	if v, err := s.Query(statement, arg); err == nil {
		defer func() {
			s.Set(key, v)
		}()
		return v, nil
	} else {
		return nil, err
	}
}

func (s Session) ExecCacheable(ctx context.Context, statement, key string, arg interface{}) (sql.Result, error) {
	if len(key) > 0 {
		s.Del(key)
	}
	s.statementHook(ctx, &statement)
	return s.Exec(statement, arg)
}
