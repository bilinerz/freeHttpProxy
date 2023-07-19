package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/bilinerz/freeHttpProxy/httpProxy"
	"io/ioutil"
	"os"
	"path"
)

var (
	// PackagePath 项目路径
	PackagePath string
	// 配置文件路径
	configPath string
	// 默认的配置文件名称
	configFilename string
	//随机生成的32位密钥
	randKey string
)

type Config struct {
	ServerAddr string `json:"serverAddr"`
	ClientAddr string `json:"clientAddr"`
	Key        string `json:"key"`
}

func init() {
	PackagePath, _ = os.Getwd()
	configFilename = "config.json"
	configPath = path.Join(PackagePath, configFilename)
	randKey = httpProxy.RandKey()
}

// ReadConfig 读取配置
func (config *Config) ReadConfig() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("配置文件 %s 不存在，创建默认配置\n", configPath)
		conf := Config{
			ServerAddr: ":1992",
			ClientAddr: ":1080",
			Key:        randKey,
		}
		conf.SaveConfig()
	}

	fmt.Printf("从路径 %s 中读取配置\n", configPath)
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("打开配置文件 %s 出错:%s", configPath, err)
		return
	}
	defer file.Close()

	//创建Json编码器
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		fmt.Printf("JSON 配置文件格式不正确:\n%s", file.Name())
		return
	}
}

// SaveConfig 保存配置到配置文件
func (config *Config) SaveConfig() {
	configJson, _ := json.MarshalIndent(config, "", "	")
	err := ioutil.WriteFile(configPath, configJson, 0644)
	if err != nil {
		fmt.Printf("保存配置到文件 %s 出错: %s", configPath, err)
		return
	}
	fmt.Printf("保存配置到文件 %s 成功\n", configPath)
}
