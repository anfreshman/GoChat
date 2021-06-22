package tinterface

// 服务器接口
type IServer interface {
	// 启动Server
	Start()
	// 停止Server
	Stop()
	// 运行Server，调用启动Server，且可以执行一些其他业务
	Serve()
}
