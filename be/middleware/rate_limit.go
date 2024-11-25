package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips   map[string]*rate.Limiter
	mu    *sync.RWMutex
	rate  rate.Limit
	burst int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:   make(map[string]*rate.Limiter),
		mu:    &sync.RWMutex{},
		rate:  r,
		burst: b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.ips[ip] = limiter
	}

	return limiter
}

// 限流中间件
func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)
		if !l.Allow() {
			c.JSON(429, gin.H{
				"code":    429,
				"message": "Too Many Requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
