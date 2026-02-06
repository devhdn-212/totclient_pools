package repository

import (
	"database/sql"
	"github.com/doug-martin/goqu/v9"
)

type GoquExecutor struct {
	*goqu.Database
}

func NewGoquExecutor(db *sql.DB) *GoquExecutor {
	return &GoquExecutor{
		Database: goqu.New("default", db),
	}
}

type DBExecutor interface {
	From(table ...interface{}) *goqu.SelectDataset
	Insert(table interface{}) *goqu.InsertDataset
	Update(table interface{}) *goqu.UpdateDataset
	Delete(table interface{}) *goqu.DeleteDataset
}
