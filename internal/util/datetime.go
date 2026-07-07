package util

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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

func PgTimeToString(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}
	total := t.Microseconds / 1e6
	return fmt.Sprintf("%02d:%02d:%02d", total/3600, (total%3600)/60, total%60)
}
func StringToPgTime(s string) pgtype.Time {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		if t, err = time.Parse("15:04", s); err != nil {
			return pgtype.Time{} // invalid → Valid: false
		}
	}
	us := int64(t.Hour())*3600e6 + int64(t.Minute())*60e6 + int64(t.Second())*1e6
	return pgtype.Time{Microseconds: us, Valid: true}
}
