package controller

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	pb "task-srv/proto/task"
	"task-srv/repository"
)

type TaskController struct {
	TaskModel repository.TaskModel
}

func (this *TaskController) Create(ctx context.Context, req *pb.Task, resp *pb.ResponseObj) error {
	if req.Body == "" || req.StartTime <= 0 || req.EndTime <= 0 {
		return errors.New("bad param")
	}
	if err := this.TaskModel.InsertOne(ctx, req); err != nil {
		return err
	}
	*resp = pb.ResponseObj{
		Result: 1,
		Code:   200,
		Msg:    "success",
	}
	return nil
}

func (this *TaskController) Delete(ctx context.Context, req *pb.Task, resp *pb.ResponseObj) error {
	if req.Id == "" {
		return errors.New("bad param")
	}
	if err := this.TaskModel.Delete(ctx, req.Id); err != nil {
		return err
	}
	*resp = pb.ResponseObj{
		Result: 1,
		Code:   200,
		Msg:    "success",
	}
	return nil
}

func (this *TaskController) Modify(ctx context.Context, req *pb.Task, resp *pb.ResponseObj) error {
	if req.Id == "" || req.Body == "" || req.StartTime <= 0 || req.EndTime <= 0 {
		return errors.New("bad param")
	}
	if err := this.TaskModel.Modify(ctx, req); err != nil {
		return err
	}
	*resp = pb.ResponseObj{
		Result: 1,
		Code:   200,
		Msg:    "success",
	}
	return nil
}

func (this *TaskController) Finished(ctx context.Context, req *pb.Task, resp *pb.ResponseObj) error {
	if req.Id == "" || req.IsFinished != repository.UnFinished && req.IsFinished != repository.Finished {
		return errors.New("bad param")
	}
	if err := this.TaskModel.Finished(ctx, req); err != nil {
		return err
	}
	*resp = pb.ResponseObj{
		Result: 1,
		Code:   200,
		Msg:    "success",
	}
	return nil
}

func (this *TaskController) Search(ctx context.Context, req *pb.SearchRequest, resp *pb.ResponseArr) error {
	count, err := this.TaskModel.Count(ctx, req.Keyword)
	if err != nil {
		return errors.New("count row number")
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "createTime"
	}
	if req.Order == 0 {
		req.Order = -1
	}
	if req.Limit*(req.Page-1) > count {
		return errors.New("There's not that much data")
	}
	tmp, err := this.TaskModel.Search(ctx, req)
	if err != nil {
		return err
	}
	var rows []*any.Any
	for _, item := range tmp {
		t, _ := ptypes.MarshalAny(item)
		rows = append(rows, t)
	}
	*resp = pb.ResponseArr{
		Result: 1,
		Code:   200,
		Msg:    "success",
		Data:   rows,
	}
	return nil
}
