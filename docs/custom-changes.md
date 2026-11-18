# 二次开发修改记录

本文档记录本仓库相对上游 [wyx2685/v2node](https://github.com/wyx2685/v2node) 的差异。同步上游前先看本表，避免把本地补丁盖掉。

代码里凡本仓库加的逻辑都标了 `LOCAL:` 注释，合并冲突时搜这个标记即可重放。

## 记录规则

- 每个功能变更必须新增一条记录。
- `类型` 使用：新增、修改、删除。
- `上游影响` 使用：无、低、中、高。
- `合并策略` 说明同步上游时应保留、重放、或可丢弃。

## 变更清单

| 序号 | 日期 | 类型 | 文件 | 功能域 | 说明 | 上游影响 | 合并策略 |
|------|------|------|------|--------|------|----------|----------|
| 1 | 2026-08-21 | 新增 | `docs/custom-changes.md` | 维护规范 | 建立相对上游的二开变更记录 | 无 | 保留 |
| 2 | 2026-08-21 | 新增 | `docs/node-speed-limit.md` | 节点限速 | 节点级 `SpeedLimit` 运维说明 | 无 | 保留 |
| 3 | 2026-08-21 | 修改 | `README.md` | 节点限速 | 增加节点限速文档链接 | 低 | 保留链接；上游 README 大改时重放这一行 |
| 4 | 2026-08-21 | 修改 | `conf/conf.go` | 节点限速 | `NodeConfig` 增加 `SpeedLimit`（Mbps）；加载后负值钳为 0。搜 `LOCAL:` | 中 | 重放字段和钳制；不要改 `AddLimiter` 签名 |
| 5 | 2026-08-21 | 修改 | `node/controller.go` | 节点限速 | `AddLimiter` 之后 `l.SpeedLimit = c.conf.SpeedLimit`。搜 `LOCAL:` | 中 | 重放这一行赋值，不要并进 `AddLimiter` 参数 |
| 6 | 2026-08-21 | 修改 | `script/install.sh` | 节点限速 | 生成的 `config.json` 带 `"SpeedLimit": 0` | 低 | 重放该字段 |
| 7 | 2026-08-21 | 修改 | `script/v2node.sh` | 节点限速 | 同上，交互生成配置带 `"SpeedLimit": 0` | 低 | 重放该字段 |
| 8 | 2026-08-21 | 新增 | `conf/conf_test.go` | 节点限速 | 配置解析：有值 / 缺省 0 / 负值钳 0 | 无 | 保留 |
| 9 | 2026-08-21 | 新增 | `limiter/speedlimit_test.go` | 节点限速 | 节点天花板、用户更低时取 min、同 UUID 共享桶 | 无 | 保留 |
| 10 | 2026-08-21 | 新增 | `docs/superpowers/specs/2026-08-21-node-speed-limit-design.md` | 节点限速 | 设计说明 | 无 | 保留 |
| 11 | 2026-08-21 | 新增 | `docs/superpowers/plans/2026-08-21-node-speed-limit.md` | 节点限速 | 实现计划 | 无 | 保留 |
| 12 | 2026-08-21 | 修改 | `script/install.sh`、`script/v2node.sh` | 发布源 | `GITHUB_REPO=Anikato/v2node`，一键安装/更新从本 fork 的 Releases 拉包 | 中 | 保留变量；不要改回 wyx2685 |
| 13 | 2026-08-21 | 修改 | `README.md` | 发布源 | 一键安装 URL 改为 Anikato/v2node `main` | 低 | 保留 |
| 14 | 2026-08-21 | 修改 | `.github/workflows/release.yml` | 发布源 | push 触发增加 `main`（本 fork 默认分支） | 低 | 保留 `main` |
| 15 | 2026-11-18 | 修改 | Go 依赖和 Xray-core | 内核 | 从既有 V2Board 兼容版同步到上游 5c2360f，保留 UniProxy 和 SpeedLimit | 中 | 维护 V2Board 配置及用户同步契约 |
| 16 | 2026-11-18 | 修改 | script、README.md | 安装源 | 锁定 Anikato V2Board 发行标签，校验安装包，移除原版安装计数和防火墙菜单 | 中 | 安装和升级只读取明确的兼容标签 |
| 17 | 2026-11-18 | 修改 | 发版工作流 | 发版 | 手动可指定已发布标签，检出并上传对应安装包 | 低 | 发布源与当前分支独立 |

## Git 远程

```
origin     https://github.com/Anikato/v2node.git
upstream   https://github.com/wyx2685/v2node.git
```

兼容版从 `v2board-compat` 分支打标签。创建对应发行页后，可用 `gh workflow run release.yml --ref v2board-compat -f release_tag=v0.4.7.anikato.v2board.1` 构建并上传安装包。面板安装脚本使用固定标签和对应的包校验值。

## 上游同步检查清单

1. `git fetch upstream`（上游 `wyx2685/v2node`），`origin` 是本 fork `Anikato/v2node`
2. 合并或 rebase 前先看本表
3. 冲突文件若含 `LOCAL:`，保留本仓库那几行，再接入上游其余改动
4. 高风险文件（上游常改）：
   - `node/controller.go`：只重放 `l.SpeedLimit = c.conf.SpeedLimit`
   - `conf/conf.go`：只重放 `SpeedLimit` 字段和负值钳制
   - `limiter/limiter.go`：**本次未改**。若上游给 `Limiter.SpeedLimit` 接上了面板下发值，需要再确认与本地配置取 min，而不是丢掉本地字段
5. 合并后跑：

```bash
GOEXPERIMENT=jsonv2 go test ./conf/ ./limiter/ -count=1
GOEXPERIMENT=jsonv2 go build -o /tmp/v2node .
```

6. 已部署节点的 `/etc/v2node/config.json` 不会被 git 合并改写。新字段缺省为 0（不限速），要限速需手改后等 Watch 重载或 `v2node restart`。

## 刻意没改的上游文件

| 文件 | 原因 |
|------|------|
| `limiter/limiter.go` | `Limiter.SpeedLimit` 和 `determineSpeedLimit` 已按节点天花板工作，只是从未赋值。改签名会放大合并冲突 |
| `limiter/dynamic.go` | 已有 min(非零) 语义，无需改 |
| v2board 面板 | 限速只写在节点本地配置，不增加后台字段 |
