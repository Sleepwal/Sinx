package global

import (
	"fmt"
	"os"

	"SleepXLink/iface"

	"gopkg.in/yaml.v3"
)

// SystemConfig 存储全局参数
// 一些参数通过配置文件读取，由用户进行配置
type SystemConfig struct {
	//======= Sever ==========
	TcpServer iface.IServer // 全局Server对象
	Host      string        `yaml:"host"` // 服务器监听的IP
	Port      int           `yaml:"port"` // 监听的端口
	Name      string        `yaml:"name"` // 服务器名称

	//======= Connection ==========
	MaxMsgChanLen uint32 `yaml:"maxMsgChanLen"` // 消息缓冲最大长度

	//======= SleepXLink ==========
	Version          string `yaml:"version"`          // Sinx版本号
	MaxConn          int    `yaml:"maxConn"`          // 服务器允许的最大连接数
	MaxPackageSize   uint32 `yaml:"maxPackageSize"`   // 数据包的最大值
	WorkerPoolSize   uint32 `yaml:"workerPoolSize"`   // 业务工作worker池的Goroutine数量
	MaxWorkerTaskLen uint32 `yaml:"maxWorkerTaskLen"` // 允许用户最多开辟多少个worker（限定条件）
}

// InitSystemConfig 提供init方法，初始化SystemConfig
func InitSystemConfig() {
	// 默认值
	SXL_CONFIG = &SystemConfig{
		Host:             "0.0.0.0",
		Port:             8888,
		Name:             "SleepXLinkServerApp",
		MaxMsgChanLen:    10,
		Version:          "V1.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 加载配置文件conf/config.yaml
	SXL_CONFIG.LoadConfig()
}

// LoadConfig 从config.yaml读取配置参数
func (sg *SystemConfig) LoadConfig() {
	// 读取 YAML 文件
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("[GlobalObj]use default config")
		return
		//panic(err)
	}

	// 解析 YAML 文件
	err = yaml.Unmarshal(file, SXL_CONFIG)
	if err != nil {
		panic(err)
	}
}
