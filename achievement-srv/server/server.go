package main

import (
	"achievement-srv/common/tracer"
	"achievement-srv/repository"
	"achievement-srv/subscriber"
	"achievement-srv/utils"
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/broker/nats"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/pkg/errors"
	"log"
	"time"
)

const (
	MONGO_URL  = "mongodb://localhost:27017"
	ServerName = "go.micro.service.achievement"
	JaegerAddr = "127.0.0.1:6831"
)

func main() {
	log.SetFlags(log.Llongfile)

	conn, err := utils.ConnectMongo(MONGO_URL, time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect(context.Background())

	// 配置jaeger连接
	jaegerTracer, closer, err := tracer.NewJaegerTracer(ServerName, JaegerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	service := micro.NewService(
		micro.Name(ServerName),
		micro.Version("lastest"),
		// 配置nats作为消息中间件，
		// 这里没有将该服务注册到etcd中，表明消息的订阅只需要知道自己订阅的事件名称以及接收事件通知的地址，
		// 发布者只需要知道事件的名称以及通知地址即可
		micro.Broker(
			nats.NewBroker(
				broker.Addrs("nats://127.0.0.1:4222"),
			),
		),
		micro.WrapSubscriber(
			// 配置链路追踪为jaeger
			opentracing.NewSubscriberWrapper(jaegerTracer),
		),
	)
	service.Init()

	handler := &subscriber.AchievementSub{
		Model: &repository.AchievementRepoImpl{
			Conn: conn,
		},
	}

	if err := micro.RegisterSubscriber("go.micro.service.task.finished", service.Server(), handler); err != nil {
		log.Fatal(errors.WithMessage(err, "subscribe"))
	}

	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}

}
