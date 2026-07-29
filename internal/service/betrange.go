package service

import (
	"fmt"

	"github.com/devhdn-212/totclient_api/dto"
	"github.com/shopspring/decimal"
)

// fourDMaxBet picks the max-bet ceiling for a 4D-family item's tipetoto —
// FULL/BB each have their own dedicated field, DISC (or anything else) uses
// the base one. Mirrors getPasaranMaxBet in web2026/src/lib/utils.ts: when
// the FULL/BB-specific field is itself zero/missing, that mode is genuinely
// uncapped, it does NOT fall back to the base field.
func fourDMaxBet(tipetoto string, base, full, bb decimal.Decimal) decimal.Decimal {
	switch tipetoto {
	case "FULL":
		return full
	case "BB":
		return bb
	default:
		return base
	}
}

// resolveMinMaxBet mirrors getPasaranMinBet/getPasaranMaxBet in utils.ts —
// the minbet/maxbet ceilings for one bet, straight off the live pasaran
// config. A value <= 0 means that side is uncapped (no minimum / no
// maximum), same convention already used for limittotal/limitglobal
// (resolveLimits) and CheckAndIncrementLimitRedis elsewhere in this package.
func resolveMinMaxBet(p dto.PasaranData, typegame, tipetoto string) (minBet, maxBet decimal.Decimal) {
	switch typegame {
	case "4D":
		return p.AngkaMinbet, fourDMaxBet(tipetoto, p.AngkaMaxbet4d, p.AngkaMaxbet4dFull, p.AngkaMaxbet4dBb)
	case "3D":
		return p.AngkaMinbet, fourDMaxBet(tipetoto, p.AngkaMaxbet3d, p.AngkaMaxbet3dFull, p.AngkaMaxbet3dBb)
	case "3DD":
		return p.AngkaMinbet, fourDMaxBet(tipetoto, p.AngkaMaxbet3dd, p.AngkaMaxbet3ddFull, p.AngkaMaxbet3ddBb)
	case "2D":
		return p.AngkaMinbet, fourDMaxBet(tipetoto, p.AngkaMaxbet2d, p.AngkaMaxbet2dFull, p.AngkaMaxbet2dBb)
	case "2DD":
		return p.AngkaMinbet, fourDMaxBet(tipetoto, p.AngkaMaxbet2dd, p.AngkaMaxbet2ddFull, p.AngkaMaxbet2ddBb)
	case "2DT":
		return p.AngkaMinbet, fourDMaxBet(tipetoto, p.AngkaMaxbet2dt, p.AngkaMaxbet2dtFull, p.AngkaMaxbet2dtBb)
	case "COLOK_BEBAS":
		return p.CbMinbet, p.CbMaxbet
	case "COLOK_MACAU":
		return p.CmacauMinbet, p.CmacauMaxbet
	case "COLOK_NAGA":
		return p.CnagaMinbet, p.CnagaMaxbet
	case "COLOK_JITU":
		return p.CjituMinbet, p.CjituMaxbet
	case "50_50_UMUM":
		return p.Umum5050Minbet, p.Umum5050Maxbet
	case "50_50_SPECIAL":
		return p.Special5050Minbet, p.Special5050Maxbet
	case "50_50_KOMBINASI":
		return p.Kombinasi5050Minbet, p.Kombinasi5050Maxbet
	case "MACAU_KOMBINASI":
		return p.MacaukombinasiMinbet, p.MacaukombinasiMaxbet
	case "DASAR":
		return p.DasarMinbet, p.DasarMaxbet
	case "SHIO":
		return p.ShioMinbet, p.ShioMaxbet
	}
	return decimal.Zero, decimal.Zero
}

// checkoutItemRejectReason validates one submitted item's nomor format and
// bet size against the live pasaran config, and reports why it's invalid (or
// "" if it's fine) as a multi-line, player-facing message — checked well
// before any Redis/DB work, see isValidNomor's doc comment in betformat.go
// for why this can't be left to calculatePayout to catch.
func checkoutItemRejectReason(p dto.PasaranData, item dto.CheckoutItem) string {
	if !isValidNomor(item.Permainan, item.Nomor) {
		return fmt.Sprintf(
			"Format nomor tidak valid\nNomor : %s\nPermainan : %s\nTipe : %s\nBet : %d",
			item.Nomor, item.Permainan, item.Tipetoto, item.Bet,
		)
	}

	minBet, maxBet := resolveMinMaxBet(p, item.Permainan, item.Tipetoto)
	betDecimal := decimal.NewFromInt(int64(item.Bet))
	if minBet.IsPositive() && betDecimal.LessThan(minBet) {
		return fmt.Sprintf(
			"Bet dibawah minimal bet\nNomor : %s\nPermainan : %s\nTipe : %s\nBet : %d\nMinimal Bet : %s",
			item.Nomor, item.Permainan, item.Tipetoto, item.Bet, minBet.String(),
		)
	}
	if maxBet.IsPositive() && betDecimal.GreaterThan(maxBet) {
		return fmt.Sprintf(
			"Bet melebihi maximal bet\nNomor : %s\nPermainan : %s\nTipe : %s\nBet : %d\nMaximal Bet : %s",
			item.Nomor, item.Permainan, item.Tipetoto, item.Bet, maxBet.String(),
		)
	}
	return ""
}