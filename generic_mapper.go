package batis

import (
	"database/sql"
)

type Finder SqlMapper
type Searcher SqlMapper
type Updater SqlMapper

func (f Finder) Exec(arg interface{}) interface{} {
	panic("implement me")
}

func (s Searcher) Exec(arg interface{})[]interface{}  {
	panic("implement me")
}

func (u Updater) Exec(arg interface{}) (sql.Result, error) {
	panic("implement me")
}

