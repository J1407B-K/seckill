package init

import (
	"github.com/spf13/viper"
	"seckill/rpc/user/global"
)

func SetupViper() {
	//先指定文件
	viper.SetConfigType("yaml")
	viper.SetConfigName("user")
	viper.SetConfigFile("./rpc/manifest/user.yaml")

	//读取
	err := viper.ReadInConfig()
	if err != nil {
		panic("Read config file failed, err: " + err.Error())
	}

	//数据类型转换
	err = viper.Unmarshal(&global.Config)
	if err != nil {
		panic("Unmarshal config file failed, err: " + err.Error())
	}

}
