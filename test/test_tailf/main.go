package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

// 测试tailf的用法
func main() {
	filename := "./my.log" //需要读取的日志文件
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
	tails, err := tail.TailFile(filename, config) //按照上面的配置读取日志文件
	if err != nil {
		fmt.Println("tail file failed , err:", err)
		return
	}
	var (
		line *tail.Line //按行读取日志
		ok   bool
	)
	for {
		line, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen , filename = %s\n", filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("line:", line.Text) //打印读取到的内容
	}
}
