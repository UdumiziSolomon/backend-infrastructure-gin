package config

import (
	"log"

	"github.com/spf13/viper"
)


type Config struct{
	ServerPort     	   string	  `mapstructure:"PORT"`
	ClientUrl 	   	 	string		`mapstructure:"CLIENT_URL"`
	LocalDBUri   	 string	     `mapstructure:"LOCAL_DB_URI"`
	DBName  	 	 string 	`mapstructure:"DB_NAME"`
}

func LoadENVFile(path string) (config Config, err error){
	viper.AddConfigPath(path)

	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil{
		return
	}

	if err := viper.Unmarshal(&config); err != nil{
		log.Fatal(err)
	}
	return
}