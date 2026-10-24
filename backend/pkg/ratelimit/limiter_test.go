package ratelimit

import (
	"testing"
	"time"
)

func TestLimiterAllowBurst(t *testing.T) {
	lim := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !lim.Allow("a") {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	if lim.Allow("a") {
		t.Fatal("4th request should be denied")
	}
	if !lim.Allow("b") {
		t.Fatal("other key should still be allowed")
	}
}

func TestLimiterWindowReset(t *testing.T) {
	lim := New(1, 40*time.Millisecond)
	if !lim.Allow("x") {
		t.Fatal("first should allow")
	}
	if lim.Allow("x") {
		t.Fatal("second should deny")
	}
	time.Sleep(50 * time.Millisecond)
	if !lim.Allow("x") {
		t.Fatal("after window should allow again")
	}
}

func TestDailyQuota(t *testing.T) {
	q := NewDailyQuota(2)
	if !q.Allow("u:1") || !q.Allow("u:1") {
		t.Fatal("first two should allow")
	}
	if q.Allow("u:1") {
		t.Fatal("third should deny")
	}
	if !q.Allow("u:2") {
		t.Fatal("other user should allow")
	}
	if q.Max() != 2 {
		t.Fatalf("Max() = %d", q.Max())
	}
}
