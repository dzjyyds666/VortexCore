package vortexUtil

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io"
	"os"
	"path"
	"time"
)

// 提供基础的配置文件加载功能，上层业务获取到map之后，自行进行解析即可
func LoadConfigFromToml(configPath string) (map[string]interface{}, error) {
	if path.Ext(configPath) != ".toml" {
		return nil, errors.New("config file must be toml format")
	}

	config := make(map[string]interface{})

	open, err := os.Open(configPath)
	if nil != err {
		return nil, err
	}

	configBytes, err := io.ReadAll(open)
	if nil != err {
		return nil, err
	}

	_, err = toml.Decode(string(configBytes), &config)
	if nil != err {
		return nil, err
	}
	return config, nil
}

// 从json中解析
func LoadConfigFromJson(configPath string) (map[string]interface{}, error) {
	if path.Ext(configPath) != ".json" {
		return nil, errors.New("config file must be json format")
	}
	config := make(map[string]interface{})
	open, err := os.Open(configPath)
	if nil != err {
		return nil, err
	}
	configBytes, err := io.ReadAll(open)
	if nil != err {
		return nil, err
	}
	err = json.Unmarshal(configBytes, &config)
	if nil != err {
		return nil, err
	}
	return config, nil
}

// 从etcd中获取配置文件
func LoadConfigFromEtcd(ctx context.Context, endpoint, configKey string) (map[string]interface{}, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 5 * time.Second,
	})
	if nil != err {
		return nil, err
	}
	defer func() {
		err = client.Close()
		if nil != err {
			return
		}
	}()
	config, err := client.Get(ctx, configKey)
	if nil != err {
		return nil, err
	}

	if len(config.Kvs) > 0 {
		configBytes := config.Kvs[0].Value
		return LoadConfigFromJson(string(configBytes))
	}
	return nil, errors.New("no config found")
}
