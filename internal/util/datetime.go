package util

import (
	"time"
)

var (
	// Lokasi Jakarta didefinisikan secara global agar bisa dipakai di mana saja
	LocJakarta *time.Location
)

func init() {
	var err error
	LocJakarta, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback jika timezone Asia/Jakarta tidak ditemukan di OS
		// Biasanya pakai offset +7 jam
		LocJakarta = time.FixedZone("Asia/Jakarta", 7*3600)
	}
}

// GetNowJakarta adalah helper untuk mendapatkan waktu saat ini langsung dalam zona Jakarta
func GetNowJakarta() time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now()
	}
	// Kirim langsung, pgx akan strip timezone info, simpan angka mentah
	return time.Now().In(loc)
}

// FormatToJakarta membantu merubah time.Time mentah ke zona Jakarta
func FormatToJakarta(t time.Time) time.Time {
	return t.In(LocJakarta)
}
