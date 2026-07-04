package util

import (
	"database/sql"
	"math/rand"
	"strconv"
	s "strings"
	"time"

	"github.com/devhdn-212/totmaster_api/internal/config"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(bytes), err
}
func Encryption(datatext string) (string, int) {
	min := 0
	max := 149
	rand.Seed(time.Now().UnixNano())
	// keymap := rand.Intn(max-min) + min
	keymap := rand.Intn(max-min) + min
	var key string = config.Keymap[keymap]
	var source string = config.Sourcechar
	result := ""
	for i := 0; i < len(datatext); i++ {
		temp_indexsource := s.Index(source, string(datatext[i]))
		temp_indexkey := s.Index(key, string(key[temp_indexsource]))
		result += string(key[temp_indexkey])
	}
	return result, keymap
}
func Decryption(dataencrypt string) string {
	temp := s.Split(dataencrypt, "|")
	keymap, _ := strconv.Atoi(temp[1])
	var key string = config.Keymap[keymap]
	var source string = config.Sourcechar
	result := ""
	for i := 0; i < len(temp[0]); i++ {
		temp_indexkey := s.Index(key, string(dataencrypt[i]))
		temp_indexsource := s.Index(source, string(source[temp_indexkey]))
		result += string(source[temp_indexsource])
	}
	return result
}

func Parsing_final(datatoken string) string {
	temp_decp := Decryption(datatoken)
	return temp_decp
}

func NsToStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// Helper convert sql.NullTime → string dengan prefix (misal nama admin/creator)
func NtToStr(nt sql.NullTime, prefix string, loc *time.Location) string {
	if nt.Valid {
		if prefix != "" {
			return prefix + ", " + nt.Time.In(loc).Format("2006-01-02 15:04:05")
		}
		return nt.Time.In(loc).Format("2006-01-02 15:04:05")
	}
	return ""
}
