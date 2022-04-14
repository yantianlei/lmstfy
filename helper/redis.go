package helper

import (
	pool "github.com/bitleak/lmstfy/go-redis-pool"
	"strings"

	"github.com/bitleak/lmstfy/config"
	"github.com/go-redis/redis"
)
//NewRedisPool return a redis pool
func NewRedisPool(cfg *config.RedisConf,opt *redis.Options) *pool.Pool {
	addrs := strings.Split(cfg.Addr, ",")
	p, _ := pool.NewHA(&pool.HAConfig{
		Master: addrs[0],
		Slaves: addrs[1:],
		PollType: pool.PollByWeight,
	})
	return p
}
// NewRedisClient wrap the standalone and sentinel client
func NewRedisClient(conf *config.RedisConf, opt *redis.Options) *redis.Client {
	if opt == nil {
		opt = &redis.Options{}
	}
	opt.Addr = conf.Addr
	opt.Password = conf.Password
	opt.PoolSize = conf.PoolSize
	opt.DB = conf.DB
	if conf.IsSentinel() {
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    conf.MasterName,
			SentinelAddrs: strings.Split(opt.Addr, ","),
			Password:      opt.Password,
			PoolSize:      opt.PoolSize,
			ReadTimeout:   opt.ReadTimeout,
			WriteTimeout:  opt.WriteTimeout,
			MinIdleConns:  opt.MinIdleConns,
			DB:            opt.DB,
		})
	}
	return redis.NewClient(opt)
}
