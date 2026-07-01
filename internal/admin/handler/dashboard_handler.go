package handler

import (
	"os"
	"runtime"
	"time"

	"bico-admin/internal/core/config"
	"bico-admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DashboardHandler 工作台处理器，负责汇总服务器运行状态与基础监控数据。
type DashboardHandler struct {
	cfg       *config.Config
	db        *gorm.DB
	startedAt time.Time
}

// DashboardOverview 工作台概览响应，聚合页面首屏需要的服务器、运行时和数据库状态。
type DashboardOverview struct {
	Server   DashboardServerInfo   `json:"server"`
	Runtime  DashboardRuntimeInfo  `json:"runtime"`
	Database DashboardDatabaseInfo `json:"database"`
	Monitor  DashboardMonitorInfo  `json:"monitor"`
}

// DashboardServerInfo 服务器基础信息，避免暴露账号、密码等敏感配置。
type DashboardServerInfo struct {
	Hostname      string `json:"hostname"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	GoVersion     string `json:"goVersion"`
	Mode          string `json:"mode"`
	Port          int    `json:"port"`
	StartedAt     string `json:"startedAt"`
	UptimeSeconds int64  `json:"uptimeSeconds"`
}

// DashboardRuntimeInfo Go 运行时信息，用于判断服务进程资源占用和调度规模。
type DashboardRuntimeInfo struct {
	CPUCores    int     `json:"cpuCores"`
	GOMAXPROCS  int     `json:"goMaxProcs"`
	Goroutines  int     `json:"goroutines"`
	AllocMB     float64 `json:"allocMb"`
	SysMB       float64 `json:"sysMb"`
	HeapInuseMB float64 `json:"heapInuseMb"`
	NextGCMB    float64 `json:"nextGcMb"`
	GCCycles    uint32  `json:"gcCycles"`
}

// DashboardDatabaseInfo 数据库连接池信息，只返回健康检查和容量判断需要的字段。
type DashboardDatabaseInfo struct {
	Driver              string  `json:"driver"`
	MaxOpenConnections  int     `json:"maxOpenConnections"`
	MaxIdleConnections  int     `json:"maxIdleConnections"`
	OpenConnections     int     `json:"openConnections"`
	InUse               int     `json:"inUse"`
	Idle                int     `json:"idle"`
	WaitCount           int64   `json:"waitCount"`
	WaitDurationSeconds float64 `json:"waitDurationSeconds"`
}

// DashboardMonitorInfo 监控指标列表，给前端卡片和趋势图提供统一数据源。
type DashboardMonitorInfo struct {
	CollectedAt string                   `json:"collectedAt"`
	Metrics     []DashboardMonitorMetric `json:"metrics"`
}

// DashboardMonitorMetric 单个监控指标，status 用于前端展示轻量健康状态。
type DashboardMonitorMetric struct {
	Key    string  `json:"key"`
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Unit   string  `json:"unit"`
	Status string  `json:"status"`
}

// NewDashboardHandler 创建工作台处理器。
func NewDashboardHandler(cfg *config.Config, db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{
		cfg:       cfg,
		db:        db,
		startedAt: time.Now(),
	}
}

// GetOverview 获取工作台概览数据。
// @Summary 获取工作台概览
// @Description 获取服务器、运行时、数据库和监控指标概览
// @Tags 工作台
// @Produce json
// @Security BearerAuth
// @Success 200 {object} adminResponse{data=DashboardOverview}
// @Router /dashboard/overview [get]
func (h *DashboardHandler) GetOverview(c *gin.Context) {
	overview, err := h.buildOverview()
	if err != nil {
		// 概览数据依赖数据库连接池状态，读取失败时返回统一业务错误。
		response.ErrorWithCode(c, 500, err.Error())
		return
	}
	response.SuccessWithData(c, overview)
}

// buildOverview 组装工作台概览数据，集中读取运行时快照以保持同一次响应数据一致。
func (h *DashboardHandler) buildOverview() (*DashboardOverview, error) {
	now := time.Now()
	memStats := h.readMemStats()
	serverInfo := h.buildServerInfo(now)
	runtimeInfo := h.buildRuntimeInfo(memStats)
	databaseInfo, err := h.buildDatabaseInfo()
	if err != nil {
		// 数据库连接池不可读时不拼装半完整数据，避免前端误判服务健康。
		return nil, err
	}

	return &DashboardOverview{
		Server:   serverInfo,
		Runtime:  runtimeInfo,
		Database: databaseInfo,
		Monitor:  h.buildMonitorInfo(now, runtimeInfo, databaseInfo),
	}, nil
}

// readMemStats 读取 Go 运行时内存快照，避免在多个方法里重复触发统计。
func (h *DashboardHandler) readMemStats() runtime.MemStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats
}

// buildServerInfo 构建服务器基础信息，主机名读取失败时使用可识别的占位值。
func (h *DashboardHandler) buildServerInfo(now time.Time) DashboardServerInfo {
	hostname, err := os.Hostname()
	if err != nil {
		// 某些容器环境可能禁止读取主机名，此时保留占位值便于页面识别。
		hostname = "unknown"
	}

	return DashboardServerInfo{
		Hostname:      hostname,
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		GoVersion:     runtime.Version(),
		Mode:          h.cfg.Server.Mode,
		Port:          h.cfg.Server.Port,
		StartedAt:     h.startedAt.Format(time.RFC3339),
		UptimeSeconds: int64(now.Sub(h.startedAt).Seconds()),
	}
}

// buildRuntimeInfo 构建 Go 进程运行时信息，内存统一换算成 MB 方便前端展示。
func (h *DashboardHandler) buildRuntimeInfo(memStats runtime.MemStats) DashboardRuntimeInfo {
	return DashboardRuntimeInfo{
		CPUCores:    runtime.NumCPU(),
		GOMAXPROCS:  runtime.GOMAXPROCS(0),
		Goroutines:  runtime.NumGoroutine(),
		AllocMB:     bytesToMB(memStats.Alloc),
		SysMB:       bytesToMB(memStats.Sys),
		HeapInuseMB: bytesToMB(memStats.HeapInuse),
		NextGCMB:    bytesToMB(memStats.NextGC),
		GCCycles:    memStats.NumGC,
	}
}

// buildDatabaseInfo 构建数据库连接池信息，用于判断连接池是否接近容量上限。
func (h *DashboardHandler) buildDatabaseInfo() (DashboardDatabaseInfo, error) {
	sqlDB, err := h.db.DB()
	if err != nil {
		// GORM 底层连接不可用时直接上抛，调用方负责转换为统一响应。
		return DashboardDatabaseInfo{}, err
	}
	stats := sqlDB.Stats()

	return DashboardDatabaseInfo{
		Driver:              h.cfg.Database.Driver,
		MaxOpenConnections:  h.cfg.Database.MaxOpenConns,
		MaxIdleConnections:  h.cfg.Database.MaxIdleConns,
		OpenConnections:     stats.OpenConnections,
		InUse:               stats.InUse,
		Idle:                stats.Idle,
		WaitCount:           stats.WaitCount,
		WaitDurationSeconds: stats.WaitDuration.Seconds(),
	}, nil
}

// buildMonitorInfo 将运行时和数据库核心指标转成前端可直接渲染的指标列表。
func (h *DashboardHandler) buildMonitorInfo(now time.Time, runtimeInfo DashboardRuntimeInfo, databaseInfo DashboardDatabaseInfo) DashboardMonitorInfo {
	return DashboardMonitorInfo{
		CollectedAt: now.Format(time.RFC3339),
		Metrics: []DashboardMonitorMetric{
			{
				Key:    "allocMemory",
				Label:  "内存占用",
				Value:  runtimeInfo.AllocMB,
				Unit:   "MB",
				Status: memoryStatus(runtimeInfo.AllocMB),
			},
			{
				Key:    "goroutines",
				Label:  "协程数量",
				Value:  float64(runtimeInfo.Goroutines),
				Unit:   "个",
				Status: goroutineStatus(runtimeInfo.Goroutines),
			},
			{
				Key:    "openConnections",
				Label:  "数据库连接",
				Value:  float64(databaseInfo.OpenConnections),
				Unit:   "个",
				Status: databaseStatus(databaseInfo),
			},
			{
				Key:    "gcCycles",
				Label:  "GC 次数",
				Value:  float64(runtimeInfo.GCCycles),
				Unit:   "次",
				Status: "normal",
			},
		},
	}
}

// bytesToMB 将字节数转为 MB，保留两位小数减少前端格式化负担。
func bytesToMB(value uint64) float64 {
	return float64(int((float64(value)/1024/1024)*100)) / 100
}

// memoryStatus 根据进程当前分配内存给出粗粒度状态，阈值只用于页面提示。
func memoryStatus(allocMB float64) string {
	if allocMB >= 1024 {
		// 当前只做轻量提示，超过 1GB 说明进程内存需要重点关注。
		return "warning"
	}
	// 未超过阈值时保持正常状态，避免普通波动被展示成告警。
	return "normal"
}

// goroutineStatus 根据协程数量给出粗粒度状态，避免异常膨胀时页面仍显示正常。
func goroutineStatus(count int) string {
	if count >= 10000 {
		// 协程数过高通常意味着任务堆积或泄漏，需要在工作台显著提示。
		return "warning"
	}
	// 低于阈值时只展示数量，不制造无意义告警。
	return "normal"
}

// databaseStatus 根据连接池等待和占用情况给出数据库状态。
func databaseStatus(info DashboardDatabaseInfo) string {
	if info.MaxOpenConnections > 0 && info.InUse >= info.MaxOpenConnections {
		// 连接池已满会直接影响请求吞吐，需要优先提示。
		return "warning"
	}
	if info.WaitCount > 0 {
		// 出现等待说明历史上发生过连接争用，状态标记为需关注。
		return "warning"
	}
	// 没有容量和等待风险时，数据库状态保持正常。
	return "normal"
}
