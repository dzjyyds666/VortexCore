## VortexCore - 多协议统一网关型服务器框架

***“开箱即用、多协议支持、可配置、开发者不需要手动管理依赖或协议接入”。***

```aiignore
 __     __                 _                 
 \ \   / /   ___    _ __  | |_    ___  __  __
  \ \ / /   / _ \  | '__| | __|  / _ \ \ \/ /
   \ V /   | (_) | | |    | |_  |  __/  >  < 
    \_/     \___/  |_|     \__|  \___| /_/\_\
```

项目旨在可以快速构建一个支持多种协议的服务器，供项目使用。

### 一、支持协议

- [x] Http
- [ ] WebSocket
- [ ] MQTT
- [ ] grpc
- [ ] .........

### 二、快速开始

``` Go
 vortex := NewVortexCore(context.Background(),
			WithListenPort("18080"),
			WithDefaultLogger(),
			WithTransport(Transport.TCP),
			WithProtocol(vortexUtil.Http1, vortexUtil.WebSocket),
		)
		vortex.Start()
```

代码解析：这段代码创建了一个简单的服务，指定端口为`18080`,并且使用`TCP`传输协议，同时支持`Http1`和`WebSocket`
协议，采用的是默认的日志配置。
