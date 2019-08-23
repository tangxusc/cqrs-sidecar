package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const debugArgName = "debug"

func InitLog() {
	if Instance.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
		logrus.Debug("[config]已开启debug模式...")
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}

func BindParameter(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&Instance.Debug, debugArgName, "v", false, "debug mod")

	cmd.PersistentFlags().StringVarP(&Instance.Db.Address, "db-address", "", "172.17.0.2", "db数据库连接地址")
	cmd.PersistentFlags().StringVarP(&Instance.Db.Port, "db-port", "", "3306", "db数据库端口")
	cmd.PersistentFlags().StringVarP(&Instance.Db.Database, "db-Database", "", "test", "db数据库实例")
	cmd.PersistentFlags().StringVarP(&Instance.Db.Username, "db-Username", "", "root", "db数据库用户名")
	cmd.PersistentFlags().StringVarP(&Instance.Db.Password, "db-Password", "", "123456", "db数据库密码")
	cmd.PersistentFlags().IntVarP(&Instance.Db.LifeTime, "db-LifeTime", "", 10, "db数据库连接最大连接周期(秒)")
	cmd.PersistentFlags().IntVarP(&Instance.Db.MaxOpen, "db-MaxOpen", "", 5, "db数据库最大连接数")
	cmd.PersistentFlags().IntVarP(&Instance.Db.MaxIdle, "db-MaxIdle", "", 5, "db数据库最大等待数量")

	cmd.PersistentFlags().StringVarP(&Instance.Pulsar.Url, "pulsar-url", "", "pulsar://localhost:6650", "pulsar消息中间件地址")
	cmd.PersistentFlags().StringVarP(&Instance.Pulsar.TopicName, "pulsar-topic-name", "", "cqrs-db", "pulsar消息中间件主题名称")

	cmd.PersistentFlags().StringVarP(&Instance.Rpc.Port, "rpc-port", "", "9999", "rpc端口")
}

type RpcConfig struct {
	Port string
}

type Config struct {
	Debug  bool
	Db     *DbConfig
	Pulsar *PulsarConfig
	Rpc    *RpcConfig
}

var Instance = &Config{
	Debug:  true,
	Db:     &DbConfig{},
	Pulsar: &PulsarConfig{},
	Rpc:    &RpcConfig{},
}

type PulsarConfig struct {
	Url       string
	TopicName string
}

type DbConfig struct {
	Address  string
	Port     string
	Database string
	Username string
	Password string

	LifeTime int
	MaxOpen  int
	MaxIdle  int
}
