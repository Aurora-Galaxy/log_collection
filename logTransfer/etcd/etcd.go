package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientV3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type TopicData struct {
	Topic string
}

// 全局连接
var cli *clientV3.Client

// Init 初始化etcd
func Init(addr []string, dialTime time.Duration) (err error) {
	cli, err = clientV3.New(clientV3.Config{
		Endpoints:   addr,     //指定连接的集群
		DialTimeout: dialTime, //设置超时
	})
	return
}

// GetConfig 从etcd中获取key的配置项
func GetConfig(key string) (topicsConfs []*TopicData, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		log.Fatalln("get to etcd failed , err = ", err)
		return
	}
	//获取get请求的相应内容（键值对）
	for _, v := range resp.Kvs {
		err = json.Unmarshal(v.Value, &topicsConfs)
		if err != nil {
			log.Fatalln("json unmarshal failed , err = ", err)
		}
	}
	return
}

// WatchConf 实时监控配置信息的变化
func WatchConf(key string, newTopicChan chan<- []*TopicData) {
	//watch操作
	// watch 会一直监视 相应的键值对的变化(新增，修改，删除)，不需要context控制超时
	resChan := cli.Watch(context.Background(), key) //返回的为一个channel
	// 从监听的通道中去取相应的值
	for wresp := range resChan {
		// 获取事件的类型和相应的key value
		for _, v := range wresp.Events {
			fmt.Printf("Type : %s Key : %s Value : %s\n", v.Type, v.Kv.Key, v.Kv.Value)
			//有新配置，通知consumer
			// 删除key会返回nil
			var newTopic []*TopicData
			if v.Type != clientV3.EventTypeDelete {
				err := json.Unmarshal(v.Kv.Value, &newTopic)
				if err != nil {
					log.Fatalln("watchConf unmarshal failed , err : ", err)
					//continue
				}
				fmt.Printf("get new topic : %s\n", newTopic)
				newTopicChan <- newTopic
			}
		}
	}
}
