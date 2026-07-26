package domain

import (
	"context"

	"github.com/devhdn-212/totclient_api/dto"
)

// RiwayatTransaksiService resolves a player's own successful bets across
// every period of a pasaran — a Transaksi-menu lookup, distinct from the
// admin TrxkeluarandetailService.All (which sees every player's rows).
type RiwayatTransaksiService interface {
	Fetch(ctx context.Context, req dto.RiwayatTransaksiRequest) (dto.RiwayatTransaksiResponse, error)
}
