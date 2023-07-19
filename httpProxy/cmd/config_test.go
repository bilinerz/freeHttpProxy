package cmd

import (
	"encoding/json"
	"os"
	"testing"
)

func clearConfigFile() {
	os.Remove(configPath)
}

//保存配置
func TestSaveConfig(t *testing.T) {
	clearConfigFile()
	config := Config{
		ServerAddr: ":1992",
		ClientAddr: ":1080",
		Key:        RandKey(),
	}
	config.SaveConfig()

	file, err := os.Open(configPath)
	if err != nil {
		t.Errorf("打开配置文件 %s 出错:%s", configPath, err)
	}
	defer file.Close()

	tmp := make(map[string]string)
	err = json.NewDecoder(file).Decode(&tmp)
	if err != nil {
		t.Error(err)
	}

	t.Log(tmp)
}

//测试读取配置
func TestReadConfig(t *testing.T) {
	config := Config{}
	config.ReadConfig()
	t.Log(config)
}
