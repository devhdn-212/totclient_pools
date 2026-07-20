package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/internal/config"
	"github.com/jackc/pgx/v5"
)

type pasaranRepository struct {
	db DBExecutor
}

func NewPasaranRepository(db DBExecutor) domain.PasaranRepository {
	return &pasaranRepository{
		db: db,
	}
}

func (u *pasaranRepository) FindAll(ctx context.Context, idcomp string) ([]domain.Pasaran, error) {
	query := `SELECT * FROM ` + config.DB_mst_company_pasaran + ` 
			  WHERE idcompany=$1 
			  ORDER BY displaypasaran ASC`

	rows, err := u.db.Query(ctx, query, idcomp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Companyconftoto
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Pasaran])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (u *pasaranRepository) FindJadwal(ctx context.Context, idcomppasaran string) ([]domain.Pasaranjadwal, error) {
	query := `SELECT * FROM ` + config.DB_mst_company_jadwaltogel + ` 
			  WHERE idcomppasaran = $1  
			  ORDER BY create_at DESC`

	rows, err := u.db.Query(ctx, query, idcomppasaran)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Companyconftoto
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Pasaranjadwal])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *pasaranRepository) FindByID(ctx context.Context, id, idcomp string) (domain.Pasaran, error) {
	var record domain.Pasaran
	query := `SELECT idcomppasaran FROM ` + config.DB_mst_company_pasaran + ` 
			WHERE idcomppasaran = $1 AND idcompany = $2 
			LIMIT 1`

	err := u.db.QueryRow(ctx, query, id, idcomp).Scan(&record.IDcomppasaran)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Pasaran{}, nil
		}
		return record, err
	}
	return record, nil
}

