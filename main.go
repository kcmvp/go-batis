package batis

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"os"
	"strings"
	"sync"
)

var (
	once sync.Once
	db   *sql.DB
	Setting Settings
)

func init() {

	viper.SetConfigName("batis")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	env := "dev"
	if envVar := os.Getenv("env"); envVar != "" {
		env = strings.ToLower(envVar)
	}
	viper.SetConfigName(env)
	viper.MergeInConfig()

}

type settings struct {
	maxOpenConns          int
	maxIdleConns          int
	maxTransactionRetries int
	mapperDir             string
}

func (s settings) MaxOpenConns() int {
	return s.maxOpenConns
}

func (s settings) MaxIdleConns() int {
	return s.maxIdleConns
}

func (s settings) MaxTransactionRetries() int {
	return s.maxTransactionRetries
}

func (s settings) MapperDir() string {
	return s.mapperDir
}

type Settings interface {
	MaxOpenConns() int
	MaxIdleConns() int
	MaxTransactionRetries() int
	MapperDir() string
}

func DB() *sql.DB {
	once.Do(func() {
		if ds, err := sql.Open("", ""); err != nil {
			db = ds
			db.SetMaxIdleConns(Setting.MaxIdleConns())
			db.SetMaxOpenConns(Setting.MaxOpenConns())
		}
	})
	return db
}
