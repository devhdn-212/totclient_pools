package repository

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/doug-martin/goqu/v9"
)

type userRepository struct {
	db *goqu.Database
}

func NewUser(con *sql.DB) domain.UserRepository {
	return &userRepository{
		db: goqu.New("default", con),
	}
}
func (u userRepository) FindByEmail(ctx context.Context, email string) (result domain.User, err error) {
	dataset := u.db.From("tbl_users").
		Where(
			goqu.C("email").Eq(email))

	_, err = dataset.ScanStructContext(ctx, &result)
	return
}
