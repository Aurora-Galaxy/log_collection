package main

import (
	"fmt"
	"log"
	"log_project/config"
	"log_project/logTransfer/Es"
	"log_project/logTransfer/etcd"
	"log_project/logTransfer/kafka"
	"log_project/utils"
	"sync"
	"time"
)

func main() {
	//config init
	config.InitConfig()
	kafkaAddr := []string{config.Conf.Kafka["consumer"].Host + ":" + config.Conf.Kafka["consumer"].Port}
	esAddr := []string{"http://" + config.Conf.Elasticsearch.Host + ":" + config.Conf.Elasticsearch.Port}
	etcdAddr := []string{config.Conf.Etcd["logtransfer"].Host + ":" + config.Conf.Etcd["logtransfer"].Port}
	timeOut := time.Duration(config.Conf.Etcd["logtransfer"].TimeOut) * time.Second

	// etcd init
	err := etcd.Init(etcdAddr, timeOut)
	if err != nil {
		log.Fatalln("connect to etcd failed , err = ", err)
		return
	}
	fmt.Println("connect to etcd success")

	// 为了实现每个logtransfer都拉取自己独有的配置，所以要以自己的ip地址作为区分
	ip, err := utils.GetOutboundIp()
	if err != nil {
		log.Fatalln("get local ip failed , err : ", err)
	}
	// 将etcd 配置文件里的%s 替换为对应的ip
	etcdConfKey := fmt.Sprintf(config.Conf.Etcd["logtransfer"].Key, ip)
	//从etcd中获取配置项,配置项中的topic需要在kafka提前创建
	topics, err := etcd.GetConfig(etcdConfKey)
	if err != nil {
		log.Fatalln("etcd get config err : ", err)
	}
	fmt.Printf("get topics : %s\n", topics)
	// es init
	err = Es.Init(esAddr)
	if err != nil {
		log.Fatalln("es init failed , err : ", err)
	}
	fmt.Println("es init success")

	//kafkaMgr init
	kafka.Init(kafkaAddr, topics)
	fmt.Println("kafkaMgr init success")

	// 监控key的变化
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etcdConfKey, kafka.NewTopicChan())
	wg.Wait()
}
