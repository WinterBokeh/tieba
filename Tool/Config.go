package Tool

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	AppPort string `json:"app_port"`
	AppHost string `json:"app_host"`
	Database DatabaseConfig `json:"database"`
}
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	User 	 string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port	 string `json:"port"`
	DbName 	 string `json:"db_name"`
	ShowSql	 bool   `json:"show_sql"`
	Charset  string `json:"charset"`
}

var cfg *Config

//获取全局配置文件
func GetCfg() *Config {
	return cfg
}

//解析配置文件
func ParseCfg(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	return nil
}
