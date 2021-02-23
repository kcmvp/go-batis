package batis

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SessionTestSuit struct {
	suite.Suite
	DbSession *Session
}

func (s SessionTestSuit) BeforeTest(suiteName, testName string) {
	//clear up the data
}

func (s SessionTestSuit) SetupSuite() {
	s.Assert().NotNil(s.DbSession)
	// init schema
}

func TestExampleTestSuite(t *testing.T) {
	cfg := DbConfig{
		Url:        "file::memory:?cache=shared",
		DriverName: "sqlite3",
	}
	if session, err := NewSession(&cfg); err == nil {
		suite.Run(t, &SessionTestSuit{DbSession: session})
	}
}

func (s *SessionTestSuit) TestName() {
	s.Equal(1,2)
}




