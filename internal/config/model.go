package config

type Config struct {
	Server     Server
	Database   Database
	Jwt        Jwt
	Redis      Redis
	Telegram   Telegram
	BalanceAPI BalanceAPI
	Kafka      Kafka
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

// BalanceAPI holds the credentials for the upstream member wallet service
// that resolves a launch token into a username + live balance (see
// service.MemberinfoService.CheckToken). The URL itself is per-agent — it
// comes from tbl_company.urlapitoto (cached in Redis, see
// MemberinfoService.getCompany), not from config/env.
type BalanceAPI struct {
	APIKey string
}

// Kafka holds the checkout event bus config this worker consumes from — see
// internal/connection/kafka.go. GroupID makes this a proper consumer group
// (offsets tracked per group on the broker), so running more than one
// totclient_pools instance splits the topic's partitions between them
// instead of every instance processing every message.
type Kafka struct {
	Brokers string
	Topic   string
	GroupID string
}
