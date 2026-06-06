package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/gofibermaster_api/domain"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db DBExecutor
}

func NewUser(db DBExecutor) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}
func (u userRepository) FindByEmail(ctx context.Context, email string) (result domain.User, err error) {
	// Query native PostgreSQL
	query := `SELECT id, email, password, name, role, created_at, updated_at 
	          FROM tbl_users 
	          WHERE email = $1 LIMIT 1`

	rows, err := u.db.Query(ctx, query, email)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	// Mapping otomatis menggunakan fitur CollectOneRow dari pgx v5
	result, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, nil
		}
		return result, err
	}

	return result, nil
}
