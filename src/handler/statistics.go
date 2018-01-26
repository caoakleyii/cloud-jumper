package handler

import (
	"net/http"
	"time"

	"github.com/caoakleyii/cloud-jumper/src/cache"
)

// Stastic structure defines the JSON model
// to be returned by GetSastics
type Stastic struct {
	Total   int     `json:"total"`
	Average float64 `json:"average"`
}

/*
	6. Statistics End-Point

	Provide a statistics endpoint to get basic information about your password hashes.
*/

// GetStastics handler function that returns
// stastics regarding the hash requests.
// Total calls and Average duration time
func GetStastics(ctx *Context) {
	var stat Stastic
	total := len(cache.InMemoryRequestLog)

	if total > 0 {
		var duration time.Duration

		for _, v := range cache.InMemoryRequestLog {
			duration += v
		}

		average := float64(duration) / float64(total)
		ms := float64(average / float64(time.Millisecond))

		stat = Stastic{total, ms}
	} else {
		stat = Stastic{total, 0}
	}

	ctx.JSON(http.StatusOK, stat)
}
