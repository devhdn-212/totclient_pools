package repository

import "github.com/doug-martin/goqu/v9"

type DBExecutor interface {
	From(table ...interface{}) *goqu.SelectDataset
	Insert(table interface{}) *goqu.InsertDataset
	Update(table interface{}) *goqu.UpdateDataset
	Delete(table interface{}) *goqu.DeleteDataset
}
