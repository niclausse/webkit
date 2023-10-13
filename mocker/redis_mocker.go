package mocker

//import (
//	"github.com/alicebob/miniredis/v2"
//	"github.com/redis/go-redis/v9"
//	"log"
//	"sync"
//)
//
//var (
//	redisMocker *RedisMocker
//	redisOnce   sync.Once
//)
//
//type RedisMocker struct {
//	RedisClient *redis.Client
//	redisServer *miniredis.Miniredis
//}
//
//func (r *RedisMocker) Close() {
//	_ = r.RedisClient.Close()
//	r.redisServer.Close()
//}
//
//func GetRedisMocker() *RedisMocker {
//	redisOnce.Do(func() {
//		if redisMocker != nil {
//			return
//		}
//
//		s, err := miniredis.Run()
//		if err != nil {
//			log.Fatalf("failed to start redis server: %+v", err)
//		}
//
//		cli := redis.NewClient(&redis.Options{
//			Network:               "",
//			Addr:                  s.Addr(),
//			ClientName:            "",
//			Dialer:                nil,
//			OnConnect:             nil,
//			Protocol:              0,
//			Username:              "",
//			Password:              "",
//			CredentialsProvider:   nil,
//			DB:                    0,
//			MaxRetries:            0,
//			MinRetryBackoff:       0,
//			MaxRetryBackoff:       0,
//			DialTimeout:           0,
//			ReadTimeout:           0,
//			WriteTimeout:          0,
//			ContextTimeoutEnabled: false,
//			PoolFIFO:              false,
//			PoolSize:              0,
//			PoolTimeout:           0,
//			MinIdleConns:          0,
//			MaxIdleConns:          0,
//			ConnMaxIdleTime:       0,
//			ConnMaxLifetime:       0,
//			TLSConfig:             nil,
//			Limiter:               nil,
//		})
//
//		redisMocker = &RedisMocker{
//			RedisClient: cli,
//			redisServer: s,
//		}
//	})
//
//	return redisMocker
//}
