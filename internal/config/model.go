package config

type Config struct {
	Server   Server
	Database Database
	Jwt      Jwt
	Redis    Redis
	Telegram Telegram
}

type Server struct {
	Host string
	Port string
}
type Jwt struct {
	Key      string
	Exp      int
	Issuer   string
	Audience string
}
type Database struct {
	Host   string
	Port   string
	Name   string
	Schema string
	User   string
	Pass   string
	Tz     string
}

type Redis struct {
	Host string
	Port string
	Pass string
	Name string
}

// Telegram holds the bot credentials used to alert on server-side errors —
// empty values disable alerting (see util.SendTelegramAlert).
type Telegram struct {
	BotToken string
	ChatID   string
}
