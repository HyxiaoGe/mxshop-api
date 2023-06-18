package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"time"
)

//如何将线上和线下的配置文件隔离
//不用改任何代码而且线上和线上的配置文件能隔离开

//type MysqlConfig struct {
//	Host string `mapstructure:"host"`
//	Port int    `mapstructure:"port"`
//}

//type ServerConfig struct {
//	ServiceName string      `mapstructure:"name"`
//	MysqlInfo   MysqlConfig `mapstructure:"mysql"`
//}

func InitConfig() {
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息为: &v", global.ServerConfig)

	//viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置信息产生变化: &v", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息为: &v", global.ServerConfig)
	})

	time.Sleep(time.Second * 300)
}
