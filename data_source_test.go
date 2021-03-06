package batis

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type Company struct {
	CmpId          int32
	CompanyName    string `fake:"{randomstring:20}"`
	CompanyAddress string `fake:"{randomstring:30}"`
	SsNo           string `fake:"{randomstring:12}"`
}

type Department struct {
	DeptId   int32
	DeptName string `fake:"{randomstring:20}"`
	Status   int    `fake:"{number:0,1}"`
}

type Employee struct {
	EmpId     int32
	FirstName sql.NullString `fake:"{randomstring:20}"`
	LastName  string         `fake:"{randomstring:20}"`
	Birthday  sql.NullTime
	Salary    sql.NullFloat64
	Gender    int `fake:"{number:0,1}"`
	Status    int `fake:"{number:0,1}"`
}

type DataSourceTestSuit struct {
	suite.Suite
	ds *DataSource
}

func (s DataSourceTestSuit) BeforeTest(suiteName, testName string) {
	// @todo
	//clear up the data
}

func (s DataSourceTestSuit) SetupSuite() {
	// @todo
	// init schema
	if entries, err := os.ReadDir("./ddl"); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				ba, _ := os.ReadFile(filepath.Join("./ddl", entry.Name()))
				if _, err = s.ds.ExecSql(string(ba)); err != nil {
					a := 1 + 1
					fmt.Sprintf("failed to init schema %+v, %v", err, a)
				}
			}
		}
	}
}

func TestSuite(t *testing.T) {

	cfg := DbConfig{
		Url:        "file::memory:?cache=shared",
		DriverName: "sqlite3",
		MapperDir:  "./mapper",
	}
	if ds, err := NewDsDefaultCache(&cfg); err == nil {
		suite.Run(t, &DataSourceTestSuit{ds: ds})
	}
}

// Test return error when the mapper dir does not exit.
func TestNoneExistsMapperDir(t *testing.T) {
	cfg := DbConfig{
		Url:        "file::memory:?cache=shared",
		DriverName: "sqlite3",
		MapperDir:  "hello",
	}
	_, err := NewDsDefaultCache(&cfg)
	assert.NotNil(t, t, err)
}

func (s *DataSourceTestSuit) TestName() {
	s.NotNil(s.ds)
}
