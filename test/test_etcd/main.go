package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func main() {
	//配置etcd客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"}, //指定连接的集群
		DialTimeout: 5 * time.Second,            //设置超时
	})
	if err != nil {
		log.Fatalln("connect to etcd failed , err = ", err)
		return
	}
	fmt.Println("connect to etcd success")
	// 延迟关闭客户端连接
	defer cli.Close()
	//put操作
	//设置context超时操作
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//如果put 1s后，没有返回相应的结果，取消操作
	//val := `[{"path":"d:/tmp/nginx.log","topic":"nginx_log"},{"path":"d:/tmp/redis.log","topic":"redis_log"},{"path":"d:/tmp/mysql.log","topic":"mysql_log"}]`
	//val := `[{"path":"d:/tmp/nginx.log","topic":"nginx_log"},{"path":"d:/tmp/redis.log","topic":"redis_log"}]`
	val := `[{"topic":"nginx_log"},{"topic":"mysql_log"}]`
	//_, err = cli.Put(ctx, "/logagent/192.168.2.104/collect_log", val)
	_, err = cli.Put(ctx, "/logtransfer/192.168.2.104/collect_log", val)
	cancel()
	if err != nil {
		log.Fatalln("put to etcd failed , err = ", err)
		return
	}
	/*
		//put操作
			//设置context超时操作
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			//如果put 1s后，没有返回相应的结果，取消操作
			_, err = cli.Put(ctx, "cyn", "123")
			cancel()
			if err != nil {
				log.Fatalln("put to etcd failed , err = ", err)
				return
			}

			//get
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			resp, err := cli.Get(ctx, "cyn")
			cancel()
			if err != nil {
				log.Fatalln("get to etcd failed , err = ", err)
				return
			}
			//获取get请求的相应内容（键值对）
			for _, v := range resp.Kvs {
				fmt.Printf("key = %s , val = %s", v.Key, v.Value)
			}
	*/

	//watch操作
	// watch 会一直监视 相应的键值对的变化(新增，修改，删除)，不需要context控制超时
	//resChan := cli.Watch(context.Background(), "cyn") //返回的为一个channel
	//// 从监听的通道中去取相应的值
	//for wresp := range resChan {
	//	// 获取事件的类型和相应的key value
	//	for _, v := range wresp.Events {
	//		fmt.Printf("Type : %s Key : %s Value : %s\n", v.Type, v.Kv.Key, v.Kv.Value)
	//	}
	//}

}