func (u *pasaranRepository) Update(ctx context.Context, pasaran *domain.Pasaran) error {
	type col struct {
		name string
		val  any
	}

	cols := []col{
		{"aliaspasaran", pasaran.Aliascomppasaran},
		{"urlpasaran", pasaran.URLpasaran},
		{"urllogo", pasaran.URLlogo},
		{"pasarandiundi", pasaran.Pasarandiundi},
		{"pasaranlibur", pasaran.Pasaranlibur},
		{"displaypasaran", pasaran.Display},
		{"angka_minbet", pasaran.AngkaMinbet},
		{"angka_maxbet4d", pasaran.AngkaMaxbet4d},
		{"angka_maxbet3d", pasaran.AngkaMaxbet3d},
		{"angka_maxbet3dd", pasaran.AngkaMaxbet3dd},
		{"angka_maxbet2d", pasaran.AngkaMaxbet2d},
		{"angka_maxbet2dd", pasaran.AngkaMaxbet2dd},
		{"angka_maxbet2dt", pasaran.AngkaMaxbet2dt},
		{"angka_win4d", pasaran.AngkaWin4d},
		{"angka_win3d", pasaran.AngkaWin3d},
		{"angka_win3dd", pasaran.AngkaWin3dd},
		{"angka_win2d", pasaran.AngkaWin2d},
		{"angka_win2dd", pasaran.AngkaWin2dd},
		{"angka_win2dt", pasaran.AngkaWin2dt},
		{"angka_disc4d", pasaran.AngkaDisc4d},
		{"angka_disc3d", pasaran.AngkaDisc3d},
		{"angka_disc3dd", pasaran.AngkaDisc3dd},
		{"angka_disc2d", pasaran.AngkaDisc2d},
		{"angka_disc2dd", pasaran.AngkaDisc2dd},
		{"angka_disc2dt", pasaran.AngkaDisc2dt},
		{"angka_limitbuang4d", pasaran.AngkaLimitbuang4d},
		{"angka_limitbuang3d", pasaran.AngkaLimitbuang3d},
		{"angka_limitbuang3dd", pasaran.AngkaLimitbuang3dd},
		{"angka_limitbuang2d", pasaran.AngkaLimitbuang2d},
		{"angka_limitbuang2dd", pasaran.AngkaLimitbuang2dd},
		{"angka_limitbuang2dt", pasaran.AngkaLimitbuang2dt},
		{"angka_limittotal4d", pasaran.AngkaLimittotal4d},
		{"angka_limittotal3d", pasaran.AngkaLimittotal3d},
		{"angka_limittotal3dd", pasaran.AngkaLimittotal3dd},
		{"angka_limittotal2d", pasaran.AngkaLimittotal2d},
		{"angka_limittotal2dd", pasaran.AngkaLimittotal2dd},
		{"angka_limittotal2dt", pasaran.AngkaLimittotal2dt},
		{"angka_maxbet4d_full", pasaran.AngkaMaxbet4dFull},
		{"angka_maxbet3d_full", pasaran.AngkaMaxbet3dFull},
		{"angka_maxbet3dd_full", pasaran.AngkaMaxbet3ddFull},
		{"angka_maxbet2d_full", pasaran.AngkaMaxbet2dFull},
		{"angka_maxbet2dd_full", pasaran.AngkaMaxbet2ddFull},
		{"angka_maxbet2dt_full", pasaran.AngkaMaxbet2dtFull},
		{"angka_maxbet4d_bb", pasaran.AngkaMaxbet4dBb},
		{"angka_maxbet3d_bb", pasaran.AngkaMaxbet3dBb},
		{"angka_maxbet3dd_bb", pasaran.AngkaMaxbet3ddBb},
		{"angka_maxbet2d_bb", pasaran.AngkaMaxbet2dBb},
		{"angka_maxbet2dd_bb", pasaran.AngkaMaxbet2ddBb},
		{"angka_maxbet2dt_bb", pasaran.AngkaMaxbet2dtBb},
		{"angka_maxbet4d_bbdisc", pasaran.AngkaMaxbet4dBbdisc},
		{"angka_maxbet3d_bbdisc", pasaran.AngkaMaxbet3dBbdisc},
		{"angka_maxbet3dd_bbdisc", pasaran.AngkaMaxbet3ddBbdisc},
		{"angka_maxbet2d_bbdisc", pasaran.AngkaMaxbet2dBbdisc},
		{"angka_maxbet2dd_bbdisc", pasaran.AngkaMaxbet2ddBbdisc},
		{"angka_maxbet2dt_bbdisc", pasaran.AngkaMaxbet2dtBbdisc},
		{"angka_win4dnodisc", pasaran.AngkaWin4dnodisc},
		{"angka_win3dnodisc", pasaran.AngkaWin3dnodisc},
		{"angka_win3ddnodisc", pasaran.AngkaWin3ddnodisc},
		{"angka_win2dnodisc", pasaran.AngkaWin2dnodisc},
		{"angka_win2ddnodisc", pasaran.AngkaWin2ddnodisc},
		{"angka_win2dtnodisc", pasaran.AngkaWin2dtnodisc},
		{"angka_win4dbb_kena", pasaran.AngkaWin4dbbKena},
		{"angka_win3dbb_kena", pasaran.AngkaWin3dbbKena},
		{"angka_win3ddbb_kena", pasaran.AngkaWin3ddbbKena},
		{"angka_win2dbb_kena", pasaran.AngkaWin2dbbKena},
		{"angka_win2ddbb_kena", pasaran.AngkaWin2ddbbKena},
		{"angka_win2dtbb_kena", pasaran.AngkaWin2dtbbKena},
		{"angka_win4dbb", pasaran.AngkaWin4dbb},
		{"angka_win3dbb", pasaran.AngkaWin3dbb},
		{"angka_win3ddbb", pasaran.AngkaWin3ddbb},
		{"angka_win2dbb", pasaran.AngkaWin2dbb},
		{"angka_win2ddbb", pasaran.AngkaWin2ddbb},
		{"angka_win2dtbb", pasaran.AngkaWin2dtbb},
		{"angka_maxbuy4d", pasaran.AngkaMaxbuy4d},
		{"angka_maxbuy3d", pasaran.AngkaMaxbuy3d},
		{"angka_maxbuy3dd", pasaran.AngkaMaxbuy3dd},
		{"angka_maxbuy2d", pasaran.AngkaMaxbuy2d},
		{"angka_maxbuy2dd", pasaran.AngkaMaxbuy2dd},
		{"angka_maxbuy2dt", pasaran.AngkaMaxbuy2dt},
		{"angka_maxbet4d_fullbb", pasaran.AngkaMaxbet4dFullbb},
		{"angka_maxbet3d_fullbb", pasaran.AngkaMaxbet3dFullbb},
		{"angka_maxbet3dd_fullbb", pasaran.AngkaMaxbet3ddFullbb},
		{"angka_maxbet2d_fullbb", pasaran.AngkaMaxbet2dFullbb},
		{"angka_maxbet2dd_fullbb", pasaran.AngkaMaxbet2ddFullbb},
		{"angka_maxbet2dt_fullbb", pasaran.AngkaMaxbet2dtFullbb},
		{"angka_limitbuang4d_fullbb", pasaran.AngkaLimitbuang4dFullbb},
		{"angka_limitbuang3d_fullbb", pasaran.AngkaLimitbuang3dFullbb},
		{"angka_limitbuang3dd_fullbb", pasaran.AngkaLimitbuang3ddFullbb},
		{"angka_limitbuang2d_fullbb", pasaran.AngkaLimitbuang2dFullbb},
		{"angka_limitbuang2dd_fullbb", pasaran.AngkaLimitbuang2ddFullbb},
		{"angka_limitbuang2dt_fullbb", pasaran.AngkaLimitbuang2dtFullbb},
		{"angka_limittotal4d_fullbb", pasaran.AngkaLimittotal4dFullbb},
		{"angka_limittotal3d_fullbb", pasaran.AngkaLimittotal3dFullbb},
		{"angka_limittotal3dd_fullbb", pasaran.AngkaLimittotal3ddFullbb},
		{"angka_limittotal2d_fullbb", pasaran.AngkaLimittotal2dFullbb},
		{"angka_limittotal2dd_fullbb", pasaran.AngkaLimittotal2ddFullbb},
		{"angka_limittotal2dt_fullbb", pasaran.AngkaLimittotal2dtFullbb},
		{"angka_limitline_4d", pasaran.AngkaLimitline4d},
		{"angka_limitline_3d", pasaran.AngkaLimitline3d},
		{"angka_limitline_2d", pasaran.AngkaLimitline2d},
		{"angka_limitline_2dd", pasaran.AngkaLimitline2dd},
		{"angka_limitline_2dt", pasaran.AngkaLimitline2dt},
		{"angka_limitline_3dd", pasaran.AngkaLimitline3dd},
		{"angka_bbfs", pasaran.AngkaBbfs},
		{"cb_minbet", pasaran.CbMinbet},
		{"cb_maxbet", pasaran.CbMaxbet},
		{"cb_maxbuy", pasaran.CbMaxbuy},
		{"cb_win", pasaran.CbWin},
		{"cb_disc", pasaran.CbDisc},
		{"cb_limitbuang", pasaran.CbLimitbuang},
		{"cb_limitotal", pasaran.CbLimitotal},
		{"cmacau_minbet", pasaran.CmacauMinbet},
		{"cmacau_maxbet", pasaran.CmacauMaxbet},
		{"cmacau_maxbuy", pasaran.CmacauMaxbuy},
		{"cmacau_win2digit", pasaran.CmacauWin2digit},
		{"cmacau_win3digit", pasaran.CmacauWin3digit},
		{"cmacau_win4digit", pasaran.CmacauWin4digit},
		{"cmacau_disc", pasaran.CmacauDisc},
		{"cmacau_limitbuang", pasaran.CmacauLimitbuang},
		{"cmacau_limittotal", pasaran.CmacauLimittotal},
		{"cnaga_minbet", pasaran.CnagaMinbet},
		{"cnaga_maxbet", pasaran.CnagaMaxbet},
		{"cnaga_maxbuy", pasaran.CnagaMaxbuy},
		{"cnaga_win3digit", pasaran.CnagaWin3digit},
		{"cnaga_win4digit", pasaran.CnagaWin4digit},
		{"cnaga_disc", pasaran.CnagaDisc},
		{"cnaga_limitbuang", pasaran.CnagaLimitbuang},
		{"cnaga_limittotal", pasaran.CnagaLimittotal},
		{"cjitu_minbet", pasaran.CjituMinbet},
		{"cjitu_maxbet", pasaran.CjituMaxbet},
		{"cjitu_maxbuy", pasaran.CjituMaxbuy},
		{"cjitu_winas", pasaran.CjituWinas},
		{"cjitu_winkop", pasaran.CjituWinkop},
		{"cjitu_winkepala", pasaran.CjituWinkepala},
		{"cjitu_winekor", pasaran.CjituWinekor},
		{"cjitu_desic", pasaran.CjituDesic},
		{"cjitu_limitbuang", pasaran.CjituLimitbuang},
		{"cjitu_limitotal", pasaran.CjituLimitotal},
		{"umum5050_minbet", pasaran.Umum5050Minbet},
		{"umum5050_maxbet", pasaran.Umum5050Maxbet},
		{"umum5050_maxbuy", pasaran.Umum5050Maxbuy},
		{"umum5050_keibesar", pasaran.Umum5050Keibesar},
		{"umum5050_keikecil", pasaran.Umum5050Keikecil},
		{"umum5050_keigenap", pasaran.Umum5050Keigenap},
		{"umum5050_keiganjil", pasaran.Umum5050Keiganjil},
		{"umum5050_keitengah", pasaran.Umum5050Keitengah},
		{"umum5050_keitepi", pasaran.Umum5050Keitepi},
		{"umum5050_discbesar", pasaran.Umum5050Discbesar},
		{"umum5050_disckecil", pasaran.Umum5050Disckecil},
		{"umum5050_discgenap", pasaran.Umum5050Discgenap},
		{"umum5050_discganjil", pasaran.Umum5050Discganjil},
		{"umum5050_disctengah", pasaran.Umum5050Disctengah},
		{"umum5050_disctepi", pasaran.Umum5050Disctepi},
		{"umum5050_limitbuang", pasaran.Umum5050Limitbuang},
		{"umum5050_limittotal", pasaran.Umum5050Limittotal},
		{"special5050_minbet", pasaran.Special5050Minbet},
		{"special5050_maxbet", pasaran.Special5050Maxbet},
		{"special5050_maxbuy", pasaran.Special5050Maxbuy},
		{"special5050_keiasganjil", pasaran.Special5050Keiasganjil},
		{"special5050_keiasgenap", pasaran.Special5050Keiasgenap},
		{"special5050_keiasbesar", pasaran.Special5050Keiasbesar},
		{"special5050_keiaskecil", pasaran.Special5050Keiaskecil},
		{"special5050_keikopganjil", pasaran.Special5050Keikopganjil},
		{"special5050_keikopgenap", pasaran.Special5050Keikopgenap},
		{"special5050_keikopbesar", pasaran.Special5050Keikopbesar},
		{"special5050_keikopkecil", pasaran.Special5050Keikopkecil},
		{"special5050_keikepalaganjil", pasaran.Special5050Keikepalaganjil},
		{"special5050_keikepalagenap", pasaran.Special5050Keikepalagenap},
		{"special5050_keikepalabesar", pasaran.Special5050Keikepalabesar},
		{"special5050_keikepalakecil", pasaran.Special5050Keikepalakecil},
		{"special5050_keiekorganjil", pasaran.Special5050Keiekorganjil},
		{"special5050_keiekorgenap", pasaran.Special5050Keiekorgenap},
		{"special5050_keiekorbesar", pasaran.Special5050Keiekorbesar},
		{"special5050_keiekorkecil", pasaran.Special5050Keiekorkecil},
		{"special5050_discasganjil", pasaran.Special5050Discasganjil},
		{"special5050_discasgenap", pasaran.Special5050Discasgenap},
		{"special5050_discasbesar", pasaran.Special5050Discasbesar},
		{"special5050_discaskecil", pasaran.Special5050Discaskecil},
		{"special5050_disckopganjil", pasaran.Special5050Disckopganjil},
		{"special5050_disckopgenap", pasaran.Special5050Disckopgenap},
		{"special5050_disckopbesar", pasaran.Special5050Disckopbesar},
		{"special5050_disckopkecil", pasaran.Special5050Disckopkecil},
		{"special5050_disckepalaganjil", pasaran.Special5050Disckepalaganjil},
		{"special5050_disckepalagenap", pasaran.Special5050Disckepalagenap},
		{"special5050_disckepalabesar", pasaran.Special5050Disckepalabesar},
		{"special5050_disckepalakecil", pasaran.Special5050Disckepalakecil},
		{"special5050_discekorganjil", pasaran.Special5050Discekorganjil},
		{"special5050_discekorgenap", pasaran.Special5050Discekorgenap},
		{"special5050_discekorbesar", pasaran.Special5050Discekorbesar},
		{"special5050_discekorkecil", pasaran.Special5050Discekorkecil},
		{"special5050_limitbuang", pasaran.Special5050Limitbuang},
		{"special5050_limittotal", pasaran.Special5050Limittotal},
		{"kombinasi5050_minbet", pasaran.Kombinasi5050Minbet},
		{"kombinasi5050_maxbet", pasaran.Kombinasi5050Maxbet},
		{"kombinasi5050_maxbuy", pasaran.Kombinasi5050Maxbuy},
		{"kombinasi5050_belakangkeimono", pasaran.Kombinasi5050Belakangkeimono},
		{"kombinasi5050_belakangkeistereo", pasaran.Kombinasi5050Belakangkeistereo},
		{"kombinasi5050_belakangkeikembang", pasaran.Kombinasi5050Belakangkeikembang},
		{"kombinasi5050_belakangkeikempis", pasaran.Kombinasi5050Belakangkeikempis},
		{"kombinasi5050_belakangkeikembar", pasaran.Kombinasi5050Belakangkeikembar},
		{"kombinasi5050_tengahkeimono", pasaran.Kombinasi5050Tengahkeimono},
		{"kombinasi5050_tengahkeistereo", pasaran.Kombinasi5050Tengahkeistereo},
		{"kombinasi5050_tengahkeikembang", pasaran.Kombinasi5050Tengahkeikembang},
		{"kombinasi5050_tengahkeikempis", pasaran.Kombinasi5050Tengahkeikempis},
		{"kombinasi5050_tengahkeikembar", pasaran.Kombinasi5050Tengahkeikembar},
		{"kombinasi5050_depankeimono", pasaran.Kombinasi5050Depankeimono},
		{"kombinasi5050_depankeistereo", pasaran.Kombinasi5050Depankeistereo},
		{"kombinasi5050_depankeikembang", pasaran.Kombinasi5050Depankeikembang},
		{"kombinasi5050_depankeikempis", pasaran.Kombinasi5050Depankeikempis},
		{"kombinasi5050_depankeikembar", pasaran.Kombinasi5050Depankeikembar},
		{"kombinasi5050_belakangdiscmono", pasaran.Kombinasi5050Belakangdiscmono},
		{"kombinasi5050_belakangdiscstereo", pasaran.Kombinasi5050Belakangdiscstereo},
		{"kombinasi5050_belakangdisckembang", pasaran.Kombinasi5050Belakangdisckembang},
		{"kombinasi5050_belakangdisckempis", pasaran.Kombinasi5050Belakangdisckempis},
		{"kombinasi5050_belakangdisckembar", pasaran.Kombinasi5050Belakangdisckembar},
		{"kombinasi5050_tengahdiscmono", pasaran.Kombinasi5050Tengahdiscmono},
		{"kombinasi5050_tengahdiscstereo", pasaran.Kombinasi5050Tengahdiscstereo},
		{"kombinasi5050_tengahdisckembang", pasaran.Kombinasi5050Tengahdisckembang},
		{"kombinasi5050_tengahdisckempis", pasaran.Kombinasi5050Tengahdisckempis},
		{"kombinasi5050_tengahdisckembar", pasaran.Kombinasi5050Tengahdisckembar},
		{"kombinasi5050_depandiscmono", pasaran.Kombinasi5050Depandiscmono},
		{"kombinasi5050_depandiscstereo", pasaran.Kombinasi5050Depandiscstereo},
		{"kombinasi5050_depandisckembang", pasaran.Kombinasi5050Depandisckembang},
		{"kombinasi5050_depandisckempis", pasaran.Kombinasi5050Depandisckempis},
		{"kombinasi5050_depandisckembar", pasaran.Kombinasi5050Depandisckembar},
		{"kombinasi5050_limitbuang", pasaran.Kombinasi5050Limitbuang},
		{"kombinasi5050_limittotal", pasaran.Kombinasi5050Limittotal},
		{"macaukombinasi_minbet", pasaran.MacaukombinasiMinbet},
		{"macaukombinasi_maxbet", pasaran.MacaukombinasiMaxbet},
		{"macaukombinasi_maxbuy", pasaran.MacaukombinasiMaxbuy},
		{"macaukombinasi_win", pasaran.MacaukombinasiWin},
		{"macaukombinasi_discount", pasaran.MacaukombinasiDiscount},
		{"macaukombinasi_limitbuang", pasaran.MacaukombinasiLimitbuang},
		{"macaukombinasi_limittotal", pasaran.MacaukombinasiLimittotal},
		{"dasar_minbet", pasaran.DasarMinbet},
		{"dasar_maxbet", pasaran.DasarMaxbet},
		{"dasar_maxbuy", pasaran.DasarMaxbuy},
		{"dasar_keibesar", pasaran.DasarKeibesar},
		{"dasar_keikecil", pasaran.DasarKeikecil},
		{"dasar_keigenap", pasaran.DasarKeigenap},
		{"dasar_keiganjil", pasaran.DasarKeiganjil},
		{"dasar_discbesar", pasaran.DasarDiscbesar},
		{"dasar_disckecil", pasaran.DasarDisckecil},
		{"dasar_discigenap", pasaran.DasarDiscigenap},
		{"dasar_discganjil", pasaran.DasarDiscganjil},
		{"dasar_limitbuang", pasaran.DasarLimitbuang},
		{"dasar_limittotal", pasaran.DasarLimittotal},
		{"shio_referal", pasaran.ShioReferal},
		{"shio_shiotahunini", pasaran.ShioShiotahunini},
		{"shio_minbet", pasaran.ShioMinbet},
		{"shio_maxbet", pasaran.ShioMaxbet},
		{"shio_maxbuy", pasaran.ShioMaxbuy},
		{"shio_win", pasaran.ShioWin},
		{"shio_disc", pasaran.ShioDisc},
		{"shio_limitbuang", pasaran.ShioLimitbuang},
		{"shio_limittotal", pasaran.ShioLimittotal},
		{"statuscompasaran", pasaran.Status},
		{"update_by", pasaran.UpdateBy},
		{"update_at", pasaran.UpdateAt},
	}

	setClauses := make([]string, 0, len(cols))
	args := make([]any, 0, len(cols)+1)
	for i, cl := range cols {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", cl.name, i+1))
		args = append(args, cl.val)
	}

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE idcomppasaran = $%d AND idcompany= $%d`,
		config.DB_mst_company_pasaran,
		strings.Join(setClauses, ", "),
		len(cols)+1,
	)
	args = append(args, pasaran.IDcomppasaran, pasaran.IDcompany)

	res, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func (u *pasaranRepository) Savejadwal(ctx context.Context, compjadwal *domain.Pasaranjadwal) error {
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
func (u *pasaranRepository) DeleteJadwal(ctx context.Context, idcomppasaran string) error {
	query := `DELETE FROM ` + config.DB_mst_company_jadwaltogel + ` WHERE idcomppasaran = $1`
	_, err := u.db.Exec(ctx, query, idcomppasaran)
	return err
}
