CREATE TABLE public.tbl_counter (
	idcounter int4 DEFAULT nextval('idcounter_seq'::regclass) NOT NULL,
	nmcounter varchar(70) NULL,
	counter int8 NOT NULL,
	CONSTRAINT tbl_counter_pk PRIMARY KEY (idcounter)
);
CREATE SEQUENCE public.idcounter_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;

CREATE TABLE tbl_admin (
	username varchar(30) NOT NULL,
	"password" varchar(250) NULL,
	idadmin varchar(30) NULL,
	"name" varchar(50) NULL,
	statuslogin varchar(1) NULL,
	lastlogin timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	joindate date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	ipaddress varchar(20) DEFAULT ''::character varying NULL,
	timezone varchar(30) DEFAULT ''::character varying NULL,
	createadmin varchar(30) NULL,
	createdateadmin timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updateadmin varchar(30) DEFAULT ''::character varying NULL,
	updatedateadmin timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_admin_pk PRIMARY KEY (username)
);




CREATE TABLE public.tbl_adminrole (
	idadminrole varchar(30) NOT NULL,
	nmadminrole varchar(50) NULL,
	ruleadmin text NULL,
	createadminrole varchar(30) DEFAULT ''::character varying NULL,
	createdateadminrole timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updateadminrole varchar(30) DEFAULT ''::character varying NULL,
	updatedateadminrole timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_adminrole_unique UNIQUE (idadminrole)
);



CREATE TABLE public.tbl_bank (
	idbank varchar(10) NOT NULL,
	typebank varchar(20) NOT NULL,
	nmbank varchar(50) NULL,
	bankstatus varchar(1) DEFAULT 'Y'::character varying NULL,
	createbank varchar(30) DEFAULT ''::character varying NULL,
	createdatebank timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updatebank varchar(30) DEFAULT ''::character varying NULL,
	updatedatebank timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_tbl_bank_unique UNIQUE (idbank)
);




CREATE TABLE public.tbl_clientrule (
	idclientrule varchar(30) NOT NULL,
	nmclientrule varchar(50) NULL,
	ruleclient text NULL,
	createclientrule varchar(30) DEFAULT ''::character varying NULL,
	createdateclientrule timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updateclientrule varchar(30) DEFAULT ''::character varying NULL,
	updatedateclientrule timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_clientrule_unique UNIQUE (idclientrule)
);


CREATE TABLE tbl_groupcompany (
	idgroupcomp varchar(20) NOT NULL,
	nmgroupcomp varchar(100) DEFAULT ''::character varying NOT NULL,
	statusgroupcomp varchar(1) DEFAULT 'Y'::character varying NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_groupcompany_unique UNIQUE (idgroupcomp)
);


CREATE TABLE public.tbl_company (
	idcompany varchar(10) NOT NULL,
	idcurrdef varchar(20) NOT NULL,
	compname varchar(50) NULL,
	endjoin timestamp NULL,
	amountcomp numeric(36, 18) DEFAULT 0 NOT NULL,
	compstatus varchar(1) DEFAULT 'Y'::character varying NULL,
	idgroupcomp varchar(20) DEFAULT ''::character varying NULL,
	telegramid varchar(50) DEFAULT ''::character varying NULL,
	urlapitoto varchar(150) DEFAULT ''::character varying NULL,
	urlapislot varchar(150) DEFAULT ''::character varying NULL,
	compactivetoto varchar DEFAULT 'N'::character varying NULL,
	compactiveslot varchar DEFAULT 'N'::character varying NULL,
	createcomp varchar(30) DEFAULT ''::character varying NULL,
	createdatecomp timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updatecomp varchar(30) DEFAULT ''::character varying NULL,
	updatedatecomp timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_company_unique UNIQUE (idcompany)
);

CREATE TABLE public.tbl_company_admin (
	idcompadmin varchar(64) NOT NULL,
	idcompany varchar(10) NOT NULL,
	idclientrule varchar(30) NOT NULL,
	usernamecompadmin varchar(30) NOT NULL,
	namecompadmin varchar(50) NULL,
	passcompadmin varchar(250) NULL,
	lastlogincompadmin timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	ipaddresscompadmin varchar(20) DEFAULT ''::character varying NULL,
	compadminstatus varchar(1) DEFAULT 'Y'::character varying NULL,
	createcompadmin varchar(30) DEFAULT ''::character varying NULL,
	createdatecompadmin timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updatecompadmin varchar(30) DEFAULT ''::character varying NULL,
	updatedatecompadmin timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_company_admin_unique UNIQUE (idcompadmin)
);

