# WeCross-Go-SDK
WeCross Go SDK提供操作跨链资源的Go API，开发者通过SDK可以方便快捷地基于[WeCross](https://github.com/WeBankBlockchain/WeCross)开发自己的跨链应用。

## 引用
    import "github.com/WeBankBlockchain/WeCross-Go-SDK"

## 关键特性
- 提供了一个简单易用的日志系统方便Go程序开发，用户可自定义配置使用
- 提供调用WeCross[RPC接口](https://wecross.readthedocs.io/zh_CN/latest/docs/manual/api.html)的Go API
- 封装了跨链资源操作接口

## 快速开始
### 配置使用日志系统
```
func main() {
	// 将标准输出添加进日志系统的输出
	logger.AddStdOutLogSystem(logger.Info)
	// 自定义一个日志输出系统
	logger.AddNewLogSystem("./", "test.log", log.LstdFlags, logger.Debug)

	// 在你的go程序中定义不同的日志标签
	testLogTag := logger.NewLogger("quickStart")
	// 使用这些标签进行填写日志
	testLogTag.Infoln("Log here as you wish.")
	testLogTag.Warnf("Use the log level you like, warn level is: %d", logger.Warn)

	// 日志消息可能需要等待一定时间才能刷入标准输出， 使用flush可以强制刷新
	// 此处为了在标准输出中显示测试结果使用Flush，实际使用时无需用户自己去Flush
	logger.Flush()
}
```

### RPC API调用
```
func main() {
	// 首先创建RPC服务并设置配置文件的classpath
	// classpath下应该放置application.toml
	rpcService := service.NewWeCrossRPCService()
	rpcService.SetClassPath("./tomldir") // 如不配置classpath,默认为当前程序运行目录

	err := rpcService.Init()
	if err != nil {
		panic(err)
	}

	weCrossRPC := rpc.NewWeCrossRPCModel(rpcService)
	call, err := weCrossRPC.Login("username", "password")
	if err != nil {
		panic(err)
	}

	res, err := call.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("The response is: %s\n", res.ToString())

	// 对response更加复杂的处理,需要知道不同RPI指令返回的response data的数据类型
	// 更多RPI指令以及所对应的response data类型可查阅官方文档中的WeCross-Go-SDK说明
	data, ok := res.Data.(*response.UAReceipt)
	if !ok {
		panic("type is not right")
	}
	fmt.Printf("Universal Account info: %s\n", data.UniversalAccount.ToString())
}
```

### 资源操作接口

```
func main() {
	// 首先创建RPC服务并设置配置文件的classpath
	// classpath下应该放置application.toml
	rpcService := service.NewWeCrossRPCService()
	rpcService.SetClassPath("./tomldir") // 如不配置classpath,默认为当前程序运行目录

	err := rpcService.Init()
	if err != nil {
		panic(err)
	}

	weCrossRPC := rpc.NewWeCrossRPCModel(rpcService)

	//使用resource前需要先登录
	call, err := weCrossRPC.Login("username", "password")
	if err != nil {
		panic(err)
	}
	_, err = call.Send()
	if err != nil {
		panic(err)
	}

	// 使用NewResource绑定resource
	testResource := resource.NewResource(weCrossRPC, "test.test.test")

	// 更多资源相关API请查阅官方文档
	if !testResource.IsActive() {
		panic("the path is not active")
	}

	res, err := testResource.SendTransaction("XXXmethod", "arg1", "arg...")
	if err != nil {
		panic(err)
	}
	
	fmt.Println("Result:")
	fmt.Println(res)
}

```


## 环境依赖
```
module github.com/WeBankBlockchain/WeCross-Go-SDK

go 1.18

require github.com/pelletier/go-toml v1.9.5
```



## 贡献说明

欢迎参与WeCross社区的维护和建设：

- 提交代码(Pull requests)，可参考[代码贡献流程](./CONTRIBUTING.md)以及[wiki指南](https://github.com/WeBankBlockchain/WeCross/wiki/贡献代码)
- [提问和提交BUG](https://github.com/WeBankBlockchain/WeCross-Go-SDK/issues/new)

感谢以下贡献者的付出（此处自动显示所有本项目的代码贡献者）

<img src="https://contrib.rocks/image?repo=WeBankBlockchain/WeCross-Go-SDK" alt="https://github.com/WeBankBlockchain/github.com/WeBankBlockchain/WeCross-Go-SDK/graphs/contributors" style="zoom:100%;" />

## 开源社区

参与meetup：[活动日历](https://github.com/WeBankBlockchain/WeCross/wiki#%E6%B4%BB%E5%8A%A8%E6%97%A5%E5%8E%86)

学习知识、讨论方案、开发新特性：[联系微信小助手，加入跨链兴趣小组（CC-SIG）](https://wecross.readthedocs.io/zh_CN/latest/docs/community/cc-sig.html#id3)

## License

WeCross Go SDK的开源协议为Apache License 2.0，详情参考[LICENSE](./LICENSE)。
