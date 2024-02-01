package main

import (
	"fmt"
	"log"
	"log_project/config"
	"log_project/log_Agent/etcd"
	"log_project/log_Agent/kafka"
	"log_project/log_Agent/tailLog"
	"log_project/utils"
	"sync"
	"time"
)

func main() {
	//config init
	config.InitConfig()
	fmt.Println(config.Conf)
	kafkaAddr := []string{config.Conf.Kafka["producer"].Host + ":" + config.Conf.Kafka["producer"].Port}
	etcdAddr := []string{config.Conf.Etcd["logagent"].Host + ":" + config.Conf.Etcd["logagent"].Port}
	timeOut := time.Duration(config.Conf.Etcd["logagent"].TimeOut) * time.Second
	fmt.Println(etcdAddr)
	//kafka init
	err := kafka.Init(kafkaAddr, config.Conf.Kafka["producer"].ChanMaxSize)
	if err != nil {
		log.Fatalln("kafka init err : ", err)
	}
	fmt.Println("kafka init success")

	// etcd init
	err = etcd.Init(etcdAddr, timeOut)
	if err != nil {
		log.Fatalln("etcd init err : ", err)
	}
	fmt.Println("etcd init success")

	// 为了实现每个logagent都拉取自己独有的配置，所以要以自己的ip地址作为区分
	ip, err := utils.GetOutboundIp()
	if err != nil {
		log.Fatalln("get local ip failed , err : ", err)
	}
	// 将etcd 配置文件里的%s 替换为对应的ip
	etcdConfKey := fmt.Sprintf(config.Conf.Etcd["logagent"].Key, ip)
	//从etcd中获取配置项,配置项中的topic需要在kafka提前创建
	logEntryConf, err := etcd.GetConfig(etcdConfKey)
	if err != nil {
		log.Fatalln("etcd get config err : ", err)
	}

	// tailLog init
	tailLog.Init(logEntryConf)
	var wg sync.WaitGroup
	wg.Add(1)
	//派一个哨兵监视日志收集项的变化（有变化及时通知logAgent实现热加载配置）
	go etcd.WatchConf(etcdConfKey, tailLog.NewConfigChan())
	wg.Wait()

}
