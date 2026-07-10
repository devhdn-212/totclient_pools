package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/internal/config"
	"github.com/jackc/pgx/v5"
)

type companypasaranRepository struct {
	db DBExecutor
}

func NewCompanypasaranRepository(db DBExecutor) domain.CompanypasaranRepository {
	return &companypasaranRepository{
		db: db,
	}
}

func (u *companypasaranRepository) FindAll(ctx context.Context, idcomp string) ([]domain.Companypasaran, error) {
	query := `SELECT * FROM ` + config.DB_mst_company_pasaran + ` 
			  WHERE idcompany=$1 
			  ORDER BY displaypasaran ASC`

	rows, err := u.db.Query(ctx, query, idcomp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Companyconftoto
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Companypasaran])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (u *companypasaranRepository) FindJadwal(ctx context.Context, idcomppasaran string) ([]domain.Companyjadwalpasaran, error) {
	query := `SELECT * FROM ` + config.DB_mst_company_jadwaltogel + ` 
			  WHERE idcomppasaran = $1  
			  ORDER BY create_at DESC`

	rows, err := u.db.Query(ctx, query, idcomppasaran)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Companyconftoto
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Companyjadwalpasaran])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *companypasaranRepository) FindByID(ctx context.Context, id, idcomp string) (domain.Companypasaran, error) {
	var record domain.Companypasaran
	query := `SELECT idcomppasaran FROM ` + config.DB_mst_company_pasaran + ` 
			WHERE idcomppasaran = $1 AND idcompany = $2 
			LIMIT 1`

	err := u.db.QueryRow(ctx, query, id, idcomp).Scan(&record.IDcomppasaran)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Companypasaran{}, nil
		}
		return record, err
	}
	return record, nil
}

