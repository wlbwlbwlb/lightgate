package config

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Toml struct {
	Port int
}

var TOML Toml

func init() {
	pflag.Int("port", 8081, "listen port")
	f := pflag.StringP("file", "f", "./configs/config.toml", "config path")

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
