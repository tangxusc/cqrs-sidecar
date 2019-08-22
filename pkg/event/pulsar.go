package event

import (
	"context"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	"github.com/tangxusc/cqrs-sidecar/pkg/config"
	"log"
	"os"
	"runtime"
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
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            config.Instance.Pulsar.TopicName,
		SubscriptionName: name,
		Type:             pulsar.KeyShared,
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
		msg, err := p.consumer.Receive(ctx)
		if err != nil {
			log.Fatal(err)
		}

		//1,存储到db
		//2,grpc 推送
		//3,ack
		err = processMessage(msg)

		if err == nil {
			p.consumer.Ack(msg)
		} else {
			p.consumer.Nack(msg)
		}
	}
}
