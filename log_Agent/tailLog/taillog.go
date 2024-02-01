package tailLog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"log_project/log_Agent/kafka"
)

//var TailClient *tail.Tail

/*
每一个日志文件初始化一个tailObj去读取日志，所以不能使用全局初始化
*/

// TailTask 管理不同的taillog
type TailTask struct {
	path    string
	topic   string
	tailObj *tail.Tail //创建一个读取日志的实例
	// 使用context控制 TailTask 的goroutine
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewTailTask 创建TailTask实例
func NewTailTask(path, topic string) (tailTask *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailTask = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	_ = tailTask.init()
	return
}

func (t *TailTask) init() (err error) {
	config := tail.Config{
		Location: &tail.SeekInfo{
			//从文件末尾开始读取
			Offset: 0,
			Whence: 2,
		}, //指定文件的读取位置
		ReOpen:      true,  //如果文件被重命名或移动，这个选项决定了是否需要重新打开文件
		MustExist:   false, //当读取的文件不存在时不报错
		Poll:        false, //是否使用轮询去检查文件变化
		Pipe:        false, //说明该文件是否是命名管道，命名管道允许两个进程相互通信
		RateLimiter: nil,   //在读取文件中的数据时，它可以限制数据的流动速度
		Follow:      true,  // 允许持续地获取新的数据行，而不是仅获取已经存在的数据行
		MaxLineSize: 0,     //允许持续地获取新的数据行，而不是仅获取已经存在的数据行。
		Logger:      nil,   //用于记录过程中发生的各种事件或错误
	}
	t.tailObj, err = tail.TailFile(t.path, config) //按照上面的配置读取日志文件
	if err != nil {
		fmt.Println("tail file failed , err:", err)
		return err
	}
	//起协程，读取日志并发送
	go t.run() // 当run函数执行完成，其对应的goroutine也会退出
	return
}

// 收集日志送入kafka
func (t *TailTask) run() {
	for {
		select {
		case <-t.ctx.Done():
			fmt.Printf("tailTask : %s_%s 结束了。。。\n", t.topic, t.path)
			//终止 run  函数
			return
		case line := <-t.tailObj.Lines:
			//将数据发送到kafka，但是该方法读取日志的速度会受到发送给kafka的速度的影响
			//err := kafka.SendToKafka(t.topic, line.Text)
			//优化：将日志读取后放入channel，发送时从chan中取出发送
			fmt.Printf("file : %v   message : %v\n", t.tailObj.Filename, line.Text)
			kafka.SendToChan(t.topic, line.Text)
			//发送数据在kafka初始化时由goroutine负责
		}
	}
}
