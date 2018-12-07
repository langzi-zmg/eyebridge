package rpcserver

import (
	"gitlab.wallstcn.com/wscnbackend/ivankaprotocol/xinge"
	"fmt"
	"time"
	"context"
	"github.com/micro/go-micro"
)

var push xinge.PushApiClient

func Init(svc micro.Service) {
	push = xinge.NewPushApiClient("gitlab.wallstcn.com.xinge", svc.Client())
}

func SendPanicMail(Params *xinge.EmailParms,)  {
	fmt.Println(Params.Content)
	fmt.Println(Params.Receivers)
	//fmt.Println(Params.Titile)
	ctx, _ := context.WithTimeout(context.Background(), (10 * time.Second))
	_, err := push.SendEmail(ctx,Params)
	if err != nil {
		fmt.Println("email-sending err: ", err.Error())
	}
	//fmt.Println(rsp.Status)
}