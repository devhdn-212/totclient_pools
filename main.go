package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devhdn-212/totclient_pools/domain"
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
	memberInvoiceRepository := repository.NewMemberInvoiceRepository(pgxExec)

	// public.tbl_trx_member_invoice (bagian 9.8 DOKUMENTASI.md) is brand new -
	// make sure the current AND next month's partition exist before any
	// checkout tries to insert into it. Ongoing month rollovers (if this
	// process runs for 2+ months without a restart) are handled by the same
	// 1-minute ticker below that already runs FlushDirtyTotals.
	{
		startupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		now := util.GetNowJakarta()
		thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		nextMonth := thisMonth.AddDate(0, 1, 0)
		if err := memberInvoiceRepository.EnsureMonthPartition(startupCtx, thisMonth); err != nil {
			logger.Error("gagal bikin partisi tbl_trx_member_invoice bulan ini", zap.Error(err))
		}
		if err := memberInvoiceRepository.EnsureMonthPartition(startupCtx, nextMonth); err != nil {
			logger.Error("gagal bikin partisi tbl_trx_member_invoice bulan depan", zap.Error(err))
		}
		cancel()
	}

	kafkaReader := connection.NewKafkaReader(cnf.Kafka)
	defer kafkaReader.Close()
	checkoutConsumer := service.NewCheckoutConsumer(dbPool, kafkaReader, memberinfoService, memberInvoiceRepository)

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
				ensureNextMonthMemberInvoicePartition(context.Background(), memberInvoiceRepository)
			}
		}
	}()

	logger.Info("totclient_pools: checkout consumer starting", zap.String("topic", cnf.Kafka.Topic))
	checkoutConsumer.Run(ctx)

	logger.Info("Gracefully shutting down...")
}

// ensureNextMonthMemberInvoicePartition mirrors the agen-side
// ensureNextMonthPartitions gating (see totagen_api/totagen_pools
// documentasi.md bagian 63) - only actually attempts the (idempotent) DDL
// once the current month is in its last 7 days, so this doesn't run a
// CREATE TABLE statement every single ticker tick for no reason.
func ensureNextMonthMemberInvoicePartition(ctx context.Context, repo domain.MemberInvoiceRepository) {
	now := util.GetNowJakarta()
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()
	if now.Day() < daysInMonth-6 {
		return
	}
	nextMonthStart := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	if err := repo.EnsureMonthPartition(ctx, nextMonthStart); err != nil {
		connection.Log.Error("gagal bikin partisi tbl_trx_member_invoice bulan depan", zap.Error(err))
	}
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