CREATE TABLE tbl_company_conf_toto (
	idcompconftoto varchar(80) NOT NULL,
	idcompany varchar(10) NOT NULL,
	angka_max_minbasket NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_4d NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_3d NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_3dd NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_2d NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_2dd NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_2dt NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_4d_bbdisc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_3d_bbdisc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_3dd_bbdisc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_2d_bbdisc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_2dd_bbdisc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_maxbet_2dt_bbdisc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win4d_full NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win3d_full NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win3dd_full NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2d_full NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2dd_full NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2dt_full NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win4d_disc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win3d_disc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win3dd_disc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2d_disc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2dd_disc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2dt_disc NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win4d_bb NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win3d_bb NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win3dd_bb NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2d_bb NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2dd_bb NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2dt_bb NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win4d_bb_kena NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win3d_bb_kena NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win3dd_bb_kena NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2d_bb_kena NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2dd_bb_kena NUMERIC(15,2) DEFAULT 0 NULL,
	angka_max_win2dt_bb_kena NUMERIC(15,2) DEFAULT 0 NULL,
	cbebas_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	cbebas_max_maxbet NUMERIC(15,2) DEFAULT 0 NULL,
	cbebas_max_win NUMERIC(5,2) DEFAULT 0 NULL,
	cmacau_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	cmacau_max_maxbet NUMERIC(15,2) DEFAULT 0 NULL,
	cmacau_max_win2 NUMERIC(5,2) DEFAULT 0 NULL,
	cmacau_max_win3 NUMERIC(5,2) DEFAULT 0 NULL,
	cmacau_max_win4 NUMERIC(5,2) DEFAULT 0 NULL,
	cnaga_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	cnaga_max_maxbet NUMERIC(15,2) DEFAULT 0 NULL,
	cnaga_max_win3 NUMERIC(5,2) DEFAULT 0 NULL,
	cnaga_max_win4 NUMERIC(5,2) DEFAULT 0 NULL,
	cjitu_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	cjitu_max_maxbet NUMERIC(15,2) DEFAULT 0 NULL,
	cjitu_max_winas NUMERIC(5,2) DEFAULT 0 NULL,
	cjitu_max_winkop NUMERIC(5,2) DEFAULT 0 NULL,
	cjitu_max_winkepala NUMERIC(5,2) DEFAULT 0 NULL,
	cjitu_max_winekor NUMERIC(5,2) DEFAULT 0 NULL,
	umum50_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	umum50_max_maxbet NUMERIC(15,2) DEFAULT 0 NULL,
	special50_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	special50_max_maxbet NUMERIC(15,2) DEFAULT 0 NULL,
	kombinasi50_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	kombinasi50_max_maxbet NUMERIC(15,2) DEFAULT 0 NULL,
	macau_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	macau_max_maxbet NUMERIC(15,2) DEFAULT 0 NULL,
	macau_max_win NUMERIC(5,2) DEFAULT 0 NULL,
	dasar_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	dasar_max_maxbet NUMERIC(15,2) DEFAULT 0 NULL,
	shio_max_minbet NUMERIC(15,2) DEFAULT 0 NULL,
	shio_max_maxbet NUMERIC(15,2) DEFAULT 0 NULL,
	shio_max_win NUMERIC(5,2) DEFAULT 0 NULL,
	shio_parent int4 DEFAULT 0 NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_company_conf_toto_unique UNIQUE (idcompconftoto)
);

CREATE TABLE public.tbl_company_wallet (
	idcompwallet varchar(64) NOT NULL,
	idcompany varchar(10) NOT NULL,
	idcurr varchar(20) NOT NULL,
	amountcompwallet numeric(36, 18) DEFAULT 0 NOT NULL,
	compwalletstatus varchar(1) DEFAULT 'Y'::character varying NULL,
	createcompwallet varchar(30) DEFAULT ''::character varying NULL,
	createdatecompwallet timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updatecompwallet varchar(30) DEFAULT ''::character varying NULL,
	updatedatecompwallet timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_company_wallet_unique UNIQUE (idcompwallet)
);

CREATE TABLE db_tot.tbl_trx_log (
	idlog int8 NOT NULL,
	datetimelog timestamptz NOT NULL,
	yearlog int4 NOT NULL,
	idcompany varchar(50) NULL,
	username varchar(100) NULL,
	pagelog varchar(50) NULL,
	tipelog varchar(50) NULL,
	notebefore text NULL,
	noteafter text NULL
);
-- db_tot.tbl_mst_pasaran_togel definition

-- Drop table

-- DROP TABLE db_tot.tbl_mst_pasaran_togel;

