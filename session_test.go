package batis

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSession(t *testing.T) {
	cfg := DbConfig{
		Url:        "file::memory:?cache=shared",
		DriverName: "sqlite3",
	}
	assertion := assert.New(t)
	session, err := NewSession(&cfg)
	assertion.NotNil(session)
	assertion.Nil(err)
}
