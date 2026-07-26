package domain

import (
	"context"

	"github.com/devhdn-212/totclient_api/dto"
)

// ResultService resolves a pasaran's past draw results (keluarantogel) for
// a given month — the "Result" menu. Distinct from RiwayatTransaksiService
// (a player's own bet history): draw results are public information for
// the whole pasaran, not scoped to any one username.
type ResultService interface {
	Fetch(ctx context.Context, req dto.ResultRequest) (dto.ResultResponse, error)
}
