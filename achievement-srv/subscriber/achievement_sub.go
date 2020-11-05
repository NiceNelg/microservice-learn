package subscriber

import (
	pb "achievement-srv/proto/task"
	"achievement-srv/repository"
	"context"
	"errors"
	"log"
	"strings"
	"time"
)

// 定义实现类
type AchievementSub struct {
	Model repository.AchievementRepo
}

func (this *AchievementSub) Finished(ctx context.Context, task *pb.Task) error {
	log.Println("Fininshed")
	log.Printf("Handler Received message: %v\\n", task)
	if task.UserId == "" || strings.TrimSpace(task.UserId) == "" {
		return errors.New("userId is blank")
	}
	entity, err := this.Model.FindByUserId(ctx, task.UserId)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	if entity == nil {
		entity = &repository.Achievement{
			UserId:        task.UserId,
			Total:         1,
			Finished1Time: now,
			UpdateTime:    now,
		}
		return this.Model.Insert(ctx, entity)
	}
	entity.Total++
	entity.UpdateTime = now
	switch entity.Total {
	case 100:
		entity.Finished100Time = now
	case 1000:
		entity.Finished1000Time = now
	}
	return this.Model.Update(ctx, entity)
}

func (this *AchievementSub) Finished2(ctx context.Context, task *pb.Task) error {
	log.Println("Fininshed2")
	return nil
}

func (this *AchievementSub) Finished3(ctx context.Context, task *pb.Task) error {
	log.Println("Fininshed3")
	return nil
}