func (u *companypasaranRepository) Save(ctx context.Context, comppasaran *domain.Companypasaran) error {
	query := `INSERT INTO ` + config.DB_mst_company_pasaran + ` (
		idcomppasaran, idcompany, aliascomppasaran, idpasarantogel,
		urlpasaran, pasarandiundi, pasaranlibur, angka_minbasket, angka_minbet, angka_maxbet4d,
		angka_maxbet3d, angka_maxbet3dd, angka_maxbet2d, angka_maxbet2dd, angka_maxbet2dt,
		angka_win4d, angka_win3d, angka_win3dd, angka_win2d, angka_win2dd, angka_win2dt, angka_disc4d,
		angka_disc3d, angka_disc3dd, angka_disc2d, angka_disc2dd, angka_disc2dt, angka_limitbuang4d,
		angka_limitbuang3d, angka_limitbuang3dd, angka_limitbuang2d, angka_limitbuang2dd,
		angka_limitbuang2dt, angka_limittotal4d, angka_limittotal3d, angka_limittotal3dd,
		angka_limittotal2d, angka_limittotal2dd, angka_limittotal2dt, angka_maxbet4d_full,
		angka_maxbet3d_full, angka_maxbet3dd_full, angka_maxbet2d_full, angka_maxbet2dd_full,
		angka_maxbet2dt_full, angka_maxbet4d_bb, angka_maxbet3d_bb, angka_maxbet3dd_bb,
		angka_maxbet2d_bb, angka_maxbet2dd_bb, angka_maxbet2dt_bb, angka_maxbet4d_bbdisc,
		angka_maxbet3d_bbdisc, angka_maxbet3dd_bbdisc, angka_maxbet2d_bbdisc, angka_maxbet2dd_bbdisc,
		angka_maxbet2dt_bbdisc, angka_win4dnodisc, angka_win3dnodisc, angka_win3ddnodisc,
		angka_win2dnodisc, angka_win2ddnodisc, angka_win2dtnodisc, angka_win4dbb_kena,
		angka_win3dbb_kena, angka_win3ddbb_kena, angka_win2dbb_kena, angka_win2ddbb_kena,
		angka_win2dtbb_kena, angka_win4dbb, angka_win3dbb, angka_win3ddbb, angka_win2dbb,
		angka_win2ddbb, angka_win2dtbb, angka_maxbuy4d, angka_maxbuy3d, angka_maxbuy3dd,
		angka_maxbuy2d, angka_maxbuy2dd, angka_maxbuy2dt, angka_maxbet4d_fullbb,
		angka_maxbet3d_fullbb, angka_maxbet3dd_fullbb, angka_maxbet2d_fullbb, angka_maxbet2dd_fullbb,
		angka_maxbet2dt_fullbb, angka_limitbuang4d_fullbb, angka_limitbuang3d_fullbb,
		angka_limitbuang3dd_fullbb, angka_limitbuang2d_fullbb, angka_limitbuang2dd_fullbb,
		angka_limitbuang2dt_fullbb, angka_limittotal4d_fullbb, angka_limittotal3d_fullbb,
		angka_limittotal3dd_fullbb, angka_limittotal2d_fullbb, angka_limittotal2dd_fullbb,
		angka_limittotal2dt_fullbb, angka_limitline_4d, angka_limitline_3d, angka_limitline_2d,
		angka_limitline_2dd, angka_limitline_2dt, angka_limitline_3dd, angka_bbfs, cb_minbet,
		cb_maxbet, cb_maxbuy, cb_win, cb_disc, cb_limitbuang, cb_limitotal, cmacau_minbet,
		cmacau_maxbet, cmacau_maxbuy, cmacau_win2digit, cmacau_win3digit, cmacau_win4digit,
		cmacau_disc, cmacau_limitbuang, cmacau_limittotal, cnaga_minbet, cnaga_maxbet, cnaga_maxbuy,
		cnaga_win3digit, cnaga_win4digit, cnaga_disc, cnaga_limitbuang, cnaga_limittotal,
		cjitu_minbet, cjitu_maxbet, cjitu_maxbuy, cjitu_winas, cjitu_winkop, cjitu_winkepala,
		cjitu_winekor, cjitu_desic, cjitu_limitbuang, cjitu_limitotal, umum5050_minbet,
		umum5050_maxbet, umum5050_maxbuy, umum5050_keibesar, umum5050_keikecil, umum5050_keigenap,
		umum5050_keiganjil, umum5050_keitengah, umum5050_keitepi, umum5050_discbesar,
		umum5050_disckecil, umum5050_discgenap, umum5050_discganjil, umum5050_disctengah,
		umum5050_disctepi, umum5050_limitbuang, umum5050_limittotal, special5050_minbet,
		special5050_maxbet, special5050_maxbuy, special5050_keiasganjil, special5050_keiasgenap,
		special5050_keiasbesar, special5050_keiaskecil, special5050_keikopganjil,
		special5050_keikopgenap, special5050_keikopbesar, special5050_keikopkecil,
		special5050_keikepalaganjil, special5050_keikepalagenap, special5050_keikepalabesar,
		special5050_keikepalakecil, special5050_keiekorganjil, special5050_keiekorgenap,
		special5050_keiekorbesar, special5050_keiekorkecil, special5050_discasganjil,
		special5050_discasgenap, special5050_discasbesar, special5050_discaskecil,
		special5050_disckopganjil, special5050_disckopgenap, special5050_disckopbesar,
		special5050_disckopkecil, special5050_disckepalaganjil, special5050_disckepalagenap,
		special5050_disckepalabesar, special5050_disckepalakecil, special5050_discekorganjil,
		special5050_discekorgenap, special5050_discekorbesar, special5050_discekorkecil,
		special5050_limitbuang, special5050_limittotal, kombinasi5050_minbet, kombinasi5050_maxbet,
		kombinasi5050_maxbuy, kombinasi5050_belakangkeimono, kombinasi5050_belakangkeistereo,
		kombinasi5050_belakangkeikembang, kombinasi5050_belakangkeikempis,
		kombinasi5050_belakangkeikembar, kombinasi5050_tengahkeimono, kombinasi5050_tengahkeistereo,
		kombinasi5050_tengahkeikembang, kombinasi5050_tengahkeikempis, kombinasi5050_tengahkeikembar,
		kombinasi5050_depankeimono, kombinasi5050_depankeistereo, kombinasi5050_depankeikembang,
		kombinasi5050_depankeikempis, kombinasi5050_depankeikembar, kombinasi5050_belakangdiscmono,
		kombinasi5050_belakangdiscstereo, kombinasi5050_belakangdisckembang,
		kombinasi5050_belakangdisckempis, kombinasi5050_belakangdisckembar,
		kombinasi5050_tengahdiscmono, kombinasi5050_tengahdiscstereo, kombinasi5050_tengahdisckembang,
		kombinasi5050_tengahdisckempis, kombinasi5050_tengahdisckembar, kombinasi5050_depandiscmono,
		kombinasi5050_depandiscstereo, kombinasi5050_depandisckembang, kombinasi5050_depandisckempis,
		kombinasi5050_depandisckembar, kombinasi5050_limitbuang, kombinasi5050_limittotal,
		macaukombinasi_minbet, macaukombinasi_maxbet, macaukombinasi_maxbuy, macaukombinasi_win,
		macaukombinasi_discount, macaukombinasi_limitbuang, macaukombinasi_limittotal, dasar_minbet,
		dasar_maxbet, dasar_maxbuy, dasar_keibesar, dasar_keikecil, dasar_keigenap, dasar_keiganjil,
		dasar_discbesar, dasar_disckecil, dasar_discigenap, dasar_discganjil, dasar_limitbuang,
		dasar_limittotal, shio_referal, shio_shiotahunini, shio_minbet, shio_maxbet, shio_maxbuy,
		shio_win, shio_disc, shio_limitbuang, shio_limittotal, create_by, create_at
	) SELECT
		$1, $2, $3,  p.idpasarantogel, p.urlpasaran, p.pasarandiundi, p.pasaranlibur,
		p.angka_minbasket, p.angka_minbet, p.angka_maxbet4d, p.angka_maxbet3d, p.angka_maxbet3dd,
		p.angka_maxbet2d, p.angka_maxbet2dd, p.angka_maxbet2dt, p.angka_win4d, p.angka_win3d,
		p.angka_win3dd, p.angka_win2d, p.angka_win2dd, p.angka_win2dt, p.angka_disc4d, p.angka_disc3d,
		p.angka_disc3dd, p.angka_disc2d, p.angka_disc2dd, p.angka_disc2dt, p.angka_limitbuang4d,
		p.angka_limitbuang3d, p.angka_limitbuang3dd, p.angka_limitbuang2d, p.angka_limitbuang2dd,
		p.angka_limitbuang2dt, p.angka_limittotal4d, p.angka_limittotal3d, p.angka_limittotal3dd,
		p.angka_limittotal2d, p.angka_limittotal2dd, p.angka_limittotal2dt, p.angka_maxbet4d_full,
		p.angka_maxbet3d_full, p.angka_maxbet3dd_full, p.angka_maxbet2d_full, p.angka_maxbet2dd_full,
		p.angka_maxbet2dt_full, p.angka_maxbet4d_bb, p.angka_maxbet3d_bb, p.angka_maxbet3dd_bb,
		p.angka_maxbet2d_bb, p.angka_maxbet2dd_bb, p.angka_maxbet2dt_bb, p.angka_maxbet4d_bbdisc,
		p.angka_maxbet3d_bbdisc, p.angka_maxbet3dd_bbdisc, p.angka_maxbet2d_bbdisc,
		p.angka_maxbet2dd_bbdisc, p.angka_maxbet2dt_bbdisc, p.angka_win4dnodisc, p.angka_win3dnodisc,
		p.angka_win3ddnodisc, p.angka_win2dnodisc, p.angka_win2ddnodisc, p.angka_win2dtnodisc,
		p.angka_win4dbb_kena, p.angka_win3dbb_kena, p.angka_win3ddbb_kena, p.angka_win2dbb_kena,
		p.angka_win2ddbb_kena, p.angka_win2dtbb_kena, p.angka_win4dbb, p.angka_win3dbb,
		p.angka_win3ddbb, p.angka_win2dbb, p.angka_win2ddbb, p.angka_win2dtbb, p.angka_maxbuy4d,
		p.angka_maxbuy3d, p.angka_maxbuy3dd, p.angka_maxbuy2d, p.angka_maxbuy2dd, p.angka_maxbuy2dt,
		p.angka_maxbet4d_fullbb, p.angka_maxbet3d_fullbb, p.angka_maxbet3dd_fullbb,
		p.angka_maxbet2d_fullbb, p.angka_maxbet2dd_fullbb, p.angka_maxbet2dt_fullbb,
		p.angka_limitbuang4d_fullbb, p.angka_limitbuang3d_fullbb, p.angka_limitbuang3dd_fullbb,
		p.angka_limitbuang2d_fullbb, p.angka_limitbuang2dd_fullbb, p.angka_limitbuang2dt_fullbb,
		p.angka_limittotal4d_fullbb, p.angka_limittotal3d_fullbb, p.angka_limittotal3dd_fullbb,
		p.angka_limittotal2d_fullbb, p.angka_limittotal2dd_fullbb, p.angka_limittotal2dt_fullbb,
		p.angka_limitline_4d, p.angka_limitline_3d, p.angka_limitline_2d, p.angka_limitline_2dd,
		p.angka_limitline_2dt, p.angka_limitline_3dd, p.angka_bbfs, p.cb_minbet, p.cb_maxbet,
		p.cb_maxbuy, p.cb_win, p.cb_disc, p.cb_limitbuang, p.cb_limitotal, p.cmacau_minbet,
		p.cmacau_maxbet, p.cmacau_maxbuy, p.cmacau_win2digit, p.cmacau_win3digit, p.cmacau_win4digit,
		p.cmacau_disc, p.cmacau_limitbuang, p.cmacau_limittotal, p.cnaga_minbet, p.cnaga_maxbet,
		p.cnaga_maxbuy, p.cnaga_win3digit, p.cnaga_win4digit, p.cnaga_disc, p.cnaga_limitbuang,
		p.cnaga_limittotal, p.cjitu_minbet, p.cjitu_maxbet, p.cjitu_maxbuy, p.cjitu_winas,
		p.cjitu_winkop, p.cjitu_winkepala, p.cjitu_winekor, p.cjitu_desic, p.cjitu_limitbuang,
		p.cjitu_limitotal, p.umum5050_minbet, p.umum5050_maxbet, p.umum5050_maxbuy,
		p.umum5050_keibesar, p.umum5050_keikecil, p.umum5050_keigenap, p.umum5050_keiganjil,
		p.umum5050_keitengah, p.umum5050_keitepi, p.umum5050_discbesar, p.umum5050_disckecil,
		p.umum5050_discgenap, p.umum5050_discganjil, p.umum5050_disctengah, p.umum5050_disctepi,
		p.umum5050_limitbuang, p.umum5050_limittotal, p.special5050_minbet, p.special5050_maxbet,
		p.special5050_maxbuy, p.special5050_keiasganjil, p.special5050_keiasgenap,
		p.special5050_keiasbesar, p.special5050_keiaskecil, p.special5050_keikopganjil,
		p.special5050_keikopgenap, p.special5050_keikopbesar, p.special5050_keikopkecil,
		p.special5050_keikepalaganjil, p.special5050_keikepalagenap, p.special5050_keikepalabesar,
		p.special5050_keikepalakecil, p.special5050_keiekorganjil, p.special5050_keiekorgenap,
		p.special5050_keiekorbesar, p.special5050_keiekorkecil, p.special5050_discasganjil,
		p.special5050_discasgenap, p.special5050_discasbesar, p.special5050_discaskecil,
		p.special5050_disckopganjil, p.special5050_disckopgenap, p.special5050_disckopbesar,
		p.special5050_disckopkecil, p.special5050_disckepalaganjil, p.special5050_disckepalagenap,
		p.special5050_disckepalabesar, p.special5050_disckepalakecil, p.special5050_discekorganjil,
		p.special5050_discekorgenap, p.special5050_discekorbesar, p.special5050_discekorkecil,
		p.special5050_limitbuang, p.special5050_limittotal, p.kombinasi5050_minbet,
		p.kombinasi5050_maxbet, p.kombinasi5050_maxbuy, p.kombinasi5050_belakangkeimono,
		p.kombinasi5050_belakangkeistereo, p.kombinasi5050_belakangkeikembang,
		p.kombinasi5050_belakangkeikempis, p.kombinasi5050_belakangkeikembar,
		p.kombinasi5050_tengahkeimono, p.kombinasi5050_tengahkeistereo,
		p.kombinasi5050_tengahkeikembang, p.kombinasi5050_tengahkeikempis,
		p.kombinasi5050_tengahkeikembar, p.kombinasi5050_depankeimono, p.kombinasi5050_depankeistereo,
		p.kombinasi5050_depankeikembang, p.kombinasi5050_depankeikempis,
		p.kombinasi5050_depankeikembar, p.kombinasi5050_belakangdiscmono,
		p.kombinasi5050_belakangdiscstereo, p.kombinasi5050_belakangdisckembang,
		p.kombinasi5050_belakangdisckempis, p.kombinasi5050_belakangdisckembar,
		p.kombinasi5050_tengahdiscmono, p.kombinasi5050_tengahdiscstereo,
		p.kombinasi5050_tengahdisckembang, p.kombinasi5050_tengahdisckempis,
		p.kombinasi5050_tengahdisckembar, p.kombinasi5050_depandiscmono,
		p.kombinasi5050_depandiscstereo, p.kombinasi5050_depandisckembang,
		p.kombinasi5050_depandisckempis, p.kombinasi5050_depandisckembar, p.kombinasi5050_limitbuang,
		p.kombinasi5050_limittotal, p.macaukombinasi_minbet, p.macaukombinasi_maxbet,
		p.macaukombinasi_maxbuy, p.macaukombinasi_win, p.macaukombinasi_discount,
		p.macaukombinasi_limitbuang, p.macaukombinasi_limittotal, p.dasar_minbet, p.dasar_maxbet,
		p.dasar_maxbuy, p.dasar_keibesar, p.dasar_keikecil, p.dasar_keigenap, p.dasar_keiganjil,
		p.dasar_discbesar, p.dasar_disckecil, p.dasar_discigenap, p.dasar_discganjil,
		p.dasar_limitbuang, p.dasar_limittotal, p.shio_referal, p.shio_shiotahunini, p.shio_minbet,
		p.shio_maxbet, p.shio_maxbuy, p.shio_win, p.shio_disc, p.shio_limitbuang, p.shio_limittotal,
		$4, $5
	FROM ` + config.DB_mst_pasaran_togel + ` p
	WHERE p.idpasarantogel = $6`

	_, err := u.db.Exec(ctx, query,
		comppasaran.IDcomppasaran,
		comppasaran.IDcompany,
		comppasaran.IDpasarantogel,
		comppasaran.CreateBy,
		comppasaran.CreateAt,
		comppasaran.IDpasarantogel,
	)
	return err
}