CREATE TABLE db_tot.tbl_mst_pasaran_togel (
	idpasarantogel bpchar(10) NOT NULL,
	nmpasarantogel varchar(70) NULL,
	tipepasaran varchar(10) NULL,
	urlpasaran varchar(350) NULL,
	pasarandiundi varchar(150) NULL,
	pasaranlibur varchar(150) NULL,
	jamtutup time NOT NULL,
	jamjadwal time NOT NULL,
	jamopen time NOT NULL,
	angka_referal float8 DEFAULT 0 NOT NULL,
	angka_minbet NUMERIC(15,2) DEFAULT 0 NOT NULL,
	angka_maxbet4d NUMERIC(15,2) DEFAULT 0 NOT NULL,
	angka_maxbet3d NUMERIC(15,2) DEFAULT 0 NOT NULL,
	angka_maxbet3dd NUMERIC(15,2) DEFAULT 0 NOT NULL,
	angka_maxbet2d NUMERIC(15,2) DEFAULT 0 NOT NULL,
	angka_maxbet2dd NUMERIC(15,2) DEFAULT 0 NOT NULL,
	angka_maxbet2dt NUMERIC(15,2) DEFAULT 0 NOT NULL,
	angka_maxbuy4d float8 DEFAULT 0 NOT NULL,
	angka_maxbuy3d float8 DEFAULT 0 NOT NULL,
	angka_maxbuy3dd float8 DEFAULT 0 NOT NULL,
	angka_maxbuy2d float8 DEFAULT 0 NOT NULL,
	angka_maxbuy2dd float8 DEFAULT 0 NOT NULL,
	angka_maxbuy2dt float8 DEFAULT 0 NOT NULL,
	angka_maxbet4d_fullbb float8 DEFAULT 0 NOT NULL,
	angka_maxbet3d_fullbb float8 DEFAULT 0 NOT NULL,
	angka_maxbet3dd_fullbb float8 DEFAULT 0 NOT NULL,
	angka_maxbet2d_fullbb float8 DEFAULT 0 NOT NULL,
	angka_maxbet2dd_fullbb float8 DEFAULT 0 NOT NULL,
	angka_maxbet2dt_fullbb float8 DEFAULT 0 NOT NULL,
	angka_win4d float8 DEFAULT 0 NOT NULL,
	angka_win3d float8 DEFAULT 0 NOT NULL,
	angka_win3dd float8 DEFAULT 0 NOT NULL,
	angka_win2d float8 DEFAULT 0 NOT NULL,
	angka_win2dd float8 DEFAULT 0 NOT NULL,
	angka_win2dt float8 DEFAULT 0 NOT NULL,
	angka_win4dnodisc float8 DEFAULT 0 NOT NULL,
	angka_win3dnodisc float8 DEFAULT 0 NOT NULL,
	angka_win3ddnodisc float8 DEFAULT 0 NOT NULL,
	angka_win2dnodisc float8 DEFAULT 0 NOT NULL,
	angka_win2ddnodisc float8 DEFAULT 0 NOT NULL,
	angka_win2dtnodisc float8 DEFAULT 0 NOT NULL,
	angka_win4dbb_kena float8 DEFAULT 0 NOT NULL,
	angka_win3dbb_kena float8 DEFAULT 0 NOT NULL,
	angka_win3ddbb_kena float8 DEFAULT 0 NOT NULL,
	angka_win2dbb_kena float8 DEFAULT 0 NOT NULL,
	angka_win2ddbb_kena float8 DEFAULT 0 NOT NULL,
	angka_win2dtbb_kena float8 DEFAULT 0 NOT NULL,
	angka_win4dbb float8 DEFAULT 0 NOT NULL,
	angka_win3dbb float8 DEFAULT 0 NOT NULL,
	angka_win3ddbb float8 DEFAULT 0 NOT NULL,
	angka_win2dbb float8 DEFAULT 0 NOT NULL,
	angka_win2ddbb float8 DEFAULT 0 NOT NULL,
	angka_win2dtbb float8 DEFAULT 0 NOT NULL,
	angka_disc4d float8 DEFAULT 0 NOT NULL,
	angka_disc3d float8 DEFAULT 0 NOT NULL,
	angka_disc3dd float8 DEFAULT 0 NOT NULL,
	angka_disc2d float8 DEFAULT 0 NOT NULL,
	angka_disc2dd float8 DEFAULT 0 NOT NULL,
	angka_disc2dt float8 DEFAULT 0 NOT NULL,
	angka_limitbuang4d float8 DEFAULT 0 NOT NULL,
	angka_limitbuang3d float8 DEFAULT 0 NOT NULL,
	angka_limitbuang3dd float8 DEFAULT 0 NOT NULL,
	angka_limitbuang2d float8 DEFAULT 0 NOT NULL,
	angka_limitbuang2dd float8 DEFAULT 0 NOT NULL,
	angka_limitbuang2dt float8 DEFAULT 0 NOT NULL,
	angka_limittotal4d float8 DEFAULT 0 NOT NULL,
	angka_limittotal3d float8 DEFAULT 0 NOT NULL,
	angka_limittotal3dd float8 DEFAULT 0 NOT NULL,
	angka_limittotal2d float8 DEFAULT 0 NOT NULL,
	angka_limittotal2dd float8 DEFAULT 0 NOT NULL,
	angka_limittotal2dt float8 DEFAULT 0 NOT NULL,
	angka_limitline_4d int4 DEFAULT 0 NOT NULL,
	angka_limitline_3d int4 DEFAULT 0 NOT NULL,
	angka_limitline_3dd int4 DEFAULT 0 NOT NULL,
	angka_limitline_2d int4 DEFAULT 0 NOT NULL,
	angka_limitline_2dd int4 DEFAULT 0 NOT NULL,
	angka_limitline_2dt int4 DEFAULT 0 NOT NULL,
	angka_limitbuang4d_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limitbuang3d_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limitbuang3dd_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limitbuang2d_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limitbuang2dd_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limitbuang2dt_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limittotal4d_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limittotal3d_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limittotal3dd_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limittotal2d_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limittotal2dd_fullbb float8 DEFAULT 0 NOT NULL,
	angka_limittotal2dt_fullbb float8 DEFAULT 0 NOT NULL,
	angka_bbfs int4 DEFAULT 6 NOT NULL,
	"2_minbet" float8 DEFAULT 0 NOT NULL,
	"2_referal" float8 DEFAULT 0 NOT NULL,
	"2_maxbet" float8 DEFAULT 0 NOT NULL,
	"2_win" float8 DEFAULT 0 NOT NULL,
	"2_disc" float8 DEFAULT 0 NOT NULL,
	"2_limitbuang" float8 DEFAULT 0 NOT NULL,
	"2_limitotal" float8 DEFAULT 0 NOT NULL,
	"3_referal" float8 DEFAULT 0 NOT NULL,
	"3_minbet" float8 DEFAULT 0 NOT NULL,
	"3_maxbet" float8 DEFAULT 0 NOT NULL,
	"3_win2digit" float8 DEFAULT 0 NOT NULL,
	"3_win3digit" float8 DEFAULT 0 NOT NULL,
	"3_win4digit" float8 DEFAULT 0 NOT NULL,
	"3_disc" float8 DEFAULT 0 NOT NULL,
	"3_limitbuang" float8 DEFAULT 0 NOT NULL,
	"3_limittotal" float8 DEFAULT 0 NOT NULL,
	"4_referal" float8 DEFAULT 0 NOT NULL,
	"4_minbet" float8 DEFAULT 0 NOT NULL,
	"4_maxbet" float8 DEFAULT 0 NOT NULL,
	"4_win3digit" float8 DEFAULT 0 NOT NULL,
	"4_win4digit" float8 DEFAULT 0 NOT NULL,
	"4_disc" float8 DEFAULT 0 NOT NULL,
	"4_limitbuang" float8 DEFAULT 0 NOT NULL,
	"4_limittotal" float8 DEFAULT 0 NOT NULL,
	"5_referal" float8 DEFAULT 0 NOT NULL,
	"5_minbet" float8 DEFAULT 0 NOT NULL,
	"5_maxbet" float8 DEFAULT 0 NOT NULL,
	"5_winas" float8 DEFAULT 0 NOT NULL,
	"5_winkop" float8 DEFAULT 0 NOT NULL,
	"5_winkepala" float8 DEFAULT 0 NOT NULL,
	"5_winekor" float8 DEFAULT 0 NOT NULL,
	"5_desic" float8 DEFAULT 0 NOT NULL,
	"5_limitbuang" float8 DEFAULT 0 NOT NULL,
	"5_limitotal" float8 DEFAULT 0 NOT NULL,
	"6_referal" float8 DEFAULT 0 NOT NULL,
	"6_minbet" float8 DEFAULT 0 NOT NULL,
	"6_maxbet" float8 DEFAULT 0 NOT NULL,
	"6_keibesar" float8 DEFAULT 0 NOT NULL,
	"6_keikecil" float8 DEFAULT 0 NOT NULL,
	"6_keigenap" float8 DEFAULT 0 NOT NULL,
	"6_keiganjil" float8 DEFAULT 0 NOT NULL,
	"6_keitengah" float8 DEFAULT 0 NOT NULL,
	"6_keitepi" float8 DEFAULT 0 NOT NULL,
	"6_discbesar" float8 DEFAULT 0 NOT NULL,
	"6_disckecil" float8 DEFAULT 0 NOT NULL,
	"6_discgenap" float8 DEFAULT 0 NOT NULL,
	"6_discganjil" float8 DEFAULT 0 NOT NULL,
	"6_disctengah" float8 DEFAULT 0 NOT NULL,
	"6_disctepi" float8 DEFAULT 0 NOT NULL,
	"6_limitbuang" float8 DEFAULT 0 NOT NULL,
	"6_limittotal" float8 DEFAULT 0 NOT NULL,
	"7_referal" float8 DEFAULT 0 NOT NULL,
	"7_minbet" float8 DEFAULT 0 NOT NULL,
	"7_maxbet" float8 DEFAULT 0 NOT NULL,
	"7_keiasganjil" float8 DEFAULT 0 NOT NULL,
	"7_keiasgenap" float8 DEFAULT 0 NOT NULL,
	"7_keiasbesar" float8 DEFAULT 0 NOT NULL,
	"7_keiaskecil" float8 DEFAULT 0 NOT NULL,
	"7_keikopganjil" float8 DEFAULT 0 NOT NULL,
	"7_keikopgenap" float8 DEFAULT 0 NOT NULL,
	"7_keikopbesar" float8 DEFAULT 0 NOT NULL,
	"7_keikopkecil" float8 DEFAULT 0 NOT NULL,
	"7_keikepalaganjil" float8 DEFAULT 0 NOT NULL,
	"7_keikepalagenap" float8 DEFAULT 0 NOT NULL,
	"7_keikepalabesar" float8 DEFAULT 0 NOT NULL,
	"7_keikepalakecil" float8 DEFAULT 0 NOT NULL,
	"7_keiekorganjil" float8 DEFAULT 0 NOT NULL,
	"7_keiekorgenap" float8 DEFAULT 0 NOT NULL,
	"7_keiekorbesar" float8 DEFAULT 0 NOT NULL,
	"7_keiekorkecil" float8 DEFAULT 0 NOT NULL,
	"7_discasganjil" float8 DEFAULT 0 NOT NULL,
	"7_discasgenap" float8 DEFAULT 0 NOT NULL,
	"7_discasbesar" float8 DEFAULT 0 NOT NULL,
	"7_discaskecil" float8 DEFAULT 0 NOT NULL,
	"7_disckopganjil" float8 DEFAULT 0 NOT NULL,
	"7_disckopgenap" float8 DEFAULT 0 NOT NULL,
	"7_disckopbesar" float8 DEFAULT 0 NOT NULL,
	"7_disckopkecil" float8 DEFAULT 0 NOT NULL,
	"7_disckepalaganjil" float8 DEFAULT 0 NOT NULL,
	"7_disckepalagenap" float8 DEFAULT 0 NOT NULL,
	"7_disckepalabesar" float8 DEFAULT 0 NOT NULL,
	"7_disckepalakecil" float8 DEFAULT 0 NOT NULL,
	"7_discekorganjil" float8 DEFAULT 0 NOT NULL,
	"7_discekorgenap" float8 DEFAULT 0 NOT NULL,
	"7_discekorbesar" float8 DEFAULT 0 NOT NULL,
	"7_discekorkecil" float8 DEFAULT 0 NOT NULL,
	"7_limitbuang" float8 DEFAULT 0 NOT NULL,
	"7_limittotal" float8 DEFAULT 0 NOT NULL,
	"8_referal" float8 DEFAULT 0 NOT NULL,
	"8_minbet" float8 DEFAULT 0 NOT NULL,
	"8_maxbet" float8 DEFAULT 0 NOT NULL,
	"8_belakangkeimono" float8 DEFAULT 0 NOT NULL,
	"8_belakangkeistereo" float8 DEFAULT 0 NOT NULL,
	"8_belakangkeikembang" float8 DEFAULT 0 NOT NULL,
	"8_belakangkeikempis" float8 DEFAULT 0 NOT NULL,
	"8_belakangkeikembar" float8 DEFAULT 0 NOT NULL,
	"8_tengahkeimono" float8 DEFAULT 0 NOT NULL,
	"8_tengahkeistereo" float8 DEFAULT 0 NOT NULL,
	"8_tengahkeikembang" float8 DEFAULT 0 NOT NULL,
	"8_tengahkeikempis" float8 DEFAULT 0 NOT NULL,
	"8_tengahkeikembar" float8 DEFAULT 0 NOT NULL,
	"8_depankeimono" float8 DEFAULT 0 NOT NULL,
	"8_depankeistereo" float8 DEFAULT 0 NOT NULL,
	"8_depankeikembang" float8 DEFAULT 0 NOT NULL,
	"8_depankeikempis" float8 DEFAULT 0 NOT NULL,
	"8_depankeikembar" float8 DEFAULT 0 NOT NULL,
	"8_belakangdiscmono" float8 DEFAULT 0 NOT NULL,
	"8_belakangdiscstereo" float8 DEFAULT 0 NOT NULL,
	"8_belakangdisckembang" float8 DEFAULT 0 NOT NULL,
	"8_belakangdisckempis" float8 DEFAULT 0 NOT NULL,
	"8_belakangdisckembar" float8 DEFAULT 0 NOT NULL,
	"8_tengahdiscmono" float8 DEFAULT 0 NOT NULL,
	"8_tengahdiscstereo" float8 DEFAULT 0 NOT NULL,
	"8_tengahdisckembang" float8 DEFAULT 0 NOT NULL,
	"8_tengahdisckempis" float8 DEFAULT 0 NOT NULL,
	"8_tengahdisckembar" float8 DEFAULT 0 NOT NULL,
	"8_depandiscmono" float8 DEFAULT 0 NOT NULL,
	"8_depandiscstereo" float8 DEFAULT 0 NOT NULL,
	"8_depandisckembang" float8 DEFAULT 0 NOT NULL,
	"8_depandisckempis" float8 DEFAULT 0 NOT NULL,
	"8_depandisckembar" float8 DEFAULT 0 NOT NULL,
	"8_limitbuang" float8 DEFAULT 0 NOT NULL,
	"8_limittotal" float8 DEFAULT 0 NOT NULL,
	"9_referal" float8 DEFAULT 0 NOT NULL,
	"9_minbet" float8 DEFAULT 0 NOT NULL,
	"9_maxbet" float8 DEFAULT 0 NOT NULL,
	"9_win" float8 DEFAULT 0 NOT NULL,
	"9_discount" float8 DEFAULT 0 NOT NULL,
	"9_limitbuang" float8 DEFAULT 0 NOT NULL,
	"9_limittotal" float8 DEFAULT 0 NOT NULL,
	"10_referal" float8 DEFAULT 0 NOT NULL,
	"10_minbet" float8 DEFAULT 0 NOT NULL,
	"10_maxbet" float8 DEFAULT 0 NOT NULL,
	"10_keibesar" float8 DEFAULT 0 NOT NULL,
	"10_keikecil" float8 DEFAULT 0 NOT NULL,
	"10_keigenap" float8 DEFAULT 0 NOT NULL,
	"10_keiganjil" float8 DEFAULT 0 NOT NULL,
	"10_discbesar" float8 DEFAULT 0 NOT NULL,
	"10_disckecil" float8 DEFAULT 0 NOT NULL,
	"10_discigenap" float8 DEFAULT 0 NOT NULL,
	"10_discganjil" float8 DEFAULT 0 NOT NULL,
	"10_limitbuang" float8 DEFAULT 0 NOT NULL,
	"10_limittotal" float8 DEFAULT 0 NOT NULL,
	"shio_referal" float8 DEFAULT 0 NOT NULL,
	"shio_shiotahunini" varchar DEFAULT ''::character varying NULL,
	"shio_minbet" float8 DEFAULT 0 NOT NULL,
	"shio_maxbet" float8 DEFAULT 0 NOT NULL,
	"11_win" float8 DEFAULT 0 NOT NULL,
	"11_disc" float8 DEFAULT 0 NOT NULL,
	"11_limitbuang" float8 DEFAULT 0 NOT NULL,
	"11_limittotal" float8 DEFAULT 0 NOT NULL,
	"12_master4d" float8 DEFAULT 0 NOT NULL,
	"12_agent4d" float8 DEFAULT 0 NOT NULL,
	"12_master3d" float8 DEFAULT 0 NOT NULL,
	"12_agent3d" float8 DEFAULT 0 NOT NULL,
	"12_master2d" float8 DEFAULT 0 NOT NULL,
	"12_agent2d" float8 DEFAULT 0 NOT NULL,
	"12_master2dd" float8 DEFAULT 0 NOT NULL,
	"12_agent2dd" float8 DEFAULT 0 NOT NULL,
	"12_master2dt" float8 DEFAULT 0 NOT NULL,
	"12_agent2dt" float8 DEFAULT 0 NOT NULL,
	"12_mastercb" float8 DEFAULT 0 NOT NULL,
	"12_agentcb" float8 DEFAULT 0 NOT NULL,
	"12_mastercm" float8 DEFAULT 0 NOT NULL,
	"12_agentcm" float8 DEFAULT 0 NOT NULL,
	"12_mastercn" float8 DEFAULT 0 NOT NULL,
	"12_agentcn" float8 DEFAULT 0 NOT NULL,
	"12_mastercj" float8 DEFAULT 0 NOT NULL,
	"12_agentcj" float8 DEFAULT 0 NOT NULL,
	"12_master5050umum" float8 DEFAULT 0 NOT NULL,
	"12_agent5050umum" float8 DEFAULT 0 NOT NULL,
	"12_master5050special" float8 DEFAULT 0 NOT NULL,
	"12_agent5050special" float8 DEFAULT 0 NOT NULL,
	"12_master5050kombinasi" float8 DEFAULT 0 NOT NULL,
	"12_agent5050kombinasi" float8 DEFAULT 0 NOT NULL,
	"12_masterkombinasi" float8 DEFAULT 0 NOT NULL,
	"12_agentkombinasi" float8 DEFAULT 0 NOT NULL,
	"12_masterdasar" float8 DEFAULT 0 NOT NULL,
	"12_agentdasar" float8 DEFAULT 0 NOT NULL,
	"12_mastershio" float8 DEFAULT 0 NOT NULL,
	"12_agentshio" float8 DEFAULT 0 NOT NULL,
	
	created_by varchar(30) DEFAULT ''::character varying NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_by varchar(30) DEFAULT ''::character varying NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	
	
	"2_maxbuy" float8 DEFAULT 0 NOT NULL,
	"3_maxbuy" float8 DEFAULT 0 NOT NULL,
	"4_maxbuy" float8 DEFAULT 0 NOT NULL,
	"5_maxbuy" float8 DEFAULT 0 NOT NULL,
	"6_maxbuy" float8 DEFAULT 0 NOT NULL,
	"7_maxbuy" float8 DEFAULT 0 NOT NULL,
	"8_maxbuy" float8 DEFAULT 0 NOT NULL,
	"9_maxbuy" float8 DEFAULT 0 NOT NULL,
	"10_maxbuy" float8 DEFAULT 0 NOT NULL,
	"11_maxbuy" float8 DEFAULT 0 NOT NULL,
	"1_maxbet4d_full" float8 DEFAULT 0 NOT NULL,
	"1_maxbet4d_bb" float8 DEFAULT 0 NOT NULL,
	"1_maxbet3d_full" float8 DEFAULT 0 NOT NULL,
	"1_maxbet3d_bb" float8 DEFAULT 0 NOT NULL,
	"1_maxbet3dd_full" float8 DEFAULT 0 NOT NULL,
	"1_maxbet3dd_bb" float8 DEFAULT 0 NOT NULL,
	"1_maxbet2d_full" float8 DEFAULT 0 NOT NULL,
	"1_maxbet2d_bb" float8 DEFAULT 0 NOT NULL,
	"1_maxbet2dd_full" float8 DEFAULT 0 NOT NULL,
	"1_maxbet2dd_bb" float8 DEFAULT 0 NOT NULL,
	"1_maxbet2dt_full" float8 DEFAULT 0 NOT NULL,
	"1_maxbet2dt_bb" float8 DEFAULT 0 NOT NULL,
	"1_maxbet4d_bbdisc" float8 DEFAULT 0 NOT NULL,
	"1_maxbet3d_bbdisc" float8 DEFAULT 0 NOT NULL,
	"1_maxbet3dd_bbdisc" float8 DEFAULT 0 NOT NULL,
	"1_maxbet2d_bbdisc" float8 DEFAULT 0 NOT NULL,
	"1_maxbet2dd_bbdisc" float8 DEFAULT 0 NOT NULL,
	"1_maxbet2dt_bbdisc" float8 DEFAULT 0 NOT NULL,
	"1_minbasket" float8 DEFAULT 1000 NULL,
	CONSTRAINT tbl_mst_pasaran_togel_pk PRIMARY KEY (idpasarantogel)
);
CREATE TABLE public.tbl_currency (
	idcurr varchar(20) NOT NULL,
	typecurr varchar(10) DEFAULT ''::character varying NOT NULL,
	statuscurr varchar(1) NOT NULL,
	createcurr varchar(30) DEFAULT ''::character varying NULL,
	createdatecurr timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updatecurr varchar(30) DEFAULT ''::character varying NULL,
	updatedatecurr timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_currency_unique UNIQUE (idcurr)
);



