package api

import (
	"github.com/labstack/echo"
	"gitlab.wallstcn.com/operation/eyebridge/models"
	"gitlab.wallstcn.com/operation/eyebridge/dao"
	"gitlab.wallstcn.com/operation/eyebridge/helper"
)

// @Title loginfo list
// @Description 获取loginfo list
// @Accept  json
// @Param page query int false "页数|默认1"
// @Param limit query int false "每页条目数|默认10"
// @Resource loginfo
// @Router /v1/loginfo   [get]

func HTTPGetLogInfo(ctx echo.Context) error {
	args := &models.Pagination{}

	if err := ctx.Bind(args); err != nil {
		return helper.ErrorResponse(ctx, err)
	}
	//fmt.Println("HTTPGetLogInfo==========args.Page", args.Page)
	//fmt.Println("HTTPGetLogInfo==========args.Limit", args.Limit)
	if args.Page <= 0 {
		args.Page = 1
	}
	if args.Limit <= 0 {
		args.Limit = 10
	}

	logInfo := dao.GetLogInfo(args.Page, args.Limit)

	return helper.SuccessResponse(ctx, logInfo)
}

// @Title loginfo list
// @Description 获取loginfo list
// @Accept  json
// @Param page query int false "页数|默认1"
// @Param limit query int false "每页条目数|默认10"
// @Resource loginfo
// @Router /v1/loginfo/{containerName}  [get]

func HTTPGetLogInfoByContainerName(ctx echo.Context) error {
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

	logInfo := dao.GetLogInfoByContainerName(containerName, args.Page, args.Limit)

	return helper.SuccessResponse(ctx, logInfo)
}
