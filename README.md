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
logAgent 和 LogTransfer 应该分开部署

# 程序执行图
![image](https://github.com/Aurora-Galaxy/log_collection/assets/84388261/b374a76c-c478-4689-9a99-6b92b5e81c06)
![image](https://github.com/Aurora-Galaxy/log_collection/assets/84388261/ada6e53e-0b2e-48a5-8f06-ccfa7e969e95)
