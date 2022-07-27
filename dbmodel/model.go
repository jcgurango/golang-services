package dbmodel

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel

	ID       int64  `bun:",pk,autoincrement"`
	Username string `bun:",unique,notnull"`
	Password string `bun:",notnull"`
}

type Account struct {
	bun.BaseModel

	ID   int64  `bun:",pk,autoincrement" json:"id"`
	Name string `bun:",notnull" json:"name"`
	User int64  `bun:",notnull"`
}

type Transaction struct {
	bun.BaseModel

	ID            int64  `bun:",pk,autoincrement"`
	Detail        string `bun:",notnull"`
	Amount        string `bun:"type:money"`
	DebitAccount  int64  `bun:",notnull"`
	CreditAccount int64  `bun:",notnull"`
}
