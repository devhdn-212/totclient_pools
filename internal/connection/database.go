package connection

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"log"
	"time"

	"github.com/devhdn-212/gofibergoqu_master/internal/config"
)

func GetDatabase(conf config.Database) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable Timezone=Asia/Jakarta",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Pass,
		conf.Name,
		conf.Schema)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		Log.Fatal("Failed to connect to database", zap.Error(err))
	} else {
		Log.Info("Successfully connected to database")
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database", zap.Error(err))
	}
	return db
}
