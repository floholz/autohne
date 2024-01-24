package src

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	Download struct {
		TwitchConfig TwitchConfig `mapstructure:"twitch"`
	} `mapstructure:"download"`
	Upload struct {
		YoutubeConfig   YoutubeConfig   `mapstructure:"youtube"`
		TiktokConfig    TiktokConfig    `mapstructure:"tiktok"`
		InstagramConfig InstagramConfig `mapstructure:"instagram"`
	} `mapstructure:"upload"`
}

type TwitchConfig struct {
	BroadcasterId string `mapstructure:"broadcaster_id"`
	ClientId      string `mapstructure:"client_id"`
	BearerToken   string `mapstructure:"bearer_token"`
}

type YoutubeConfig struct {
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type TiktokConfig struct {
}

type InstagramConfig struct {
}

func NewAppConfig(configFile ...string) AppConfig {
	config := AppConfig{}
	if len(configFile) == 0 {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(configFile[0])
	}
	// read config
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
