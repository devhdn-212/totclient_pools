package repository

import (
	"context"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"gofibergocu/domain"
	"time"
)

type customerRepository struct {
	exec DBExecutor
}

func NewCustomer(exec DBExecutor) domain.CustomerRepository {
	return &customerRepository{exec: exec}
}

func (cr customerRepository) FindAll(ctx context.Context) (result []domain.Customer, err error) {
	dataset := cr.exec.From("tbl_customer").Where(goqu.C("deleted_at").IsNull())
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

func (cr customerRepository) FindByID(ctx context.Context, id string) (result domain.Customer, err error) {
	dataset := cr.exec.From("tbl_customer").
		Where(
			goqu.C("deleted_at").IsNull(),
			goqu.C("id").Eq(id))

	_, err = dataset.ScanStructContext(ctx, &result)
	return
}
func (cr customerRepository) FindByCode(ctx context.Context, code string) (result domain.Customer, err error) {
	dataset := cr.exec.From("tbl_customer").
		Where(
			goqu.C("deleted_at").IsNull(),
			goqu.C("code").Eq(code))

	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

func (cr customerRepository) Save(ctx context.Context, c *domain.Customer) error {
	exec := cr.exec.Insert("tbl_customer").Rows(c).Executor()
	_, err := exec.ExecContext(ctx)
	return err
}

func (cr customerRepository) Update(ctx context.Context, c *domain.Customer) error {
	exec := cr.exec.Update("tbl_customer").
		Where(goqu.C("id").Eq(c.ID)).Set(c).Executor()
	_, err := exec.ExecContext(ctx)
	return err
}

func (cr customerRepository) Delete(ctx context.Context, id string) error {
	exec := cr.exec.Update("tbl_customer").
		Where(goqu.C("id").Eq(id)).
		Set(goqu.Record{"deleted_at": sql.NullTime{Valid: true, Time: time.Now()}}).
		Executor()
	_, err := exec.ExecContext(ctx)
	return err
}
