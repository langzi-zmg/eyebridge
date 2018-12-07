package business

import (
	"flag"
	"strings"
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"fmt"
	"sync"
	"encoding/json"
	"log"
	"gitlab.wallstcn.com/operation/eyebridge/dao"
	"gitlab.wallstcn.com/wscnbackend/ivankaprotocol/xinge"
	"gitlab.wallstcn.com/operation/eyebridge/rpcserver"
	"strconv"
	"time"
	"gitlab.wallstcn.com/operation/eyebridge/models"
	"github.com/goinggo/mapstructure"
	"text/template"
	"io"
)

type KubernetesInfo struct {
	ContainerName    string      `json:"container_name"`
	NamespaceName    string      `json:"namespace_name"`
	Pod_Name         string      `json:"pod_name"`
	ContainerImage   string      `json:"container_image"`
	ContainerImageId string      `json:"container_image_id"`
	PodId            string      `json:"pod_id"`
	Labels           interface{} `json:"labels"`
	Host             string      `json:"host"`
	MasterUrl        string      `json:"master_url"`
	NamespaceId      string      `json:"namespace_id"`
	NamespaceLabels  interface{} `json:"namespace_labels"`
}

type EmailInventory struct {
	ENV             string
	Container_Name  string
	Namespace_Name  string
	Container_Image string
	Host            string
	Message         string
}

var (
	wg            sync.WaitGroup
	brokerList    = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The comma separated list of brokers in the Kafka cluster. You can also set the KAFKA_PEERS environment variable")
	topic         = flag.String("topic", os.Getenv("KAFKA_TOPIC"), "REQUIRED: the topic to produce to")
	contains      = flag.String("containes", os.Getenv("CONTAINES"), "CONTAINES")
	receivers     = flag.String("receivers", os.Getenv("RECEIVERS"), "receivers")
	ReceiverEmail string
	ReceiversByDb *models.Receiver
)

func ReadMessage() {

	flag.Parse()
	// init
	brokers := strings.Split(*brokerList, ",")
	topics := strings.Split(*topic, ",")
	// int cluster config
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Group.PartitionStrategy = cluster.StrategyRoundRobin

	consumer, err := cluster.NewConsumer(brokers, "applog_to_kafka", topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg, more := <-consumer.Messages():
			if more {

				consumer.MarkOffset(msg, "") // mark message as processed
				msgValue := string(msg.Value)
				env := os.Getenv("CONFIGOR_ENV")
				MessageTo(env, msgValue)
			}

		case err, more := <-consumer.Errors():
			if more {
				log.Printf("Error: %s\n", err.Error())
			}
		case ntf, more := <-consumer.Notifications():
			if more {
				log.Printf("Rebalanced: %+v\n", ntf)
			}
		case <-signals:
			return
		}
	}
}

