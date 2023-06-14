package rabbit

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/netx/tcpx"
)

func NewRabbitServer(cfg server.CfgRabbitServer) server.IRabbitServer {
	container := NewServerExtensionContainer()
	server := tcpx.NewTCPServer()
	statusDetail := NewServerStatusDetail(cfg.Name, DefaultStatsInterval)
	mgr := NewServerExtensionManager(statusDetail)
	logger := logx.NewLogger()
	return &RabbitServer{
		Config:       cfg,
		SockServer:   server,
		ExtContainer: container,
		ExtManager:   mgr,
		StatusDetail: statusDetail,
		Logger:       logger,
	}
}

type RabbitServer struct {
	Config       server.CfgRabbitServer
	SockServer   tcpx.ITCPServer
	ExtContainer IServerExtensionContainer
	ExtManager   IServerExtensionManager
	StatusDetail *ServerStatusDetail
	Logger       logx.ILogger
}

func (s *RabbitServer) GetId() string {
	return s.Config.Id
}

func (s *RabbitServer) GetName() string {
	return s.Config.Name
}

func (s *RabbitServer) GetLogger() logx.ILogger {
	return s.Logger
}

func (s *RabbitServer) Init() {
	// 注入Extension
	ForeachExtensionConstructor(func(constructor FuncServerExtension) {
		s.ExtContainer.AppendExtension(constructor())
	})
	// 设置SockServer信息
	s.SockServer.SetName(s.Config.FromUser.Name)
	s.SockServer.SetMax(100)
	s.SockServer.SetLogger(s.Logger)
	// 初始化ExtensionManager
	s.ExtManager.InitManager(s.SockServer.GetPackHandlerContainer(), s.ExtContainer, s.SockServer)
	s.ExtManager.SetLogger(s.Logger)
	s.ExtManager.SetAddressProxy(AddressProxy)
	// 初始化Logger
	cfgLog := s.Config.Log
	if nil != cfgLog {
		s.Logger.SetConfig(logx.LogConfig{Type: cfgLog.LogType, Level: cfgLog.LogLevel,
			FilePath: cfgLog.GetLogPath(), MaxSize: cfgLog.MaxSize()})
	}
}

func (s *RabbitServer) Start() {
	s.StatusDetail.Start()
	s.SockServer.AddEventListener(netx.ServerEventConnOpened, s.onConnOpened)
	s.SockServer.AddEventListener(netx.ServerEventConnClosed, s.onConnClosed)
	s.ExtManager.StartManager()
	s.SockServer.StartServer(netx.SockParams{
		Network: netx.ParseSockNetwork(s.Config.FromUser.Network), LocalAddress: s.Config.FromUser.Addr}) //这里会阻塞
}

func (s *RabbitServer) Stop() {
	s.SockServer.StopServer()
	s.ExtManager.StopManager()
	s.SockServer.RemoveEventListener(netx.ServerEventConnOpened, s.onConnOpened)
	s.SockServer.RemoveEventListener(netx.ServerEventConnClosed, s.onConnClosed)
	s.StatusDetail.ReStats()
}

func (s *RabbitServer) Restart() {
	s.Stop()
	s.Start()
}

func (s *RabbitServer) Save() {
	//TODO implement me
	panic("implement me")
}

func (s *RabbitServer) onConnOpened(evd *eventx.EventData) {
	s.StatusDetail.AddLinkCount()
}

func (s *RabbitServer) onConnClosed(evd *eventx.EventData) {
	address := evd.Data.(string)
	AddressProxy.RemoveByAddress(address)
	s.StatusDetail.RemoveLinkCount()
}
