package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/devhdn-212/totclient_pools/domain"
	"github.com/devhdn-212/totclient_pools/dto"
	"github.com/devhdn-212/totclient_pools/internal/connection"
	"github.com/devhdn-212/totclient_pools/internal/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

const (
	RedisPasaran = "client:pasaran"
)

type pasaranService struct {
	db      *pgxpool.Pool
	repo    domain.PasaranRepository
	trxRepo domain.TrxkeluaranRepository
	sf      singleflight.Group
}

func NewPasaranService(db *pgxpool.Pool, repo domain.PasaranRepository, trxRepo domain.TrxkeluaranRepository) domain.PasaranService {
	return &pasaranService{
		db:      db,
		repo:    repo,
		trxRepo: trxRepo,
	}
}

func (u *pasaranService) FindID(ctx context.Context, idcomp, codepasaran string) (dto.PasaranData, error) {
	redisKey := RedisPasaran + ":" + strings.ToLower(idcomp) + ":" + strings.ToLower(codepasaran)

	cached, found, err := connection.GetRedis(redisKey)
	if err != nil {
		return dto.PasaranData{}, err
	}
	var record dto.PasaranData
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			connection.Log.Info("Returning data from Redis - Pasaran")
			return record, nil
		}
	}

	result, err, _ := u.sf.Do(redisKey, func() (any, error) {
		return u.fetchPasaranData(ctx, idcomp, codepasaran, redisKey)
	})
	if err != nil {
		return dto.PasaranData{}, err
	}
	return result.(dto.PasaranData), nil
}

