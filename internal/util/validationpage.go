package util

import (
	"encoding/json"
	"strings"

	"github.com/devhdn-212/totclient_pools/internal/connection"
)

const (
	RedisAuth = "agen:client:"
)

func Validpage(username, page string) bool {
	RedisAuthPage := RedisAuth + username
	type Authredis struct {
		Username string `json:"username"`
		IDUle    string `json:"idule"`
		Rule     string `json:"rule"`
	}
	cached, found, err := connection.GetRedis(RedisAuthPage)
	if err != nil {
		return false
	}
	if found {
		var auth Authredis
		if err := json.Unmarshal([]byte(cached), &auth); err == nil {
			ruleMap := ParseRules(auth.Rule)
			return ruleMap[page]
		}
	}
	return found
}
func ParseRules(ruleStr string) map[string]bool {
	rules := make(map[string]bool)

	for _, r := range strings.Split(ruleStr, ",") {
		rules[strings.TrimSpace(r)] = true
	}

	return rules
}
func GetDataRedisClient(username string) (bool, string, string) {
	RedisAuthPage := RedisAuth + username
	type Authredis struct {
		Username string `json:"username"`
		ID       string `json:"id"`
		IDComp   string `json:"idcomp"`
		IDrule   string `json:"idrule"`
		Rule     string `json:"rule"`
	}
	cached, found, err := connection.GetRedis(RedisAuthPage, 0)
	if err != nil {
		return false, "", ""
	}
	if found {
		var auth Authredis
		if err := json.Unmarshal([]byte(cached), &auth); err == nil {
			return found, auth.ID, auth.IDComp
		}
	}
	return found, "", ""
}
