package service

import (
	"context"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"

	"go.uber.org/zap"
)

const (
	RedisMemberInfo = "client:memberinfo"
)

type MemberinfoService struct {
}

func NewMemberinfoService() domain.MemberinfoService {
	return &MemberinfoService{}
}
func (d *MemberinfoService) CheckToken(ctx context.Context, req dto.MemberinfoResponse) (*domain.Memberinfo, error) {
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
