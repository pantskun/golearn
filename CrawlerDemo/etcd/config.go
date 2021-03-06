package etcd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/pantskun/golearn/CrawlerDemo/xlogutil"
)

// GetClientConfig
// 获取etcd配置文件
func GetClientConfig(filePath string) clientv3.Config {
	const timeSecond = 5.0

	defaultConfig := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: timeSecond * time.Second,
	}

	f, err := os.Open(filePath)
	if err != nil {
		xlogutil.Warning("can not find config file, use default config\n", defaultConfig)
		return defaultConfig
	}

	fb, _ := ioutil.ReadAll(f)

	type EtcdConfigs struct {
		Endpoint    string
		DialTimeout time.Duration
	}

	var fileConfigs EtcdConfigs
	if err := json.Unmarshal(fb, &fileConfigs); err != nil {
		xlogutil.Error(err)
		return defaultConfig
	}

	if fileConfigs.Endpoint == "" || fileConfigs.DialTimeout == 0 {
		return defaultConfig
	}

	defaultConfig.Endpoints = []string{fileConfigs.Endpoint}
	defaultConfig.DialTimeout = fileConfigs.DialTimeout * time.Second

	return defaultConfig
}
