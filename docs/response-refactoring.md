# 统一响应说明

## 目标

后端统一使用 `internal/pkg/response/response.go` 输出响应，避免各处手写 `c.JSON`。

## 标准结构

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

## 常用函数

- `SuccessWithData(c, data)`
- `SuccessWithMessage(c, msg, data)`
- `ErrorWithCode(c, code, msg)`
- `BadRequest(c, msg)`
- `NotFound(c, msg)`
- `TooManyRequests(c, msg)`
- `SuccessWithPagination(c, list, total)`

## HTTP 状态码约定（当前实现）

1. 成功：HTTP `200`
2. 业务错误（`ErrorWithCode`）：HTTP `200`
3. 参数错误（`BadRequest`）：HTTP `400`
4. 限流（`TooManyRequests`）：HTTP `429`
5. `NotFound`：当前实现为 HTTP `200` + `code=404`

## 典型示例

```go
if err := c.ShouldBindJSON(&req); err != nil {
    response.BadRequest(c, "参数错误: "+err.Error())
    return
}

if err := svc.Do(); err != nil {
    response.ErrorWithCode(c, 400, err.Error())
    return
}

response.SuccessWithMessage(c, "保存成功", data)
```

## 分页示例

```go
response.SuccessWithPagination(c, list, total)
```

返回：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [],
    "total": 0
  }
}
```

## 建议

1. 新接口默认使用 `response` 包，不再新增手写响应格式。
2. 若要调整全局响应规范，只改 `response` 包即可。
