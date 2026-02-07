package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局变量，存放所有配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID uint16 `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	Log   *LogConfig   `mapstructure:"log"`
	MySQL *MySQLConfig `mapstructure:"mysql"`
	Redis *RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// InitViper 仅负责加载配置文件到内存中
func InitViper() (err error) {

	viper.SetConfigFile("./conf/config.yaml")
	viper.SetConfigName("config")
	//viper.SetConfigType("yaml")//配置文件可以有多个类型，如json，yaml，toml
	viper.AddConfigPath("./conf")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("读取配置失败: %v", err)
	}

	if err := viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("viper.Unmarshal failed, err:%v", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("检测到配置文件变动: %s\n", in.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}

	})
	return
}
