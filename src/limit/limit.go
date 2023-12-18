package limit

import (
	"net/http"
	"time"

	"github.com/Cosmin2410/proxy-backend-test/src/model"
)

func RateLimiter(handler http.Handler) http.Handler {
	var rateLimitMap = make(map[string]*model.RateLimitInfo)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if info, exists := rateLimitMap[ip]; exists {
			if time.Since(info.StartTime) > 10*time.Second {
				info.RequestCount = 1
				info.StartTime = time.Now()
			} else {
				info.RequestCount++
				if info.RequestCount > 10 {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)

					return
				}
			}
		} else {
			rateLimitMap[ip] = &model.RateLimitInfo{RequestCount: 1, StartTime: time.Now()}
		}

		handler.ServeHTTP(w, r)
	})
}
