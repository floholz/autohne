package src

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
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

var AUTOHNE_APP_CONTEXT = os.Getenv("AUTOHNE_APP_CONTEXT")
var AUTOHNE_CONFIG_PATH = os.Getenv("AUTOHNE_CONFIG_PATH")
var AUTOHNE_VIDEOS_DIR = os.Getenv("AUTOHNE_VIDEOS_DIR")

func NewAppConfig() AppConfig {
	if AUTOHNE_APP_CONTEXT == "" {
		home, err := os.UserHomeDir()
		AUTOHNE_APP_CONTEXT = filepath.Join(home, "autohne")
		err = os.Setenv("AUTOHNE_APP_CONTEXT", AUTOHNE_APP_CONTEXT)
		if err != nil {
			log.Fatal(err)
		}
	}

	if AUTOHNE_VIDEOS_DIR == "" {
		AUTOHNE_VIDEOS_DIR = filepath.Join(AUTOHNE_APP_CONTEXT, "videos")
		err := os.Setenv("AUTOHNE_VIDEOS_DIR", AUTOHNE_VIDEOS_DIR)
		if err != nil {
			log.Fatal(err)
		}
	}

	if AUTOHNE_CONFIG_PATH == "" {
		AUTOHNE_CONFIG_PATH = filepath.Join(AUTOHNE_APP_CONTEXT, "config.yml")
		err := os.Setenv("AUTOHNE_CONFIG_PATH", AUTOHNE_CONFIG_PATH)
		if err != nil {
			log.Fatal(err)
		}
	}

	// read config
	viper.SetConfigFile(AUTOHNE_CONFIG_PATH)
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	config := AppConfig{}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
