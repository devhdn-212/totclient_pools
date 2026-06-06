package connection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"time"

	"github.com/devhdn-212/gofibermaster_api/internal/config"
)

func GetDatabase(conf config.Database) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=%s&sslmode=disable&timezone=Asia/Jakarta",
		conf.User,
		conf.Pass,
		conf.Host,
		conf.Port,
		conf.Name,
		conf.Schema,
	)

	// 2. Parsing konfigurasi pool
	configPool, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		Log.Fatal("Gagal parsing config database", zap.Error(err))
	}

	// 3. Pengaturan Connection Pool (Mirip dengan sql.DB)
	configPool.MaxConns = 100                     // SetMaxOpenConns
	configPool.MinConns = 10                      // SetMaxIdleConns (pendekatan terdekat)
	configPool.MaxConnIdleTime = 5 * time.Minute  // SetConnMaxIdleTime
	configPool.MaxConnLifetime = 60 * time.Minute // SetConnMaxLifetime

	// 4. Membuat koneksi pool
	// Context background digunakan karena ini inisialisasi awal
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
