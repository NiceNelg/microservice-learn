package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/web"
	"log"
	pb "task-api/proto/task"
)

func main() {
	g := gin.Default()

	service := web.NewService(
		web.Name("go.micro.api.task"),
		web.Address(":8888"),
		web.Handler(g),
	)

	cli := pb.NewTaskService("go.micro.service.task", client.DefaultClient)

	v1 := g.Group("/task")
	v1.GET("/search", func(c *gin.Context) {
		req := new(pb.SearchRequest)
		if err := c.BindQuery(req); err != nil {
			c.JSON(200, gin.H{
				"code": "500",
				"msg":  "bad param",
			})
			return
		}
		if resp, err := cli.Search(c, req); err != nil {
			c.JSON(200, gin.H{
				"code": "500",
				"msg":  err.Error(),
				"data": nil,
			})
		} else {
			d := resp.GetData()
			var data []*pb.Task
			for _, i2 := range d {
				t := new(pb.Task)
				ptypes.UnmarshalAny(i2, t)
				data = append(data, t)
			}
			c.JSON(200, gin.H{
				"code": "200",
				"msg":  "success",
				"data": data,
			})
		}
	})
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
