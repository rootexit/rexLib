package rexUlid

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func New() ulid.ULID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id
}

func NewString() string {
	return New().String()
}

func MonotonicNew(inc uint64) ulid.ULID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), inc)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id
}

func MonotonicNewString(inc uint64) string {
	return MonotonicNew(inc).String()
}