func (u *companypasaranRepository) Update(ctx context.Context, comppasaran *domain.Companypasaran) error {
	type col struct {
		name string
		val  any
	}

	cols := []col{
		{"aliascomppasaran", comppasaran.Aliascomppasaran},
		{"urlpasaran", comppasaran.URLpasaran},
		{"pasarandiundi", comppasaran.Pasarandiundi},
		{"pasaranlibur", comppasaran.Pasaranlibur},
		{"displaypasaran", comppasaran.Display},
		{"angka_minbasket", comppasaran.AngkaMinbasket},
		{"angka_minbet", comppasaran.AngkaMinbet},
		{"angka_maxbet4d", comppasaran.AngkaMaxbet4d},
		{"angka_maxbet3d", comppasaran.AngkaMaxbet3d},
		{"angka_maxbet3dd", comppasaran.AngkaMaxbet3dd},
		{"angka_maxbet2d", comppasaran.AngkaMaxbet2d},
		{"angka_maxbet2dd", comppasaran.AngkaMaxbet2dd},
		{"angka_maxbet2dt", comppasaran.AngkaMaxbet2dt},
		{"angka_win4d", comppasaran.AngkaWin4d},
		{"angka_win3d", comppasaran.AngkaWin3d},
		{"angka_win3dd", comppasaran.AngkaWin3dd},
		{"angka_win2d", comppasaran.AngkaWin2d},
		{"angka_win2dd", comppasaran.AngkaWin2dd},
		{"angka_win2dt", comppasaran.AngkaWin2dt},
		{"angka_disc4d", comppasaran.AngkaDisc4d},
		{"angka_disc3d", comppasaran.AngkaDisc3d},
		{"angka_disc3dd", comppasaran.AngkaDisc3dd},
		{"angka_disc2d", comppasaran.AngkaDisc2d},
		{"angka_disc2dd", comppasaran.AngkaDisc2dd},
		{"angka_disc2dt", comppasaran.AngkaDisc2dt},
		{"angka_limitbuang4d", comppasaran.AngkaLimitbuang4d},
		{"angka_limitbuang3d", comppasaran.AngkaLimitbuang3d},
		{"angka_limitbuang3dd", comppasaran.AngkaLimitbuang3dd},
		{"angka_limitbuang2d", comppasaran.AngkaLimitbuang2d},
		{"angka_limitbuang2dd", comppasaran.AngkaLimitbuang2dd},
		{"angka_limitbuang2dt", comppasaran.AngkaLimitbuang2dt},
		{"angka_limittotal4d", comppasaran.AngkaLimittotal4d},
		{"angka_limittotal3d", comppasaran.AngkaLimittotal3d},
		{"angka_limittotal3dd", comppasaran.AngkaLimittotal3dd},
		{"angka_limittotal2d", comppasaran.AngkaLimittotal2d},
		{"angka_limittotal2dd", comppasaran.AngkaLimittotal2dd},
		{"angka_limittotal2dt", comppasaran.AngkaLimittotal2dt},
		{"angka_maxbet4d_full", comppasaran.AngkaMaxbet4dFull},
		{"angka_maxbet3d_full", comppasaran.AngkaMaxbet3dFull},
		{"angka_maxbet3dd_full", comppasaran.AngkaMaxbet3ddFull},
		{"angka_maxbet2d_full", comppasaran.AngkaMaxbet2dFull},
		{"angka_maxbet2dd_full", comppasaran.AngkaMaxbet2ddFull},
		{"angka_maxbet2dt_full", comppasaran.AngkaMaxbet2dtFull},
		{"angka_maxbet4d_bb", comppasaran.AngkaMaxbet4dBb},
		{"angka_maxbet3d_bb", comppasaran.AngkaMaxbet3dBb},
		{"angka_maxbet3dd_bb", comppasaran.AngkaMaxbet3ddBb},
		{"angka_maxbet2d_bb", comppasaran.AngkaMaxbet2dBb},
		{"angka_maxbet2dd_bb", comppasaran.AngkaMaxbet2ddBb},
		{"angka_maxbet2dt_bb", comppasaran.AngkaMaxbet2dtBb},
		{"angka_maxbet4d_bbdisc", comppasaran.AngkaMaxbet4dBbdisc},
		{"angka_maxbet3d_bbdisc", comppasaran.AngkaMaxbet3dBbdisc},
		{"angka_maxbet3dd_bbdisc", comppasaran.AngkaMaxbet3ddBbdisc},
		{"angka_maxbet2d_bbdisc", comppasaran.AngkaMaxbet2dBbdisc},
		{"angka_maxbet2dd_bbdisc", comppasaran.AngkaMaxbet2ddBbdisc},
		{"angka_maxbet2dt_bbdisc", comppasaran.AngkaMaxbet2dtBbdisc},
		{"angka_win4dnodisc", comppasaran.AngkaWin4dnodisc},
		{"angka_win3dnodisc", comppasaran.AngkaWin3dnodisc},
		{"angka_win3ddnodisc", comppasaran.AngkaWin3ddnodisc},
		{"angka_win2dnodisc", comppasaran.AngkaWin2dnodisc},
		{"angka_win2ddnodisc", comppasaran.AngkaWin2ddnodisc},
		{"angka_win2dtnodisc", comppasaran.AngkaWin2dtnodisc},
		{"angka_win4dbb_kena", comppasaran.AngkaWin4dbbKena},
		{"angka_win3dbb_kena", comppasaran.AngkaWin3dbbKena},
		{"angka_win3ddbb_kena", comppasaran.AngkaWin3ddbbKena},
		{"angka_win2dbb_kena", comppasaran.AngkaWin2dbbKena},
		{"angka_win2ddbb_kena", comppasaran.AngkaWin2ddbbKena},
		{"angka_win2dtbb_kena", comppasaran.AngkaWin2dtbbKena},
		{"angka_win4dbb", comppasaran.AngkaWin4dbb},
		{"angka_win3dbb", comppasaran.AngkaWin3dbb},
		{"angka_win3ddbb", comppasaran.AngkaWin3ddbb},
		{"angka_win2dbb", comppasaran.AngkaWin2dbb},
		{"angka_win2ddbb", comppasaran.AngkaWin2ddbb},
		{"angka_win2dtbb", comppasaran.AngkaWin2dtbb},
		{"angka_maxbuy4d", comppasaran.AngkaMaxbuy4d},
		{"angka_maxbuy3d", comppasaran.AngkaMaxbuy3d},
		{"angka_maxbuy3dd", comppasaran.AngkaMaxbuy3dd},
		{"angka_maxbuy2d", comppasaran.AngkaMaxbuy2d},
		{"angka_maxbuy2dd", comppasaran.AngkaMaxbuy2dd},
		{"angka_maxbuy2dt", comppasaran.AngkaMaxbuy2dt},
		{"angka_maxbet4d_fullbb", comppasaran.AngkaMaxbet4dFullbb},
		{"angka_maxbet3d_fullbb", comppasaran.AngkaMaxbet3dFullbb},
		{"angka_maxbet3dd_fullbb", comppasaran.AngkaMaxbet3ddFullbb},
		{"angka_maxbet2d_fullbb", comppasaran.AngkaMaxbet2dFullbb},
		{"angka_maxbet2dd_fullbb", comppasaran.AngkaMaxbet2ddFullbb},
		{"angka_maxbet2dt_fullbb", comppasaran.AngkaMaxbet2dtFullbb},
		{"angka_limitbuang4d_fullbb", comppasaran.AngkaLimitbuang4dFullbb},
		{"angka_limitbuang3d_fullbb", comppasaran.AngkaLimitbuang3dFullbb},
		{"angka_limitbuang3dd_fullbb", comppasaran.AngkaLimitbuang3ddFullbb},
		{"angka_limitbuang2d_fullbb", comppasaran.AngkaLimitbuang2dFullbb},
		{"angka_limitbuang2dd_fullbb", comppasaran.AngkaLimitbuang2ddFullbb},
		{"angka_limitbuang2dt_fullbb", comppasaran.AngkaLimitbuang2dtFullbb},
		{"angka_limittotal4d_fullbb", comppasaran.AngkaLimittotal4dFullbb},
		{"angka_limittotal3d_fullbb", comppasaran.AngkaLimittotal3dFullbb},
		{"angka_limittotal3dd_fullbb", comppasaran.AngkaLimittotal3ddFullbb},
		{"angka_limittotal2d_fullbb", comppasaran.AngkaLimittotal2dFullbb},
		{"angka_limittotal2dd_fullbb", comppasaran.AngkaLimittotal2ddFullbb},
		{"angka_limittotal2dt_fullbb", comppasaran.AngkaLimittotal2dtFullbb},
		{"angka_limitline_4d", comppasaran.AngkaLimitline4d},
		{"angka_limitline_3d", comppasaran.AngkaLimitline3d},
		{"angka_limitline_2d", comppasaran.AngkaLimitline2d},
		{"angka_limitline_2dd", comppasaran.AngkaLimitline2dd},
		{"angka_limitline_2dt", comppasaran.AngkaLimitline2dt},
		{"angka_limitline_3dd", comppasaran.AngkaLimitline3dd},
		{"angka_bbfs", comppasaran.AngkaBbfs},
		{"cb_minbet", comppasaran.CbMinbet},
		{"cb_maxbet", comppasaran.CbMaxbet},
		{"cb_maxbuy", comppasaran.CbMaxbuy},
		{"cb_win", comppasaran.CbWin},
		{"cb_disc", comppasaran.CbDisc},
		{"cb_limitbuang", comppasaran.CbLimitbuang},
		{"cb_limitotal", comppasaran.CbLimitotal},
		{"cmacau_minbet", comppasaran.CmacauMinbet},
		{"cmacau_maxbet", comppasaran.CmacauMaxbet},
		{"cmacau_maxbuy", comppasaran.CmacauMaxbuy},
		{"cmacau_win2digit", comppasaran.CmacauWin2digit},
		{"cmacau_win3digit", comppasaran.CmacauWin3digit},
		{"cmacau_win4digit", comppasaran.CmacauWin4digit},
		{"cmacau_disc", comppasaran.CmacauDisc},
		{"cmacau_limitbuang", comppasaran.CmacauLimitbuang},
		{"cmacau_limittotal", comppasaran.CmacauLimittotal},
		{"cnaga_minbet", comppasaran.CnagaMinbet},
		{"cnaga_maxbet", comppasaran.CnagaMaxbet},
		{"cnaga_maxbuy", comppasaran.CnagaMaxbuy},
		{"cnaga_win3digit", comppasaran.CnagaWin3digit},
		{"cnaga_win4digit", comppasaran.CnagaWin4digit},
		{"cnaga_disc", comppasaran.CnagaDisc},
		{"cnaga_limitbuang", comppasaran.CnagaLimitbuang},
		{"cnaga_limittotal", comppasaran.CnagaLimittotal},
		{"cjitu_minbet", comppasaran.CjituMinbet},
		{"cjitu_maxbet", comppasaran.CjituMaxbet},
		{"cjitu_maxbuy", comppasaran.CjituMaxbuy},
		{"cjitu_winas", comppasaran.CjituWinas},
		{"cjitu_winkop", comppasaran.CjituWinkop},
		{"cjitu_winkepala", comppasaran.CjituWinkepala},
		{"cjitu_winekor", comppasaran.CjituWinekor},
		{"cjitu_desic", comppasaran.CjituDesic},
		{"cjitu_limitbuang", comppasaran.CjituLimitbuang},
		{"cjitu_limitotal", comppasaran.CjituLimitotal},
		{"umum5050_minbet", comppasaran.Umum5050Minbet},
		{"umum5050_maxbet", comppasaran.Umum5050Maxbet},
		{"umum5050_maxbuy", comppasaran.Umum5050Maxbuy},
		{"umum5050_keibesar", comppasaran.Umum5050Keibesar},
		{"umum5050_keikecil", comppasaran.Umum5050Keikecil},
		{"umum5050_keigenap", comppasaran.Umum5050Keigenap},
		{"umum5050_keiganjil", comppasaran.Umum5050Keiganjil},
		{"umum5050_keitengah", comppasaran.Umum5050Keitengah},
		{"umum5050_keitepi", comppasaran.Umum5050Keitepi},
		{"umum5050_discbesar", comppasaran.Umum5050Discbesar},
		{"umum5050_disckecil", comppasaran.Umum5050Disckecil},
		{"umum5050_discgenap", comppasaran.Umum5050Discgenap},
		{"umum5050_discganjil", comppasaran.Umum5050Discganjil},
		{"umum5050_disctengah", comppasaran.Umum5050Disctengah},
		{"umum5050_disctepi", comppasaran.Umum5050Disctepi},
		{"umum5050_limitbuang", comppasaran.Umum5050Limitbuang},
		{"umum5050_limittotal", comppasaran.Umum5050Limittotal},
		{"special5050_minbet", comppasaran.Special5050Minbet},
		{"special5050_maxbet", comppasaran.Special5050Maxbet},
		{"special5050_maxbuy", comppasaran.Special5050Maxbuy},
		{"special5050_keiasganjil", comppasaran.Special5050Keiasganjil},
		{"special5050_keiasgenap", comppasaran.Special5050Keiasgenap},
		{"special5050_keiasbesar", comppasaran.Special5050Keiasbesar},
		{"special5050_keiaskecil", comppasaran.Special5050Keiaskecil},
		{"special5050_keikopganjil", comppasaran.Special5050Keikopganjil},
		{"special5050_keikopgenap", comppasaran.Special5050Keikopgenap},
		{"special5050_keikopbesar", comppasaran.Special5050Keikopbesar},
		{"special5050_keikopkecil", comppasaran.Special5050Keikopkecil},
		{"special5050_keikepalaganjil", comppasaran.Special5050Keikepalaganjil},
		{"special5050_keikepalagenap", comppasaran.Special5050Keikepalagenap},
		{"special5050_keikepalabesar", comppasaran.Special5050Keikepalabesar},
		{"special5050_keikepalakecil", comppasaran.Special5050Keikepalakecil},
		{"special5050_keiekorganjil", comppasaran.Special5050Keiekorganjil},
		{"special5050_keiekorgenap", comppasaran.Special5050Keiekorgenap},
		{"special5050_keiekorbesar", comppasaran.Special5050Keiekorbesar},
		{"special5050_keiekorkecil", comppasaran.Special5050Keiekorkecil},
		{"special5050_discasganjil", comppasaran.Special5050Discasganjil},
		{"special5050_discasgenap", comppasaran.Special5050Discasgenap},
		{"special5050_discasbesar", comppasaran.Special5050Discasbesar},
		{"special5050_discaskecil", comppasaran.Special5050Discaskecil},
		{"special5050_disckopganjil", comppasaran.Special5050Disckopganjil},
		{"special5050_disckopgenap", comppasaran.Special5050Disckopgenap},
		{"special5050_disckopbesar", comppasaran.Special5050Disckopbesar},
		{"special5050_disckopkecil", comppasaran.Special5050Disckopkecil},
		{"special5050_disckepalaganjil", comppasaran.Special5050Disckepalaganjil},
		{"special5050_disckepalagenap", comppasaran.Special5050Disckepalagenap},
		{"special5050_disckepalabesar", comppasaran.Special5050Disckepalabesar},
		{"special5050_disckepalakecil", comppasaran.Special5050Disckepalakecil},
		{"special5050_discekorganjil", comppasaran.Special5050Discekorganjil},
		{"special5050_discekorgenap", comppasaran.Special5050Discekorgenap},
		{"special5050_discekorbesar", comppasaran.Special5050Discekorbesar},
		{"special5050_discekorkecil", comppasaran.Special5050Discekorkecil},
		{"special5050_limitbuang", comppasaran.Special5050Limitbuang},
		{"special5050_limittotal", comppasaran.Special5050Limittotal},
		{"kombinasi5050_minbet", comppasaran.Kombinasi5050Minbet},
		{"kombinasi5050_maxbet", comppasaran.Kombinasi5050Maxbet},
		{"kombinasi5050_maxbuy", comppasaran.Kombinasi5050Maxbuy},
		{"kombinasi5050_belakangkeimono", comppasaran.Kombinasi5050Belakangkeimono},
		{"kombinasi5050_belakangkeistereo", comppasaran.Kombinasi5050Belakangkeistereo},
		{"kombinasi5050_belakangkeikembang", comppasaran.Kombinasi5050Belakangkeikembang},
		{"kombinasi5050_belakangkeikempis", comppasaran.Kombinasi5050Belakangkeikempis},
		{"kombinasi5050_belakangkeikembar", comppasaran.Kombinasi5050Belakangkeikembar},
		{"kombinasi5050_tengahkeimono", comppasaran.Kombinasi5050Tengahkeimono},
		{"kombinasi5050_tengahkeistereo", comppasaran.Kombinasi5050Tengahkeistereo},
		{"kombinasi5050_tengahkeikembang", comppasaran.Kombinasi5050Tengahkeikembang},
		{"kombinasi5050_tengahkeikempis", comppasaran.Kombinasi5050Tengahkeikempis},
		{"kombinasi5050_tengahkeikembar", comppasaran.Kombinasi5050Tengahkeikembar},
		{"kombinasi5050_depankeimono", comppasaran.Kombinasi5050Depankeimono},
		{"kombinasi5050_depankeistereo", comppasaran.Kombinasi5050Depankeistereo},
		{"kombinasi5050_depankeikembang", comppasaran.Kombinasi5050Depankeikembang},
		{"kombinasi5050_depankeikempis", comppasaran.Kombinasi5050Depankeikempis},
		{"kombinasi5050_depankeikembar", comppasaran.Kombinasi5050Depankeikembar},
		{"kombinasi5050_belakangdiscmono", comppasaran.Kombinasi5050Belakangdiscmono},
		{"kombinasi5050_belakangdiscstereo", comppasaran.Kombinasi5050Belakangdiscstereo},
		{"kombinasi5050_belakangdisckembang", comppasaran.Kombinasi5050Belakangdisckembang},
		{"kombinasi5050_belakangdisckempis", comppasaran.Kombinasi5050Belakangdisckempis},
		{"kombinasi5050_belakangdisckembar", comppasaran.Kombinasi5050Belakangdisckembar},
		{"kombinasi5050_tengahdiscmono", comppasaran.Kombinasi5050Tengahdiscmono},
		{"kombinasi5050_tengahdiscstereo", comppasaran.Kombinasi5050Tengahdiscstereo},
		{"kombinasi5050_tengahdisckembang", comppasaran.Kombinasi5050Tengahdisckembang},
		{"kombinasi5050_tengahdisckempis", comppasaran.Kombinasi5050Tengahdisckempis},
		{"kombinasi5050_tengahdisckembar", comppasaran.Kombinasi5050Tengahdisckembar},
		{"kombinasi5050_depandiscmono", comppasaran.Kombinasi5050Depandiscmono},
		{"kombinasi5050_depandiscstereo", comppasaran.Kombinasi5050Depandiscstereo},
		{"kombinasi5050_depandisckembang", comppasaran.Kombinasi5050Depandisckembang},
		{"kombinasi5050_depandisckempis", comppasaran.Kombinasi5050Depandisckempis},
		{"kombinasi5050_depandisckembar", comppasaran.Kombinasi5050Depandisckembar},
		{"kombinasi5050_limitbuang", comppasaran.Kombinasi5050Limitbuang},
		{"kombinasi5050_limittotal", comppasaran.Kombinasi5050Limittotal},
		{"macaukombinasi_minbet", comppasaran.MacaukombinasiMinbet},
		{"macaukombinasi_maxbet", comppasaran.MacaukombinasiMaxbet},
		{"macaukombinasi_maxbuy", comppasaran.MacaukombinasiMaxbuy},
		{"macaukombinasi_win", comppasaran.MacaukombinasiWin},
		{"macaukombinasi_discount", comppasaran.MacaukombinasiDiscount},
		{"macaukombinasi_limitbuang", comppasaran.MacaukombinasiLimitbuang},
		{"macaukombinasi_limittotal", comppasaran.MacaukombinasiLimittotal},
		{"dasar_minbet", comppasaran.DasarMinbet},
		{"dasar_maxbet", comppasaran.DasarMaxbet},
		{"dasar_maxbuy", comppasaran.DasarMaxbuy},
		{"dasar_keibesar", comppasaran.DasarKeibesar},
		{"dasar_keikecil", comppasaran.DasarKeikecil},
		{"dasar_keigenap", comppasaran.DasarKeigenap},
		{"dasar_keiganjil", comppasaran.DasarKeiganjil},
		{"dasar_discbesar", comppasaran.DasarDiscbesar},
		{"dasar_disckecil", comppasaran.DasarDisckecil},
		{"dasar_discigenap", comppasaran.DasarDiscigenap},
		{"dasar_discganjil", comppasaran.DasarDiscganjil},
		{"dasar_limitbuang", comppasaran.DasarLimitbuang},
		{"dasar_limittotal", comppasaran.DasarLimittotal},
		{"shio_referal", comppasaran.ShioReferal},
		{"shio_shiotahunini", comppasaran.ShioShiotahunini},
		{"shio_minbet", comppasaran.ShioMinbet},
		{"shio_maxbet", comppasaran.ShioMaxbet},
		{"shio_maxbuy", comppasaran.ShioMaxbuy},
		{"shio_win", comppasaran.ShioWin},
		{"shio_disc", comppasaran.ShioDisc},
		{"shio_limitbuang", comppasaran.ShioLimitbuang},
		{"shio_limittotal", comppasaran.ShioLimittotal},
		{"statuscompasaran", comppasaran.Status},
		{"update_by", comppasaran.UpdateBy},
		{"update_at", comppasaran.UpdateAt},
	}

	setClauses := make([]string, 0, len(cols))
	args := make([]any, 0, len(cols)+1)
	for i, cl := range cols {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", cl.name, i+1))
		args = append(args, cl.val)
	}

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE idcomppasaran = $%d`,
		config.DB_mst_company_pasaran,
		strings.Join(setClauses, ", "),
		len(cols)+1,
	)
	args = append(args, comppasaran.IDcomppasaran)

	res, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func (u *companypasaranRepository) SavejadwalCopyPasaran(ctx context.Context, compjadwal *domain.Companyjadwalpasaran, idpasarantogel string) error {
	query := `INSERT INTO ` + config.DB_mst_company_jadwaltogel + ` (
		idjadwalcomppasaran, idcomppasaran, haripasaran, jamtutup, jamjadwal, jamopen,
		create_by, create_at
	) SELECT
		to_char($4::timestamp, 'YYMM') || replace(gen_random_uuid()::text, '-', ''),
		$1, j.haripasaran, j.jamtutup, j.jamjadwal, j.jamopen, $3, $4
	FROM ` + config.DB_mst_pasaran_jadwaltogel + ` j
	WHERE j.idpasarantogel = $2`

	_, err := u.db.Exec(ctx, query,
		compjadwal.IDcomppasaran,
		idpasarantogel,
		compjadwal.CreateBy,
		compjadwal.CreateAt,
	)
	return err
}
func (u *companypasaranRepository) Savejadwal(ctx context.Context, compjadwal *domain.Companyjadwalpasaran) error {
	query := `INSERT INTO ` + config.DB_mst_company_jadwaltogel + ` (
		idjadwalcomppasaran,
		idcomppasaran,
		haripasaran,
		jamtutup,
		jamjadwal,
		jamopen,
		create_by,
		create_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	)`

	_, err := u.db.Exec(ctx, query,
		compjadwal.IDjadwalcomppasaran,
		compjadwal.IDcomppasaran,
		compjadwal.Haripasaran,
		compjadwal.Jamtutup,
		compjadwal.Jamjadwal,
		compjadwal.Jamopen,
		compjadwal.CreateBy,
		compjadwal.CreateAt,
	)
	return err
}
func (u *companypasaranRepository) DeleteJadwal(ctx context.Context, idcomppasaran string) error {
	query := `DELETE FROM ` + config.DB_mst_company_jadwaltogel + ` WHERE idcomppasaran = $1`
	_, err := u.db.Exec(ctx, query, idcomppasaran)
	return err
}
