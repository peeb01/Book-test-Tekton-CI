package middlewares

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// block IP
var blockedIPs = sync.Map{}

func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,              
		Expiration: 1 * time.Minute,  
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // use IP as Key
		},
		LimitReached: func(c *fiber.Ctx) error {
			ip := c.IP()

			// บล็อก IP ไว้ 5 นาที
			blockedIPs.Store(ip, time.Now().Add(5*time.Minute))

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too many requests",
				"message": "Your IP is blocked for 5 minutes.",
			})
		},
	})
}


func BlockMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()

		// ตรวจสอบว่า IP อยู่ในรายการบล็อก
		if unblockTime, ok := blockedIPs.Load(ip); ok {
			// ถ้าถึงเวลาปลดบล็อกแล้ว ให้ลบ IP ออก
			if time.Now().After(unblockTime.(time.Time)) {
				blockedIPs.Delete(ip)
			} else {
				// ถ้ายังถูกบล็อก ให้ส่ง Forbidden 
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error":   "Forbidden",
					"message": "Your IP is blocked. Please try again later.",
				})
			}
		}

		return c.Next()
	}
}
