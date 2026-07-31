package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devhdn-212/totclient_pools/internal/config"
	"github.com/devhdn-212/totclient_pools/internal/connection"
	"github.com/devhdn-212/totclient_pools/internal/repository"
	"github.com/devhdn-212/totclient_pools/internal/service"
	"github.com/devhdn-212/totclient_pools/internal/util"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// totclient_pools is a pure background worker — no HTTP server. It consumes
// checkout events totclient_api publishes to Kafka, debits the wallet, and
// persists the bet. See internal/service/checkoutconsumer.go.
func main() {
	cnf := config.Get()
	logger := NewGCPLogger(cnf.Telegram)

	connection.SetLogger(logger)

	dbPool := connection.GetDatabase(cnf.Database)
	defer dbPool.Close()

	if err := connection.InitRedis(cnf.Redis); err != nil {
		logger.Fatal("Failed to init Redis", zap.Error(err))
	}
	defer connection.RDB.Close()
	defer connection.RDBLimit.Close()

	if !connection.RedisHealth() {
		logger.Fatal("Redis is not healthy")
	}

	pgxExec := repository.NewPGXExecutor(dbPool)
	companyRepository := repository.NewCompanyRepository(pgxExec)
	trxkeluaranRepository := repository.NewTrxkeluaranRepository(pgxExec)
	memberinfoService := service.NewMemberinfoService(companyRepository, cnf.BalanceAPI)

	kafkaReader := connection.NewKafkaReader(cnf.Kafka)
	defer kafkaReader.Close()
	checkoutConsumer := service.NewCheckoutConsumer(dbPool, kafkaReader, memberinfoService)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Drains periods the consumer marked dirty (service.MarkTotalsDirty) and
	// recomputes total_member/total_bet/total_pairs/total_payout — same
	// batching rationale as totclient_api's own ticker (see trxkeluaran.go):
	// one recompute per busy period per tick instead of one UPDATE per
	// checkout event.
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				service.FlushDirtyTotals(context.Background(), trxkeluaranRepository)
			}
		}
	}()

	logger.Info("totclient_pools: checkout consumer starting", zap.String("topic", cnf.Kafka.Topic))
	checkoutConsumer.Run(ctx)

	logger.Info("Gracefully shutting down...")
}

func NewGCPLogger(tg config.Telegram) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		enc.AppendString(t.In(loc).Format("2006-01-02 15:04:05"))
	}
	encoderConfig.LevelKey = "severity"
	encoderConfig.MessageKey = "message"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.Hooks(func(entry zapcore.Entry) error {
		// Every Error/Fatal log anywhere in the app flows through here —
		// entry.Caller pinpoints the exact file:line so the alert says
		// where to look, not just what broke. Fire-and-forget: never let a
		// slow Telegram API delay the consumer loop.
		if entry.Level < zapcore.ErrorLevel {
			return nil
		}
		go util.SendTelegramAlert(tg, entry.Level.CapitalString(), entry.Caller.String(), entry.Message)
		return nil
	}))
}
