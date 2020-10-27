package batis

import (
	"sync"
	"sync/atomic"
	"time"
)

// Settings defines methods to get or set configuration values.
type Settings interface {
	// SetPreparedStatementCache enables or disables the prepared statement
	// internal.
	SetPreparedStatementCache(bool)

	// PreparedStatementCacheEnabled returns true if the prepared statement internal
	// is enabled, false otherwise.
	PreparedStatementCacheEnabled() bool

	// SetConnMaxLifetime sets the default maximum amount of time a connection
	// may be reused.
	SetConnMaxLifetime(time.Duration)

	// ConnMaxLifetime returns the default maximum amount of time a connection
	// may be reused.
	ConnMaxLifetime() time.Duration

	// SetMaxIdleConns sets the default maximum number of connections in the idle
	// connection pool.
	SetMaxIdleConns(int)

	// MaxIdleConns returns the default maximum number of connections in the idle
	// connection pool.
	MaxIdleConns() int

	// SetMaxOpenConns sets the default maximum number of open connections to the
	// database.
	SetMaxOpenConns(int)

	// MaxOpenConns returns the default maximum number of open connections to the
	// database.
	MaxOpenConns() int

	// SetMaxTransactionRetries sets the number of times a transaction can
	// be retried.
	SetMaxTransactionRetries(int)

	// MaxTransactionRetries returns the maximum number of times a
	// transaction can be retried.
	MaxTransactionRetries() int

	SetMapperDir(string)
}

type settings struct {
	sync.RWMutex

	preparedStatementCacheEnabled uint32

	connMaxLifetime time.Duration
	maxOpenConns    int
	maxIdleConns    int

	maxTransactionRetries int
	mapperDir             string
}

func (c *settings) binaryOption(opt *uint32) bool {
	return atomic.LoadUint32(opt) == 1
}

func (c *settings) setBinaryOption(opt *uint32, value bool) {
	if value {
		atomic.StoreUint32(opt, 1)
		return
	}
	atomic.StoreUint32(opt, 0)
}

func (c *settings) SetPreparedStatementCache(value bool) {
	c.setBinaryOption(&c.preparedStatementCacheEnabled, value)
}

func (c *settings) PreparedStatementCacheEnabled() bool {
	return c.binaryOption(&c.preparedStatementCacheEnabled)
}

func (c *settings) SetConnMaxLifetime(t time.Duration) {
	c.Lock()
	c.connMaxLifetime = t
	c.Unlock()
}

func (c *settings) ConnMaxLifetime() time.Duration {
	c.RLock()
	defer c.RUnlock()
	return c.connMaxLifetime
}

func (c *settings) SetMaxIdleConns(n int) {
	c.Lock()
	c.maxIdleConns = n
	c.Unlock()
}

func (c *settings) MaxIdleConns() int {
	c.RLock()
	defer c.RUnlock()
	return c.maxIdleConns
}

func (c *settings) SetMaxTransactionRetries(n int) {
	c.Lock()
	c.maxTransactionRetries = n
	c.Unlock()
}

func (c *settings) MaxTransactionRetries() int {
	c.RLock()
	defer c.RUnlock()
	if c.maxTransactionRetries < 1 {
		return 1
	}
	return c.maxTransactionRetries
}

func (c *settings) SetMaxOpenConns(n int) {
	c.Lock()
	c.maxOpenConns = n
	c.Unlock()
}

func (c *settings) MaxOpenConns() int {
	c.RLock()
	defer c.RUnlock()
	return c.maxOpenConns
}

func (c *settings) SetMapperDir(mapperDir string) {
	c.Lock()
	c.mapperDir = mapperDir
	c.Unlock()
}

func (c *settings) MapperDir() string {
	c.RLock()
	defer c.RUnlock()
	return c.mapperDir
}

func NewSettings() Settings {
	def := DefaultSettings.(*settings)
	return &settings{
		preparedStatementCacheEnabled: def.preparedStatementCacheEnabled,
		connMaxLifetime:               def.connMaxLifetime,
		maxIdleConns:                  def.maxIdleConns,
		maxOpenConns:                  def.maxOpenConns,
		maxTransactionRetries:         def.maxTransactionRetries,
		mapperDir: def.mapperDir,
	}
}

var DefaultSettings Settings = &settings{
	preparedStatementCacheEnabled: 0,
	connMaxLifetime:               time.Duration(0),
	maxIdleConns:                  10,
	maxOpenConns:                  0,
	maxTransactionRetries:         1,
	mapperDir:                     "mappers",
}
