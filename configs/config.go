package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigType struct {
	Port      int
	Prefix    string
	StaticDir string `yaml:"staticDir"`

	LogDir         string `yaml:"logDir"`
	LogReserveTime int    `yaml:"logReserveTime"`

	TokenValidTime         int64 `yaml:"tokenValidTime"`
	TokenExceedRefreshTime int64 `yaml:"tokenExceedRefreshTime"`
}

var Config ConfigType

const template = `
prefix: "/trail"  # 路由前缀
port: 9528  # 启动端口
staticDir: "static"  # 前端静态页面

logDir: "./logs"        # 日志目录
logReserveTime: 30  # 日志保留时间(d)

tokenValidTime: 7200  # 令牌有效时间(s)
tokenExceedRefreshTime: 86400  # 令牌超过刷新时间(s)
`

func init() {
	configFile := "./config.yml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		os.Create(configFile)
		os.WriteFile(configFile, []byte(template), 0777)
		data, _ = os.ReadFile(configFile)
	}

	if err := yaml.Unmarshal([]byte(data), &Config); err != nil {
		panic(err)
	}
}
