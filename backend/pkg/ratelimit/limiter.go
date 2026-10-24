package ratelimit

import (
	"strconv"
	"sync"
	"time"
)

// Limiter is a thread-safe sliding-window counter keyed by string.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	max     int
	window  time.Duration
}

type bucket struct {
	timestamps []time.Time
}

func New(max int, window time.Duration) *Limiter {
	if max < 1 {
		max = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &Limiter{
		buckets: make(map[string]*bucket),
		max:     max,
		window:  window,
	}
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{}
		l.buckets[key] = b
	}

	filtered := b.timestamps[:0]
	for _, ts := range b.timestamps {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	b.timestamps = filtered

	if len(b.timestamps) >= l.max {
		return false
	}
	b.timestamps = append(b.timestamps, now)
	return true
}

func UserKey(userID uint) string {
	return "u:" + strconv.FormatUint(uint64(userID), 10)
}

// DailyQuota limits how many Allow() calls succeed per key per UTC day.
type DailyQuota struct {
	mu      sync.Mutex
	counts  map[string]dayCount
	maxPerDay int
}

type dayCount struct {
	day   string // YYYY-MM-DD UTC
	count int
}

func NewDailyQuota(maxPerDay int) *DailyQuota {
	if maxPerDay < 1 {
		maxPerDay = 1
	}
	return &DailyQuota{
		counts:    make(map[string]dayCount),
		maxPerDay: maxPerDay,
	}
}

func (q *DailyQuota) Max() int {
	return q.maxPerDay
}

func (q *DailyQuota) Allow(key string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	today := time.Now().UTC().Format("2006-01-02")
	cur, ok := q.counts[key]
	if !ok || cur.day != today {
		q.counts[key] = dayCount{day: today, count: 1}
		return true
	}
	if cur.count >= q.maxPerDay {
		return false
	}
	cur.count++
	q.counts[key] = cur
	return true
}
