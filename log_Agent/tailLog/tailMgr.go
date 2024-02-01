package tailLog

import (
	"fmt"
	"log_project/log_Agent/etcd"
	"time"
)

// TailMsg 管理TailTask
type TailMsg struct {
	logConfigs  []*etcd.LogEntry
	taskMap     map[string]*TailTask  //保存已经正在使用的tailtask
	newConfChan chan []*etcd.LogEntry //等待最新配置，将新的配置先放入chan缓冲
}

// 全局变量，管理所有的TailTask
var tskMsg *TailMsg

func Init(logEntry []*etcd.LogEntry) {
	// 初始化
	tskMsg = &TailMsg{
		logConfigs:  logEntry,
		taskMap:     make(map[string]*TailTask, 16), //一般一台机器只收集16个配置项对应的日志
		newConfChan: make(chan []*etcd.LogEntry),    //无缓冲区chan，没有值时一直堵塞
	}
	for _, logEntry := range tskMsg.logConfigs {
		//对于每一个日志路径都创建一个tailTask去读取日志
		tailTask := NewTailTask(logEntry.Path, logEntry.Topic)
		// 创建的tailTask需要被记录，方便后续管理
		//使用path和topic作为key,因为日志配置项，这两项都有可能被更改
		mk := fmt.Sprintf("%s_%s", logEntry.Topic, logEntry.Path)
		tskMsg.taskMap[mk] = tailTask
	}
	//初始化时，使用一个后台goroutine去监听自己的chan
	go tskMsg.MonitorChan()

}

// MonitorChan 监听自己的通道，查看chan内是否有值，并作相应的处理
func (t *TailMsg) MonitorChan() {
	for {
		select {
		case newConf := <-t.newConfChan:
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Topic, conf.Path)
				_, ok := t.taskMap[mk]
				if ok {
					// 配置原本就存在，没有变化
					continue
				} else {
					// 出现新增的配置,加入map方便管理
					tsk := NewTailTask(conf.Path, conf.Topic)
					t.taskMap[mk] = tsk
				}
			}
			//找到原本存在，但是新的conf中不存在的配置项，将对应的tailtask协程 kill
			for _, c1 := range t.logConfigs { //t.logConfigs 旧配置
				isDelete := true
				for _, c2 := range newConf { //newConf 新配置
					if c2.Path == c1.Path && c2.Topic == c2.Topic {
						isDelete = false
						//配置相同，不做处理
						continue
					}
				}
				if isDelete {
					//新配置中没有，需要将该配置项对应的tailTask协程 kill
					delKey := fmt.Sprintf("%s_%s", c1.Topic, c1.Path)
					t.taskMap[delKey].cancelFunc()
					// 将map中的键值对也一并删除
					delete(t.taskMap, delKey)
				}
			}
			//	配置修改
			//  配置删除
			fmt.Println("配置更改，新的配置出现", newConf)

		default:
			time.Sleep(time.Second)
		}
	}
}

// NewConfigChan 向外部暴露自己的私有字段
func NewConfigChan() chan<- []*etcd.LogEntry {
	return tskMsg.newConfChan
}
