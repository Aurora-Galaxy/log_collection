package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var Conf *Config

type Config struct {
	// yaml字段要和配置文件中的命名对应,和结构体内的字段名也要一致
	Kafka         map[string]*KafkaConfig `yaml:"kafka"` //可以配置多个kafka
	TailLog       *TailConfig             `yaml:"tailLog"`
	Etcd          map[string]*EtcdConfig  `yaml:"etcd"` // logagent，logtransfer分别配置etcd
	Elasticsearch *ESConfig               `yaml:"elasticsearch"`
}

type KafkaConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	//Topic string `yaml:"topic"`
	ChanMaxSize int `yaml:"chanMaxSize"`
}

type TailConfig struct {
	FilePath string `yaml:"filePath"`
}

type EtcdConfig struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	TimeOut int    `yaml:"timeOut"`
	Key     string `yaml:"key"`
	Size    int    `yaml:"size"`
}

type ESConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func InitConfig() {
	workDir, _ := os.Getwd()          // 获取当前根路径
	workPath := filepath.Dir(workDir) // 获取上一级目录
	workPath = filepath.Join(workPath, "config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// 将读取的配置文件反序列化
	err = viper.Unmarshal(&Conf)
	//fmt.Println(*Conf)
	if err != nil {
		panic(err)
	}
}
