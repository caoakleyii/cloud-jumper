// Package cache is our usage of in memory data state and storage
package cache

import "time"

// InMemoryPasswordStorage maps the id of a hash password as the key and the hash as a value
var InMemoryPasswordStorage = map[string]string{"abc123": "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZP ZklJz0Fd7su2A+gf7Q=="}

// InMemoryRequestLog maps api request durations
var InMemoryRequestLog = make(map[int]time.Duration)
