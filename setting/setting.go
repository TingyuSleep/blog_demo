package setting

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局变量，用来保存程序的所有配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"` //注意，不要使用yaml json，统一使用mapstructure标签（使用了第三方库viper）
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
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
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init() (err error) {
	//使用viper管理配置
	/*
		方式一：直接指定配置文件路径（相对路径或绝对路径）
			viper.SetConfigFile("./conf/config.yaml")//相对路径
			viper.SetConfigFile("/users/xxx/web_frame02/conf/config.yaml")
	*/

	//方法二：指定配置文件名和配置文件的位置，viper自行查找可用的配置文件
	//viper.SetConfigFile("config.yaml")
	viper.SetConfigName("config") // 指定配置文件名称（不需要带后缀）
	viper.SetConfigType("yaml")   // 指定配置文件类型(专用于从远程获取配置信息是指定配置文件类型的)
	viper.AddConfigPath(".")      // 指定查找配置文件的路径（这里使用相对路径，.表示在当前目录（即，根目录）下找）
	viper.AddConfigPath("./conf") //再从conf目录下找配置文件

	//基本上是配合远程控制中心使用的，告诉viper当前的数据使用的什么格式去解析
	//viper.SetConfigType("json")

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}

	viper.WatchConfig() //实时监控
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil { //配置文件修改了，需要再把viper中的配置信息反序列化到结构体中
			fmt.Printf("viper.Unmarshal failed,err:%v", err)
		}
	})
	return
}
