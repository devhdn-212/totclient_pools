package service

import (
	"reflect"
	"strings"

	"github.com/devhdn-212/totclient_pools/dto"
	"github.com/shopspring/decimal"
)

// This file recomputes diskon/kei/win/bayar server-side from the live
// pasaran config — a Go mirror of web2026/src/lib/utils.ts's
// getPasaranRates/calculatePayout/getPasaranWinValue. The client sends the
// same numbers for display purposes, but they're never trusted for
// persistence: a modified request could otherwise inflate its own payout.

var fourDDiscField = map[string]func(dto.PasaranData) decimal.Decimal{
	"4D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaDisc4d },
	"3D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaDisc3d },
	"3DD": func(p dto.PasaranData) decimal.Decimal { return p.AngkaDisc3dd },
	"2D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaDisc2d },
	"2DD": func(p dto.PasaranData) decimal.Decimal { return p.AngkaDisc2dd },
	"2DT": func(p dto.PasaranData) decimal.Decimal { return p.AngkaDisc2dt },
}

var flatDiscField = map[string]func(dto.PasaranData) decimal.Decimal{
	"COLOK_BEBAS":     func(p dto.PasaranData) decimal.Decimal { return p.CbDisc },
	"COLOK_MACAU":     func(p dto.PasaranData) decimal.Decimal { return p.CmacauDisc },
	"COLOK_NAGA":      func(p dto.PasaranData) decimal.Decimal { return p.CnagaDisc },
	"COLOK_JITU":      func(p dto.PasaranData) decimal.Decimal { return p.CjituDesic },
	"MACAU_KOMBINASI": func(p dto.PasaranData) decimal.Decimal { return p.MacaukombinasiDiscount },
	"SHIO":            func(p dto.PasaranData) decimal.Decimal { return p.ShioDisc },
}

var fourDWinSuffix = map[string]string{
	"4D": "4d", "3D": "3d", "3DD": "3dd", "2D": "2d", "2DD": "2dd", "2DT": "2dt",
}

var fourDWinDisc = map[string]func(dto.PasaranData) decimal.Decimal{
	"4D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin4d },
	"3D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin3d },
	"3DD": func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin3dd },
	"2D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin2d },
	"2DD": func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin2dd },
	"2DT": func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin2dt },
}
var fourDWinFull = map[string]func(dto.PasaranData) decimal.Decimal{
	"4D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin4dnodisc },
	"3D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin3dnodisc },
	"3DD": func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin3ddnodisc },
	"2D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin2dnodisc },
	"2DD": func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin2ddnodisc },
	"2DT": func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin2dtnodisc },
}
var fourDWinBB = map[string]func(dto.PasaranData) decimal.Decimal{
	"4D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin4dbbKena },
	"3D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin3dbbKena },
	"3DD": func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin3ddbbKena },
	"2D":  func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin2dbbKena },
	"2DD": func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin2ddbbKena },
	"2DT": func(p dto.PasaranData) decimal.Decimal { return p.AngkaWin2dtbbKena },
}

var cjituWinField = map[string]func(dto.PasaranData) decimal.Decimal{
	"as":     func(p dto.PasaranData) decimal.Decimal { return p.CjituWinas },
	"kop":    func(p dto.PasaranData) decimal.Decimal { return p.CjituWinkop },
	"kepala": func(p dto.PasaranData) decimal.Decimal { return p.CjituWinkepala },
	"ekor":   func(p dto.PasaranData) decimal.Decimal { return p.CjituWinekor },
}

