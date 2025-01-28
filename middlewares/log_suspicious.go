package middlewares

import (
	"log"
	"time"
	"sync"
	"github.com/gofiber/fiber/v2"
)


var requestCounts = sync.Map{}

func LogSuspiciousRequests() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		method := c.Method()
		path := c.Path()

		count, _ := requestCounts.LoadOrStore(ip, 0)
		requestCounts.Store(ip, count.(int)+1)

		// 100 req/min
		if count.(int) > 100 {
			log.Printf("Suspicious request from IP => %s | Method => %s | Path => %s\n", ip, method, path)
		}

		// reset every 1 minute
		go func(ip string) {
			time.Sleep(1 * time.Minute)
			requestCounts.Store(ip, 0)
		}(ip)

		return c.Next()
	}
}