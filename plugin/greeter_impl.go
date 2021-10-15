package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	example "go-plugin-demo/commons"
	"os"
)


//业务接口的真正实现
type GreeterHello struct {
	logger hclog.Logger
}

//之前暴露的插件业务接口，此处必须实现，供宿主机进程RPC调用
func (g *GreeterHello) Greet() string {
	g.logger.Debug("message from GreeterHello.Greet")
	return "Hello!"
}

//握手配置，插件进程和宿主机进程，都需要保持一致
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Trace,
		Output: os.Stderr,
		JSONFormat: true,
	})

	//实例化一个greeter
	greeter := &GreeterHello{
		logger: logger,
	}

	// pluginMap is the map of plugins we can dispense.
	// 插件进程必须指定Impl,此处赋值为greeter对象
	var pluginMap = map[string]plugin.Plugin{
		"greeter": &example.GreeterPlugin{Impl: greeter},
	}

	logger.Debug("message from plugin", "foo", "bar")

	//调用plugin.Serve() 启动侦听，并提供服务
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: pluginMap,
	})
}