CREATE TABLE public.tbl_mst_uom (
	iduom varchar(10) NOT NULL,
	nmuom varchar(100) DEFAULT ''::character varying NOT NULL,
	status varchar(1) NOT NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_mst_uom_unique UNIQUE (iduom)
);



CREATE TABLE db_bbca.tbl_account_balance_log (
	idaccbalancelog varchar(150) NOT NULL,
	ref_idtrx varchar(150) NOT NULL,
	ref_table varchar(150) NOT NULL,
	typeaccbalancelog varchar(10) NOT NULL,
	dateaccbalancelog date NOT NULL,
	idcurr varchar(20) NOT NULL,
	amount_credit numeric(36, 18) DEFAULT 0 NOT NULL,
	amount_debit numeric(36, 18) DEFAULT 0 NOT NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	idwallet varchar(150) DEFAULT ''::character varying NULL,
	CONSTRAINT tbl_ledger_company_unique UNIQUE (idaccbalancelog)
);
CREATE INDEX tbl_ledger_company_idcompwalletbank_idx ON db_bbca.tbl_account_balance_log USING btree (idwallet);




CREATE TABLE db_bbca.tbl_mst_grouptrx (
	idgroup varchar(4) NOT NULL,
	nmgroup varchar(100) DEFAULT ''::character varying NOT NULL,
	status varchar(1) NOT NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_mst_grouptrx_unique UNIQUE (idgroup)
);




