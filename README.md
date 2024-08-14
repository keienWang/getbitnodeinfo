# 比特币版本消息客户端

这个 Go 程序演示了如何创建一个基本的比特币客户端，它连接到比特币节点，发送版本消息，并处理响应。

## 概述

该程序执行以下操作：
1. 通过 TCP 连接到比特币节点。
2. 构造一个 `VersionMessage`，包含多个字段。
3. 将 `VersionMessage` 序列化为字节流。
4. 将序列化的版本消息发送到比特币节点。
5. 接收并处理节点的响应。

## 代码说明

### `VersionMessage` 结构体

此结构体定义了比特币版本消息的格式：

- `Version`: 协议版本
- `Services`: 发送方支持的服务
- `Timestamp`: 当前时间（Unix 时间戳）
- `AddrRecv`: 接收节点的地址
- `AddrFrom`: 发送节点的地址
- `Nonce`: 连接的随机数
- `UserAgent`: 客户端版本字符串
- `StartHeight`: 起始区块高度
- `Relay`: 继承标志

### 方法

- `Serialize() ([]byte, error)`: 将 `VersionMessage` 转换为字节流以供传输。
- `parseBitcoinResponse(response []byte)`: 解析并处理比特币节点的响应。
- `parseVersionMessage(body []byte) (*VersionMessage, error)`: 从响应体中解析 `VersionMessage`。

### 主函数

1. 连接到比特币节点 `203.11.72.110:8333`。
2. 创建一个 `VersionMessage`。
3. 序列化并发送版本消息。
4. 等待响应并打印结果。

## 使用方法

1. 确保系统上安装了 Go 环境。
2. 克隆此仓库或将代码复制到本地 Go 工作区。
3. 运行 `go run main.go` 来执行程序。

## 注意事项

- 地址 `203.11.72.110:8333` 是一个占位符，实际使用时应替换为真实的比特币节点地址。
- 程序目前仅演示了基本功能。要实现完整的功能，需要根据比特币协议处理各种类型的消息和响应。

## 许可证

//TODO
