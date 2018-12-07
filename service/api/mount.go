package api

import "github.com/labstack/echo"
// MountAPI registers API
func MountAPI(group *echo.Group) {
	MountLogInfoAPI(group)
	MountReceiverAPI(group)
}

func MountLogInfoAPI(group *echo.Group) {
	task := group.Group("/loginfo")
	task.GET("", HTTPGetLogInfo)
	task.GET("/:containerName",HTTPGetLogInfoByContainerName)

}

func MountReceiverAPI(group *echo.Group)  {
	task := group.Group("/receivers")
	task.GET("", HTTPGetReceivers)
	task.POST("",HTTPAddReceiver)
	task.PUT("/:id",HTTPUpdateReceiver)
	task.GET("/del/:id",HTTPDelReceiver)
	task.GET("/:containerName",HTTPGetReceiversByContainerName)

}

