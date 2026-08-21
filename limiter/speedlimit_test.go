package limiter

import (
	"testing"

	panel "github.com/wyx2685/v2node/api/v2board"
	"github.com/wyx2685/v2node/common/format"
)

func mbpsToBytes(mbps int) float64 {
	return float64(mbps) * 1000000 / 8
}

func newTestLimiter(t *testing.T, nodeMbps, userMbps int) (*Limiter, string) {
	t.Helper()
	Init()
	tag := "test-node"
	uuid := "user-uuid"
	users := []panel.UserInfo{{
		Id:         1,
		Uuid:       uuid,
		SpeedLimit: userMbps,
	}}
	l := AddLimiter("vmess", tag, users, map[int]int{})
	l.SpeedLimit = nodeMbps
	return l, format.UserTag(tag, uuid)
}

func TestCheckLimit_NodeSpeedLimitCapsHigherUserLimit(t *testing.T) {
	l, key := newTestLimiter(t, 50, 100)
	bucket, reject := l.CheckLimit(key, "1.2.3.4", true)
	if reject {
		t.Fatal("unexpected reject")
	}
	if bucket == nil {
		t.Fatal("expected rate bucket")
	}
	got := bucket.Get().Rate()
	want := mbpsToBytes(50)
	if got != want {
		t.Fatalf("rate = %v, want %v (50 Mbps ceiling over 100 Mbps user)", got, want)
	}
}

func TestCheckLimit_NodeSpeedLimitAppliesWhenUserUnlimited(t *testing.T) {
	l, key := newTestLimiter(t, 50, 0)
	bucket, reject := l.CheckLimit(key, "1.2.3.4", true)
	if reject {
		t.Fatal("unexpected reject")
	}
	if bucket == nil {
		t.Fatal("expected rate bucket")
	}
	got := bucket.Get().Rate()
	want := mbpsToBytes(50)
	if got != want {
		t.Fatalf("rate = %v, want %v", got, want)
	}
}

func TestCheckLimit_LowerUserLimitWinsOverNode(t *testing.T) {
	l, key := newTestLimiter(t, 50, 20)
	bucket, reject := l.CheckLimit(key, "1.2.3.4", true)
	if reject {
		t.Fatal("unexpected reject")
	}
	if bucket == nil {
		t.Fatal("expected rate bucket")
	}
	got := bucket.Get().Rate()
	want := mbpsToBytes(20)
	if got != want {
		t.Fatalf("rate = %v, want %v", got, want)
	}
}

func TestCheckLimit_SameUUIDSharesOneBucket(t *testing.T) {
	l, key := newTestLimiter(t, 50, 0)
	b1, reject1 := l.CheckLimit(key, "1.2.3.4", true)
	b2, reject2 := l.CheckLimit(key, "5.6.7.8", true)
	if reject1 || reject2 {
		t.Fatal("unexpected reject")
	}
	if b1 == nil || b2 == nil {
		t.Fatal("expected buckets")
	}
	if b1 != b2 {
		t.Fatal("same UUID must share one speed bucket across devices")
	}
}

func TestCheckLimit_ZeroNodeAndUserMeansUnlimited(t *testing.T) {
	l, key := newTestLimiter(t, 0, 0)
	bucket, reject := l.CheckLimit(key, "1.2.3.4", true)
	if reject {
		t.Fatal("unexpected reject")
	}
	if bucket != nil {
		t.Fatal("expected no rate bucket when both limits are 0")
	}
}