CREATE TABLE db_bbca.tbl_mst_gudang (
	idgudang varchar(10) NOT NULL,
	nmgudang varchar(100) DEFAULT ''::character varying NOT NULL,
	status varchar(1) NOT NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_mst_gudang_unique UNIQUE (idgudang)
);



CREATE TABLE db_bbca.tbl_mst_item (
	iditem varchar(20) NOT NULL,
	iditemcategory int4 NOT NULL,
	item_type varchar(20) NOT NULL,
	nmitem varchar(100) DEFAULT ''::character varying NOT NULL,
	description text DEFAULT ''::text NOT NULL,
	status varchar(1) NOT NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz NULL,
	CONSTRAINT tbl_mst_item_item_type_check CHECK (((item_type)::text = ANY ((ARRAY['STOCK'::character varying, 'NON_STOCK'::character varying, 'SERVICE'::character varying])::text[]))),
	CONSTRAINT tbl_mst_item_pkey PRIMARY KEY (iditem),
	CONSTRAINT tbl_mst_item_status_check CHECK (((status)::text = ANY ((ARRAY['Y'::character varying, 'N'::character varying])::text[])))
);




CREATE TABLE db_bbca.tbl_mst_member (
	idmember varchar(150) NOT NULL,
	usernamemember varchar(30) NOT NULL,
	passmember varchar(250) NULL,
	namemember varchar(50) NULL,
	hpmember varchar(30) DEFAULT ''::character varying NULL,
	emailmember varchar(100) DEFAULT ''::character varying NULL,
	lastloginmember timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	ipaddressmember varchar(20) DEFAULT ''::character varying NULL,
	statusmember varchar(1) DEFAULT 'Y'::character varying NULL,
	createmember varchar(30) DEFAULT ''::character varying NULL,
	createdatemember timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updatemember varchar(30) DEFAULT ''::character varying NULL,
	updatedatemember timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_member_unique UNIQUE (idmember)
);




