package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"time"
)

type LogData struct {
	topic string
	data  string
}

var (
	kafkaClient sarama.SyncProducer
	LogChan     chan *LogData
)

// Init 初始化kafka连接
func Init(addr []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //等待leader和follower全部ack后继续进行下一步
	config.Producer.Partitioner = sarama.NewRandomPartitioner //随机选择一个分区
	config.Producer.Return.Successes = true                   //ack通过Return的success返回

	//连接kafka
	kafkaClient, err = sarama.NewSyncProducer(addr, config) //可以连接多个
	if err != nil {
		return err
	}
	//初始化LogChan
	LogChan = make(chan *LogData, maxSize)
	//开启后台协程从chan中取数据发送给kafka
	go sendToKafka()
	return
}

// SendToChan 构造chan，将日志先写入chan
func SendToChan(topic, data string) {
	msg := &LogData{
		topic: topic,
		data:  data,
	}
	//将msg写入chan
	LogChan <- msg
}

// sendToKafka 向kafka发送信息
func sendToKafka() {
	//循环从通道内取数据
	for {
		select {
		case val := <-LogChan:
			//构造需要发送的消息
			msg := &sarama.ProducerMessage{
				Topic: val.topic,
				//val需要经过特殊编码
				Value: sarama.StringEncoder(val.data),
			}
			//发送
			partition, offset, err := kafkaClient.SendMessage(msg)
			if err != nil {
				log.Fatalln("send message failed! err : ", err)
			}
			fmt.Println("partition :", partition, " offset: ", offset)
			fmt.Println("send success")
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}
