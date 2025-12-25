package secure_tools

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type clientInfo struct {
	timestamps []time.Time
}

const (
	maxRequests = 20
	window      = 10 * time.Minute
)

var ipRequests sync.Map

func RateLimiter(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !utils.GetBoolValue(utils.Config.RATE_LIMITING) {
			log.Printf("[WEB]: Rate limiting skipped\n")
			next.ServeHTTP(w, r)
			return
		}
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}

		now := time.Now()

		val, _ := ipRequests.LoadOrStore(ip, &clientInfo{})
		info := val.(*clientInfo)

		// filter out old timestamps older than the window
		filtered := info.timestamps[:0]
		for _, t := range info.timestamps {
			if now.Sub(t) < window {
				filtered = append(filtered, t)
			}
		}
		info.timestamps = filtered

		// Check if too many requests
		if len(info.timestamps) >= maxRequests {
			log.Printf("Rate limit exceeded for IP %s", ip)
			http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
			return
		}

		// Add current timeStamp & store back
		info.timestamps = append(info.timestamps, now)
		ipRequests.Store(ip, info)

		next.ServeHTTP(w, r)
	}
}
