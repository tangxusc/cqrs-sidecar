package rpc

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-sidecar/pkg/config"
	"github.com/tangxusc/cqrs-sidecar/pkg/db"
	"github.com/tangxusc/cqrs-sidecar/pkg/event"
	"github.com/tangxusc/cqrs-sidecar/pkg/proxy"
	_ "github.com/tangxusc/cqrs-sidecar/pkg/proxy/handler"
	"google.golang.org/grpc"
	"io"
	"os"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	config.Instance = &config.Config{
		Debug: false,
		Db: &config.DbConfig{
			Address:  "172.17.0.2",
			Port:     "3306",
			Database: "test",
			Username: "root",
			Password: "123456",
			LifeTime: 10,
			MaxOpen:  5,
			MaxIdle:  5,
		},
		Pulsar: &config.PulsarConfig{
			Url:       "pulsar://localhost:6650",
			TopicName: "cqrs-db",
		},
		Rpc:     &config.RpcConfig{Port: `9999`},
		Sidecar: &config.SidecarConfig{Port: `3307`},
	}
	todo := context.TODO()
	db.InitConn(todo)

	Start(todo)
	go proxy.Start(todo)

	time.Sleep(time.Second * 10)
	startClient()

	e := event.ProcessEvent(todo, event.NewEventImpl("1", "E1", "A1", "AT1", time.Now(), `{"name":"test"}`), "E1:1")
	if e != nil {
		panic(e)
	}
	select {}
}

func startClient() {
	conn, e := grpc.Dial(fmt.Sprintf("127.0.0.1:%s", config.Instance.Rpc.Port), grpc.WithInsecure())
	if e != nil {
		panic(e)
	}

	client := NewConsumerClient(conn)
	consumeClient, e := client.Consume(context.TODO(), &ConsumeRequest{})
	if e != nil {
		panic(e)
	}
	go func() {
		for {
			response, e := consumeClient.Recv()
			if e == io.EOF {
				println("读取完成...")
				return
			}
			if e != nil {
				println("客户端错误6,", e.Error())
				return
			}
			println("客户端收到消息7,", response.String())

			//ack
			fmt.Println("开始ack")
			ack(response)
		}
	}()

}

func ack(response *ConsumeResponse) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=true", config.Instance.Db.Username, config.Instance.Db.Password,
		"tcp", `127.0.0.1`, config.Instance.Sidecar.Port, config.Instance.Db.Database)
	conn, e := sql.Open("mysql", dsn)
	if e != nil {
		logrus.Errorf("[proxy]连接出现错误,url:%v,错误:%v", dsn, e.Error())
		os.Exit(1)
	}
	defer conn.Close()
	conn.SetConnMaxLifetime(time.Duration(config.Instance.Db.LifeTime) * time.Second)
	conn.SetMaxOpenConns(config.Instance.Db.MaxOpen)
	conn.SetMaxIdleConns(config.Instance.Db.MaxIdle)

	tx, e := conn.Begin()
	if e != nil {
		panic(e)
	}
	result, e := tx.Exec(`call ack('` + response.Id + `')`)
	if e != nil {
		panic(e)
	}
	tx.Commit()
	fmt.Println(result)
}