// fieldDecimal reads a decimal.Decimal field off Pasaran by Go field name —
// used only for 50/50 Special/Kombinasi, whose rate field names are built
// from a posisi+jenis/kondisi combination (dozens of them) that's not worth
// spelling out as explicit maps the way the other types are above.
func fieldDecimal(p dto.PasaranData, fieldName string) decimal.Decimal {
	v := reflect.ValueOf(p).FieldByName(fieldName)
	if !v.IsValid() {
		return decimal.Zero
	}
	d, ok := v.Interface().(decimal.Decimal)
	if !ok {
		return decimal.Zero
	}
	return d
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

type payoutRates struct {
	DiscRate decimal.Decimal
	KeiRate  decimal.Decimal
	HasKei   bool
}

// resolveRates reads disc/kei straight off the pasaran config — those
// columns already hold the fraction (0.66 means 66%), not a raw percentage
// number, so no /100 conversion happens here.
func resolveRates(p dto.PasaranData, typegame, nomor string) payoutRates {
	if f, ok := fourDDiscField[typegame]; ok {
		return payoutRates{DiscRate: f(p)}
	}
	if f, ok := flatDiscField[typegame]; ok {
		return payoutRates{DiscRate: f(p)}
	}

	if typegame == "DASAR" {
		switch nomor {
		case "BESAR":
			return payoutRates{DiscRate: p.DasarDiscbesar, KeiRate: p.DasarKeibesar, HasKei: true}
		case "KECIL":
			return payoutRates{DiscRate: p.DasarDisckecil, KeiRate: p.DasarKeikecil, HasKei: true}
		case "GENAP":
			return payoutRates{DiscRate: p.DasarDiscigenap, KeiRate: p.DasarKeigenap, HasKei: true}
		case "GANJIL":
			return payoutRates{DiscRate: p.DasarDiscganjil, KeiRate: p.DasarKeiganjil, HasKei: true}
		}
		return payoutRates{}
	}

	if typegame == "50_50_UMUM" {
		switch nomor {
		case "BESAR":
			return payoutRates{DiscRate: p.Umum5050Discbesar, KeiRate: p.Umum5050Keibesar, HasKei: true}
		case "KECIL":
			return payoutRates{DiscRate: p.Umum5050Disckecil, KeiRate: p.Umum5050Keikecil, HasKei: true}
		case "GENAP":
			return payoutRates{DiscRate: p.Umum5050Discgenap, KeiRate: p.Umum5050Keigenap, HasKei: true}
		case "GANJIL":
			return payoutRates{DiscRate: p.Umum5050Discganjil, KeiRate: p.Umum5050Keiganjil, HasKei: true}
		case "TENGAH":
			return payoutRates{DiscRate: p.Umum5050Disctengah, KeiRate: p.Umum5050Keitengah, HasKei: true}
		case "TEPI":
			return payoutRates{DiscRate: p.Umum5050Disctepi, KeiRate: p.Umum5050Keitepi, HasKei: true}
		}
		return payoutRates{}
	}

	// Special5050 numbers are "{POS}_{KONDISI}" (AS/KOP/KEPALA/EKOR x
	// BESAR/KECIL/GENAP/GANJIL); every field is
	// Special5050{Disc,Kei}{pos}{kondisi} lowercase.
	if typegame == "50_50_SPECIAL" {
		pos, kondisi, ok := strings.Cut(nomor, "_")
		if !ok {
			return payoutRates{}
		}
		suffix := strings.ToLower(pos) + strings.ToLower(kondisi)
		disc := fieldDecimal(p, "Special5050Disc"+suffix)
		kei := fieldDecimal(p, "Special5050Kei"+suffix)
		return payoutRates{DiscRate: disc, KeiRate: kei, HasKei: true}
	}

	// Kombinasi5050 numbers are "{POSISI}_{JENIS}" (BELAKANG/TENGAH/DEPAN x
	// MONO/STEREO/KEMBANG/KEMPIS/KEMBAR); fields are
	// Kombinasi5050{Posisi}{disc,kei}{jenis} — posisi capitalized, jenis
	// lowercase.
	if typegame == "50_50_KOMBINASI" {
		posisi, jenis, ok := strings.Cut(nomor, "_")
		if !ok {
			return payoutRates{}
		}
		p1 := capitalize(strings.ToLower(posisi))
		j := strings.ToLower(jenis)
		disc := fieldDecimal(p, "Kombinasi5050"+p1+"disc"+j)
		kei := fieldDecimal(p, "Kombinasi5050"+p1+"kei"+j)
		return payoutRates{DiscRate: disc, KeiRate: kei, HasKei: true}
	}

	return payoutRates{}
}

// resolveWinValue mirrors getPasaranWinValue in utils.ts — a rate/multiplier
// recorded alongside the bet for reference; the real payout is only settled
// once the draw result is in.
func resolveWinValue(p dto.PasaranData, typegame, nomor, tipetoto string) decimal.Decimal {
	if _, ok := fourDWinSuffix[typegame]; ok {
		switch tipetoto {
		case "FULL":
			return fourDWinFull[typegame](p)
		case "BB":
			return fourDWinBB[typegame](p)
		default:
			return fourDWinDisc[typegame](p)
		}
	}

	switch typegame {
	case "COLOK_BEBAS":
		return p.CbWin
	case "COLOK_MACAU":
		switch len(nomor) {
		case 2:
			return p.CmacauWin2digit
		case 3:
			return p.CmacauWin3digit
		case 4:
			return p.CmacauWin4digit
		}
		return decimal.Zero
	case "COLOK_NAGA":
		switch len(nomor) {
		case 3:
			return p.CnagaWin3digit
		case 4:
			return p.CnagaWin4digit
		}
		return decimal.Zero
	case "COLOK_JITU":
		_, pos, ok := strings.Cut(nomor, "_")
		if !ok {
			return decimal.Zero
		}
		if f, ok := cjituWinField[strings.ToLower(pos)]; ok {
			return f(p)
		}
		return decimal.Zero
	case "MACAU_KOMBINASI":
		return p.MacaukombinasiWin
	case "SHIO":
		return p.ShioWin
	case "50_50_UMUM", "50_50_SPECIAL", "50_50_KOMBINASI", "DASAR":
		return decimal.NewFromInt(1)
	}
	return decimal.Zero
}

type payoutResult struct {
	// Disc/Kei are the computed money amounts (bet * rate) — used for the
	// invoice's totaldiscount/totalkei/payout math.
	Disc decimal.Decimal
	Kei  decimal.Decimal
	// DiscRate/KeiRate are the raw rate straight from the pasaran config
	// (e.g. 0.66, -0.015) — this is what gets persisted onto the detail
	// row's diskon/kei columns, not the computed amount, so the row records
	// what rate was in effect rather than a number that only makes sense
	// combined with that bet's size.
	DiscRate decimal.Decimal
	KeiRate  decimal.Decimal
	Win      decimal.Decimal
	Payout   decimal.Decimal
}

// calculatePayout mirrors calculatePayout in utils.ts: DISC applies
// discRate (and keiRate, for the types that have one); FULL/BB pay out the
// bet as-is, no discount.
func calculatePayout(p dto.PasaranData, typegame, nomor string, bet int, tipetoto string) payoutResult {
	betDecimal := decimal.NewFromInt(int64(bet))
	win := resolveWinValue(p, typegame, nomor, tipetoto)

	if tipetoto != "DISC" {
		return payoutResult{Payout: betDecimal, Win: win}
	}

	rates := resolveRates(p, typegame, nomor)
	disc := betDecimal.Mul(rates.DiscRate)
	if rates.HasKei {
		kei := betDecimal.Mul(rates.KeiRate).Neg()
		return payoutResult{
			Disc:     disc,
			Kei:      kei,
			DiscRate: rates.DiscRate,
			KeiRate:  rates.KeiRate,
			Win:      win,
			Payout:   betDecimal.Add(kei).Sub(disc),
		}
	}
	return payoutResult{Disc: disc, DiscRate: rates.DiscRate, Win: win, Payout: betDecimal.Sub(disc)}
}
