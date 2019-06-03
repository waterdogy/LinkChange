package service

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"strings"
	"time"
	"utils"
)

var confFilename = "D:/LinkChange/src/conf/app.ini"
var addr []string
var topic string

func init(){
	fmt.Println("初始化kafka生产者")
	iniParser := utils.IniParser{}
	if err:=iniParser.Load(confFilename); err!=nil{
		fmt.Println("loading ini file fail",err)
		return
	}
	addr = strings.Split(iniParser.GetString("kafka","Host"),",")
	topic = iniParser.GetString("kafka","Topic")
}

//同步消息模式
func SyncProducer(message string)  {
	config := sarama.NewConfig()
	//是否等待成功和失败后的响应
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	//新建一个同步生产者
	p, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	defer p.Close()
	msg := &sarama.ProducerMessage{
		Topic:topic,
		Value:sarama.ByteEncoder(message),
	}
	part, offset, err := p.SendMessage(msg)
	if err != nil {
		log.Printf("send message(%s) err=%s \n", message, err)
	}else {
		fmt.Fprintf(os.Stdout, message + "发送成功，partition=%d, offset=%d \n", part, offset)
	}
}

// asyncProducer 异步生产者
// 并发量大时，必须采用这种方式
func AsyncProducer(message string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true //必须有这个选项
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewAsyncProducer(addr, config)
	defer p.Close()
	if err != nil {
		return
	}

	//必须有这个匿名函数内容（轮询错误信息）
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					log.Println(err)
				}
			case <-success:
			}
		}
	}(p)

	msg := &sarama.ProducerMessage{
		Topic:topic,
		Value:sarama.ByteEncoder(message),
	}
	p.Input() <- msg

}
