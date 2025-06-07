package vortexUtil

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 提供基础的配置文件加载功能，上层业务获取到map之后，自行进行解析即可
func LoadConfigFromToml(configPath string, config any) error {
	if path.Ext(configPath) != ".toml" {
		return errors.New("config file must be toml format")
	}

	open, err := os.Open(configPath)
	if nil != err {
		return err
	}

	configBytes, err := io.ReadAll(open)
	if nil != err {
		return err
	}

	_, err = toml.Decode(string(configBytes), config)
	if nil != err {
		return err
	}
	return nil
}

// 从json中解析
func LoadConfigFromJson(configPath string, config any) error {
	if path.Ext(configPath) != ".json" {
		return errors.New("config file must be json format")
	}
	open, err := os.Open(configPath)
	if nil != err {
		return err
	}
	configBytes, err := io.ReadAll(open)
	if nil != err {
		return err
	}
	err = json.Unmarshal(configBytes, config)
	if nil != err {
		return err
	}
	return nil
}

// 从etcd中获取配置文件
func LoadConfigFromEtcd(ctx context.Context, endpoint, configKey string, config any) error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 5 * time.Second,
	})
	if nil != err {
		return err
	}
	defer func() {
		err = client.Close()
		if nil != err {
			return
		}
	}()
	configInfo, err := client.Get(ctx, configKey)
	if nil != err {
		return err
	}

	if len(configInfo.Kvs) > 0 {
		configBytes := configInfo.Kvs[0].Value
		err = json.Unmarshal(configBytes, config)
		if nil != err {
			return err
		}
		return nil
	}
	return errors.New("no config found")
}
