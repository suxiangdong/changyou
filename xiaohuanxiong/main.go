package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/thinkingdata/logbus/plugin/parser2"
	proto "github.com/thinkingdata/logbus/plugin/prot"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var separator = []byte("::")

// Parser 自定义解析器
type Parser interface {
	Parse([]byte) ([]byte, error)
}

// GRPCServer 插件服务端
type GRPCServer struct {
	Impl Parser
}

type GRPCPlugin struct {
	plugin.Plugin
	Impl Parser
}

// GRPCClient GRPC客户端
func (_ *GRPCPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return nil, nil
}

func (g *GRPCPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterParserServer(s, &GRPCServer{Impl: g.Impl})
	return nil
}

// Parse 服务端解析
func (g *GRPCServer) Parse(_ context.Context, req *proto.Request) (*proto.Response, error) {
	data, err := g.Impl.Parse(req.Content)
	return &proto.Response{
		Content: data,
	}, err
}

type Par struct{}

// Parse 服务端解析
func (g *Par) Parse(b []byte) ([]byte, error) {
	return parser2.Parse(b)
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "LogBus",
	MagicCookieValue: "v2",
}

func main() {
	if err := parser2.InitConfig(); err != nil {
		panic(err)
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"parser": &GRPCPlugin{Impl: &Par{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
