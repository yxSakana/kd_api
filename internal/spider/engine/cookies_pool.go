package engine

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CookiesPool struct {
	rdb *redis.Client
}

func NewCookiesPool(addr, pw string, db int) *CookiesPool {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
	return &CookiesPool{rdb}
}

func (cp *CookiesPool) AddCookies(username string, cookies []*http.Cookie, ttl time.Duration) error {
	var cookiesStrs []string
	for _, cookie := range cookies {
		cs := fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
		cookiesStrs = append(cookiesStrs, cs)
	}
	key := fmt.Sprintf("cookies:%s", username)
	value := strings.Join(cookiesStrs, "; ")
	return cp.rdb.Set(ctx, key, value, ttl).Err()
}

func (cp *CookiesPool) GetCookies(username string) ([]*http.Cookie, error) {
	key := fmt.Sprintf("%s:%s", "cookies", username)
	cookieValue, err := cp.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var cookies []*http.Cookie
	pairs := strings.Split(cookieValue, "; ")
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			cookies = append(cookies, &http.Cookie{
				Name:  parts[0],
				Value: parts[1],
			})
		}
	}
	return cookies, nil
}

func (cp *CookiesPool) DeleteCookies(username string) error {
	return cp.rdb.Del(ctx, username).Err()
}

func (cp *CookiesPool) SyncToCookieJar(username string, jar http.CookieJar, domain string) error {
	cookies, err := cp.GetCookies(username)
	if err != nil {
		return err
	}

	u, err := url.Parse(domain)
	if err != nil {
		return err
	}
	jar.SetCookies(u, cookies)
	return nil
}

func (cp *CookiesPool) ClearAll() error {
	return cp.rdb.FlushDB(ctx).Err()
}

var ctx = context.Background()