func MessageTo(env string, msgValue string) {
	//get containsInfo
	containsInfo := strings.Split(*contains, ",")
	//receiversFromEnv := strings.Split(*receivers, ",")

	//define var
	var dat map[string]interface{}
	var content string
	var logMap map[string]string
	var kubernetesInfo KubernetesInfo
	var containerName, namespaceName, containerImage, host string

	for _, contains := range containsInfo {
		if strings.Contains(msgValue, contains) {
			err := json.Unmarshal([]byte(msgValue), &dat)
			if err != nil {
				fmt.Println(err.Error())
			}
			res, _ := json.Marshal(dat["kubernetes"])
			err = json.Unmarshal(res, &kubernetesInfo)
			if err == nil {
				containerName = kubernetesInfo.ContainerName
				namespaceName = kubernetesInfo.NamespaceName
				containerImage = kubernetesInfo.ContainerImage
				host = kubernetesInfo.Host
				content = content + "ENV : " + env + "\n"
				content = content + "container_name : " + containerName + "\n"
				content = content + "namespace_name : " + namespaceName + "\n"
				content = content + "container_image : " + containerImage + "\n"
				content = content + "host : " + host + "\n"
				//get email receivers
				//fmt.Println("containerName===============containerName", containerName)
				ReceiversByDb = dao.GetReceiversByContainerName(containerName)
				//fmt.Println("receiversByDb.ReceiverEmail=============",ReceiversByDb.ReceiverEmail)

				if ReceiversByDb != nil {
					ReceiverEmail = ReceiversByDb.ReceiverEmail
				}
				if ReceiversByDb == nil {
					ReceiverEmail = *receivers
				}
				//fmt.Println("receivers====================", *receivers)
				//fmt.Println("receiverEmail====================", ReceiverEmail)
			}
			err = json.Unmarshal([]byte(fmt.Sprintf("%s", dat["log"])), &logMap);
			//log info if log is json
			if err == nil {
				for key, val := range logMap {
					content = content + key + " : " + val + "\n"
				}
			}
			if err != nil {
				content = content + "log : " + fmt.Sprintf("%s", dat["log"])
			}
			//ratelimit
			count := dao.CountLastLogInfo(containerName)
			if count > 0 {
				optLog, _ := dao.GetLastSendMailLogInfo(containerName, true)
				lastTime := optLog.CreatedAt
				ifSendMailOptLog := optLog.IfSendMail
				//get log time
				nowTime, _ := json.Marshal(dat["time"])
				//save db  but sometimes not send mail
				ifSendMail := RateLimitMessage(lastTime, string(nowTime), ifSendMailOptLog)

				//send mail
				SendMail(ifSendMail, content, ReceiverEmail)
				//save db
				dao.PutLogInfo(fmt.Sprintf("%s", dat["log"]), containerName, namespaceName, containerImage, host, strconv.FormatInt(time.Now().Unix(), 10), ifSendMail)

			}
			if count == 0 {
				dao.PutLogInfo(fmt.Sprintf("%s", dat["log"]), containerName, namespaceName, containerImage, host, strconv.FormatInt(time.Now().Unix(), 10), true)
				SendMail(true, content, ReceiverEmail)
			}
		}
	}

}

func RateLimitMessage(lastTime string, nowTime string, ifSendMailOptLog bool) bool {
	lastTimeInt64, _ := strconv.ParseInt(lastTime, 10, 64)
	nowTimeInt64, _ := strconv.ParseInt(nowTime, 10, 64)
	if (nowTimeInt64-lastTimeInt64)/60 >= 5 && ifSendMailOptLog {
		return true
	}
	return false
}
func SendMail(ifSendMail bool, content string, receivers string) {
	allReceiver := strings.Split(receivers, ",")
	emailContent := ContentInTemplate(content)
	if ifSendMail {
		parms := new(xinge.EmailParms)
		parms.Titile = "panic"
		parms.Receivers = allReceiver
		//parms.Content = content
		parms.Content = emailContent
		rpcserver.SendPanicMail(parms)
	}
}

func ContentInTemplate(content string) string {
	var logMap map[string]string
	var inventory EmailInventory
	var message string

	logMap = make(map[string]string)
	ss := strings.Split(content, "\n")
	for i := 0; i < 5; i++ {
		sss := strings.Split(ss[i], " : ")
		logMap[sss[0]] = sss[1]
	}
	for i := 5; i < len(ss); i++ {
		message = message + ss[i] + "\n"
	}
	logMap["Message"] = message
	mapstructure.Decode(logMap, &inventory)
	tmpl, err := template.ParseFiles("/template/email.html")

	var email  string
	resultWriter := &Result{}
	io.WriteString(resultWriter,email)
	err = tmpl.Execute(resultWriter, inventory)
	if err != nil {
		panic(err)
	}
	return resultWriter.output
}

func (p *Result) Write(b []byte) (n int, err error) {
	p.output += string(b)
	return len(b), nil
}

type Result struct {
	output string
}
