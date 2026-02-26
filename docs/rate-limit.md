# 限流中间件

## 当前实现

项目使用令牌桶限流（`golang.org/x/time/rate`），默认按客户端 IP 做全局限流。

核心代码：`internal/core/middleware/rate_limit.go`。

## 配置

```yaml
rate_limit:
  enabled: true
  rps: 100
  burst: 200
```

## 生效方式

启动时在 `BuildContext` 中创建限流器，并在 `server.NewServer` 中全局挂载。

- `enabled=false` 时，会创建一个近似无限流的限流器
- `enabled=true` 时，按 `rps/burst` 限流

## 触发响应

限流时返回：

- HTTP 状态码：`429`
- 响应体：`{"code":429,"msg":"请求过于频繁，请稍后再试"}`

## 压测示例

```bash
ab -n 1000 -c 100 http://localhost:8080/health
wrk -t 10 -c 100 -d 30s http://localhost:8080/health
```

## 注意事项

1. 当前项目默认只使用 `RateLimit()`（IP 维度）。
2. 配置文件变化不会自动重建限流器，修改 `rate_limit` 后需重启服务。
3. 反向代理场景需正确传递真实客户端 IP。
