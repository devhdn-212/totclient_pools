package util

import (
	"encoding/json"
	"gofibergocu/internal/connection"
	"strings"
)

func Validpage(username, page string) bool {
	RedisAuthPage := "client:" + username
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
