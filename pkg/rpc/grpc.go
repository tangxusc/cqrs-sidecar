package rpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-sidecar/pkg/config"
	"github.com/tangxusc/cqrs-sidecar/pkg/event"
	"google.golang.org/grpc"
	"net"
	"os"
)

type aggSender struct {
	eventChan   chan event.Event
	recoverChan chan event.Event
	sender      *grpcSender
}

func (sender *aggSender) start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-sender.recoverChan:
				sender.sender.bufChan <- e
			case e := <-sender.eventChan:
				sender.sender.bufChan <- e
			}
		}
	}()
}

type grpcSender struct {
	aggSenders map[string]*aggSender
	bufChan    chan event.Event
}

func (g *grpcSender) Consume(request *ConsumeRequest, stream Consumer_ConsumeServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case e := <-g.bufChan:
			//转换
			response := convert(e)
			//如何处理事务,避免重复发送?
			g.send(stream, response)
		}
	}
}

func convert(e event.Event) *ConsumeResponse {
	ts := &timestamp.Timestamp{
		Seconds: int64(e.CreateTime().Second()),
		Nanos:   int32(e.CreateTime().UnixNano()),
	}
	return &ConsumeResponse{
		Id:        e.Id(),
		EventType: e.EventType(),
		AggId:     e.AggId(),
		AggType:   e.AggType(),
		Data:      e.Data(),
		Create:    ts,
	}
}

func (g *grpcSender) SendEvent(ctx context.Context, e event.Event, key string) {
	sender := g.getSender(key, ctx)
	sender.eventChan <- e
}
func (g *grpcSender) SendRecoverEvent(ctx context.Context, e event.Event, key string) {
	sender := g.getSender(key, ctx)
	sender.recoverChan <- e
}

func (g *grpcSender) getSender(key string, ctx context.Context) *aggSender {
	sender, ok := g.aggSenders[key]
	if !ok {
		newSender := &aggSender{
			eventChan:   make(chan event.Event, 10),
			recoverChan: make(chan event.Event, 10),
			sender:      g,
		}
		g.aggSenders[key] = newSender
		newSender.start(ctx)
		sender = newSender
	}
	return sender
}

func (g *grpcSender) send(stream Consumer_ConsumeServer, response *ConsumeResponse) {
	//TODO:在发送时check数据,确认是否已经ack
	for {
		if err := stream.Send(response); err == nil {
			return
		}
	}
}

var server *grpc.Server

//go:generate protoc --go_out=plugins=grpc:. event.proto

func Start(ctx context.Context) {
	sender := &grpcSender{
		aggSenders: make(map[string]*aggSender),
		bufChan:    make(chan event.Event, 100),
	}
	event.SenderImpl = sender
	go func() {
		listener, e := net.Listen("tcp", fmt.Sprintf(":%s", config.Instance.Rpc.Port))
		if e != nil {
			logrus.Errorf("[rpc]监听出现错误,错误:%v", e.Error())
			os.Exit(1)
		}
		server = grpc.NewServer()
		RegisterConsumerServer(server, sender)
		e = server.Serve(listener)
		if e != nil {
			logrus.Errorf("[rpc]监听出现错误,错误:%v", e.Error())
			os.Exit(1)
		}
	}()
}

func Close() {
	if server != nil {
		server.GracefulStop()
		server = nil
	}
}
