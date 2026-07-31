package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func newTestRouter(rps, burst int, ttl time.Duration) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Limit(rps, burst, ttl))
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return router
}

func doRequest(router http.Handler, ip string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = ip
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestLimitExceedsBurst(t *testing.T) {
	defer StopAll()

	router := newTestRouter(1, 2, time.Minute)

	for i := 0; i < 2; i++ {
		if w := doRequest(router, "192.168.1.10:1234"); w.Code != http.StatusOK {
			t.Fatalf("request %d: got status %d, want %d", i+1, w.Code, http.StatusOK)
		}
	}
	if w := doRequest(router, "192.168.1.10:1234"); w.Code != http.StatusTooManyRequests {
		t.Fatalf("request 3: got status %d, want %d", w.Code, http.StatusTooManyRequests)
	}
}

func TestLimitSeparateBucketsPerIP(t *testing.T) {
	defer StopAll()

	router := newTestRouter(1, 1, time.Minute)

	// исчерпываем bucket первого IP
	if w := doRequest(router, "10.0.0.1:1234"); w.Code != http.StatusOK {
		t.Fatalf("first request: got status %d, want %d", w.Code, http.StatusOK)
	}
	if w := doRequest(router, "10.0.0.1:1234"); w.Code != http.StatusTooManyRequests {
		t.Fatalf("second request same IP: got status %d, want %d", w.Code, http.StatusTooManyRequests)
	}

	// другой IP не должен быть ограничен
	if w := doRequest(router, "10.0.0.2:1234"); w.Code != http.StatusOK {
		t.Fatalf("request from other IP: got status %d, want %d", w.Code, http.StatusOK)
	}
}

func TestLimitSetsRetryAfter(t *testing.T) {
	defer StopAll()

	router := newTestRouter(1, 0, time.Minute)

	w := doRequest(router, "10.0.0.3:1234")
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("got status %d, want %d", w.Code, http.StatusTooManyRequests)
	}
	if got := w.Header().Get("Retry-After"); got != "60" {
		t.Fatalf("Retry-After = %q, want %q", got, "60")
	}
}

func TestStopAllDoesNotPanic(t *testing.T) {
	router := newTestRouter(1, 1, time.Minute)
	doRequest(router, "10.0.0.4:1234")
	StopAll()
	StopAll() // повторный вызов не должен паниковать
}
