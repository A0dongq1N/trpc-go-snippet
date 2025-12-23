# tRPC-Go GORM 插件超时控制指南

### 背景
- `trpc-database` 的 GORM 插件默认会忽略客户端的 `timeout` 配置，导致超时控制不符合预期。

### DSN 超时配置详解
DSN 中的超时配置由 `go-sql-driver` 底层包定义，其原理与 TCP 通信机制紧密相关：
- **`timeout`**: 控制建立 TCP 连接的超时时间，对于连接池中新建连接的场景很重要。
- **`writeTimeout`**: 作用于 `net.Conn` 的 `SetWriteDeadline`，控制将数据从应用缓冲区拷贝到内核 socket 发送缓冲区的超时时间。
- **`readTimeout`**: 作用于 `net.Conn` 的 `SetReadDeadline`，控制从内核 socket 接收缓冲区读取数据的超时时间。

### 为什么需要 Context 超时？
- **协议特性**: MySQL 使用半双工协议，客户端必须等待服务器响应才能发送下一个请求，所有操作本质上都是相同的请求模型。
- **场景局限**: 在复杂业务场景下，全局接口超时与单次 MySQL 操作超时可能不一致，DSN 超时配置无法处理全局超时问题。
- **灵活控制**: GORM 支持通过 `WithContext(ctx)` 方法传入带超时时间的 context 来控制单次 SQL 操作的超时，实现更灵活的超时管理。

### 最佳实践推荐
- 推荐在 DSN 中配置 `timeout`、`readTimeout` 和 `writeTimeout`（通常建议 **3 秒**）。
- 结合 GORM 的 `WithContext(ctx)` 方法进行全面的超时控制。
