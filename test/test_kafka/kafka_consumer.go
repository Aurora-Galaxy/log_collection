package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"sync"
)

func main() {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	topic := "mysql_log"
	// //根据topic获取其所有分区
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Printf("fail to get topic partition, err:%v\n", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(partitionList))
	for _, partition := range partitionList {
		// 针对每个分区创建一个对应的分区消费者
		cons, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest) // OffsetNewest 表示从最新的位置开始消费
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer cons.AsyncClose()
		go func(partitionConsumer sarama.PartitionConsumer) {
			for msg := range cons.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(cons)
	}
	wg.Wait()
}
