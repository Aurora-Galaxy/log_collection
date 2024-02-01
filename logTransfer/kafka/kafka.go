package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"sync"
)

type TopicData struct {
	Topic string `json:"topic"`
	Text  string `json:"text"`
}

// ConsumerTask 消费者 ，每一个ConsumerTask去消费一个topic
type ConsumerTask struct {
	topic       string
	consumerObj sarama.Consumer // 消费者实例
	// 控制consumer的goroutine
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewConsumerTask(addr []string, topic string) (consumerTask *ConsumerTask, err error) {
	consumer, err := sarama.NewConsumer(addr, nil)
	ctx, cancel := context.WithCancel(context.Background())
	consumerTask = &ConsumerTask{
		topic:       topic,
		consumerObj: consumer,
		ctx:         ctx,
		cancelFunc:  cancel,
	}
	go consumerTask.ConsumeTopic()
	return
}

func (c *ConsumerTask) ConsumeTopic() {
	// 获取topic 的分区
	partitionList, err := c.consumerObj.Partitions(c.topic)
	if err != nil {
		fmt.Printf("fail to get topic partition, err:%v\n", err)
		return
	}
	fmt.Printf("获取topic:%s分区成功\n", c.topic)
	for _, partition := range partitionList {
		// 针对每个分区创建一个对应的分区消费者
		cons, err := c.consumerObj.ConsumePartition(c.topic, partition, sarama.OffsetNewest) // OffsetNewest 表示从最新的位置开始消费
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer cons.AsyncClose()
		var wg sync.WaitGroup
		wg.Add(1)
		go func(partitionConsumer sarama.PartitionConsumer) {
			for msg := range cons.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
				// 将消费的消息写入chan，然后在发送到es
				logData := &TopicData{
					Text:  string(msg.Value),
					Topic: msg.Topic,
				}
				fmt.Printf("get topic : %s msg : %s\n", msg.Topic, string(msg.Value))
				MessageChan <- logData
			}
		}(cons)
		wg.Wait()
	}
}
