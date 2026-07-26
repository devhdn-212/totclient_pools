package service

import (
	"context"
	"strings"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
)

type riwayatTransaksiService struct {
	pasaranService      domain.PasaranService
	trxkeluarandetailSv domain.TrxkeluarandetailService
	trxkeluaranmemberSv domain.TrxkeluaranmemberService
	memberinfoService   domain.MemberinfoService
}

func NewRiwayatTransaksiService(
	pasaranService domain.PasaranService,
	trxkeluarandetailSv domain.TrxkeluarandetailService,
	trxkeluaranmemberSv domain.TrxkeluaranmemberService,
	memberinfoService domain.MemberinfoService,
) domain.RiwayatTransaksiService {
	return &riwayatTransaksiService{
		pasaranService:      pasaranService,
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

	// No idtrxkeluaran: the client is opening the "Transaksi" list, which
	// only needs the period list (level 1) — trxkeluaranmember aggregated
	// per idtrxkeluaran covers every period the player has ever transacted
	// in for this pasaran, not just the currently-open one.
	if req.Idtrxkeluaran == 0 {
		periods, err := s.trxkeluaranmemberSv.Periods(ctx, idcomp, req.PasaranIdcomp, memberinfo.Username)
		if err != nil {
			return dto.RiwayatTransaksiResponse{}, err
		}
		return dto.RiwayatTransaksiResponse{Periods: periods}, nil
	}

	// idtrxkeluaran set: the client opened one specific period (level 2) —
	// fetch only that period's bet rows. statuskeluarandetail here is the
	// real, agen-settled value (RUNNING/WINNER/LOSE), read as-is.
	details, err := s.trxkeluarandetailSv.AllByUsername(ctx, idcomp, req.PasaranIdcomp, []int{req.Idtrxkeluaran}, memberinfo.Username)
	if err != nil {
		return dto.RiwayatTransaksiResponse{}, err
	}

	return dto.RiwayatTransaksiResponse{Details: details}, nil
}
