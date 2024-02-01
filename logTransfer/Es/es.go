package Es

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"log_project/logTransfer/kafka"
	"time"
)

var client *elasticsearch.TypedClient

func Init(esAddr []string) (err error) {
	// Es 配置，指定需要连接的服务地址
	cfg := elasticsearch.Config{
		Addresses: esAddr,
	}
	// 创建客户端连接
	client, err = elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Fatalln("elasticsearch.NewClient failed, err : ", err)
		return
	}

	fmt.Println("es init success")
	go SendToEs()
	return
}

// SendToEs 从chan中取出数据发送到es
func SendToEs() {
	for {
		select {
		case data := <-kafka.MessageChan:
			//var jsonData interface{}
			//err := json.Unmarshal([]byte(data.Text), &jsonData)
			//if err != nil {
			//	log.Fatalln("Failed to parse JSON document , err : ", err)
			//	continue
			//}
			resp, err := client.Index(data.Topic).Document(data).Do(context.Background())
			if err != nil {
				log.Fatalln("add document failed , err : ", err)
				continue
			}
			fmt.Printf("es add document success , index : %#v\n", resp.Result)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}
