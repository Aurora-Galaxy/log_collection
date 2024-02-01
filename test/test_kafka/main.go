package main

import (
	"github.com/IBM/sarama"
	"log"
)

// 使用sarama 连接kafka
func main() {
	//进行sarama的配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //等待leader和follower全部ack后继续进行下一步
	config.Producer.Partitioner = sarama.NewRandomPartitioner //随机选择一个分区
	config.Producer.Return.Successes = true                   //ack通过Return的success返回

	// 构造一个消息
	msg := &sarama.ProducerMessage{
		Topic: "test_log",
		//val需要经过特殊编码
		Value: sarama.StringEncoder("this is a test log"),
	}
	//构造一个生产者连接kafka
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config) //可以连接多个
	if err != nil {
		log.Fatalln(err)
	}
	// 防止关闭连接的时候出现错误
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	}

}
