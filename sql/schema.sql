CREATE TABLE public.tbl_counter (
	idcounter int4 DEFAULT nextval('idcounter_seq'::regclass) NOT NULL,
	nmcounter varchar(70) NULL,
	counter int8 NOT NULL,
	CONSTRAINT tbl_counter_pk PRIMARY KEY (idcounter)
);


CREATE TABLE public.tbl_admin (
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



CREATE TABLE public.tbl_company (
	idcompany varchar(10) NOT NULL,
	idcurrdef varchar(20) NOT NULL,
	compname varchar(50) NULL,
	endjoin timestamp NULL,
	amountcomp numeric(36, 18) DEFAULT 0 NOT NULL,
	compstatus varchar(1) DEFAULT 'Y'::character varying NULL,
	createcomp varchar(30) DEFAULT ''::character varying NULL,
	createdatecomp timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updatecomp varchar(30) DEFAULT ''::character varying NULL,
	updatedatecomp timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_company_unique UNIQUE (idcompany)
);


CREATE TABLE public.tbl_company_admin (
	idcompadmin varchar(64) NOT NULL,
	idcompany varchar(10) NOT NULL,
	idclientrule varchar(30) NOT NULL,
	usernamecompadmin varchar(30) NOT NULL,
	namecompadmin varchar(50) NULL,
	passcompadmin varchar(250) NULL,
	lastlogincompadmin timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	ipaddresscompadmin varchar(20) DEFAULT ''::character varying NULL,
	compadminstatus varchar(1) DEFAULT 'Y'::character varying NULL,
	createcompadmin varchar(30) DEFAULT ''::character varying NULL,
	createdatecompadmin timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updatecompadmin varchar(30) DEFAULT ''::character varying NULL,
	updatedatecompadmin timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_company_admin_unique UNIQUE (idcompadmin)
);


CREATE TABLE public.tbl_company_wallet (
	idcompwallet varchar(64) NOT NULL,
	idcompany varchar(10) NOT NULL,
	idcurr varchar(20) NOT NULL,
	amountcompwallet numeric(36, 18) DEFAULT 0 NOT NULL,
	compwalletstatus varchar(1) DEFAULT 'Y'::character varying NULL,
	createcompwallet varchar(30) DEFAULT ''::character varying NULL,
	createdatecompwallet timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updatecompwallet varchar(30) DEFAULT ''::character varying NULL,
	updatedatecompwallet timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_company_wallet_unique UNIQUE (idcompwallet)
);



CREATE TABLE public.tbl_currency (
	idcurr varchar(20) NOT NULL,
	typecurr varchar(10) DEFAULT ''::character varying NOT NULL,
	status varchar(1) NOT NULL,
	createcurr varchar(30) DEFAULT ''::character varying NULL,
	createdatecurr timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updatecurr varchar(30) DEFAULT ''::character varying NULL,
	updatedatecurr timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT tbl_currency_status_check CHECK (((status)::bpchar = ANY (ARRAY['Y'::bpchar, 'N'::bpchar]))),
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