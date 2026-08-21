# Node Speed Limit Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Wire a per-node `SpeedLimit` (Mbps) from `/etc/v2node/config.json` into the existing limiter so it becomes a hard ceiling for every user on that node.

**Architecture:** `Limiter.SpeedLimit` and `determineSpeedLimit` already take `min(non-zero node, user)`. Add the config field, clamp negatives to 0, assign after `AddLimiter` without changing that function's signature. Document the fork delta for upstream merges.

**Tech Stack:** Go, viper JSON config, existing `github.com/juju/ratelimit` buckets.

---

### File map

- Create: `conf/conf_test.go`
- Create: `limiter/speedlimit_test.go`
- Modify: `conf/conf.go`
- Modify: `node/controller.go`
- Modify: `script/install.sh`, `script/v2node.sh`
- Create: `docs/custom-changes.md`, `docs/node-speed-limit.md`
- Modify: `README.md`

### Task 1: Config field tests (RED)

**Files:**
- Create: `conf/conf_test.go`

- [ ] Write failing tests that `LoadFromPath` reads `SpeedLimit`, defaults missing to 0, clamps negative to 0
- [ ] Run `go test ./conf/ -count=1` and confirm FAIL (unknown field / zero value)
- [ ] Add `SpeedLimit int` + clamp in `LoadFromPath`
- [ ] Re-run tests, expect PASS

### Task 2: Limiter ceiling tests + wiring

**Files:**
- Create: `limiter/speedlimit_test.go`
- Modify: `node/controller.go`

- [ ] Tests: node 50 / user 100 → 50 Mbps; node 50 / user 0 → 50; node 50 / user 20 → 20; same UUID shares one bucket
- [ ] Run `go test ./limiter/ -count=1`
- [ ] Assign `l.SpeedLimit = c.conf.SpeedLimit` after `AddLimiter` with `LOCAL:` comment

### Task 3: Install templates + docs

- [ ] Add `"SpeedLimit": 0` to generated JSON in `script/install.sh` and `script/v2node.sh`
- [ ] Write `docs/custom-changes.md` and `docs/node-speed-limit.md`
- [ ] Link from `README.md`
- [ ] Run `go test ./conf/ ./limiter/ -count=1` and `go build -o /tmp/v2node .`
