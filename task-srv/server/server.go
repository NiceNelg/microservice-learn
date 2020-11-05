package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/pkg/errors"
	"log"
	"task-srv/controller"
	pb "task-srv/proto/task"
	"task-srv/repository"
	"task-srv/utils"
	"time"
)

const MONGO_URL = "mongodb://127.0.0.1:27017"

func main() {
	log.SetFlags(log.Llongfile)

	conn, err := utils.ConnectMongo(MONGO_URL, time.Second*60)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect(context.Background())

	service := micro.NewService(
		micro.Name("go.micro.service.task"),
		micro.Version("lastest"),
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
