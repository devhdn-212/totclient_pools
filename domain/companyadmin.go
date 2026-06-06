package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/gofibermaster_api/dto"
)

type Companyadmin struct {
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
type CompanyadminRepository interface {
	FindAll(ctx context.Context, id string) ([]Companyadmin, error)
	FindByID(ctx context.Context, id, username string) (Companyadmin, error)
	Save(ctx context.Context, compadmin *Companyadmin) error
	Update(ctx context.Context, compadmin *Companyadmin, flag bool) error
}
type CompanyadminService interface {
	All(ctx context.Context, id string) ([]dto.CompanyadminData, error)
	Save(ctx context.Context, req dto.CompanyadminSave, client string) error
}
