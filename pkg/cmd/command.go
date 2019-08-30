package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tangxusc/cqrs-sidecar/pkg/config"
	"github.com/tangxusc/cqrs-sidecar/pkg/db"
	"github.com/tangxusc/cqrs-sidecar/pkg/event"
	"github.com/tangxusc/cqrs-sidecar/pkg/proxy"
	"github.com/tangxusc/cqrs-sidecar/pkg/rpc"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func NewCommand(ctx context.Context) *cobra.Command {
	var command = &cobra.Command{
		Use:   "start",
		Short: "start sidecar",
		Run: func(cmd *cobra.Command, args []string) {
			rand.Seed(time.Now().Unix())
			config.InitLog()

			db.InitConn(ctx)
			defer db.CloseConn()

			go proxy.Start(ctx)
			//2,启动grpc
			rpc.Start(ctx)
			defer rpc.Close()
			//1,启动消息监听
			go event.StartConsumer(ctx)
			defer event.StopConsumer()
			//3,metrics

			<-ctx.Done()
		},
	}
	logrus.SetFormatter(&logrus.TextFormatter{})
	config.BindParameter(command)

	return command
}

func HandlerNotify(cancel context.CancelFunc) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		cancel()
	}()
}
