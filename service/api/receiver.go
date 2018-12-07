package api

import (
	"github.com/labstack/echo"
	"gitlab.wallstcn.com/operation/eyebridge/models"
	"gitlab.wallstcn.com/operation/eyebridge/dao"
	"gitlab.wallstcn.com/operation/eyebridge/helper"
	"fmt"
	"strconv"
)

// @Title Receiver list
// @Description Receiver list
// @Accept  json
// @Param page query int false "页数|默认1"
// @Param limit query int false "每页条目数|默认10"
// @Resource Receiver
// @Router /v1/receivers   [get]

func HTTPGetReceivers(ctx echo.Context) error {
	args := &models.Pagination{}

	if err := ctx.Bind(args); err != nil {
		return helper.ErrorResponse(ctx, err)
	}
	fmt.Println("HTTPGetReceivers==========args.Page", args.Page)
	fmt.Println("HTTPGetReceivers==========args.Limit", args.Limit)
	if args.Page <= 0 {
		args.Page = 1
	}
	if args.Limit <= 0 {
		args.Limit = 10
	}

	receivers := dao.GetReceivers(args.Page, args.Limit)

	return helper.SuccessResponse(ctx, receivers)
}

// @Title Receiver list
// @Description Receiver list
// @Accept  json
// @Param page query int false "页数|默认1"
// @Param limit query int false "每页条目数|默认10"
// @Resource Receiver
// @Router /v1/receivers/{containerName}   [get]
func HTTPGetReceiversByContainerName(ctx echo.Context) error {

	containerName := ctx.Param("containerName")

	args := &models.Pagination{}

	if err := ctx.Bind(args); err != nil {
		return helper.ErrorResponse(ctx, err)
	}

	//fmt.Println("HTTPGetLogInfoByContainerName==========args.Page", args.Page)
	//fmt.Println("HTTPGetLogInfoByContainerName==========args.Limit", args.Limit)
	if args.Page <= 0 {
		args.Page = 1
	}
	if args.Limit <= 0 {
		args.Limit = 10
	}
	receiver := dao.GetReceiversByContainerName(containerName)

	return helper.SuccessResponse(ctx, receiver)
}

// @Title 新增Receiver
// @Description 新增Receiver
// @Accept  json
// @Param   args
// @Resource Receiver
// @Router /v1/receivers [post]
func HTTPAddReceiver(ctx echo.Context) error {

	args := new(models.Receiver)
	if err := ctx.Bind(args); err != nil {
		return helper.ErrorResponse(ctx, err)
	}
	fmt.Println("HTTPAddReceiver==========args.ReceiverEmail", args.ReceiverEmail)
	fmt.Println("HTTPAddReceiver==========args.ContainerName", args.ContainerName)
	rec := models.Receiver{
		ReceiverEmail: args.ReceiverEmail,
		ContainerName: args.ContainerName,
	}
	fmt.Println("HTTPAddReceiver==========rec.ReceiverEmail", rec.ReceiverEmail)
	fmt.Println("HTTPAddReceiver==========rec.ContainerName", rec.ContainerName)
	err := dao.SaveReceiver(rec.ReceiverEmail, rec.ContainerName)
	if err != nil {
		return helper.ErrorResponse(ctx, err)
	}
	return helper.SuccessResponse(ctx, rec)
}

// @Title 更新Receiver
// @Description 更新Receiver
// @Accept  json
// @Param   args
// @Resource Receiver
// @Router /v1/receivers/{id}  [put]

func HTTPUpdateReceiver(ctx echo.Context) error {
	id, parseIntErr := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if parseIntErr != nil {
		return helper.ErrorResponse(ctx, parseIntErr)
	}

	args := new(models.Receiver)
	if err := ctx.Bind(args); err != nil {
		return helper.ErrorResponse(ctx, err)
	}

	err := dao.UpdateReceiver(id, args.ReceiverEmail, args.ContainerName)
	if err != nil {
		return helper.ErrorResponse(ctx, err)
	}
	return helper.SuccessResponse(ctx, err)
}

// @Title delete Receiver by id
// @Description  delete Receiver by id
// @Accept  json
// @Param page query int false "页数|默认1"
// @Param limit query int false "每页条目数|默认10"
// @Resource Receiver
// @Router /v1/receivers/del/{id}   [get]

func HTTPDelReceiver(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return helper.ErrorResponse(ctx, err)
	}

	error := dao.DelReceiver(id)
	if error != nil {
		return helper.ErrorResponse(ctx, error)
	}
	return helper.SuccessResponse(ctx, error)

}
