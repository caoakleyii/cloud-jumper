/*
Package middleware contains handlers to be used pre and post route request.
*/
package middleware

import (
	"strings"
	"time"

	"github.com/caoakleyii/cloud-jumper/src/cache"

	"github.com/caoakleyii/cloud-jumper/src/handler"
)

var t time.Time

// PreStatistics is pre middleware that runs before the route
// to set the time for logging purposes
func PreStatistics(ctx *handler.Context) {
	t = time.Now()
}

// Statistics is middleware that runs after the route is called
// if the route is for hashing, it logs the duration
func Statistics(ctx *handler.Context) {
	if !strings.Contains(ctx.Request.URL.Path, "hash") {
		return
	}

	elap := time.Since(t)
	cache.InMemoryRequestLog[len(cache.InMemoryRequestLog)] = elap
}
