package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/devhdn-212/totclient_api/internal/util"
)

const RedisResult = "client:result"

type resultService struct {
	pasaranService    domain.PasaranService
	trxkeluaranRepo   domain.TrxkeluaranRepository
	memberinfoService domain.MemberinfoService
}

func NewResultService(
	pasaranService domain.PasaranService,
	trxkeluaranRepo domain.TrxkeluaranRepository,
	memberinfoService domain.MemberinfoService,
) domain.ResultService {
	return &resultService{
		pasaranService:    pasaranService,
		trxkeluaranRepo:   trxkeluaranRepo,
		memberinfoService: memberinfoService,
	}
}

func (s *resultService) Fetch(ctx context.Context, req dto.ResultRequest) (dto.ResultResponse, error) {
	idcomp := strings.ToUpper(req.Company)
	pasaranCode := strings.ToUpper(req.PasaranCode)

	if _, err := s.memberinfoService.CheckToken(ctx, dto.MemberinfoResponse{
		Agen:   idcomp,
		Market: pasaranCode,
		Token:  req.Token,
	}); err != nil {
		return dto.ResultResponse{}, err
	}

	// Cached (24h TTL + singleflight) — just confirming the pasaran exists
	// here, no need for a fresh DB hit every call.
	if _, err := s.pasaranService.FindID(ctx, idcomp, pasaranCode); err != nil {
		return dto.ResultResponse{}, err
	}

	// Month is always within the CURRENT year — never client-controlled —
	// only which month (1-12) is. Out-of-range/omitted defaults to the
	// current month.
	now := util.GetNowJakarta()
	month := req.Bulan
	if month < 1 || month > 12 {
		month = int(now.Month())
	}
	year := now.Year()
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, util.LocJakarta)
	end := start.AddDate(0, 1, 0)

	cacheKey := fmt.Sprintf("%s:%s:%s:%d-%02d", RedisResult, strings.ToLower(idcomp), req.PasaranIdcomp, year, month)
	cached, found, err := connection.GetRedis(cacheKey)
	if err != nil {
		return dto.ResultResponse{}, err
	}
	var results []dto.ResultItemData
	if found {
		if err := json.Unmarshal([]byte(cached), &results); err == nil {
			connection.Log.Info("Returning data from Redis - Result")
			return dto.ResultResponse{Results: results}, nil
		}
	}

	rows, err := s.trxkeluaranRepo.FindResultsByMonth(ctx, idcomp, req.PasaranIdcomp, start, end)
	if err != nil {
		return dto.ResultResponse{}, err
	}
	for _, v := range rows {
		results = append(results, dto.ResultItemData{
			Periode:          fmt.Sprintf("#%d-%d", v.IDtrxkeluaran, v.Keluaranperiode),
			Datekeluaran:     v.Datekeluaran.In(util.LocJakarta).Format("2006-01-02"),
			Keluarantogel:    v.Keluarantogel,
			Aliascomppasaran: v.Aliascomppasaran,
		})
	}
	go connection.SetRedis(cacheKey, results, 1*time.Hour)
	connection.Log.Info("Returning data Database - Result")
	return dto.ResultResponse{Results: results}, nil
}
