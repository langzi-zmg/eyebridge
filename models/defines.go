package models

import (
	"github.com/jinzhu/gorm"
	"gitlab.wallstcn.com/wscnbackend/ivankastd"
	"gitlab.wallstcn.com/wscnbackend/ivankastd/toolkit"
)

var (
	db *gorm.DB
)

func InitModel(config ivankastd.ConfigMysql) {
	db = toolkit.CreateDB(config)
	ivankastd.LogInfoLn("start init mysql model")
	db.AutoMigrate(&LogsInfo{}, &Receiver{})
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func DB() *gorm.DB {
	return db
}

type LogsInfo struct {
	Id              int64  `gorm:"primary_key"`
	LogContent      string `gorm:"size:10000"`
	ContainerName  string `json:"container_name"`
	NamespaceName  string `json:"namespace_name"`
	ContainerImage string `json:"container_image"`
	Host            string `json:"host"`
	CreatedAt       string `json:"created_at"`
	IfSendMail      bool   `json:"if_send_mail"`
}

type Receiver struct {
	Id             int64  `gorm:"primary_key"`
	ReceiverEmail string `json:"receiver_email"`
	ContainerName string `json:"container_name"`
	//KeyWords       string `json:"key_words"`
}

type Pagination struct {
	Page  int64 `json:"page" query:"page"`
	Limit int64 `json:"limit" query:"limit"`
}

type PageAndO struct {
	Page  int64 `json:"page" query:"page"`
	Limit int64 `json:"limit" query:"limit"`
	ContainerName  string `json:"container_name"`
}
