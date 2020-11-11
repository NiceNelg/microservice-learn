package main

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"log"
	pb "task-srv/proto/task"
	"task-srv/repository"
	"time"
)

func main() {
	log.SetFlags(log.Llongfile)

	server := micro.NewService(
		micro.Name("go.micro.client.task"),
		// 配置etcd为注册中心，配置etcd路径，默认端口是2379
		micro.Registry(
			etcd.NewRegistry(
				registry.Addrs("127.0.0.1:2379"),
			),
		),
	)
	server.Init()

	taskService := pb.NewTaskService("go.micro.service.task", server.Client())

	// 调用服务生成三条任务
	now := time.Now()
	insertTask(taskService, "完成学习笔记（一）", now.Unix(), now.Add(time.Hour*24).Unix())
	insertTask(taskService, "完成学习笔记（二）", now.Add(time.Hour*24).Unix(), now.Add(time.Hour*48).Unix())
	insertTask(taskService, "完成学习笔记（三）", now.Add(time.Hour*48).Unix(), now.Add(time.Hour*72).Unix())

	// 分页查询任务列表
	page, err := taskService.Search(context.Background(), &pb.SearchRequest{
		Page:  1,
		Limit: 20,
	})
	if err != nil {
		log.Fatal("search1：", err)
	}
	log.Println(page)

	// 更新第一条记录为完成
	var list []*pb.Task
	for _, item := range page.Data {
		tmp := new(pb.Task)
		ptypes.UnmarshalAny(item, tmp)
		list = append(list, tmp)
	}
	log.Println(list)

	row := list[0]
	if _, err = taskService.Finished(context.Background(), &pb.Task{
		Id:         row.Id,
		IsFinished: repository.Finished,
	}); err != nil {
		log.Fatal("finished", row.Id, err)
	}

	// 修改查询到的第二条数据,延长截至日期
	row = list[1]
	if _, err = taskService.Modify(context.Background(), &pb.Task{
		Id:        row.Id,
		Body:      row.Body,
		StartTime: row.StartTime,
		EndTime:   now.Add(time.Hour * 72).Unix(),
	}); err != nil {
		log.Fatal("modify", row.Id, err)
	}

	// 删除第三条记录
	row = list[2]
	if _, err = taskService.Delete(context.Background(), &pb.Task{
		Id: row.Id,
	}); err != nil {
		log.Fatal("delete", row.Id, err)
	}

	// 再次分页查询，校验修改结果
	page, err = taskService.Search(context.Background(), &pb.SearchRequest{})
	if err != nil {
		log.Fatal("search2", err)
	}
	log.Println(page.String())
}

func insertTask(taskService pb.TaskService, body string, start, end int64) {
	_, err := taskService.Create(context.Background(), &pb.Task{
		UserId:    "1000",
		Body:      body,
		StartTime: start,
		EndTime:   end,
	})
	if err != nil {
		log.Fatal("create", err)
	}
	log.Println("create task success! ")
}
