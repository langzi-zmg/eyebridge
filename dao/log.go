package dao

import (
	"gitlab.wallstcn.com/operation/eyebridge/models"
	std "gitlab.wallstcn.com/wscnbackend/ivankastd"
)

func PutLogInfo(logContent string, container_Name string, namespace_Name string, container_Image string, host string, createdAt string, ifSendMail bool) {
	var opLog models.LogsInfo
	{
		opLog.LogContent = logContent
		opLog.ContainerName = container_Name
		opLog.NamespaceName = namespace_Name
		opLog.ContainerImage = container_Image
		opLog.Host = host
		opLog.CreatedAt = createdAt
		opLog.IfSendMail = ifSendMail
	}
	db := models.DB().Create(&opLog)
	if err := db.Error; err != nil {
		std.LogErrorc("mysql", err, "failed to create opLog")
	}
}
func GetLogInfo(page int64, limit int64) (*[]models.LogsInfo) {
	var optLog []models.LogsInfo
	models.DB().Limit(limit).Offset((page - 1) * limit).Find(&optLog).Order("CreatedAt desc")
	return &optLog
}

func GetLastSendMailLogInfo(container_name string, ifSendMail bool) (models.LogsInfo, error) {
	var optLog models.LogsInfo
	db := models.DB().Last(&optLog, "container_name = ? AND if_send_mail = ?", container_name, ifSendMail)
	std.LogInfoLn("%s", db)
	if err := db.Error; err != nil {
		std.LogErrorc("mysql", err, "failed to select ")
	}
	return optLog, nil
}

func CountLastLogInfo(container_name string) int64 {
	var opLog models.LogsInfo
	var count int64
	db := models.DB().Model(opLog).Where("container_name = ?", container_name).Count(&count)
	std.LogInfoLn("%s", db)
	return count
}

func GetLogInfoByContainerName(containerName string,page int64, limit int64) (*[]models.LogsInfo) {
	var optLog []models.LogsInfo
	if models.DB().Where("container_name = ? ", containerName).Limit(limit).Offset((page - 1) * limit).Find(&optLog).RecordNotFound() {
		return nil
	}
	return &optLog
}
