package batis

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"os"
	"sync"
)

const cfgName = "application"

var (
	once   sync.Once
	db     *sql.DB
	Config Settings
)

func initialization() {
	env := os.Getenv("env")
	//@todo init from etcd
	v := viper.New()
	v.SetConfigName(cfgName)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("../")
	err := v.ReadInConfig() // Find and read the initialization file
	if err != nil {         // Handle errors reading the initialization file
		panic(fmt.Errorf("Fatal error initialization file: %s \n", err))
	}
	v.SetConfigName(fmt.Sprintf("%v-%v", cfgName, env))
	if err := v.MergeInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error initialization file: %s \n", err))
	}
	v = v.Sub("batis")

	s := settings{}
	v.Unmarshal(&s)
	Config = &s
}

type settings struct {
	DriverName            string `mapstructure:"driverName"`
	Url                   string `mapstructure:"url"`
	MaxOpenConns          int    `mapstructure:"maxOpenConns"`
	MaxIdleConns          int    `mapstructure:"maxIdleConns"`
	MaxTransactionRetries int    `mapstructure:"maxTransactionRetries"`
	Mappers               string `mapstructure:"mapperDir"`
}

func (s settings) Driver() string {
	return s.DriverName
}

func (s settings) MaxOpen() int {
	return s.MaxOpenConns
}

func (s settings) MaxIdle() int {
	return s.MaxIdleConns
}

func (s settings) MaxRetries() int {
	return s.MaxTransactionRetries
}

func (s settings) DBUrl() string {
	return s.Url
}

func (s settings) MapperDir() string {
	return s.Mappers
}

type Settings interface {
	MaxOpen() int
	MaxIdle() int
	MaxRetries() int
	DBUrl() string
	MapperDir() string
	Driver() string
}

func DB() *sql.DB {
	once.Do(func() {
		initialization()
		var err error
		if db, err = sql.Open(Config.Driver(), Config.DBUrl()); err != nil {
			panic(fmt.Sprintf("failed to connect to database :%v", err.Error()))
		} else {
			db.SetMaxIdleConns(Config.MaxIdle())
			db.SetMaxOpenConns(Config.MaxOpen())
		}
	})
	return db
}
