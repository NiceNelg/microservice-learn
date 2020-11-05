package main

import (
	"achievement-srv/repository"
	"achievement-srv/subscriber"
	"achievement-srv/utils"
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/pkg/errors"
	"log"
	"time"
)

const (
	MONGO_URL = "mongodb://localhost:27017"
)

func main() {
	log.SetFlags(log.Llongfile)

	conn, err := utils.ConnectMongo(MONGO_URL, time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect(context.Background())

	service := micro.NewService(
		micro.Name("go.micro.service.achievement"),
		micro.Version("lastest"),
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
