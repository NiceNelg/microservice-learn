package main

import (
	"gateway/plugins/auth"
	"github.com/micro/micro/v2/client/api"
	"github.com/micro/micro/v2/cmd"
	"log"
)

func main()  {
	err := api.Register(auth.NewPlugin())
	if err != nil {
		log.Fatal("auth register")
	}

	cmd.Init()
}
