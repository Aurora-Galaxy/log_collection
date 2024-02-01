package main

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	//导入 search 包时 注意辨别 是/core/search 还是 /eql/search
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"log"
	"strconv"
	"time"
)

// Review 评价数据
type Review struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userID"`
	Score       uint8     `json:"score"`
	Content     string    `json:"content"`
	Tags        []Tag     `json:"tags"`
	Status      int       `json:"status"`
	PublishTime time.Time `json:"publishDate"`
}

type Tag struct {
	Code  int    `json:"code"`
	Title string `json:"title"`
}

func main() {
	// Es 配置，指定需要连接的服务地址
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	// 创建客户端连接
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Fatalln("elasticsearch.NewTypedClient failed, err : ", err)
		return
	}
	fmt.Println("ES connect success")
	createIndex(client, "mysql_log")
	//createDocument(client, "review")
	//getDocument(client, "review", "1")
	//searchContent(client, "review", "好评")

}

func createIndex(client *elasticsearch.TypedClient, index string) {
	resp, err := client.Indices.Create(index).Do(context.Background())
	if err != nil {
		log.Fatalln("create index failed , err : ", err)
		return
	}
	// %#v 会将类型信息和详细的结构都打印出来
	fmt.Printf("index : %#v\n", resp.Index)
}

// 创建document
func createDocument(client *elasticsearch.TypedClient, index string) {
	//定义需要创建的document的实例
	review := Review{
		ID:      1,
		UserID:  147982601,
		Score:   5,
		Content: "这是一个好评！",
		Tags: []Tag{
			{1000, "好评"},
			{1100, "物超所值"},
			{9000, "有图"},
		},
		Status:      2,
		PublishTime: time.Now(),
	}
	// 添加文档
	resp, err := client.Index(index).Id(strconv.Itoa(int(review.ID))).Document(review).Do(context.Background())
	if err != nil {
		log.Fatalln("add document failed , err : ", err)
		return
	}
	fmt.Printf("index : %#v\n", resp.Result)
}

func getDocument(client *elasticsearch.TypedClient, index, id string) {
	resp, err := client.Get(index, id).Do(context.Background())
	if err != nil {
		log.Fatalln("get document failed , err : ", err)
		return
	}
	fmt.Printf("document : %s\n", resp.Source_)
}

// 搜索包含相应内容的文档
func searchContent(client *elasticsearch.TypedClient, index, content string) {
	resp, err := client.Search().Index(index).Request(&search.Request{
		Query: &types.Query{
			MatchPhrase: map[string]types.MatchPhraseQuery{
				"content": {Query: content},
			},
		},
	}).Do(context.Background())
	if err != nil {
		log.Fatalln("search document failed , err : ", err)
		return
	}
	// 打印到符合条件的个数
	fmt.Printf("document total: %d\n", resp.Hits.Total.Value)
	// 遍历所有结果
	for _, v := range resp.Hits.Hits {
		fmt.Printf("%s\n", v.Source_)
	}
}