CREATE TABLE db_bbca.tbl_mst_merek (
	idmerek varchar(10) NOT NULL,
	nmmerek varchar(100) DEFAULT ''::character varying NOT NULL,
	status varchar(1) NOT NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_mst_merek_unique UNIQUE (idmerek)
);


-- db_bbca.tbl_mst_supplier definition

-- Drop table

-- DROP TABLE db_bbca.tbl_mst_supplier;

CREATE TABLE db_bbca.tbl_mst_supplier (
	idsupplier varchar(20) NOT NULL,
	nmsupplier varchar(100) DEFAULT ''::character varying NOT NULL,
	hp1 varchar(25) DEFAULT ''::character varying NOT NULL,
	hp2 varchar(25) DEFAULT ''::character varying NOT NULL,
	email varchar(150) DEFAULT ''::character varying NOT NULL,
	tempo_pembayaran varchar(50) DEFAULT ''::character varying NOT NULL,
	tipe_transaksi varchar(20) DEFAULT ''::character varying NOT NULL,
	idbank varchar(10) NOT NULL,
	norek varchar(50) DEFAULT ''::character varying NOT NULL,
	nmrek varchar(100) DEFAULT ''::character varying NOT NULL,
	alamat text DEFAULT ''::text NOT NULL,
	status varchar(1) NOT NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_mst_supplier_unique UNIQUE (idsupplier)
);




