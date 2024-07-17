package middleware

import (
	"github.com/asb1302/innopolis_go_hw11/internal/config"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	requests int
	lastSeen time.Time
}

var clients = make(map[string]*Client)
var mu sync.Mutex

func getIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

func RateLimiterMiddleware(cfg config.RateLimiterConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)
			log.Printf("Запрос от IP: %s", ip)

			mu.Lock()
			defer mu.Unlock()

			client, exists := clients[ip]
			if !exists {
				log.Printf("Новый клиент: %s", ip)

				clients[ip] = &Client{requests: 1, lastSeen: time.Now()}
			} else {
				now := time.Now()
				if now.Sub(client.lastSeen) > cfg.Duration {
					log.Printf("Сброс счетчика запросов для клиента %s из-за превышения длительности", ip)

					client.requests = 1
					client.lastSeen = now
				} else {
					client.requests++

					if client.requests > cfg.Requests {
						log.Printf("Клиент %s превысил лимит запросов", ip)

						w.WriteHeader(http.StatusTooManyRequests)
						w.Write([]byte("Слишком много запросов"))
						return
					}
				}
			}

			if client != nil {
				log.Printf("Счетчик запросов клиента %s: %d, lastSeen: %s", ip, client.requests, client.lastSeen.Format(time.RFC3339))
			} else {
				log.Printf("Клиент %s пока не инициализирован", ip)
			}

			next.ServeHTTP(w, r)
		})
	}
}
