package consumer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"log"
	"os"
	"strings"
	"time"
	"utils"
)

var confFilename = "D:/LinkChange/src/conf/app.ini"
var addr []string
var topic [] string
var VisitCount = 0

func init(){
	fmt.Println("初始化kafka消费者")
	iniParser := utils.IniParser{}
	if err:=iniParser.Load(confFilename); err!=nil{
		fmt.Println("loading ini file fail",err)
		return
	}
	//消费者集群地址
	addr = strings.Split(iniParser.GetString("kafka","Host"),",")
	//消费者topic名称
	topic = strings.Split(iniParser.GetString("kafka","Topic"), ",")
}

// consumer 消费者
func Consumer() {
	groupID := "group-1"
	config := cluster.NewConfig()
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest//初始从最新的offset开始

	c, err := cluster.NewConsumer(addr, groupID, topic, config)
	if err != nil {
		log.Printf("fail create consumer, message=%s \n", err)
		return
	}
	defer c.Close()
	//轮询错误线程
	go func(c *cluster.Consumer) {
		errors := c.Errors()
		noti := c.Notifications()
		for {
			select {
			case err := <-errors:
				log.Println(err)
			case <-noti:
			}
		}
	}(c)

	for msg := range c.Messages() {
		fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
		VisitCount ++
		c.MarkOffset(msg, "") //MarkOffset 并不是实时写入kafka，有可能在程序crash时丢掉未提交的offset
	}
}
