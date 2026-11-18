# v2node for V2Board

基于 [wyx2685/v2node](https://github.com/wyx2685/v2node) 和修改版 Xray-core 的 V2Board 节点。节点通过 UniProxy 使用现有 ApiHost、NodeID 和 ApiKey。支持节点级 SpeedLimit，详见 [限速说明](docs/node-speed-limit.md)。

## 安装和更新

优先从 V2Board 后台复制该节点的一键命令。在 Linux x86-64 或 ARM64 上也可下载此发行标签的脚本：

```bash
curl -fLSs -o install.sh https://raw.githubusercontent.com/Anikato/v2node/v0.4.7.anikato.v2board.1/script/install.sh
bash install.sh
```

脚本使用此标签发布的兼容包，检查 SHA-256 后才替换二进制。已有 V2Board 配置会保留。安装后运行 `v2node update` 使用同一版本，运行 `v2node version` 查看节点版本。

同步和发版记录见 [定制记录](docs/custom-changes.md)。
