package utils

import (
	"fmt"
	"os"

	"github.com/SleepWalker/sinx/iface"
	"gopkg.in/yaml.v3"
)

/**
* 存储全局参数
* 一些参数通过配置文件读取，由用户进行配置
**/
type GlobalObj struct {
	//======= Sever ==========
	TcpServer iface.IServer // 全局Server对象
	Host      string        `yaml:"host"` // 服务器监听的IP
	Port      int           `yaml:"port"` // 监听的端口
	Name      string        `yaml:"name"` // 服务器名称

	//======= Sinx ==========
	Version        string `yaml:"version"`        // Sinx版本号
	MaxConn        int    `yaml:"maxConn"`        // 服务器允许的最大连接数
	MaxPackageSize uint32 `yaml:"maxPackageSize"` // 数据包的最大值
}

// 定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

/**
* 提供init方法，初始化GlobalObject
**/
func init() {
	// 默认值
	GlobalObject = &GlobalObj{
		Host:           "0.0.0.0",
		Port:           8888,
		Name:           "SinxServerApp",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	// 加载配置文件conf/sinx.yaml
	GlobalObject.LoadConfig()
}

// 从conf/sinx.yaml配置参数
func (g *GlobalObj) LoadConfig() {
	// 读取 YAML 文件
	file, err := os.ReadFile("conf/sinx.yaml")
	if err != nil {
		fmt.Println("[GlobalObj]use default config")
		return
		// panic(err)
	}

	// 解析 YAML 文件
	err = yaml.Unmarshal(file, GlobalObject)
	if err != nil {
		panic(err)
	}
}
