package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/util"
)

type riwayatTransaksiService struct {
	pasaranService      domain.PasaranService
	trxkeluaranRepo     domain.TrxkeluaranRepository
	trxkeluarandetailSv domain.TrxkeluarandetailService
	trxkeluaranmemberSv domain.TrxkeluaranmemberService
	memberinfoService   domain.MemberinfoService
}

func NewRiwayatTransaksiService(
	pasaranService domain.PasaranService,
	trxkeluaranRepo domain.TrxkeluaranRepository,
	trxkeluarandetailSv domain.TrxkeluarandetailService,
	trxkeluaranmemberSv domain.TrxkeluaranmemberService,
	memberinfoService domain.MemberinfoService,
) domain.RiwayatTransaksiService {
	return &riwayatTransaksiService{
		pasaranService:      pasaranService,
		trxkeluaranRepo:     trxkeluaranRepo,
		trxkeluarandetailSv: trxkeluarandetailSv,
		trxkeluaranmemberSv: trxkeluaranmemberSv,
		memberinfoService:   memberinfoService,
	}
}

func (s *riwayatTransaksiService) Fetch(ctx context.Context, req dto.RiwayatTransaksiRequest) (dto.RiwayatTransaksiResponse, error) {
	idcomp := strings.ToUpper(req.Company)
	pasaranCode := strings.ToUpper(req.PasaranCode)

	memberinfo, err := s.memberinfoService.CheckToken(ctx, dto.MemberinfoResponse{
		Agen:   idcomp,
		Market: pasaranCode,
		Token:  req.Token,
	})
	if err != nil {
		return dto.RiwayatTransaksiResponse{}, err
	}

	// Cached (24h TTL + singleflight) — just confirming the pasaran exists
	// here, no need for a fresh DB hit every call.
	if _, err := s.pasaranService.FindID(ctx, idcomp, pasaranCode); err != nil {
		return dto.RiwayatTransaksiResponse{}, err
	}

	trxkeluaran, err := s.trxkeluaranRepo.FindByID(ctx, idcomp, req.PasaranIdcomp)
	if err != nil {
		return dto.RiwayatTransaksiResponse{}, err
	}
	if trxkeluaran.ID == 0 {
		return dto.RiwayatTransaksiResponse{}, fmt.Errorf("trxkeluaran %w", util.ErrNotFound)
	}

	members, err := s.trxkeluaranmemberSv.AllByUsername(ctx, idcomp, trxkeluaran.ID, memberinfo.Username)
	if err != nil {
		return dto.RiwayatTransaksiResponse{}, err
	}
	details, err := s.trxkeluarandetailSv.AllByUsername(ctx, idcomp, trxkeluaran.ID, memberinfo.Username)
	if err != nil {
		return dto.RiwayatTransaksiResponse{}, err
	}

	return dto.RiwayatTransaksiResponse{Members: members, Details: details}, nil
}
