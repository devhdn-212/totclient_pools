package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/dto"
	"github.com/devhdn-212/totagen_api/internal/connection"
	"github.com/devhdn-212/totagen_api/internal/repository"
	"github.com/devhdn-212/totagen_api/internal/util"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

const (
	RedisPasaran         = "agen:pasaran:all"
	RedisPasaranconftoto = "agen:pasaran:conftoto"
)

type pasaranService struct {
	db   *pgxpool.Pool
	repo domain.PasaranRepository
}

func NewPasaranService(db *pgxpool.Pool, repo domain.PasaranRepository) domain.PasaranService {
	return &pasaranService{
		db:   db,
		repo: repo,
	}
}

func (u *pasaranService) All(ctx context.Context, idcomp string) ([]dto.PasaranData, error) {
	cached, found, err := connection.GetRedis(RedisPasaran + ":" + strings.ToLower(idcomp))
	if err != nil {
		return nil, err
	}
	var record []dto.PasaranData
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			connection.Log.Info("Returning data from Redis - Pasaran")
			return record, nil
		}
	}

	pasaran, err := u.repo.FindAll(ctx, idcomp)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range pasaran {
		var createdAt, updateAt string

		if v.CreateAt.Valid {
			if v.CreateBy != "" {
				createdAt = v.CreateBy + ", " + v.CreateAt.Time.In(util.LocJakarta).Format("2006-01-02 15:04:05")
			} else {
				createdAt = ""
			}
		}
		if v.UpdateAt.Valid {
			if v.UpdateBy != "" {
				updateAt = v.UpdateBy + ", " + v.UpdateAt.Time.In(util.LocJakarta).Format("2006-01-02 15:04:05")
			} else {
				updateAt = ""
			}
		}

		jadwal, err := u.repo.FindJadwal(ctx, v.IDcomppasaran)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		var record_jadwal []dto.Companyjadwalpasaran
		for _, d := range jadwal {
			record_jadwal = append(record_jadwal, dto.Companyjadwalpasaran{
				IDjadwalpasaran: d.IDjadwalcomppasaran,
				Haripasaran:     d.Haripasaran,
				Jamtutup:        util.PgTimeToString(d.Jamtutup),
				Jamjadwal:       util.PgTimeToString(d.Jamjadwal),
				Jamopen:         util.PgTimeToString(d.Jamopen),
			})
		}

		record = append(record, dto.PasaranData{
			IDcomppasaran:                    v.IDcomppasaran,
			IDcompany:                        v.IDcompany,
			Aliascomppasaran:                 v.Aliascomppasaran,
			URLpasaran:                       v.URLpasaran,
			URLlogo:                          v.URLlogo,
			Pasarandiundi:                    v.Pasarandiundi,
			Pasaranlibur:                     v.Pasaranlibur,
			Display:                          v.Display,
			Status:                           v.Status,
			AngkaMinbasket:                   v.AngkaMinbasket,
			AngkaMinbet:                      v.AngkaMinbet,
			AngkaMaxbet4d:                    v.AngkaMaxbet4d,
			AngkaMaxbet3d:                    v.AngkaMaxbet3d,
			AngkaMaxbet3dd:                   v.AngkaMaxbet3dd,
			AngkaMaxbet2d:                    v.AngkaMaxbet2d,
			AngkaMaxbet2dd:                   v.AngkaMaxbet2dd,
			AngkaMaxbet2dt:                   v.AngkaMaxbet2dt,
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
			AngkaMaxbet4dBbdisc:              v.AngkaMaxbet4dBbdisc,
			AngkaMaxbet3dBbdisc:              v.AngkaMaxbet3dBbdisc,
			AngkaMaxbet3ddBbdisc:             v.AngkaMaxbet3ddBbdisc,
			AngkaMaxbet2dBbdisc:              v.AngkaMaxbet2dBbdisc,
			AngkaMaxbet2ddBbdisc:             v.AngkaMaxbet2ddBbdisc,
			AngkaMaxbet2dtBbdisc:             v.AngkaMaxbet2dtBbdisc,
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
			AngkaMaxbet4dFullbb:              v.AngkaMaxbet4dFullbb,
			AngkaMaxbet3dFullbb:              v.AngkaMaxbet3dFullbb,
			AngkaMaxbet3ddFullbb:             v.AngkaMaxbet3ddFullbb,
			AngkaMaxbet2dFullbb:              v.AngkaMaxbet2dFullbb,
			AngkaMaxbet2ddFullbb:             v.AngkaMaxbet2ddFullbb,
			AngkaMaxbet2dtFullbb:             v.AngkaMaxbet2dtFullbb,
			AngkaLimitbuang4dFullbb:          v.AngkaLimitbuang4dFullbb,
			AngkaLimitbuang3dFullbb:          v.AngkaLimitbuang3dFullbb,
			AngkaLimitbuang3ddFullbb:         v.AngkaLimitbuang3ddFullbb,
			AngkaLimitbuang2dFullbb:          v.AngkaLimitbuang2dFullbb,
			AngkaLimitbuang2ddFullbb:         v.AngkaLimitbuang2ddFullbb,
			AngkaLimitbuang2dtFullbb:         v.AngkaLimitbuang2dtFullbb,
			AngkaLimittotal4dFullbb:          v.AngkaLimittotal4dFullbb,
			AngkaLimittotal3dFullbb:          v.AngkaLimittotal3dFullbb,
			AngkaLimittotal3ddFullbb:         v.AngkaLimittotal3ddFullbb,
			AngkaLimittotal2dFullbb:          v.AngkaLimittotal2dFullbb,
			AngkaLimittotal2ddFullbb:         v.AngkaLimittotal2ddFullbb,
			AngkaLimittotal2dtFullbb:         v.AngkaLimittotal2dtFullbb,
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
			ShioReferal:                      v.ShioReferal,
			ShioShiotahunini:                 v.ShioShiotahunini,
			ShioMinbet:                       v.ShioMinbet,
			ShioMaxbet:                       v.ShioMaxbet,
			ShioMaxbuy:                       v.ShioMaxbuy,
			ShioWin:                          v.ShioWin,
			ShioDisc:                         v.ShioDisc,
			ShioLimitbuang:                   v.ShioLimitbuang,
			ShioLimittotal:                   v.ShioLimittotal,
			Jadwalpasaran:                    record_jadwal,
			Created:                          createdAt,
			Update:                           updateAt,
		})
	}
	go connection.SetRedis(RedisPasaran+":"+strings.ToLower(idcomp), record, 24*time.Hour)
	connection.Log.Info("Returning data Database - Pasaran")
	return record, nil
}

