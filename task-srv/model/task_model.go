package model

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	pb "task-srv/proto/task"
	"time"
)

const (
	DbName         = "todolist"
	TaskCollection = "task"
	UnFinished     = 0
	Finished       = 1
)

type TaskModel interface {
	InsertOne(ctx context.Context, task *pb.Task) error
	Delete(ctx context.Context, id string) error
	Modify(ctx context.Context, task *pb.Task) error
	Finished(ctx context.Context, task *pb.Task) error
	Count(ctx context.Context, keyword string) (int64, error)
	Search(ctx context.Context, req *pb.SearchRequest) ([]*pb.Task, error)
	// 根据ID搜索
	FindById(ctx context.Context, id string) (*pb.Task, error)
}

type TaskModelImpl struct {
	Conn *mongo.Client
}

func (this *TaskModelImpl) collection() *mongo.Collection {
	return this.Conn.Database(DbName).Collection(TaskCollection)
}

func (this *TaskModelImpl) InsertOne(ctx context.Context, task *pb.Task) error {
	_, err := this.collection().InsertOne(ctx, bson.M{
		"body":       task.Body,
		"startTime":  task.StartTime,
		"endTime":    task.EndTime,
		"isFinished": UnFinished,
		"createTime": time.Now().Unix(),
		// 新增用户id
		"userId": task.UserId,
	})
	return err
}

func (this *TaskModelImpl) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = this.collection().DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (this *TaskModelImpl) Modify(ctx context.Context, task *pb.Task) error {
	id, err := primitive.ObjectIDFromHex(task.Id)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	update := bson.M{
		"isFinished": int32(task.IsFinished),
		"updateTime": now,
	}
	if task.IsFinished == Finished {
		update["finishTime"] = now
	}
	log.Print(task)
	log.Println(update)
	_, err = this.collection().UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (this *TaskModelImpl) Finished(ctx context.Context, task *pb.Task) error {
	id, err := primitive.ObjectIDFromHex(task.Id)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	update := bson.M{
		"isFinished": int32(task.IsFinished),
		"updateTime": now,
	}
	if task.IsFinished == Finished {
		update["finishTime"] = now
	}
	log.Print(task)
	log.Println(update)
	_, err = this.collection().UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (this *TaskModelImpl) Count(ctx context.Context, keyword string) (int64, error) {
	filter := bson.M{}
	if keyword != "" && strings.TrimSpace(keyword) != "" {
		filter = bson.M{
			"body": bson.M{"$regex": keyword},
		}
	}
	count, err := this.collection().CountDocuments(ctx, filter)
	return count, err
}

func (this *TaskModelImpl) Search(ctx context.Context, req *pb.SearchRequest) ([]*pb.Task, error) {
	filter := bson.M{}
	if req.Keyword != "" && strings.TrimSpace(req.Keyword) != "" {
		filter = bson.M{
			"body": bson.M{"$regex": req.Keyword},
		}
	}
	cursor, err := this.collection().Find(
		ctx,
		filter,
		options.Find().SetSkip((req.Page-1)*req.Limit),
		options.Find().SetLimit(req.Limit),
		options.Find().SetSort(bson.M{req.SortBy: req.Order}),
	)
	if err != nil {
		return nil, errors.WithMessage(err, "search mongo")
	}
	var rows []*pb.Task
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, errors.WithMessage(err, "parse data")
	}
	return rows, err
}

func (this *TaskModelImpl) FindById(ctx context.Context, id string) (*pb.Task, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.WithMessage(err, "parse ID")
	}
	result := this.collection().FindOne(ctx, bson.M{"_id": objectId})
	task := new(pb.Task)
	if err := result.Decode(task); err != nil {
		return nil, errors.WithMessage(err, "search mongo")
	}
	return task, nil
}
