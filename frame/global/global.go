package global

import (
	"log/slog"
)

/****************************************
@Author : SleepWalker
@Description: 全局变量
@File : global
@Time : 2024/1/2 18:12
****************************************/

var (
	// SXL_LOG 日志
	SXL_LOG *slog.Logger
	// SXL_CONFIG 系统配置
	SXL_CONFIG *SystemConfig
)

func init() {
	InitLog()
	InitSystemConfig()
}

func InitLog() {
	var logger = slog.Default()
	SXL_LOG = logger
	SXL_LOG.Info("Log initialized")
}