func (u *pasaranService) Save(ctx context.Context, req dto.PasaranSave, client, idcomp string) error {
	// Start Transaction native pgx v5
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer rollback jika terjadi panic atau error sebelum commit
	defer tx.Rollback(ctx)

	// Executor transaksi native pgx
	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewCompanypasaranRepository(txExec)
	txRepoconftoto := repository.NewCompanyconftotoRepository(txExec)

	flag, err := txRepo.FindByID(ctx, req.IDcomppasaran, req.IDcompany)
	if err != nil {
		return err
	}

	now := util.GetNowJakarta()

	if req.Type == "Edit" {
		if flag.IDcomppasaran == "" {
			return fmt.Errorf("Pasaran %w", util.ErrNotFound)
		}

		cached_conftoto, found_conftoto, err_conftoto := connection.GetRedis(RedisPasaranconftoto + ":" + strings.ToLower(idcomp))
		if err_conftoto != nil {
			return err_conftoto
		}

		var record domain.Companyconftoto
		if found_conftoto {
			if err := json.Unmarshal([]byte(cached_conftoto), &record); err != nil {
				found_conftoto = false
			} else {
				connection.Log.Info("Returning data from Redis - Comp Conf Toto")
			}
		}

		if !found_conftoto {
			flagconfoto, errconfoto := txRepoconftoto.FindByID(ctx, idcomp)
			if errconfoto != nil {
				return errconfoto
			}
			go connection.SetRedis(RedisPasaranconftoto+":"+strings.ToLower(idcomp), flagconfoto, 24*time.Hour)
			record = flagconfoto
		}
		if err := _validatePasaranLimits(req, record); err != nil {
			return err
		}

		flag.IDcomppasaran = req.IDcomppasaran
		flag.IDcompany = idcomp
		flag.Aliascomppasaran = req.Aliascomppasaran
		flag.Display = req.Display
		flag.URLpasaran = req.URLpasaran
		flag.URLlogo = req.URLlogo
		flag.Pasarandiundi = req.Pasarandiundi
		flag.Pasaranlibur = req.Pasaranlibur
		flag.AngkaMinbasket = req.AngkaMinbasket
		flag.AngkaMinbet = req.AngkaMinbet
		flag.AngkaMaxbet4d = req.AngkaMaxbet4d
		flag.AngkaMaxbet3d = req.AngkaMaxbet3d
		flag.AngkaMaxbet3dd = req.AngkaMaxbet3dd
		flag.AngkaMaxbet2d = req.AngkaMaxbet2d
		flag.AngkaMaxbet2dd = req.AngkaMaxbet2dd
		flag.AngkaMaxbet2dt = req.AngkaMaxbet2dt
		flag.AngkaWin4d = req.AngkaWin4d
		flag.AngkaWin3d = req.AngkaWin3d
		flag.AngkaWin3dd = req.AngkaWin3dd
		flag.AngkaWin2d = req.AngkaWin2d
		flag.AngkaWin2dd = req.AngkaWin2dd
		flag.AngkaWin2dt = req.AngkaWin2dt
		flag.AngkaDisc4d = req.AngkaDisc4d
		flag.AngkaDisc3d = req.AngkaDisc3d
		flag.AngkaDisc3dd = req.AngkaDisc3dd
		flag.AngkaDisc2d = req.AngkaDisc2d
		flag.AngkaDisc2dd = req.AngkaDisc2dd
		flag.AngkaDisc2dt = req.AngkaDisc2dt
		flag.AngkaLimitbuang4d = req.AngkaLimitbuang4d
		flag.AngkaLimitbuang3d = req.AngkaLimitbuang3d
		flag.AngkaLimitbuang3dd = req.AngkaLimitbuang3dd
		flag.AngkaLimitbuang2d = req.AngkaLimitbuang2d
		flag.AngkaLimitbuang2dd = req.AngkaLimitbuang2dd
		flag.AngkaLimitbuang2dt = req.AngkaLimitbuang2dt
		flag.AngkaLimittotal4d = req.AngkaLimittotal4d
		flag.AngkaLimittotal3d = req.AngkaLimittotal3d
		flag.AngkaLimittotal3dd = req.AngkaLimittotal3dd
		flag.AngkaLimittotal2d = req.AngkaLimittotal2d
		flag.AngkaLimittotal2dd = req.AngkaLimittotal2dd
		flag.AngkaLimittotal2dt = req.AngkaLimittotal2dt
		flag.AngkaMaxbet4dFull = req.AngkaMaxbet4dFull
		flag.AngkaMaxbet3dFull = req.AngkaMaxbet3dFull
		flag.AngkaMaxbet3ddFull = req.AngkaMaxbet3ddFull
		flag.AngkaMaxbet2dFull = req.AngkaMaxbet2dFull
		flag.AngkaMaxbet2ddFull = req.AngkaMaxbet2ddFull
		flag.AngkaMaxbet2dtFull = req.AngkaMaxbet2dtFull
		flag.AngkaMaxbet4dBb = req.AngkaMaxbet4dBb
		flag.AngkaMaxbet3dBb = req.AngkaMaxbet3dBb
		flag.AngkaMaxbet3ddBb = req.AngkaMaxbet3ddBb
		flag.AngkaMaxbet2dBb = req.AngkaMaxbet2dBb
		flag.AngkaMaxbet2ddBb = req.AngkaMaxbet2ddBb
		flag.AngkaMaxbet2dtBb = req.AngkaMaxbet2dtBb
		flag.AngkaMaxbet4dBbdisc = req.AngkaMaxbet4dBbdisc
		flag.AngkaMaxbet3dBbdisc = req.AngkaMaxbet3dBbdisc
		flag.AngkaMaxbet3ddBbdisc = req.AngkaMaxbet3ddBbdisc
		flag.AngkaMaxbet2dBbdisc = req.AngkaMaxbet2dBbdisc
		flag.AngkaMaxbet2ddBbdisc = req.AngkaMaxbet2ddBbdisc
		flag.AngkaMaxbet2dtBbdisc = req.AngkaMaxbet2dtBbdisc
		flag.AngkaWin4dnodisc = req.AngkaWin4dnodisc
		flag.AngkaWin3dnodisc = req.AngkaWin3dnodisc
		flag.AngkaWin3ddnodisc = req.AngkaWin3ddnodisc
		flag.AngkaWin2dnodisc = req.AngkaWin2dnodisc
		flag.AngkaWin2ddnodisc = req.AngkaWin2ddnodisc
		flag.AngkaWin2dtnodisc = req.AngkaWin2dtnodisc
		flag.AngkaWin4dbbKena = req.AngkaWin4dbbKena
		flag.AngkaWin3dbbKena = req.AngkaWin3dbbKena
		flag.AngkaWin3ddbbKena = req.AngkaWin3ddbbKena
		flag.AngkaWin2dbbKena = req.AngkaWin2dbbKena
		flag.AngkaWin2ddbbKena = req.AngkaWin2ddbbKena
		flag.AngkaWin2dtbbKena = req.AngkaWin2dtbbKena
		flag.AngkaWin4dbb = req.AngkaWin4dbb
		flag.AngkaWin3dbb = req.AngkaWin3dbb
		flag.AngkaWin3ddbb = req.AngkaWin3ddbb
		flag.AngkaWin2dbb = req.AngkaWin2dbb
		flag.AngkaWin2ddbb = req.AngkaWin2ddbb
		flag.AngkaWin2dtbb = req.AngkaWin2dtbb
		flag.AngkaMaxbuy4d = req.AngkaMaxbuy4d
		flag.AngkaMaxbuy3d = req.AngkaMaxbuy3d
		flag.AngkaMaxbuy3dd = req.AngkaMaxbuy3dd
		flag.AngkaMaxbuy2d = req.AngkaMaxbuy2d
		flag.AngkaMaxbuy2dd = req.AngkaMaxbuy2dd
		flag.AngkaMaxbuy2dt = req.AngkaMaxbuy2dt
		flag.AngkaMaxbet4dFullbb = req.AngkaMaxbet4dFullbb
		flag.AngkaMaxbet3dFullbb = req.AngkaMaxbet3dFullbb
		flag.AngkaMaxbet3ddFullbb = req.AngkaMaxbet3ddFullbb
		flag.AngkaMaxbet2dFullbb = req.AngkaMaxbet2dFullbb
		flag.AngkaMaxbet2ddFullbb = req.AngkaMaxbet2ddFullbb
		flag.AngkaMaxbet2dtFullbb = req.AngkaMaxbet2dtFullbb
		flag.AngkaLimitbuang4dFullbb = req.AngkaLimitbuang4dFullbb
		flag.AngkaLimitbuang3dFullbb = req.AngkaLimitbuang3dFullbb
		flag.AngkaLimitbuang3ddFullbb = req.AngkaLimitbuang3ddFullbb
		flag.AngkaLimitbuang2dFullbb = req.AngkaLimitbuang2dFullbb
		flag.AngkaLimitbuang2ddFullbb = req.AngkaLimitbuang2ddFullbb
		flag.AngkaLimitbuang2dtFullbb = req.AngkaLimitbuang2dtFullbb
		flag.AngkaLimittotal4dFullbb = req.AngkaLimittotal4dFullbb
		flag.AngkaLimittotal3dFullbb = req.AngkaLimittotal3dFullbb
		flag.AngkaLimittotal3ddFullbb = req.AngkaLimittotal3ddFullbb
		flag.AngkaLimittotal2dFullbb = req.AngkaLimittotal2dFullbb
		flag.AngkaLimittotal2ddFullbb = req.AngkaLimittotal2ddFullbb
		flag.AngkaLimittotal2dtFullbb = req.AngkaLimittotal2dtFullbb
		flag.AngkaLimitline4d = req.AngkaLimitline4d
		flag.AngkaLimitline3d = req.AngkaLimitline3d
		flag.AngkaLimitline2d = req.AngkaLimitline2d
		flag.AngkaLimitline2dd = req.AngkaLimitline2dd
		flag.AngkaLimitline2dt = req.AngkaLimitline2dt
		flag.AngkaLimitline3dd = req.AngkaLimitline3dd
		flag.AngkaBbfs = req.AngkaBbfs
		flag.CbMinbet = req.CbMinbet
		flag.CbMaxbet = req.CbMaxbet
		flag.CbMaxbuy = req.CbMaxbuy
		flag.CbWin = req.CbWin
		flag.CbDisc = req.CbDisc
		flag.CbLimitbuang = req.CbLimitbuang
		flag.CbLimitotal = req.CbLimitotal
		flag.CmacauMinbet = req.CmacauMinbet
		flag.CmacauMaxbet = req.CmacauMaxbet
		flag.CmacauMaxbuy = req.CmacauMaxbuy
		flag.CmacauWin2digit = req.CmacauWin2digit
		flag.CmacauWin3digit = req.CmacauWin3digit
		flag.CmacauWin4digit = req.CmacauWin4digit
		flag.CmacauDisc = req.CmacauDisc
		flag.CmacauLimitbuang = req.CmacauLimitbuang
		flag.CmacauLimittotal = req.CmacauLimittotal
		flag.CnagaMinbet = req.CnagaMinbet
		flag.CnagaMaxbet = req.CnagaMaxbet
		flag.CnagaMaxbuy = req.CnagaMaxbuy
		flag.CnagaWin3digit = req.CnagaWin3digit
		flag.CnagaWin4digit = req.CnagaWin4digit
		flag.CnagaDisc = req.CnagaDisc
		flag.CnagaLimitbuang = req.CnagaLimitbuang
		flag.CnagaLimittotal = req.CnagaLimittotal
		flag.CjituMinbet = req.CjituMinbet
		flag.CjituMaxbet = req.CjituMaxbet
		flag.CjituMaxbuy = req.CjituMaxbuy
		flag.CjituWinas = req.CjituWinas
		flag.CjituWinkop = req.CjituWinkop
		flag.CjituWinkepala = req.CjituWinkepala
		flag.CjituWinekor = req.CjituWinekor
		flag.CjituDesic = req.CjituDesic
		flag.CjituLimitbuang = req.CjituLimitbuang
		flag.CjituLimitotal = req.CjituLimitotal
		flag.Umum5050Minbet = req.Umum5050Minbet
		flag.Umum5050Maxbet = req.Umum5050Maxbet
		flag.Umum5050Maxbuy = req.Umum5050Maxbuy
		flag.Umum5050Keibesar = req.Umum5050Keibesar
		flag.Umum5050Keikecil = req.Umum5050Keikecil
		flag.Umum5050Keigenap = req.Umum5050Keigenap
		flag.Umum5050Keiganjil = req.Umum5050Keiganjil
		flag.Umum5050Keitengah = req.Umum5050Keitengah
		flag.Umum5050Keitepi = req.Umum5050Keitepi
		flag.Umum5050Discbesar = req.Umum5050Discbesar
		flag.Umum5050Disckecil = req.Umum5050Disckecil
		flag.Umum5050Discgenap = req.Umum5050Discgenap
		flag.Umum5050Discganjil = req.Umum5050Discganjil
		flag.Umum5050Disctengah = req.Umum5050Disctengah
		flag.Umum5050Disctepi = req.Umum5050Disctepi
		flag.Umum5050Limitbuang = req.Umum5050Limitbuang
		flag.Umum5050Limittotal = req.Umum5050Limittotal
		flag.Special5050Minbet = req.Special5050Minbet
		flag.Special5050Maxbet = req.Special5050Maxbet
		flag.Special5050Maxbuy = req.Special5050Maxbuy
		flag.Special5050Keiasganjil = req.Special5050Keiasganjil
		flag.Special5050Keiasgenap = req.Special5050Keiasgenap
		flag.Special5050Keiasbesar = req.Special5050Keiasbesar
		flag.Special5050Keiaskecil = req.Special5050Keiaskecil
		flag.Special5050Keikopganjil = req.Special5050Keikopganjil
		flag.Special5050Keikopgenap = req.Special5050Keikopgenap
		flag.Special5050Keikopbesar = req.Special5050Keikopbesar
		flag.Special5050Keikopkecil = req.Special5050Keikopkecil
		flag.Special5050Keikepalaganjil = req.Special5050Keikepalaganjil
		flag.Special5050Keikepalagenap = req.Special5050Keikepalagenap
		flag.Special5050Keikepalabesar = req.Special5050Keikepalabesar
		flag.Special5050Keikepalakecil = req.Special5050Keikepalakecil
		flag.Special5050Keiekorganjil = req.Special5050Keiekorganjil
		flag.Special5050Keiekorgenap = req.Special5050Keiekorgenap
		flag.Special5050Keiekorbesar = req.Special5050Keiekorbesar
		flag.Special5050Keiekorkecil = req.Special5050Keiekorkecil
		flag.Special5050Discasganjil = req.Special5050Discasganjil
		flag.Special5050Discasgenap = req.Special5050Discasgenap
		flag.Special5050Discasbesar = req.Special5050Discasbesar
		flag.Special5050Discaskecil = req.Special5050Discaskecil
		flag.Special5050Disckopganjil = req.Special5050Disckopganjil
		flag.Special5050Disckopgenap = req.Special5050Disckopgenap
		flag.Special5050Disckopbesar = req.Special5050Disckopbesar
		flag.Special5050Disckopkecil = req.Special5050Disckopkecil
		flag.Special5050Disckepalaganjil = req.Special5050Disckepalaganjil
		flag.Special5050Disckepalagenap = req.Special5050Disckepalagenap
		flag.Special5050Disckepalabesar = req.Special5050Disckepalabesar
		flag.Special5050Disckepalakecil = req.Special5050Disckepalakecil
		flag.Special5050Discekorganjil = req.Special5050Discekorganjil
		flag.Special5050Discekorgenap = req.Special5050Discekorgenap
		flag.Special5050Discekorbesar = req.Special5050Discekorbesar
		flag.Special5050Discekorkecil = req.Special5050Discekorkecil
		flag.Special5050Limitbuang = req.Special5050Limitbuang
		flag.Special5050Limittotal = req.Special5050Limittotal
		flag.Kombinasi5050Minbet = req.Kombinasi5050Minbet
		flag.Kombinasi5050Maxbet = req.Kombinasi5050Maxbet
		flag.Kombinasi5050Maxbuy = req.Kombinasi5050Maxbuy
		flag.Kombinasi5050Belakangkeimono = req.Kombinasi5050Belakangkeimono
		flag.Kombinasi5050Belakangkeistereo = req.Kombinasi5050Belakangkeistereo
		flag.Kombinasi5050Belakangkeikembang = req.Kombinasi5050Belakangkeikembang
		flag.Kombinasi5050Belakangkeikempis = req.Kombinasi5050Belakangkeikempis
		flag.Kombinasi5050Belakangkeikembar = req.Kombinasi5050Belakangkeikembar
		flag.Kombinasi5050Tengahkeimono = req.Kombinasi5050Tengahkeimono
		flag.Kombinasi5050Tengahkeistereo = req.Kombinasi5050Tengahkeistereo
		flag.Kombinasi5050Tengahkeikembang = req.Kombinasi5050Tengahkeikembang
		flag.Kombinasi5050Tengahkeikempis = req.Kombinasi5050Tengahkeikempis
		flag.Kombinasi5050Tengahkeikembar = req.Kombinasi5050Tengahkeikembar
		flag.Kombinasi5050Depankeimono = req.Kombinasi5050Depankeimono
		flag.Kombinasi5050Depankeistereo = req.Kombinasi5050Depankeistereo
		flag.Kombinasi5050Depankeikembang = req.Kombinasi5050Depankeikembang
		flag.Kombinasi5050Depankeikempis = req.Kombinasi5050Depankeikempis
		flag.Kombinasi5050Depankeikembar = req.Kombinasi5050Depankeikembar
		flag.Kombinasi5050Belakangdiscmono = req.Kombinasi5050Belakangdiscmono
		flag.Kombinasi5050Belakangdiscstereo = req.Kombinasi5050Belakangdiscstereo
		flag.Kombinasi5050Belakangdisckembang = req.Kombinasi5050Belakangdisckembang
		flag.Kombinasi5050Belakangdisckempis = req.Kombinasi5050Belakangdisckempis
		flag.Kombinasi5050Belakangdisckembar = req.Kombinasi5050Belakangdisckembar
		flag.Kombinasi5050Tengahdiscmono = req.Kombinasi5050Tengahdiscmono
		flag.Kombinasi5050Tengahdiscstereo = req.Kombinasi5050Tengahdiscstereo
		flag.Kombinasi5050Tengahdisckembang = req.Kombinasi5050Tengahdisckembang
		flag.Kombinasi5050Tengahdisckempis = req.Kombinasi5050Tengahdisckempis
		flag.Kombinasi5050Tengahdisckembar = req.Kombinasi5050Tengahdisckembar
		flag.Kombinasi5050Depandiscmono = req.Kombinasi5050Depandiscmono
		flag.Kombinasi5050Depandiscstereo = req.Kombinasi5050Depandiscstereo
		flag.Kombinasi5050Depandisckembang = req.Kombinasi5050Depandisckembang
		flag.Kombinasi5050Depandisckempis = req.Kombinasi5050Depandisckempis
		flag.Kombinasi5050Depandisckembar = req.Kombinasi5050Depandisckembar
		flag.Kombinasi5050Limitbuang = req.Kombinasi5050Limitbuang
		flag.Kombinasi5050Limittotal = req.Kombinasi5050Limittotal
		flag.MacaukombinasiMinbet = req.MacaukombinasiMinbet
		flag.MacaukombinasiMaxbet = req.MacaukombinasiMaxbet
		flag.MacaukombinasiMaxbuy = req.MacaukombinasiMaxbuy
		flag.MacaukombinasiWin = req.MacaukombinasiWin
		flag.MacaukombinasiDiscount = req.MacaukombinasiDiscount
		flag.MacaukombinasiLimitbuang = req.MacaukombinasiLimitbuang
		flag.MacaukombinasiLimittotal = req.MacaukombinasiLimittotal
		flag.DasarMinbet = req.DasarMinbet
		flag.DasarMaxbet = req.DasarMaxbet
		flag.DasarMaxbuy = req.DasarMaxbuy
		flag.DasarKeibesar = req.DasarKeibesar
		flag.DasarKeikecil = req.DasarKeikecil
		flag.DasarKeigenap = req.DasarKeigenap
		flag.DasarKeiganjil = req.DasarKeiganjil
		flag.DasarDiscbesar = req.DasarDiscbesar
		flag.DasarDisckecil = req.DasarDisckecil
		flag.DasarDiscigenap = req.DasarDiscigenap
		flag.DasarDiscganjil = req.DasarDiscganjil
		flag.DasarLimitbuang = req.DasarLimitbuang
		flag.DasarLimittotal = req.DasarLimittotal
		flag.ShioReferal = req.ShioReferal
		flag.ShioShiotahunini = req.ShioShiotahunini
		flag.ShioMinbet = req.ShioMinbet
		flag.ShioMaxbet = req.ShioMaxbet
		flag.ShioMaxbuy = req.ShioMaxbuy
		flag.ShioWin = req.ShioWin
		flag.ShioDisc = req.ShioDisc
		flag.ShioLimitbuang = req.ShioLimitbuang
		flag.ShioLimittotal = req.ShioLimittotal
		flag.Status = req.Status
		flag.UpdateBy = client
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		if err = txRepo.Update(ctx, &flag); err != nil {
			fmt.Println("Error update Pasaran: ", err)
			return err
		}

		// jadwal: hapus semua, nanti di-insert ulang di loop bawah
		if err = txRepo.DeleteJadwal(ctx, req.IDcomppasaran); err != nil {
			return err
		}
		//jadwal

		for _, j := range req.Jadwalpasaran {
			jt := util.StringToPgTime(j.Jamtutup)
			jj := util.StringToPgTime(j.Jamjadwal)
			jo := util.StringToPgTime(j.Jamopen)
			if !jt.Valid || !jj.Valid || !jo.Valid {
				return fmt.Errorf("format jam tidak valid untuk hari %s", j.Haripasaran)
			}

			raw := strings.ReplaceAll(uuid.NewString(), "-", "")
			date := time.Now().Format("0601")

			jadwalpasarantoto := domain.Companyjadwalpasaran{
				IDjadwalcomppasaran: fmt.Sprintf("%s%s", date, raw),
				IDcomppasaran:       req.IDcomppasaran,
				Haripasaran:         j.Haripasaran,
				Jamtutup:            util.StringToPgTime(j.Jamtutup),
				Jamjadwal:           util.StringToPgTime(j.Jamjadwal),
				Jamopen:             util.StringToPgTime(j.Jamopen),
				CreateBy:            client,
				CreateAt:            sql.NullTime{Valid: true, Time: now},
			}

			if err = txRepo.Savejadwal(ctx, &jadwalpasarantoto); err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == "23505" {
					return util.ErrDuplicate
				}
				return err
			}
		}
	}

	// Commit transaksi
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisPasaran + ":" + strings.ToLower(idcomp))
	return nil
}
func _validatePasaranLimits(req dto.PasaranSave, record domain.Companyconftoto) error {
	checks := []struct {
		label string
		val   decimal.Decimal
		max   decimal.Decimal
	}{
		{"Total Keranjang", req.AngkaMinbasket, record.AngkaMaxMinbasket},
		{"MinBet", req.AngkaMinbet, record.AngkaMaxMinbet},
		{"MinBet Colok Bebas", req.CbMinbet, record.CbebasMaxMinbet},
		{"MinBet Colok Macau", req.CmacauMinbet, record.CmacauMaxMinbet},
		{"MinBet Colok Naga", req.CmacauMinbet, record.CmacauMaxMinbet},
		{"MinBet Colok Jitu", req.CjituMinbet, record.CjituMaxMinbet},
		{"MinBet 5050 Umum", req.Umum5050Minbet, record.Umum50MaxMinbet},
		{"MinBet 5050 Special", req.Special5050Minbet, record.Special50MaxMinbet},
		{"MinBet 5050 Kombinasi", req.Kombinasi5050Minbet, record.Kombinasi50MaxMinbet},
		{"MinBet Macau / Kombinasi", req.MacaukombinasiMinbet, record.MacauMaxMinbet},
		{"MinBet Dasar", req.DasarMinbet, record.DasarMaxMinbet},
		{"MinBet Shio", req.ShioMinbet, record.ShioMaxMinbet},
		{"MaxBet 4D", req.AngkaMaxbet4d, record.AngkaMaxMaxbet4d},
		{"MaxBet 3D", req.AngkaMaxbet3d, record.AngkaMaxMaxbet3d},
		{"MaxBet 3DD", req.AngkaMaxbet3dd, record.AngkaMaxMaxbet3dd},
		{"MaxBet 2D", req.AngkaMaxbet2d, record.AngkaMaxMaxbet2d},
		{"MaxBet 2DD", req.AngkaMaxbet2dd, record.AngkaMaxMaxbet2dd},
		{"MaxBet 2DT", req.AngkaMaxbet2dt, record.AngkaMaxMaxbet2dt},
		{"MaxBet Colok Bebas", req.CbMaxbet, record.CbebasMaxMaxbet},
		{"MaxBet Colok Macau", req.CmacauMaxbet, record.CmacauMaxMaxbet},
		{"MaxBet Colok Naga", req.CnagaMaxbet, record.CnagaMaxMaxbet},
		{"MaxBet Colok Jitu", req.CjituMaxbet, record.CjituMaxMaxbet},
		{"MaxBet 5050 Umum", req.Umum5050Maxbet, record.Umum50MaxMaxbet},
		{"MaxBet 5050 Special", req.Special5050Maxbet, record.Special50MaxMaxbet},
		{"MaxBet 5050 Kombinasi", req.Kombinasi5050Maxbet, record.Kombinasi50MaxMaxbet},
		{"MaxBet Macau / Kombinasi", req.MacaukombinasiMaxbet, record.MacauMaxMaxbet},
		{"MaxBet Dasar", req.DasarMaxbet, record.DasarMaxMaxbet},
		{"MaxBet Shio", req.ShioMaxbet, record.ShioMaxMaxbet},

		{"Win 4D FULL", req.AngkaWin4dnodisc, record.AngkaMaxWin4dFull},
		{"Win 4D DISC", req.AngkaWin4d, record.AngkaMaxWin4dDisc},
		{"Win 4D BB", req.AngkaWin4dbb, record.AngkaMaxWin4dBb},
		{"Win 4D BB Match", req.AngkaWin4dbbKena, record.AngkaMaxWin4dBbKena},

		{"Win 3D FULL", req.AngkaWin3dnodisc, record.AngkaMaxWin3dFull},
		{"Win 3D DISC", req.AngkaWin3d, record.AngkaMaxWin3dDisc},
		{"Win 3D BB", req.AngkaWin3dbb, record.AngkaMaxWin3dBb},
		{"Win 3D BB Match", req.AngkaWin3dbbKena, record.AngkaMaxWin3dBbKena},

		{"Win 3DD FULL", req.AngkaWin3ddnodisc, record.AngkaMaxWin3ddFull},
		{"Win 3DD DISC", req.AngkaWin3dd, record.AngkaMaxWin3ddDisc},
		{"Win 3DD BB", req.AngkaWin3ddbb, record.AngkaMaxWin3ddBb},
		{"Win 3DD BB Match", req.AngkaWin3ddbbKena, record.AngkaMaxWin3ddBbKena},

		{"Win 2D FULL", req.AngkaWin2dnodisc, record.AngkaMaxWin2dFull},
		{"Win 2D DISC", req.AngkaWin2d, record.AngkaMaxWin2dDisc},
		{"Win 2D BB", req.AngkaWin2dbb, record.AngkaMaxWin2dBb},
		{"Win 2D BB Match", req.AngkaWin2dbbKena, record.AngkaMaxWin2dBbKena},

		{"Win 2DD FULL", req.AngkaWin2ddnodisc, record.AngkaMaxWin2ddFull},
		{"Win 2DD DISC", req.AngkaWin2dd, record.AngkaMaxWin2ddDisc},
		{"Win 2DD BB", req.AngkaWin2ddbb, record.AngkaMaxWin2ddBb},
		{"Win 2DD BB Match", req.AngkaWin2ddbbKena, record.AngkaMaxWin2ddBbKena},

		{"Win 2DT FULL", req.AngkaWin2dtnodisc, record.AngkaMaxWin2dtFull},
		{"Win 2DT DISC", req.AngkaWin2dt, record.AngkaMaxWin2dtDisc},
		{"Win 2DT BB", req.AngkaWin2dtbb, record.AngkaMaxWin2dtBb},
		{"Win 2DT BB Match", req.AngkaWin2dtbbKena, record.AngkaMaxWin2dtBbKena},

		{"Win Colok Bebas", req.CbWin, record.CbebasMaxWin},
		{"Win 2 Digit Colok Macau", req.CmacauWin2digit, record.CmacauMaxWin2},
		{"Win 3 Digit Colok Macau", req.CmacauWin3digit, record.CmacauMaxWin3},
		{"Win 4 Digit Colok Macau", req.CmacauWin3digit, record.CmacauMaxWin4},
		{"Win 3 Digit Colok Naga", req.CnagaWin3digit, record.CnagaMaxWin3},
		{"Win 4 Digit Colok Naga", req.CnagaWin4digit, record.CnagaMaxWin4},
		{"Win As Colok Jitu", req.CjituWinas, record.CjituMaxWinas},
		{"Win Kop Colok Jitu", req.CjituWinkop, record.CjituMaxWinkop},
		{"Win Kepala Colok Jitu", req.CjituWinkepala, record.CjituMaxWinkepala},
		{"Win Ekor Colok Jitu", req.CjituWinekor, record.CjituMaxWinekor},
		{"Win Macau / Kombinasi", req.MacaukombinasiWin, record.MacauMaxWin},
		{"Win Shio", req.ShioWin, record.ShioMaxWin},
	}

	for _, c := range checks {
		if c.val.GreaterThan(c.max) {
			return fmt.Errorf("%s tidak boleh melebihi %s", c.label, c.max.String())
		}
	}
	return nil
}
