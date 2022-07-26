package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer //全局Server对象
	Host      string         //IP
	TcpPort   int            //端口号
	Name      string         //服务器名称

	Version        string //当前版本号
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前zinx框架数据包最大值

	WorkerPoolSize    uint32 //当前worker池的goroutine数量
	MaxWorkerTaskSize uint32 //最大处理消息数量
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)

	}
}

//提供init方法,初始化当前的对象
func init() {
	GlobalObject = &GlobalObj{
		Host:              "0.0.0.0",
		TcpPort:           8999,
		Name:              "ZinxServerApp",
		Version:           "V0.6",
		MaxConn:           1000,
		MaxPackageSize:    4096,
		WorkerPoolSize:    10,
		MaxWorkerTaskSize: 1024,
	}

	//尝试从conf/zinx.json去加载
	GlobalObject.Reload()
}
