package domain

import (
	"context"
	"database/sql"
	"gofibergocu/dto"
)

type Admin struct {
	Username  string       `db:"username"`
	Pass      string       `db:"password"`
	Idadmin   string       `db:"idadmin"`
	Name      string       `db:"name"`
	Status    string       `db:"statuslogin"`
	Lastlogin sql.NullTime `db:"lastlogin"`
	Joindate  sql.NullTime `db:"joindate"`
	Ipaddress string       `db:"ipaddress"`
	Timezone  string       `db:"timezone"`
	Created   string       `db:"createadmin"`
	CreatedAt sql.NullTime `db:"createdateadmin"`
	Update    string       `db:"updateadmin"`
	UpdateAt  sql.NullTime `db:"updatedateadmin"`
}
type AdminsRepository interface {
	FindAll(ctx context.Context) ([]Admin, error)
	FindByUsername(ctx context.Context, username string) (Admin, error)
	Save(ctx context.Context, admin *Admin) error
	Update(ctx context.Context, admin *Admin, flag bool) error
	UpdateLogin(ctx context.Context, admin *Admin) error
}
type AdminService interface {
	All(ctx context.Context) ([]dto.AdminData, error)
	Save(ctx context.Context, req dto.AdminSave, client string) error
}