CREATE TABLE db_bbca.tbl_wallet (
	idwallet varchar(150) NOT NULL,
	wallettype varchar(20) DEFAULT ''::character varying NULL,
	idowner varchar(150) DEFAULT ''::character varying NULL,
	idcurr varchar(20) NOT NULL,
	idbank varchar(10) NULL,
	account_number varchar(150) DEFAULT ''::character varying NULL,
	account_name varchar(150) DEFAULT ''::character varying NULL,
	networkcrypto varchar(20) DEFAULT ''::character varying NULL,
	status varchar(1) DEFAULT 'Y'::character varying NULL,
	created_by varchar(30) DEFAULT ''::character varying NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_by varchar(30) DEFAULT ''::character varying NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_wallet_unique UNIQUE (idwallet)
);




CREATE TABLE db_bbca.tbl_wallet_balance (
	idwallet varchar(150) NOT NULL,
	total_credit numeric(36, 18) DEFAULT 0 NOT NULL,
	total_debit numeric(36, 18) DEFAULT 0 NOT NULL,
	updated_by varchar(30) NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_wallet_balance_unique UNIQUE (idwallet)
);




CREATE TABLE db_bbca.tbl_wallet_company_metadata (
	idwallet varchar(150) NOT NULL,
	category varchar(20) DEFAULT ''::character varying NULL,
	note varchar(150) DEFAULT ''::character varying NULL,
	effective_from timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	effective_to timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_company_wallet_metadata_unique UNIQUE (idwallet)
);




CREATE TABLE db_bbca.tbl_wallet_external_transaction (
	idtrx varchar(150) NOT NULL,
	idmember varchar(150) DEFAULT ''::character varying NULL,
	typeact varchar(20) NOT NULL,
	typetrx varchar(20) NOT NULL,
	idcurr varchar(20) NOT NULL,
	datetrx timestamptz NOT NULL,
	source_wallet_id varchar(150) NULL,
	source_wallet_snapshot varchar(200) NULL,
	destination_wallet_id varchar(150) NULL,
	destination_wallet_snapshot varchar(200) NULL,
	amount numeric(36, 18) DEFAULT 0 NOT NULL,
	admin_fee numeric(36, 18) DEFAULT 0 NOT NULL,
	rate_trx numeric(36, 18) DEFAULT 0 NOT NULL,
	amount_final numeric(36, 18) DEFAULT 0 NOT NULL,
	amountratefinal_trx numeric(36, 18) DEFAULT 0 NOT NULL,
	balance_source_before numeric(36, 18) DEFAULT 0 NOT NULL,
	balance_destination_before numeric(36, 18) DEFAULT 0 NOT NULL,
	status varchar(20) DEFAULT 'DRAFT'::character varying NULL,
	title varchar(150) DEFAULT ''::character varying NULL,
	detail varchar(250) DEFAULT ''::character varying NULL,
	note varchar(150) DEFAULT ''::character varying NULL,
	approved_by varchar(30) NULL,
	approved_at timestamptz NULL,
	created_by varchar(30) DEFAULT ''::character varying NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_by varchar(30) DEFAULT ''::character varying NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_wallet_external_transaction_unique UNIQUE (idtrx)
);


