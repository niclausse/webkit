package mocker

//import (
//	"context"
//	"github.com/redis/go-redis/v9"
//	"testing"
//	"time"
//)
//
//func TestGetRedisMocker(t *testing.T) {
//	mocker := GetRedisMocker()
//	defer mocker.Close()
//
//	if err := set(context.Background(), mocker.RedisClient, "hello", "world", 10*time.Second); err != nil {
//		t.Errorf("failed to set redis: %+v", err)
//		return
//	}
//}
//
//func set(ctx context.Context, cli *redis.Client, k, v string, expire time.Duration) error {
//	return cli.Set(ctx, k, v, expire).Err()
//}
