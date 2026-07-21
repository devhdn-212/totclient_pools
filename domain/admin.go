package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/totclient_api/dto"
)

type Admin struct {
	ID           string       `db:"idcompadmin"`
	IDCompany    string       `db:"idcompany"`
	IDClientrule string       `db:"idclientrule"`
	Username     string       `db:"usernamecompadmin"`
	Pass         string       `db:"passcompadmin"`
	Name         string       `db:"namecompadmin"`
	Ipaddress    string       `db:"ipaddresscompadmin"`
	Lastlogin    sql.NullTime `db:"lastlogincompadmin"`
	Status       string       `db:"compadminstatus"`
	Created      string       `db:"createcompadmin"`
	CreatedAt    sql.NullTime `db:"createdatecompadmin"`
	Update       string       `db:"updatecompadmin"`
	UpdateAt     sql.NullTime `db:"updatedatecompadmin"`
}
type AdminsRepository interface {
	FindAll(ctx context.Context, idcomp string) ([]Admin, error)
	FindByUsernameComp(ctx context.Context, username, idcomp string) (Admin, error)
	FindByUsername(ctx context.Context, username string) (Admin, error)
	Save(ctx context.Context, admin *Admin) error
	Update(ctx context.Context, admin *Admin, flag bool) error
	UpdateLogin(ctx context.Context, admin *Admin) error
}
type AdminService interface {
	All(ctx context.Context, idcomp string) ([]dto.AdminData, error)
	Save(ctx context.Context, req dto.AdminSave, client, idcomp string) error
}
