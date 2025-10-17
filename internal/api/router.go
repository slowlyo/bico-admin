package api

import "github.com/gin-gonic/gin"

// Router 实现路由注册
type Router struct{}

// NewRouter 创建路由实例
func NewRouter() *Router {
	return &Router{}
}

// Register 注册路由
func (r *Router) Register(engine *gin.Engine) {
	// TODO: 实现 API 路由
}
