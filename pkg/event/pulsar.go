package event

import (
	"context"
	"encoding/json"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-sidecar/pkg/config"
	"github.com/tangxusc/cqrs-sidecar/pkg/db"
	"os"
	"runtime"
	"time"
)

type PulsarProvider struct {
}

func (p *PulsarProvider) Support(instance *config.Config) bool {
	if len(config.Instance.Pulsar.Url) <= 0 {
		return true
	}
	return false
}

func (p *PulsarProvider) Consumer() Consumer {
	return &PulsarConsumer{}
}

var PulsarProviderImpl = &PulsarProvider{}

type PulsarConsumer struct {
	client   pulsar.Client
	consumer pulsar.Consumer
}

func (p *PulsarConsumer) Start(ctx context.Context) error {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:                     config.Instance.Pulsar.Url,
		OperationTimeoutSeconds: 5,
		MessageListenerThreads:  runtime.NumCPU(),
	})
	if err != nil {
		return err
	}
	name, err := os.Hostname()
	if err != nil {
		return err
	}
	duration := time.Second * 2
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:               config.Instance.Pulsar.TopicName,
		SubscriptionName:    name,
		Type:                pulsar.KeyShared,
		NackRedeliveryDelay: &duration,
	})
	if err != nil {
		return err
	}
	p.consumer = consumer
	p.client = client

	go p.listen(ctx)

	return err
}

func (p *PulsarConsumer) Stop() error {
	e := p.consumer.Close()
	if e != nil {
		return e
	}
	e = p.client.Close()
	if e != nil {
		return e
	}
	return nil
}

func (p *PulsarConsumer) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := p.consumer.Receive(ctx)
			if err != nil {
				logrus.Errorf("[event]接收消息出现错误:%v", err)
				continue
			}
			data, err := unmarshal(msg)
			if err != nil {
				logrus.Errorf("[event]反序列化消息出现错误:%v", err)
				continue
			}
			//1,存储到db
			err = db.ConnInstance.Save(data)
			if err != nil {
				//错误处理
				err = p.consumer.Nack(msg)
				if err != nil {
					logrus.Errorf("[event]Nack出现错误:%v", err)
				}
				continue
			} else {
				err = p.consumer.Ack(msg)
				if err != nil {
					logrus.Errorf("[event]Nack出现错误:%v", err)
					continue
				}
			}
			//2,grpc 推送
			SenderImpl.SendEvent(ctx, data, msg.Key())
		}
	}
}

func unmarshal(message pulsar.Message) (Event, error) {
	impl := &Impl{}
	e := json.Unmarshal(message.Payload(), impl)
	return impl, e
}
