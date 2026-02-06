package config

type Config struct {
	Server   Server
	Database Database
	Jwt      Jwt
	Redis    Redis
}

type Server struct {
	Host string
	Port string
}
type Jwt struct {
	Key string
	Exp int
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
