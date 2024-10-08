package config

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Toml struct {
	Addr    string
	Authapi string
	Nsq     NsqConfig
	Redis   RedisConfig
}

type NsqConfig struct {
	Lookupd string
	Nsqd    string
}

type RedisConfig struct {
	Server   string
	Password string
}

var TOML Toml

func init() {
	pflag.String("addr", ":8081", "tcp addr")
	f := pflag.String("toml", "./configs/config.toml", "config path")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	//viper.SetConfigName("config")
	//viper.AddConfigPath(".")
	viper.SetConfigFile(*f)

	e := viper.ReadInConfig()
	if e != nil {
		panic(e.Error())
	}
	viper.Unmarshal(&TOML) //将配置文件绑定到config上
	log.Printf("%+v\n", TOML)
}
