package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
	"github.com/siddontang/go-mysql/test_util/test_keys"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-sidecar/pkg/config"
	"net"
	"os"
)

const serverVersion = "8.0.3"

func Start(ctx context.Context) {
	l, e := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.Instance.Sidecar.Port))
	if e != nil {
		logrus.Errorf("[proxy]监听tcp出现错误,错误:%v", e.Error())
		os.Exit(1)
	}
	provider := server.NewInMemoryProvider()
	//sidecar不需要密码
	provider.AddUser(config.Instance.Db.Username, config.Instance.Db.Password)
	var tlsConf = server.NewServerTLSConfig(test_keys.CaPem, test_keys.CertPem, test_keys.KeyPem, tls.VerifyClientCertIfGiven)
	newServer := server.NewServer(serverVersion, mysql.DEFAULT_COLLATION_ID, mysql.AUTH_SHA256_PASSWORD, test_keys.PubPem, tlsConf)
	for {
		select {
		case <-ctx.Done():
			e = l.Close()
			return
		default:
			lConn, e := l.Accept()
			if e != nil {
				logrus.Errorf("[proxy]tcp Accept出现错误,错误:%v", e.Error())
				return
			}
			handler := NewHandler()
			conn, e := server.NewCustomizedConn(lConn, newServer, provider, handler)
			if e != nil {
				logrus.Errorf("[proxy]NewCustomizedConn出现错误,错误:%v", e.Error())
				break
			}
			connHandler := handler.(*ConnHandler)
			connHandler.Conn = conn
			go startConnCommandHandler(ctx, conn)
		}
	}
}

func startConnCommandHandler(ctx context.Context, conn *server.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			e := conn.HandleCommand()
			if e != nil {
				logrus.Warningf("[proxy]HandleCommand出现错误,错误:%v", e.Error())
				return
			}
		}
	}
}
