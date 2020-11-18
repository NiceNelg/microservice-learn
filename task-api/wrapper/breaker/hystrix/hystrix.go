package hystrix

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/micro/go-micro/v2/client"
	"log"
	pb "task-api/proto/task"
)

type clientWrapper struct {
	client.Client
}

// 自定义熔断后执行的操作
func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{},
	opts ...client.CallOption) error {

	// 命令名的写法参考官方插件，服务名和方法名拼接
	name := req.Service() + "." + req.Endpoint()

	// 自定义当前命令的熔断配置，除了超时时间还有很多其他配置请自行研究
	// 这些配置在wrapper调用时才执行，因此具有最高的优先级
	// ---如果打算使用全局参数配置，请注释掉下面几行---
	config := hystrix.CommandConfig{
		Timeout: 500,
	}
	hystrix.ConfigureCommand(name, config)
	// ---如果打算使用全局参数配置，请注释掉上面几行---

	do := func() error {
		// 这里调用了真正的服务
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	demote := func(err error) error {
		// 因为是示例程序，只处理请求超时这一种错误的降级，其他错误仍抛给上级调用函数
		if err != hystrix.ErrTimeout {
			return err
		}
		switch r := rsp.(type) {
		case *pb.ResponseArr:
			log.Print("search task fail: ", err)
			*r = pb.ResponseArr{
				Result: 1,
				Code:   200,
				Msg:    "",
				Data:   []*any.Any{},
			}
		default:
			log.Print("unknown err: ", err)
		}
		return nil
	}

	return hystrix.Do(name, do, demote)
}

func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}
