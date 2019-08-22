package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tangxusc/cqrs-sidecar/pkg/config"
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
			//1,启动消息监听
			//2,启动grpc
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
