package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type clientData struct {
	lastRequest time.Time
	requests    int
}

var (
	rateLimit      = 10          // max requests
	windowDuration = time.Minute // per minute
	clientRequests = make(map[string]*clientData)
	mu             sync.Mutex
)

func RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()

		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		data, exists := clientRequests[ip]

		if !exists || now.Sub(data.lastRequest) > windowDuration {
			clientRequests[ip] = &clientData{lastRequest: now, requests: 1}
		} else {
			if data.requests >= rateLimit {
				return c.Status(fiber.StatusTooManyRequests).Render("rate_limit", fiber.Map{
					"Title": "Rate Limit",
					"Retry": int(windowDuration.Seconds()),
				}, "layouts/main")
			}
			data.requests++
			data.lastRequest = now
		}

		return c.Next()
	}
}
