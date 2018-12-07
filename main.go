package main

import (
	"gitlab.wallstcn.com/operation/eyebridge/common"
	ivksvc "gitlab.wallstcn.com/wscnbackend/ivankastd/service"
	"gitlab.wallstcn.com/operation/eyebridge/rpcserver"
	"gitlab.wallstcn.com/operation/eyebridge/models"
	"gitlab.wallstcn.com/operation/eyebridge/service"
	"gitlab.wallstcn.com/operation/eyebridge/business"
)

func main() {
	common.LoadConfig("./conf/eyebridge.yaml")
	common.Initalise()
	startService()
	go business.ReadMessage()
	service.RunServer()


	defer models.CloseDB()
}

func startService() {
	svc := ivksvc.NewService(common.GlobalConf.Micro)
	svc.Init()
	rpcserver.Init(svc)

}
