package handler

import (
	"github.com/asb1302/innopolis_go_hw11/internal/config"
	"github.com/asb1302/innopolis_go_hw11/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterMiddleware(t *testing.T) {
	cfg := config.RateLimiterConfig{
		Requests: 1,
		Duration: 3 * time.Second,
	}

	handler := http.HandlerFunc(HelloHandler)

	rateLimiter := middleware.RateLimiterMiddleware(cfg)
	testHandler := rateLimiter(handler)

	clientIP := "127.0.0.1:12345"

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	req.RemoteAddr = clientIP

	// Первый запрос должен пройти
	rr := httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Обработчик вернул неверный статус код: получили %v ожидали %v", status, http.StatusOK)
	}

	// Второй запрос должен быть ограничен
	rr = httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusTooManyRequests {
		t.Errorf("Обработчик вернул неверный статус код: получили %v ожидали %v", status, http.StatusTooManyRequests)
	}

	// Ожидание для сброса лимита
	time.Sleep(cfg.Duration)

	// Запрос после сброса лимита должен пройти
	rr = httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Обработчик вернул неверный статус код: получили %v ожидали %v", status, http.StatusOK)
	}
}