func (u *pasaranService) fetchPasaranData(ctx context.Context, idcomp, codepasaran, redisKey string) (dto.PasaranData, error) {
	var record dto.PasaranData

	v, err := u.repo.FindByID(ctx, idcomp, strings.ToUpper(codepasaran))
	if err != nil {
		connection.Log.Error("pasaranRepository.FindByID failed", zap.Error(err))
		return dto.PasaranData{}, err
	}
	if v.IDcomppasaran == "" {
		return dto.PasaranData{}, fmt.Errorf("Pasaran %w", util.ErrNotFound)
	}

	jadwal, err := u.repo.FindJadwal(ctx, v.IDcomppasaran)
	if err != nil {
		connection.Log.Error("pasaranRepository.FindJadwal failed", zap.Error(err))
		return dto.PasaranData{}, err
	}

	trxkeluaran, err := u.trxRepo.FindByID(ctx, strings.ToUpper(idcomp), v.IDcomppasaran)
	if err != nil {
		connection.Log.Error("trxkeluaranRepository.FindByID failed", zap.Error(err))
		return dto.PasaranData{}, err
	}

	now := util.GetNowJakarta()
	tglKeluaran := trxkeluaran.Datekeluaran
	tglAwal := time.Date(tglKeluaran.Year(), tglKeluaran.Month(), tglKeluaran.Day(), 0, 0, 0, 0, util.LocJakarta)
	tglHariIni := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, util.LocJakarta)

	status := "ONLINE"
	var jadwalOpen, jadwalTutup string
	var openTime, tutupTime time.Time
	var jamOpenValid, jamTutupValid bool

	hariKeluaran := util.HariIndonesia(tglKeluaran)

	// Ambil satu baris jadwal yang cocok dengan hari draw ini saja (break setelah
	// ketemu) - supaya openTime/tutupTime yang dipakai buat keputusan ONLINE/OFFLINE
	// dan yang ditampilkan (jadwalOpen/jadwalTutup di caller) selalu berasal dari
	// baris yang sama (kalau ada >1 baris utk hari yang sama, tidak nyangkut campuran
	// dua baris berbeda).
	for _, d := range jadwal {
		if d.Haripasaran != hariKeluaran {
			continue
		}
		if d.Jamopen.Valid {
			// Jamopen dan Jamtutup dua-duanya selalu di tanggal draw ini sendiri
			// (tglAwal) apa adanya - walau time-of-day-nya jamopen > jamtutup (mis.
			// buka 18:30, tutup 17:30), itu BUKAN sinyal "buka kemarin malam". Sempat
			// dicoba di-shift -1 hari (asumsi awal yang ternyata salah), dan sempat
			// juga dicoba anchor jeda maintenance ke tglHariIni (juga salah) -
			// dikonfirmasi user semuanya harus anchor ke tglAwal/Datekeluaran.
			openTime = tglAwal.Add(time.Duration(d.Jamopen.Microseconds) * time.Microsecond)
			jamOpenValid = true
		}
		if d.Jamtutup.Valid {
			tutupTime = tglAwal.Add(time.Duration(d.Jamtutup.Microseconds) * time.Microsecond)
			jamTutupValid = true
		}
		break
	}

	if jamOpenValid {
		jadwalOpen = openTime.Format("2006-01-02 15:04:05")
	}
	if jamTutupValid {
		jadwalTutup = tutupTime.Format("2006-01-02 15:04:05")
	}

	// Port dari logika lama (goment) yang sudah terbukti jalan di produksi:
	//   - tutupTime.Before(tglHariIni): draw ini (tglAwal) sudah lewat harinya sama
	//     sekali dari hari ini -> pasti OFFLINE (setara "dateNow > lastDateOpen").
	//   - !now.Before(tutupTime) && !now.After(openTime): now ada di jeda antara jam
	//     tutup dan jam buka draw ini (openTime, setara "lastScheduleOpen"/
	//     "pasaran_marketopen") -> OFFLINE (jeda maintenance).
	//   - !now.Before(openTime): now sudah lewat openTime draw ini -> ronde
	//     berikutnya sudah buka, draw ini otomatis OFFLINE (sudah digantikan).
	if tutupTime.Before(tglHariIni) {
		status = "OFFLINE"
	} else if jamOpenValid && jamTutupValid {
		if !now.Before(tutupTime) && !now.After(openTime) {
			status = "OFFLINE"
		} else if !now.Before(openTime) {
			status = "OFFLINE"
		} else {
			status = "ONLINE"
		}
	}

	record = dto.PasaranData{
		IDcomppasaran:                    v.IDcomppasaran,
		Codepasaran:                      v.Codecomppasaran,
		Aliascomppasaran:                 v.Aliascomppasaran,
		URLlogo:                          v.URLlogo,
		AngkaMinbasket:                   v.AngkaMinbasket,
		AngkaMinbet:                      v.AngkaMinbet,
		AngkaMaxbet4d:                    v.AngkaMaxbet4d,
		AngkaMaxbet3d:                    v.AngkaMaxbet3d,
		AngkaMaxbet3dd:                   v.AngkaMaxbet3dd,
		AngkaMaxbet2d:                    v.AngkaMaxbet2d,
		AngkaMaxbet2dd:                   v.AngkaMaxbet2dd,
		AngkaMaxbet2dt:                   v.AngkaMaxbet2dt,
		AngkaMaxbet4dFull:                v.AngkaMaxbet4dFull,
		AngkaMaxbet3dFull:                v.AngkaMaxbet3dFull,
		AngkaMaxbet3ddFull:               v.AngkaMaxbet3ddFull,
		AngkaMaxbet2dFull:                v.AngkaMaxbet2dFull,
		AngkaMaxbet2ddFull:               v.AngkaMaxbet2ddFull,
		AngkaMaxbet2dtFull:               v.AngkaMaxbet2dtFull,
		AngkaMaxbet4dBb:                  v.AngkaMaxbet4dBb,
		AngkaMaxbet3dBb:                  v.AngkaMaxbet3dBb,
		AngkaMaxbet3ddBb:                 v.AngkaMaxbet3ddBb,
		AngkaMaxbet2dBb:                  v.AngkaMaxbet2dBb,
		AngkaMaxbet2ddBb:                 v.AngkaMaxbet2ddBb,
		AngkaMaxbet2dtBb:                 v.AngkaMaxbet2dtBb,
		AngkaWin4d:                       v.AngkaWin4d,
		AngkaWin3d:                       v.AngkaWin3d,
		AngkaWin3dd:                      v.AngkaWin3dd,
		AngkaWin2d:                       v.AngkaWin2d,
		AngkaWin2dd:                      v.AngkaWin2dd,
		AngkaWin2dt:                      v.AngkaWin2dt,
		AngkaDisc4d:                      v.AngkaDisc4d,
		AngkaDisc3d:                      v.AngkaDisc3d,
		AngkaDisc3dd:                     v.AngkaDisc3dd,
		AngkaDisc2d:                      v.AngkaDisc2d,
		AngkaDisc2dd:                     v.AngkaDisc2dd,
		AngkaDisc2dt:                     v.AngkaDisc2dt,
		AngkaLimitbuang4d:                v.AngkaLimitbuang4d,
		AngkaLimitbuang3d:                v.AngkaLimitbuang3d,
		AngkaLimitbuang3dd:               v.AngkaLimitbuang3dd,
		AngkaLimitbuang2d:                v.AngkaLimitbuang2d,
		AngkaLimitbuang2dd:               v.AngkaLimitbuang2dd,
		AngkaLimitbuang2dt:               v.AngkaLimitbuang2dt,
		AngkaLimittotal4d:                v.AngkaLimittotal4d,
		AngkaLimittotal3d:                v.AngkaLimittotal3d,
		AngkaLimittotal3dd:               v.AngkaLimittotal3dd,
		AngkaLimittotal2d:                v.AngkaLimittotal2d,
		AngkaLimittotal2dd:               v.AngkaLimittotal2dd,
		AngkaLimittotal2dt:               v.AngkaLimittotal2dt,
		AngkaWin4dnodisc:                 v.AngkaWin4dnodisc,
		AngkaWin3dnodisc:                 v.AngkaWin3dnodisc,
		AngkaWin3ddnodisc:                v.AngkaWin3ddnodisc,
		AngkaWin2dnodisc:                 v.AngkaWin2dnodisc,
		AngkaWin2ddnodisc:                v.AngkaWin2ddnodisc,
		AngkaWin2dtnodisc:                v.AngkaWin2dtnodisc,
		AngkaWin4dbbKena:                 v.AngkaWin4dbbKena,
		AngkaWin3dbbKena:                 v.AngkaWin3dbbKena,
		AngkaWin3ddbbKena:                v.AngkaWin3ddbbKena,
		AngkaWin2dbbKena:                 v.AngkaWin2dbbKena,
		AngkaWin2ddbbKena:                v.AngkaWin2ddbbKena,
		AngkaWin2dtbbKena:                v.AngkaWin2dtbbKena,
		AngkaWin4dbb:                     v.AngkaWin4dbb,
		AngkaWin3dbb:                     v.AngkaWin3dbb,
		AngkaWin3ddbb:                    v.AngkaWin3ddbb,
		AngkaWin2dbb:                     v.AngkaWin2dbb,
		AngkaWin2ddbb:                    v.AngkaWin2ddbb,
		AngkaWin2dtbb:                    v.AngkaWin2dtbb,
		AngkaMaxbuy4d:                    v.AngkaMaxbuy4d,
		AngkaMaxbuy3d:                    v.AngkaMaxbuy3d,
		AngkaMaxbuy3dd:                   v.AngkaMaxbuy3dd,
		AngkaMaxbuy2d:                    v.AngkaMaxbuy2d,
		AngkaMaxbuy2dd:                   v.AngkaMaxbuy2dd,
		AngkaMaxbuy2dt:                   v.AngkaMaxbuy2dt,
		AngkaLimitline4d:                 v.AngkaLimitline4d,
		AngkaLimitline3d:                 v.AngkaLimitline3d,
		AngkaLimitline2d:                 v.AngkaLimitline2d,
		AngkaLimitline2dd:                v.AngkaLimitline2dd,
		AngkaLimitline2dt:                v.AngkaLimitline2dt,
		AngkaLimitline3dd:                v.AngkaLimitline3dd,
		AngkaBbfs:                        v.AngkaBbfs,
		CbMinbet:                         v.CbMinbet,
		CbMaxbet:                         v.CbMaxbet,
		CbMaxbuy:                         v.CbMaxbuy,
		CbWin:                            v.CbWin,
		CbDisc:                           v.CbDisc,
		CbLimitbuang:                     v.CbLimitbuang,
		CbLimitotal:                      v.CbLimitotal,
		CmacauMinbet:                     v.CmacauMinbet,
		CmacauMaxbet:                     v.CmacauMaxbet,
		CmacauMaxbuy:                     v.CmacauMaxbuy,
		CmacauWin2digit:                  v.CmacauWin2digit,
		CmacauWin3digit:                  v.CmacauWin3digit,
		CmacauWin4digit:                  v.CmacauWin4digit,
		CmacauDisc:                       v.CmacauDisc,
		CmacauLimitbuang:                 v.CmacauLimitbuang,
		CmacauLimittotal:                 v.CmacauLimittotal,
		CnagaMinbet:                      v.CnagaMinbet,
		CnagaMaxbet:                      v.CnagaMaxbet,
		CnagaMaxbuy:                      v.CnagaMaxbuy,
		CnagaWin3digit:                   v.CnagaWin3digit,
		CnagaWin4digit:                   v.CnagaWin4digit,
		CnagaDisc:                        v.CnagaDisc,
		CnagaLimitbuang:                  v.CnagaLimitbuang,
		CnagaLimittotal:                  v.CnagaLimittotal,
		CjituMinbet:                      v.CjituMinbet,
		CjituMaxbet:                      v.CjituMaxbet,
		CjituMaxbuy:                      v.CjituMaxbuy,
		CjituWinas:                       v.CjituWinas,
		CjituWinkop:                      v.CjituWinkop,
		CjituWinkepala:                   v.CjituWinkepala,
		CjituWinekor:                     v.CjituWinekor,
		CjituDesic:                       v.CjituDesic,
		CjituLimitbuang:                  v.CjituLimitbuang,
		CjituLimitotal:                   v.CjituLimitotal,
		Umum5050Minbet:                   v.Umum5050Minbet,
		Umum5050Maxbet:                   v.Umum5050Maxbet,
		Umum5050Maxbuy:                   v.Umum5050Maxbuy,
		Umum5050Keibesar:                 v.Umum5050Keibesar,
		Umum5050Keikecil:                 v.Umum5050Keikecil,
		Umum5050Keigenap:                 v.Umum5050Keigenap,
		Umum5050Keiganjil:                v.Umum5050Keiganjil,
		Umum5050Keitengah:                v.Umum5050Keitengah,
		Umum5050Keitepi:                  v.Umum5050Keitepi,
		Umum5050Discbesar:                v.Umum5050Discbesar,
		Umum5050Disckecil:                v.Umum5050Disckecil,
		Umum5050Discgenap:                v.Umum5050Discgenap,
		Umum5050Discganjil:               v.Umum5050Discganjil,
		Umum5050Disctengah:               v.Umum5050Disctengah,
		Umum5050Disctepi:                 v.Umum5050Disctepi,
		Umum5050Limitbuang:               v.Umum5050Limitbuang,
		Umum5050Limittotal:               v.Umum5050Limittotal,
		Special5050Minbet:                v.Special5050Minbet,
		Special5050Maxbet:                v.Special5050Maxbet,
		Special5050Maxbuy:                v.Special5050Maxbuy,
		Special5050Keiasganjil:           v.Special5050Keiasganjil,
		Special5050Keiasgenap:            v.Special5050Keiasgenap,
		Special5050Keiasbesar:            v.Special5050Keiasbesar,
		Special5050Keiaskecil:            v.Special5050Keiaskecil,
		Special5050Keikopganjil:          v.Special5050Keikopganjil,
		Special5050Keikopgenap:           v.Special5050Keikopgenap,
		Special5050Keikopbesar:           v.Special5050Keikopbesar,
		Special5050Keikopkecil:           v.Special5050Keikopkecil,
		Special5050Keikepalaganjil:       v.Special5050Keikepalaganjil,
		Special5050Keikepalagenap:        v.Special5050Keikepalagenap,
		Special5050Keikepalabesar:        v.Special5050Keikepalabesar,
		Special5050Keikepalakecil:        v.Special5050Keikepalakecil,
		Special5050Keiekorganjil:         v.Special5050Keiekorganjil,
		Special5050Keiekorgenap:          v.Special5050Keiekorgenap,
		Special5050Keiekorbesar:          v.Special5050Keiekorbesar,
		Special5050Keiekorkecil:          v.Special5050Keiekorkecil,
		Special5050Discasganjil:          v.Special5050Discasganjil,
		Special5050Discasgenap:           v.Special5050Discasgenap,
		Special5050Discasbesar:           v.Special5050Discasbesar,
		Special5050Discaskecil:           v.Special5050Discaskecil,
		Special5050Disckopganjil:         v.Special5050Disckopganjil,
		Special5050Disckopgenap:          v.Special5050Disckopgenap,
		Special5050Disckopbesar:          v.Special5050Disckopbesar,
		Special5050Disckopkecil:          v.Special5050Disckopkecil,
		Special5050Disckepalaganjil:      v.Special5050Disckepalaganjil,
		Special5050Disckepalagenap:       v.Special5050Disckepalagenap,
		Special5050Disckepalabesar:       v.Special5050Disckepalabesar,
		Special5050Disckepalakecil:       v.Special5050Disckepalakecil,
		Special5050Discekorganjil:        v.Special5050Discekorganjil,
		Special5050Discekorgenap:         v.Special5050Discekorgenap,
		Special5050Discekorbesar:         v.Special5050Discekorbesar,
		Special5050Discekorkecil:         v.Special5050Discekorkecil,
		Special5050Limitbuang:            v.Special5050Limitbuang,
		Special5050Limittotal:            v.Special5050Limittotal,
		Kombinasi5050Minbet:              v.Kombinasi5050Minbet,
		Kombinasi5050Maxbet:              v.Kombinasi5050Maxbet,
		Kombinasi5050Maxbuy:              v.Kombinasi5050Maxbuy,
		Kombinasi5050Belakangkeimono:     v.Kombinasi5050Belakangkeimono,
		Kombinasi5050Belakangkeistereo:   v.Kombinasi5050Belakangkeistereo,
		Kombinasi5050Belakangkeikembang:  v.Kombinasi5050Belakangkeikembang,
		Kombinasi5050Belakangkeikempis:   v.Kombinasi5050Belakangkeikempis,
		Kombinasi5050Belakangkeikembar:   v.Kombinasi5050Belakangkeikembar,
		Kombinasi5050Tengahkeimono:       v.Kombinasi5050Tengahkeimono,
		Kombinasi5050Tengahkeistereo:     v.Kombinasi5050Tengahkeistereo,
		Kombinasi5050Tengahkeikembang:    v.Kombinasi5050Tengahkeikembang,
		Kombinasi5050Tengahkeikempis:     v.Kombinasi5050Tengahkeikempis,
		Kombinasi5050Tengahkeikembar:     v.Kombinasi5050Tengahkeikembar,
		Kombinasi5050Depankeimono:        v.Kombinasi5050Depankeimono,
		Kombinasi5050Depankeistereo:      v.Kombinasi5050Depankeistereo,
		Kombinasi5050Depankeikembang:     v.Kombinasi5050Depankeikembang,
		Kombinasi5050Depankeikempis:      v.Kombinasi5050Depankeikempis,
		Kombinasi5050Depankeikembar:      v.Kombinasi5050Depankeikembar,
		Kombinasi5050Belakangdiscmono:    v.Kombinasi5050Belakangdiscmono,
		Kombinasi5050Belakangdiscstereo:  v.Kombinasi5050Belakangdiscstereo,
		Kombinasi5050Belakangdisckembang: v.Kombinasi5050Belakangdisckembang,
		Kombinasi5050Belakangdisckempis:  v.Kombinasi5050Belakangdisckempis,
		Kombinasi5050Belakangdisckembar:  v.Kombinasi5050Belakangdisckembar,
		Kombinasi5050Tengahdiscmono:      v.Kombinasi5050Tengahdiscmono,
		Kombinasi5050Tengahdiscstereo:    v.Kombinasi5050Tengahdiscstereo,
		Kombinasi5050Tengahdisckembang:   v.Kombinasi5050Tengahdisckembang,
		Kombinasi5050Tengahdisckempis:    v.Kombinasi5050Tengahdisckempis,
		Kombinasi5050Tengahdisckembar:    v.Kombinasi5050Tengahdisckembar,
		Kombinasi5050Depandiscmono:       v.Kombinasi5050Depandiscmono,
		Kombinasi5050Depandiscstereo:     v.Kombinasi5050Depandiscstereo,
		Kombinasi5050Depandisckembang:    v.Kombinasi5050Depandisckembang,
		Kombinasi5050Depandisckempis:     v.Kombinasi5050Depandisckempis,
		Kombinasi5050Depandisckembar:     v.Kombinasi5050Depandisckembar,
		Kombinasi5050Limitbuang:          v.Kombinasi5050Limitbuang,
		Kombinasi5050Limittotal:          v.Kombinasi5050Limittotal,
		MacaukombinasiMinbet:             v.MacaukombinasiMinbet,
		MacaukombinasiMaxbet:             v.MacaukombinasiMaxbet,
		MacaukombinasiMaxbuy:             v.MacaukombinasiMaxbuy,
		MacaukombinasiWin:                v.MacaukombinasiWin,
		MacaukombinasiDiscount:           v.MacaukombinasiDiscount,
		MacaukombinasiLimitbuang:         v.MacaukombinasiLimitbuang,
		MacaukombinasiLimittotal:         v.MacaukombinasiLimittotal,
		DasarMinbet:                      v.DasarMinbet,
		DasarMaxbet:                      v.DasarMaxbet,
		DasarMaxbuy:                      v.DasarMaxbuy,
		DasarKeibesar:                    v.DasarKeibesar,
		DasarKeikecil:                    v.DasarKeikecil,
		DasarKeigenap:                    v.DasarKeigenap,
		DasarKeiganjil:                   v.DasarKeiganjil,
		DasarDiscbesar:                   v.DasarDiscbesar,
		DasarDisckecil:                   v.DasarDisckecil,
		DasarDiscigenap:                  v.DasarDiscigenap,
		DasarDiscganjil:                  v.DasarDiscganjil,
		DasarLimitbuang:                  v.DasarLimitbuang,
		DasarLimittotal:                  v.DasarLimittotal,
		ShioMinbet:                       v.ShioMinbet,
		ShioMaxbet:                       v.ShioMaxbet,
		ShioMaxbuy:                       v.ShioMaxbuy,
		ShioWin:                          v.ShioWin,
		ShioDisc:                         v.ShioDisc,
		ShioLimitbuang:                   v.ShioLimitbuang,
		ShioLimittotal:                   v.ShioLimittotal,
		Status:                           status,
		JadwalOpen:                       jadwalOpen,
		JadwalTutup:                      jadwalTutup,
		IDtrxkeluaran:                    trxkeluaran.ID,
		Keluaranperiode:                  trxkeluaran.Keluaranperiode,
		Datekeluaran:                     trxkeluaran.Datekeluaran.Format("2006-01-02"),
	}
	go connection.SetRedis(redisKey, record, 24*time.Hour)
	connection.Log.Info("Returning data Database - Pasaran")
	return record, nil
}
