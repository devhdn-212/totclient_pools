package service

import (
	"regexp"
	"strings"
)

// fourDDigitCount is the number of real digits nomor must contain for each
// 4D-family type, once "*" wildcard markers (used by 3DD/2DD/2DT, see
// stripAsterisks in checkout.go) are stripped out.
var fourDDigitCount = map[string]int{
	"4D": 4, "3D": 3, "3DD": 3, "2D": 2, "2DD": 2, "2DT": 2,
}

var onlyDigitsAndStars = regexp.MustCompile(`^[0-9*]+$`)
var onlyDigits = regexp.MustCompile(`^[0-9]+$`)

var cjituPositions = map[string]bool{"AS": true, "KOP": true, "KEPALA": true, "EKOR": true}
var special5050Positions = map[string]bool{"AS": true, "KOP": true, "KEPALA": true, "EKOR": true}
var special5050Kondisi = map[string]bool{"BESAR": true, "KECIL": true, "GENAP": true, "GANJIL": true}
var kombinasi5050Posisi = map[string]bool{"BELAKANG": true, "TENGAH": true, "DEPAN": true}
var kombinasi5050Jenis = map[string]bool{
	"MONO": true, "STEREO": true, "KEMBANG": true, "KEMPIS": true, "KEMBAR": true,
}
var umum5050Options = map[string]bool{
	"BESAR": true, "KECIL": true, "GENAP": true, "GANJIL": true, "TENGAH": true, "TEPI": true,
}
var dasarOptions = map[string]bool{"GANJIL": true, "GENAP": true, "BESAR": true, "KECIL": true}
var macauKombinasiPosisi = map[string]bool{"BELAKANG": true, "TENGAH": true, "DEPAN": true}
var macauKombinasiBesarKecil = map[string]bool{"BESAR": true, "KECIL": true}
var macauKombinasiGenapGanjil = map[string]bool{"GENAP": true, "GANJIL": true}

// shioNames matches web2026's Formshio.svelte SHIO_OPTIONS exactly (title
// case) — payout.go never reads nomor for SHIO (one flat rate for every
// animal), so this exists purely to reject junk, not to select a rate.
var shioNames = map[string]bool{
	"Tikus": true, "Kerbau": true, "Macan": true, "Kelinci": true, "Naga": true, "Ular": true,
	"Kuda": true, "Kambing": true, "Monyet": true, "Ayam": true, "Anjing": true, "Babi": true,
}

// isValidNomor reports whether nomor is a well-formed bet number for
// typegame, mirroring exactly what web2026's Form*.svelte components ever
// produce (see DOKUMENTASI.md §3.9). calculatePayout/resolveRates never
// reject a malformed nomor on their own — an unrecognized typegame or a
// nomor that doesn't match any expected pattern just falls through to a
// zero/no-op rate (or, for 4D-family types, gets a real payout despite the
// nomor being nonsense, since typegame alone drives the rate lookup there).
// A client hitting /api/checkout directly could otherwise get a "accepted"
// bet persisted for input that can never settle against a real draw result.
// This check runs before any limit/payout logic so that case is rejected
// outright instead of silently accepted.
func isValidNomor(typegame, nomor string) bool {
	if digitCount, ok := fourDDigitCount[typegame]; ok {
		if !onlyDigitsAndStars.MatchString(nomor) {
			return false
		}
		return len(strings.ReplaceAll(nomor, "*", "")) == digitCount
	}

	switch typegame {
	case "COLOK_BEBAS":
		return len(nomor) == 1 && onlyDigits.MatchString(nomor)
	case "COLOK_MACAU":
		return len(nomor) >= 2 && len(nomor) <= 4 && onlyDigits.MatchString(nomor)
	case "COLOK_NAGA":
		return len(nomor) >= 3 && len(nomor) <= 4 && onlyDigits.MatchString(nomor)
	case "COLOK_JITU":
		digit, pos, ok := strings.Cut(nomor, "_")
		return ok && len(digit) == 1 && onlyDigits.MatchString(digit) && cjituPositions[strings.ToUpper(pos)]
	case "50_50_UMUM":
		return umum5050Options[nomor]
	case "50_50_SPECIAL":
		pos, kondisi, ok := strings.Cut(nomor, "_")
		return ok && special5050Positions[strings.ToUpper(pos)] && special5050Kondisi[strings.ToUpper(kondisi)]
	case "50_50_KOMBINASI":
		posisi, jenis, ok := strings.Cut(nomor, "_")
		return ok && kombinasi5050Posisi[strings.ToUpper(posisi)] && kombinasi5050Jenis[strings.ToUpper(jenis)]
	case "MACAU_KOMBINASI":
		parts := strings.Split(nomor, "_")
		return len(parts) == 3 &&
			macauKombinasiPosisi[parts[0]] && macauKombinasiBesarKecil[parts[1]] && macauKombinasiGenapGanjil[parts[2]]
	case "DASAR":
		return dasarOptions[nomor]
	case "SHIO":
		return shioNames[nomor]
	}
	return false
}
