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

	cmd.PersistentFlags().StringVarP(&Instance.DbConfig.Address, "db-address", "", "172.17.0.2", "db数据库连接地址")
	cmd.PersistentFlags().StringVarP(&Instance.DbConfig.Port, "db-port", "", "3306", "db数据库端口")
	cmd.PersistentFlags().StringVarP(&Instance.DbConfig.Database, "db-Database", "", "test", "db数据库实例")
	cmd.PersistentFlags().StringVarP(&Instance.DbConfig.Username, "db-Username", "", "root", "db数据库用户名")
	cmd.PersistentFlags().StringVarP(&Instance.DbConfig.Password, "db-Password", "", "123456", "db数据库密码")
	cmd.PersistentFlags().IntVarP(&Instance.DbConfig.LifeTime, "db-LifeTime", "", 10, "db数据库连接最大连接周期(秒)")
	cmd.PersistentFlags().IntVarP(&Instance.DbConfig.MaxOpen, "db-MaxOpen", "", 5, "db数据库最大连接数")
	cmd.PersistentFlags().IntVarP(&Instance.DbConfig.MaxIdle, "db-MaxIdle", "", 5, "db数据库最大等待数量")
}

type Config struct {
	Debug    bool
	DbConfig *DbConfig
}

var Instance = &Config{
	Debug:    true,
	DbConfig: &DbConfig{},
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
