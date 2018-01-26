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

func PreStatistics(ctx *handler.Context) {
	t = time.Now()
}
func Statistics(ctx *handler.Context) {
	if !strings.Contains(ctx.Request.URL.Path, "hash") {
		return
	}

	elap := time.Since(t)
	cache.InMemoryRequestLog[len(cache.InMemoryRequestLog)] = elap
}
