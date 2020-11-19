package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/broker/nats"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/pkg/errors"
	"log"
	"task-srv/common/tracer"
	"task-srv/controller"
	pb "task-srv/proto/task"
	"task-srv/repository"
	"task-srv/utils"
	"time"
)

const (
	MONGO_URL  = "mongodb://127.0.0.1:27017"
	ServerName = "go.micro.service.task"
	JaegerAddr = "127.0.0.1:6831"
)

func main() {
	log.SetFlags(log.Llongfile)

	conn, err := utils.ConnectMongo(MONGO_URL, time.Second*60)
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
		// 配置etcd作为注册中心
		micro.Registry(
			etcd.NewRegistry(
				registry.Addrs("127.0.0.1:2379"),
			),
		),
		micro.Broker(
			nats.NewBroker(
				broker.Addrs("nats://127.0.0.1:4222"),
			),
		),
		// 配置链路追踪为jaeger
		micro.WrapHandler(
			opentracing.NewHandlerWrapper(jaegerTracer),
		),
	)
	service.Init()

	ctro := &controller.TaskController{
		TaskRepo: &repository.TaskRepoImpl{
			Conn: conn,
		},
		// 注入消息发送实例,为避免消息名冲突,这里的topic我们用服务名+自定义消息名拼出
		TaskFinishedPubEvent: micro.NewEvent("go.micro.service."+controller.TaskFinishedTopic, service.Client()),
	}

	//resp := new(pb.ResponseObj)
	//now := time.Now()
	//err = ctro.Create(context.Background(), &pb.Task{
	//	Body:      "完成学习笔记（一）",
	//	StartTime: now.Unix(),
	//	EndTime:   now.Add(time.Hour * 24).Unix(),
	//}, resp)
	//fmt.Println(err)
	//fmt.Println(resp)
	//return

	if err := pb.RegisterTaskServiceHandler(service.Server(), ctro); err != nil {
		log.Fatal(errors.WithMessage(err, "register server"))
	}

	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}
}
