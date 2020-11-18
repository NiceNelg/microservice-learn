package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"log"
	"net/http"
	pb "task-api/proto/task"
)

var service pb.TaskService

func Router(g *gin.Engine, taskService pb.TaskService) {

	service = taskService
	v1 := g.Group("/task")
	{
		v1.GET("/search", Search)
		v1.POST("/finished", Finished)
	}
}

func Search(c *gin.Context) {

	req := new(pb.SearchRequest)
	if err := c.BindQuery(req); err != nil {
		log.Print("bad request param: ", err)
		return
	}
	if resp, err := service.Search(c, req); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code": 200,
				"msg":  err.Error(),
				"data": nil,
			},
		)
	} else {
		var data []*pb.Task
		d := resp.GetData()
		for _, i2 := range d {
			t := new(pb.Task)
			ptypes.UnmarshalAny(i2, t)
			data = append(data, t)
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": 200,
				"msg":  "",
				"data": data,
			},
		)
	}
}

func Finished(c *gin.Context) {

	req := new(pb.Task)
	if err := c.BindJSON(req); err != nil {
		log.Print("bad request param: ", err)
		return
	}
	if resp, err := service.Finished(c, req); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code": 200,
				"msg":  err.Error(),
				"data": nil,
			},
		)
	} else {
		d := resp.GetResult()
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": 200,
				"msg":  "",
				"data": d,
			},
		)
	}
}
