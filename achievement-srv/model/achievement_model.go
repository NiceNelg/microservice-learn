package model

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	DbName         = "todolist"
	TaskCollection = "achievement"
)

type Achievement struct {
	Id string `bson:"_id,omitempty"`
	// 用户ID
	UserId string `bson:"userId"`
	// 完成任务总数
	Total int64 `bson:"total"`
	// 完成第一个任务的时间
	Finished1Time int64 `bson:"finished1Time"`
	// 完成第一百个任务的时间
	Finished100Time int64 `bson:"finished100Time"`
	// 完成第一千个任务的时间
	Finished1000Time int64 `bson:"finished1000Time"`
	// 更新时间
	UpdateTime int64 `bson:"updateTime"`
}

type AchievementModel interface {
	FindByUserId(ctx context.Context, userId string) (*Achievement, error)
	Insert(ctx context.Context, achievement *Achievement) error
	Update(ctx context.Context, achievement *Achievement) error
}

type AchievementModelImpl struct {
	Conn *mongo.Client
}

func (this *AchievementModelImpl) collection() *mongo.Collection {
	return this.Conn.Database(DbName).Collection(TaskCollection)
}

func (this *AchievementModelImpl) FindByUserId(ctx context.Context, userId string) (*Achievement, error) {
	result := this.collection().FindOne(ctx, bson.M{"userId": userId})
	// findOne如果查不到是会报错的,这里要处理一下
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	achievement := &Achievement{}
	if err := result.Decode(achievement); err != nil {
		return nil, errors.WithMessage(err, "search mongo")
	}
	return achievement, nil
}
func (this *AchievementModelImpl) Insert(ctx context.Context, achievement *Achievement) error {
	_, err := this.collection().InsertOne(ctx, achievement)
	return err
}

func (this *AchievementModelImpl) Update(ctx context.Context, achievement *Achievement) error {
	achievement.UpdateTime = time.Now().Unix()
	oid, err := primitive.ObjectIDFromHex(achievement.Id)
	if err != nil {
		return err
	}
	achievement.Id = ""
	_, err = this.collection().UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": achievement},
	)
	return err
}
