# v2node 节点级限速

Date: 2026-08-21
Status: approved

## 问题

v2board 只下发用户/套餐 `speed_limit`。节点机房带宽不同时，无法在节点本地给所有人加一个硬天花板。

## 行为

在 `/etc/v2node/config.json` 的每个 Node 上增加 `SpeedLimit`（Mbps，整数）：

- `0` 或不写：不限，与现在相同
- `> 0`：该节点上所有用户的硬天花板。实际速度 = `min(非零的 节点限速, 用户套餐限速)`
- 同一 UUID 的多设备、上行+下行共用一个 token bucket
- 负值当 0
- 改配置走现有 Watch，整节点重载后生效

不在范围内：v2board 后台字段、整机总带宽桶、每台设备各自一份限速、Hysteria2 brutal 带宽另做一套。

## 实现

限速器已有 `Limiter.SpeedLimit` 和 `determineSpeedLimit`（取非零最小值），`AddLimiter` 从未赋值，始终为 0。

不改 `AddLimiter` 签名。启动后赋值：

```go
l := limiter.AddLimiter(...)
l.SpeedLimit = c.conf.SpeedLimit // LOCAL: node speed ceiling
```

改动文件：

| 文件 | 改动 |
|------|------|
| `conf/conf.go` | `NodeConfig.SpeedLimit`，加载后负值钳为 0 |
| `node/controller.go` | `AddLimiter` 后写入 `l.SpeedLimit` |
| `script/install.sh`、`script/v2node.sh` | 生成配置带 `"SpeedLimit": 0` |
| 测试 | 配置解析 + 节点天花板语义 |

每处本仓库补丁用 `LOCAL:` 注释标出。

## 文档与上游合并

- `docs/custom-changes.md`：变更表（文件、类型、合并策略）
- `docs/node-speed-limit.md`：运维说明
- `README.md`：一行链接

同步上游时先看 `docs/custom-changes.md`，冲突处重放 `LOCAL:` 补丁，不要用上游覆盖本字段。
