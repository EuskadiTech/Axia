package request

import (
	"axia4/bruteforce"
	"encoding/json"
)

func BruteforceGet(reqJson json.RawMessage) (any, error) {

	var res struct {
		HostsTracked int `json:"hostsTracked"`
		HostsBlocked int `json:"hostsBlocked"`
	}
	res.HostsTracked, res.HostsBlocked = bruteforce.GetCounts()
	return res, nil
}
