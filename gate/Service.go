package gate

import (
	"github.com/ssgo/hub/dock"
	"github.com/ssgo/log"
	"github.com/ssgo/redis"
	"github.com/ssgo/s"
	"github.com/ssgo/u"
)

var redisPool *redis.Redis

type gateConfig struct {
	Proxies  map[string]string
	Rewrites map[string]string
	Prefix   string
}

var logger = log.New(u.ShortUniqueId())
var proxiesKey string

func Registers() {
	//s.SetAuthChecker(dock.Auth)
	redisPool = redis.GetRedis(dock.GetDiscover(), logger)
	s.Restful(1, "GET", "/gateway", getGatewayInfo)
	s.Restful(2, "POST", "/gateway", setGatewayInfo)
}

func getPrefix() string {
	proxyKeys := redisPool.KEYS("_*proxies")
	if len(proxyKeys) < 1 {
		return "_"
	}
	proxiesKey := proxyKeys[0]
	if len(proxyKeys) <= 7 {
		return "_"
	}
	prefix := proxiesKey[0 : len(proxyKeys)-7]
	return prefix
}

func getGatewayInfo() (gatewayConfig gateConfig) {
	prefix := getPrefix()
	gatewayConfig.Proxies = redisPool.Do("HGETALL", prefix+"proxies").StringMap()
	gatewayConfig.Rewrites = redisPool.Do("HGETALL", prefix+"rewrites").StringMap()
	gatewayConfig.Prefix = prefix
	return
}

func setGatewayInfo(gatewayConfig gateConfig) bool {
	prefix := getPrefix()
	newProxies := gatewayConfig.Proxies
	newRewrites := gatewayConfig.Rewrites

	return saveMulti(prefix, "proxies", newProxies) && saveMulti(prefix, "rewrites", newRewrites)
}

func saveMulti(prefix string, key string, newList map[string]string) bool {
	oldList := redisPool.Do("HGETALL", prefix+key).StringMap()
	currentKey := prefix + key
	for index, single := range newList {
		if !redisPool.HSET(currentKey, index, single) {
			return false
		}
	}
	for index, _ := range oldList {
		_, ok := newList[index]
		if !ok {
			redisPool.HDEL(currentKey, index)
		}
	}
	return true
}