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

func (u *pasaranRepository) Update(ctx context.Context, comppasaran *domain.Pasaran) error {
	type col struct {
		name string
		val  any
	}

	cols := []col{
		{"aliascomppasaran", comppasaran.Aliascomppasaran},
		{"urlpasaran", comppasaran.URLpasaran},
		{"urllogo", comppasaran.URLlogo},
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
