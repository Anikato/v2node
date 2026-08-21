# 节点级限速

本节点所有用户的速度天花板，写在 v2node 本地配置里，不经过 v2board。套餐限速再高，也超不过这里的值。

## 配置

编辑 `/etc/v2node/config.json`，在对应 Node 上设置 `SpeedLimit`（单位 Mbps）：

```json
{
    "Log": {
        "Level": "warning",
        "Output": "",
        "Access": "none"
    },
    "Nodes": [
        {
            "ApiHost": "https://panel.example.com/",
            "NodeID": 1,
            "ApiKey": "通讯密钥",
            "Timeout": 15,
            "SpeedLimit": 50
        }
    ]
}
```

| 值 | 含义 |
|----|------|
| `0` 或不写 | 不限。只按 v2board 下发的用户/套餐 `speed_limit` |
| `> 0` | 该节点硬天花板（Mbps） |
| `< 0` | 当作 0 |

保存后现有 Watch 会整节点重载，连接会断一下。也可 `v2node restart`。

一键安装生成的配置默认 `"SpeedLimit": 0`。已在跑的旧配置没有这个字段时，行为与不限速相同。

## 和面板限速的关系

实际速度 = `min(非零的 节点 SpeedLimit, 用户 speed_limit)`。

| 节点 SpeedLimit | 用户套餐 | 实际 |
|-----------------|----------|------|
| 50 | 100 | 50 |
| 50 | 不限 (0) | 50 |
| 50 | 20 | 20 |
| 0 / 不写 | 100 | 100 |

同一 UUID 的多台设备、上行和下载，共用这一份限额。两台设备同时跑 50 Mbps 节点，合计最高 50，不是各 50。不同 UUID 互不影响。

这不是整机出口总带宽。100 个不限速用户连 `SpeedLimit: 50` 的节点，每人都能到 50。

Token bucket 容量约 1 秒的量，瞬时可能略超，按秒平均不会超过配置值。

## 一键安装

v2board 后台复制的命令会安装本 fork（[Anikato/v2node](https://github.com/Anikato/v2node)），不是上游官方包。装完后在配置里改 `SpeedLimit` 即可，不必重装。

面板代码要已部署带 Anikato 安装地址的 v2board，节点机器也要能访问 GitHub。

## 实现位置（给合并用）

限速真正卡在 dispatcher 的 token bucket。本地只做一件事：把配置写进已经存在、但从未赋值的 `Limiter.SpeedLimit`。

- 配置字段：`conf/conf.go`（`LOCAL:`）
- 赋值：`node/controller.go` 里 `AddLimiter` 之后（`LOCAL:`）
- 取 min：上游 `limiter/determineSpeedLimit`，本仓库未改

同步上游步骤见 [custom-changes.md](./custom-changes.md)。
