package sql

import (
	"github.com/kcmvp/go-batis/session"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitSession(t *testing.T) {
	viper := session.InitDefault()
	assert := assert.New(t)
	assert.NotNil(viper)
}
