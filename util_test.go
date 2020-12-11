package batis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var clauses = []struct {
	sql    string
	tables []string
}{
	{"from ta, tb where a = b", []string{"ta","tb"}},
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestGetTableList(t *testing.T) {
	assert := assert.New(t)
	for _, c := range clauses {
		t.Run(c.sql, func(t *testing.T) {
			tb := tableList(c.sql)
			assert.Equal(len(c.tables), len(tb))
			for _, table := range c.tables {
				assert.True(contains(c.tables,table))
			}
		})
	}
}

