package batis
import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)



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
}

func TestSuite(t *testing.T) {

	cfg := DbConfig{
		Url:        "file::memory:?cache=shared",
		DriverName: "sqlite3",
		MapperDir: "./mapper",
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
		MapperDir: "hello",
	}
	_, err := NewDsDefaultCache(&cfg);
	assert.NotNil(t, t, err)
}

func (s *DataSourceTestSuit) TestName() {
	s.NotNil(s.ds)
}
