package connection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"time"

	"github.com/devhdn-212/totmaster_api/internal/config"
)

func GetDatabase(conf config.Database) *pgxpool.Pool {
	// 1. Buat config dasar pakai string minimalis tanpa password yang aneh-aneh
	baseDSN := fmt.Sprintf("postgres://%s:%s/subapp?sslmode=disable", conf.Host, conf.Port)

	configPool, err := pgxpool.ParseConfig(baseDSN)
	if err != nil {
		Log.Fatal("Gagal parsing config database", zap.Error(err))
	}

	// 2. Isi data aslinya langsung ke objek struct-nya (Aman dari eror karakter @ atau !)
	configPool.ConnConfig.User = conf.User
	configPool.ConnConfig.Password = conf.Pass
	configPool.ConnConfig.Database = conf.Name

	// Jika ingin menambahkan runtime params seperti search_path atau timezone
	if configPool.ConnConfig.RuntimeParams == nil {
		configPool.ConnConfig.RuntimeParams = make(map[string]string)
	}
	configPool.ConnConfig.RuntimeParams["search_path"] = conf.Schema
	configPool.ConnConfig.RuntimeParams["timezone"] = "Asia/Jakarta"

	// 3. Pengaturan Connection Pool
	configPool.MaxConns = 100
	configPool.MinConns = 10
	configPool.MaxConnIdleTime = 5 * time.Minute
	configPool.MaxConnLifetime = 60 * time.Minute

	// 4. Membuat koneksi pool
	dbPool, err := pgxpool.NewWithConfig(context.Background(), configPool)
	if err != nil {
		Log.Fatal("Gagal membuat pool database", zap.Error(err))
	}

	// 5. Verifikasi koneksi dengan Ping
	err = dbPool.Ping(context.Background())
	if err != nil {
		Log.Fatal("Gagal ping database", zap.Error(err))
	}

	Log.Info("Berhasil terhubung ke database dengan pgxpool")

	return dbPool
}
