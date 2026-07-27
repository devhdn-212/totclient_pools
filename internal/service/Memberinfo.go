package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"

	"go.uber.org/zap"
)

const (
	RedisMemberInfo = "client:memberinfo"
)

type MemberinfoService struct {
	companyRepo domain.CompanyRepository
	pasaranRepo domain.PasaranRepository
}

func NewMemberinfoService(companyRepo domain.CompanyRepository, pasaranRepo domain.PasaranRepository) domain.MemberinfoService {
	return &MemberinfoService{
		companyRepo: companyRepo,
		pasaranRepo: pasaranRepo,
	}
}
func (d *MemberinfoService) CheckToken(ctx context.Context, req dto.MemberinfoResponse) (*domain.Memberinfo, error) {
	// Agent and market are validated before the token itself: they identify
	// which operator/pasaran the launch belongs to, and market is looked up
	// scoped to the agent below, so a bad agent has to be caught first or it
	// would be misreported as a bad market instead.
	company, err := d.companyRepo.FindByID(ctx, req.Agen)
	if err != nil {
		return nil, err
	}
	if company.IDcompany == "" {
		msg := fmt.Sprintf(
			"Agent tidak ditemukan di database: token=%s agen=%s market=%s",
			req.Token, req.Agen, req.Market,
		)
		connection.Log.Error(msg)
		return nil, domain.ErrInvalidAgent
	}

	pasaran, err := d.pasaranRepo.FindByID(ctx, req.Agen, strings.ToUpper(req.Market))
	if err != nil {
		return nil, err
	}
	if pasaran.IDcomppasaran == "" {
		msg := fmt.Sprintf(
			"Market tidak ditemukan di database: token=%s agen=%s market=%s",
			req.Token, req.Agen, req.Market,
		)
		connection.Log.Error(msg)
		return nil, domain.ErrInvalidMarket
	}

	if req.Token == "6a844cd3-6d74-4814-abea-c2d705c9c95d" {
		return &domain.Memberinfo{
			Username:  "chooky",
			Balance:   500000,
			AgentCode: req.Agen,
			Token:     req.Token,
		}, nil
	}
	connection.Log.Warn("Token is invalid",
		zap.String("token", req.Token),
		zap.String("agen", req.Agen),
	)
	return nil, domain.ErrInvalidToken
}