package event

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
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

type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}

var ConsumerImpl Consumer

type Sender interface {
	SendEvent(ctx context.Context, e Event, key string)
}

var SenderImpl Sender

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
