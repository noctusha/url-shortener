package ratelimit

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	rdb    *redis.Client
	prefix string
	limit  int
	window time.Duration
}

func NewLimiter(rdb *redis.Client, prefix string, limit int, window time.Duration) *Limiter {
	return &Limiter{
		rdb:    rdb,
		prefix: prefix,
		limit:  limit,
		window: window,
	}
}

func (l *Limiter) MiddleWare(keyFn func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := l.prefix + ":" + keyFn(r)
			allowed, retryAfter, err := l.allow(r.Context(), key)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (l *Limiter) allow(ctx context.Context, key string) (bool, time.Duration, error) {
	pipe := l.rdb.TxPipeline()

	// INCR — увеличивает значение ключа на 1 и возвращает текущее значение (count запросов).
	incr := pipe.Incr(ctx, key)
	// TTL  —  сколько осталось до автоудаления ключа Redisом.
	ttl := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	// сколько запросов уже сделано в текущем окне
	count := incr.Val()
	// сколько осталось до удаления ключа (т.е. до конца текущего окна лимита)
	t := ttl.Val()
	// установка TTL, если он еще не установлен, либо некорректен
	if t <= 0 {
		if err := l.rdb.Expire(ctx, key, l.window).Err(); err != nil {
			return false, 0, err
		}
		t = l.window
	}

	if int(count) > l.limit {
		return false, t, nil
	}

	return true, t, nil
}

func ClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}
