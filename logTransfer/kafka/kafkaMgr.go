package kafka

import (
	"fmt"
	"log"
	"log_project/config"
	"log_project/logTransfer/etcd"
	"time"
)

var MessageChan chan *TopicData

type ConsumerTaskMgr struct {
	addr         []string
	topics       []*etcd.TopicData
	taskMap      map[string]*ConsumerTask
	newTopicChan chan []*etcd.TopicData
}

// 全局变量，管理所有的TailTask
var consumerMsg *ConsumerTaskMgr

func Init(addr []string, topics []*etcd.TopicData) {
	// 初始化
	consumerMsg = &ConsumerTaskMgr{
		addr:         addr,
		topics:       topics,
		taskMap:      make(map[string]*ConsumerTask, 16), //可以同时收集16个topic
		newTopicChan: make(chan []*etcd.TopicData),       //无缓冲区chan，没有值时一直堵塞
	}
	MessageChan = make(chan *TopicData, config.Conf.Kafka["consumer"].ChanMaxSize)
	for _, topic := range consumerMsg.topics {
		//对于每一个日志路径都创建一个tailTask去读取日志
		consumerTask, err := NewConsumerTask(addr, topic.Topic)
		if err != nil {
			log.Fatalf("create topic(%s) consumer failed , err : %s\n", topic.Topic, err)
		}
		// 创建的consumerTask需要被记录，方便后续管理
		//使用topic作为key
		consumerMsg.taskMap[consumerTask.topic] = consumerTask
	}
	//初始化时，使用一个后台goroutine去监听自己的chan
	go consumerMsg.MonitorChan()

}

// MonitorChan 监听自己的通道，查看chan内是否有值，并作相应的处理
func (t *ConsumerTaskMgr) MonitorChan() {
	for {
		select {
		case newTopic := <-t.newTopicChan:
			for _, conf := range newTopic {
				_, ok := t.taskMap[conf.Topic]
				if ok {
					// 配置原本就存在，没有变化
					continue
				} else {
					// 出现新增的配置,加入map方便管理
					tsk, err := NewConsumerTask(t.addr, conf.Topic)
					if err != nil {
						log.Fatalf("create topic(%s) consumer failed , err : %s\n", conf, err)
					}
					t.taskMap[conf.Topic] = tsk
				}
			}
			//找到原本存在，但是新的topics 中不存在的topic，将对应的task协程 kill
			for _, c1 := range t.topics { //t.logConfigs 旧配置
				isDelete := true
				for _, c2 := range newTopic { //newConf 新配置
					if c1.Topic == c2.Topic {
						isDelete = false
						//配置相同，不做处理
						continue
					}
				}
				if isDelete {
					//新配置中没有，需要将该配置项对应的tailTask协程 kill
					fmt.Printf("del topic : %s\n", c1.Topic)
					t.taskMap[c1.Topic].cancelFunc()
					// 将map中的键值对也一并删除
					delete(t.taskMap, c1.Topic)
				}
			}
			//	配置修改
			//  配置删除
			fmt.Println("配置更改，新的配置出现", newTopic)

		default:
			time.Sleep(time.Second)
		}
	}
}

// NewTopicChan 向外部暴露自己的私有字段
func NewTopicChan() chan<- []*etcd.TopicData {
	return consumerMsg.newTopicChan
}