-- db_bbca.tbl_wallet_member definition

-- Drop table

-- DROP TABLE db_bbca.tbl_wallet_member;

CREATE TABLE db_bbca.tbl_wallet_member (
	idwalletmember varchar(150) NOT NULL,
	idmember varchar(150) DEFAULT ''::character varying NULL,
	idcurr varchar(20) NOT NULL,
	idbank varchar(10) NULL,
	account_number varchar(150) DEFAULT ''::character varying NULL,
	account_name varchar(150) DEFAULT ''::character varying NULL,
	networkcrypto varchar(20) DEFAULT ''::character varying NULL,
	status varchar(1) DEFAULT 'Y'::character varying NULL,
	created_by varchar(30) DEFAULT ''::character varying NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_by varchar(30) DEFAULT ''::character varying NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_wallet_member_unique UNIQUE (idwalletmember)
);


-- db_bbca.tbl_wallet_transaction definition

-- Drop table

-- DROP TABLE db_bbca.tbl_wallet_transaction;

CREATE TABLE db_bbca.tbl_wallet_transaction (
	idtrx varchar(150) NOT NULL,
	idowner varchar(150) DEFAULT ''::character varying NULL,
	typeact varchar(20) NOT NULL,
	typetrx varchar(20) NOT NULL,
	idcurr varchar(20) NOT NULL,
	datetrx timestamptz NOT NULL,
	source_wallet_id varchar(150) NULL,
	source_wallet_snapshot varchar(200) NULL,
	destination_wallet_id varchar(150) NULL,
	destination_wallet_snapshot varchar(200) NULL,
	amount numeric(36, 18) DEFAULT 0 NOT NULL,
	admin_fee numeric(36, 18) DEFAULT 0 NOT NULL,
	rate_trx numeric(36, 18) DEFAULT 0 NOT NULL,
	amount_final numeric(36, 18) DEFAULT 0 NOT NULL,
	amountratefinal_trx numeric(36, 18) DEFAULT 0 NOT NULL,
	balance_source_before numeric(36, 18) DEFAULT 0 NOT NULL,
	balance_destination_before numeric(36, 18) DEFAULT 0 NOT NULL,
	status varchar(20) DEFAULT 'DRAFT'::character varying NULL,
	title varchar(150) DEFAULT ''::character varying NULL,
	detail varchar(250) DEFAULT ''::character varying NULL,
	note varchar(150) DEFAULT ''::character varying NULL,
	approved_by varchar(30) NULL,
	approved_at timestamp NULL,
	created_by varchar(30) DEFAULT ''::character varying NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_by varchar(30) DEFAULT ''::character varying NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	idgroup varchar(4) DEFAULT ''::character varying NULL,
	CONSTRAINT tbl_wallet_transaction_unique UNIQUE (idtrx)
);




CREATE TABLE db_bbca.tbl_item_category (
	id bigserial NOT NULL,
	"name" varchar(100) NOT NULL,
	parent_id int8 NULL,
	"level" int4 DEFAULT 1 NOT NULL,
	"path" text DEFAULT ''::text NOT NULL,
	status varchar(1) NOT NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT chk_no_self_parent CHECK (((id IS NULL) OR (parent_id IS NULL) OR (id <> parent_id))),
	CONSTRAINT tbl_item_category_pkey PRIMARY KEY (id),
	CONSTRAINT uq_item_category_name_parent UNIQUE (name, parent_id),
	CONSTRAINT fk_item_category_parent FOREIGN KEY (parent_id) REFERENCES db_bbca.tbl_item_category(id) ON DELETE SET NULL
);
CREATE INDEX idx_item_category_active ON db_bbca.tbl_item_category USING btree (status);
CREATE INDEX idx_item_category_parent ON db_bbca.tbl_item_category USING btree (parent_id);
CREATE INDEX idx_item_category_path ON db_bbca.tbl_item_category USING btree (path);

-- Table Triggers

create trigger trg_after_insert_fix after
insert
    on
    db_bbca.tbl_item_category for each row execute function db_bbca.trg_fix_category();
create trigger trg_3_status_sync after
update
    of status on
    db_bbca.tbl_item_category for each row execute function db_bbca.trg_item_category_status_sync();
create trigger trg_item_category_path_logic before
insert
    or
update
    of parent_id on
    db_bbca.tbl_item_category for each row execute function db_bbca.fn_recursive_category_update();




CREATE TABLE db_bbca.tbl_mst_item_stock (
	iditemstock bigserial NOT NULL,
	iditem varchar(20) NOT NULL,
	iduom varchar(10) NOT NULL,
	total_in numeric(36, 18) DEFAULT 0 NOT NULL,
	total_out numeric(36, 18) DEFAULT 0 NOT NULL,
	create_by varchar(30) DEFAULT ''::character varying NULL,
	create_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_by varchar(30) DEFAULT ''::character varying NULL,
	update_at timestamptz NULL,
	CONSTRAINT tbl_mst_item_stock_pkey PRIMARY KEY (iditemstock),
	CONSTRAINT tbl_mst_item_stock_total_in_check CHECK ((total_in >= (0)::numeric)),
	CONSTRAINT tbl_mst_item_stock_total_out_check CHECK ((total_out >= (0)::numeric)),
	CONSTRAINT uq_item_uom UNIQUE (iditem, iduom),
	CONSTRAINT fk_item_stock_item FOREIGN KEY (iditem) REFERENCES db_bbca.tbl_mst_item(iditem) ON DELETE RESTRICT ON UPDATE CASCADE
);
CREATE INDEX idx_item_stock_item ON db_bbca.tbl_mst_item_stock USING btree (iditem);
CREATE INDEX idx_item_stock_uom ON db_bbca.tbl_mst_item_stock USING btree (iduom);