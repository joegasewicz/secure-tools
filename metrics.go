package secure_tools

import (
	"log"
	"net/http"
	"sync"
	"time"

	identity_client "github.com/joegasewicz/identity-client"
)

// Only use "start" key
var requestTimeMap = make(map[string]time.Time)
var requestTimeMutex sync.Mutex

var domain = "..."    // TODO
var requestURL = ".." // TODO

func RequestTime(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		requestTimeMutex.Lock()
		defer requestTimeMutex.Unlock()

		firstTime, ok := requestTimeMap["start"]

		if !ok {
			// first time ever
			requestTimeMap["start"] = start
			go storeFirstRequestOfDay(r, duration)
			return
		}

		if time.Since(firstTime) >= 24*time.Hour {
			// 24 hours has elapsed so add to metrics
			go storeFirstRequestOfDay(r, duration)
			// start a new 24-hour timer
			requestTimeMap["start"] = start
		}
	}
}

func storeFirstRequestOfDay(r *http.Request, duration time.Duration) {
	// Send metrics
	log.Println("[WEB]: sending request metrics")
	apiURL := requestURL
	log.Printf("[WEB]: Making request on %s\n", apiURL)

	token, err := utils.GetBearerTokenForAdmin(r) // TODO
	if err != nil {
		log.Printf("[WEB] %s\n", err.Error())
		return
	}

	identity := identity_client.Identity{
		URL:    apiURL,
		Client: nil,
		Token:  identity_client.AddTokenPrefix(token),
	}
	jsonData := map[string]interface{}{
		"RequestTime":     time.Now().Format(time.RFC3339),
		"RequestDuration": duration.Seconds(),
		"Domain":          domain,
	}
	data, err := identity.Post(r, jsonData)
	if err != nil {
		log.Println("[WEB]: error sending metrics " + err.Error())
		return
	}
	log.Printf("[WEB]: successfully created metrics: %s\n", data)
}
