# 架构图

![image-20240201142048642](https://raw.githubusercontent.com/Aurora-Galaxy/image/main/image-20240201142048642.png)

# 目录结构图

```
├─config                 
├─logTransfer      # 负责消费kafka内日志，并发送到es        
│  ├─Es                  
│  ├─etcd                
│  └─kafka               
├─log_Agent        #按照配置项收集日志送入kafka     
│  ├─etcd                
│  ├─kafka               
│  └─tailLog
├─test             # 工具的基本使用
│  ├─test_es
│  ├─test_etcd
│  ├─test_kafka
│  └─test_tailf
└─utils
```

