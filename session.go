package sqlx



type Session interface {

	Exec(mapper Mapper)

}

