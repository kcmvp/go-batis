package plugin

import (
	"database/sql"
	"time"
)

type BaseEntity struct {
	CreatedAt sql.NullTime `json:"created_at" db:"name=createdAt,PK"`
	CreatedBy string       `json:"created_by"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	UpdatedBy string       `json:"updated_by"`
	When      time.Time
}

type OrderHeader struct {
	Sku     string
	BatchNo *string
	Seq     uint8
	BaseEntity
}

type Order struct {
	Id        int16 `json:"id" db:"name=id, PK"`
	CustNo    string
	BasicInfo OrderHeader `db:"join=Sku, BatchNo"`
	OrderNum  string
	OrderQty  int
	Price     float32
	BaseEntity
}
