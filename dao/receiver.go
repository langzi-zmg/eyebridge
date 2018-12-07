package dao

import (
	"gitlab.wallstcn.com/operation/eyebridge/models"
	std "gitlab.wallstcn.com/wscnbackend/ivankastd"
)

func SaveReceiver(receiverEmail string, containerName string) error {
	var receiver models.Receiver
	{
		receiver.ReceiverEmail = receiverEmail
		receiver.ContainerName = containerName
	}

	db := models.DB().Save(&receiver)
	if err := db.Error; err != nil {
		std.LogErrorc("mysql", err, "failed to create opLog")
		return err
	}
	return nil
}
func GetReceivers(page int64, limit int64) (*[]models.Receiver) {
	var receiver []models.Receiver
	models.DB().Limit(limit).Offset((page - 1) * limit).Find(&receiver)
	return &receiver
}

func GetReceiversByContainerName(containerName string) *models.Receiver {
	var receiver models.Receiver
	if models.DB().Where("container_name = ? ", containerName).Find(&receiver).RecordNotFound() {
		return nil
	}
	return &receiver
}

func UpdateReceiver(Id int64,receiverEmail string, containerName string) error{
	var receiver models.Receiver
	models.DB().Where("id = ?", Id).Find(&receiver)
	receiver.ReceiverEmail = receiverEmail
	receiver.ContainerName = containerName
	db := models.DB().Save(&receiver)
	if err := db.Error ;err != nil {
		std.LogErrorc("mysql", err, "failed to update Watcher")
		return err
	}
	return  nil
}

func DelReceiver(Id int64) error {
	var receiver models.Receiver
	db := models.DB().Where("id = ?", Id).Delete(&receiver)
	if err := db.Error; err != nil {
		std.LogErrorc("mysql", err, "failed to delete Watcher")
		return err
	}
	return nil
}
