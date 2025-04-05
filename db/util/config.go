package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

//stores all configurations of the application
//values are read by viper from a config file or env file
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}


//It reads config from file or environment var
func LoadConfig(path string)(config Config, err error){
	viper.SetConfigName("app") // Name without extension
	viper.SetConfigType("env") // File type
	viper.AddConfigPath(path)   // Path to the directory containing the config file

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Error while unmarshalling response",err)
	}
	return
}
