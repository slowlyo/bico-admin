package crud

import "errors"

// BizError 业务错误，消息可以安全地原样展示给客户端。
//
// 说明：CRUD 框架只会把 BizError 的消息返回给前端；其余错误（数据库驱动错误等）
// 可能带有 SQL、表名或约束名，统一记日志后返回通用提示。
type BizError struct {
	Msg string
}

// Error 实现 error 接口。
func (e *BizError) Error() string {
	return e.Msg
}

// NewBizError 创建业务错误。
//
// 说明：返回 error 接口，可直接作为哨兵错误使用，errors.Is 按指针比较。
func NewBizError(msg string) error {
	return &BizError{Msg: msg}
}

// bizMessage 提取错误链中的业务错误消息，不是业务错误时返回 false。
func bizMessage(err error) (string, bool) {
	var biz *BizError
	if errors.As(err, &biz) {
		return biz.Msg, true
	}
	return "", false
}
