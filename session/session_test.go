package session

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type BaseEntity struct {
	Id int32 `db:"pk;seq"`
	CreatedAt sql.NullTime
	CreatedBy string `db:"size=20"`
	UpdatedAt sql.NullTime
	UpdatedBy string `db:"size=20"`
}

type Company struct {
	CompanyName    string
	CompanyAddress string
	SsNo           string
	BaseEntity
}

type Department struct {
	DeptName string
	Status   int
}

type Employee struct {
	FirstName string `db:"size=20"`
	LastName  string `db:"size=20"`
	Birthday  sql.NullTime
	Salary    sql.NullFloat64
	Gender    int
	Status    int
	BaseEntity
}

type DataSourceTestSuit struct {
	suite.Suite
	ds *Session
}

func (s DataSourceTestSuit) BeforeTest(suiteName, testName string) {
	// @todo
	//clear up the data
}

func (s DataSourceTestSuit) SetupSuite() {
	// @todo
	// init schema
	if entries, err := os.ReadDir("./schema/sqlite"); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				/*
				ba, _ := os.ReadFile(filepath.Join("./ddl", entry.Id()))
				if _, err = s.ds.ExecSql(string(ba)); err != nil {
					a := 1 + 1
					fmt.Sprintf("failed to init schema %+v, %v", err, a)
				}
				 */
			}
		}
	}
}

func TestSuite(t *testing.T) {
	cfg := Configuration{
		Url:        "file::memory:?cache=shared",
		DriverName: "sqlite3",
		//Id:       "./mapper",
	}
	suite.Run(t, &DataSourceTestSuit{ds: InitSessionDefault(&cfg)})
}


func (s *DataSourceTestSuit) TestName() {
	s.NotNil(s.ds)
}


