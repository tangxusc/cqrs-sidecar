package event

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-sidecar/pkg/db"
	"os"
	"time"
)

type Impl struct {
	ImplId         string    `json:"id"`
	ImplEventType  string    `json:"event_type"`
	ImplAggId      string    `json:"agg_id"`
	ImplAggType    string    `json:"agg_type"`
	ImplCreateTime time.Time `json:"create_time"`
	ImplData       string    `json:"data"`
}

func (impl *Impl) Data() string {
	return impl.ImplData
}

func (impl *Impl) ToJson() ([]byte, error) {
	return json.Marshal(impl)
}

func (impl *Impl) CreateTime() time.Time {
	return impl.ImplCreateTime
}

func (impl *Impl) AggType() string {
	return impl.ImplAggType
}

func (impl *Impl) AggId() string {
	return impl.ImplAggId
}

func (impl *Impl) EventType() string {
	return impl.ImplEventType
}

func (i *Impl) Id() string {
	return i.ImplId
}

func NewEventImpl(id string, eventType string, aggId string, aggType string, createTime time.Time, data string) *Impl {
	return &Impl{ImplId: id, ImplEventType: eventType, ImplAggId: aggId, ImplAggType: aggType, ImplCreateTime: createTime, ImplData: data}
}

type Event interface {
	Id() string
	EventType() string
	AggId() string
	AggType() string
	CreateTime() time.Time
	Data() string
	ToJson() ([]byte, error)
}

var ConsumerImpl Consumer

type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}

func StartConsumer(ctx context.Context) {
	ConsumerImpl = GetConsumerImpl()
	if ConsumerImpl == nil {
		logrus.Errorf("[event]没找到可支持的mq")
		os.Exit(1)
	}
	e := ConsumerImpl.Start(ctx)
	if e != nil {
		logrus.Errorf("[event]连接消息中间件出现错误,错误:%v", e.Error())
		os.Exit(1)
	}
}

func StopConsumer() {
	e := ConsumerImpl.Stop()
	if e != nil {
		logrus.Errorf("[event]关闭消息中间件连接出现错误,错误:%v", e.Error())
	}
}

var SenderImpl Sender

type Sender interface {
	SendEvent(ctx context.Context, e Event, key string)
}

var Confirmed = `Confirmed`
var NotConfirmed = `NotConfirmed`

/*
处理事件
*/
func ProcessEvent(ctx context.Context, eve Event, key string) error {
	//1,存储到db
	e := db.ConnInstance.Save(eve.Id(), eve.EventType(), eve.AggId(), eve.AggType(), eve.CreateTime(), eve.Data(), NotConfirmed)
	if e != nil {
		return e
	}
	//2,grpc 推送
	go SenderImpl.SendEvent(ctx, eve, key)
	return e
}
