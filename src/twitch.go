package src

import (
	"cmp"
	"encoding/json"
	"github.com/flytam/filenamify"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

type TwitchClip struct {
	Id              string    `json:"id"`
	Url             string    `json:"url"`
	EmbedUrl        string    `json:"embed_url"`
	BroadcasterId   string    `json:"broadcaster_id"`
	BroadcasterName string    `json:"broadcaster_name"`
	CreatorId       string    `json:"creator_id"`
	CreatorName     string    `json:"creator_name"`
	VideoId         string    `json:"video_id"`
	GameId          string    `json:"game_id"`
	Language        string    `json:"language"`
	Title           string    `json:"title"`
	ViewCount       int       `json:"view_count"`
	CreatedAt       time.Time `json:"created_at"`
	ThumbnailUrl    string    `json:"thumbnail_url"`
	Duration        float32   `json:"duration"`
	VodOffset       int       `json:"vod_offset"`
	IsFeatured      bool      `json:"is_featured"`
}

type GetClipsResponse struct {
	Data       []TwitchClip `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

type TwitchConfig struct {
	Streamer struct {
		BroadcasterId string `mapstructure:"broadcaster_id"`
	} `yaml:"streamer"`
	Me struct {
		BearerToken string `mapstructure:"bearer_token"`
		ClientId    string `mapstructure:"client_id"`
	} `yaml:"me"`
}

type TwitchApi struct {
	BaseUrl string
	Config  TwitchConfig
}

func MakeTwitchApi() TwitchApi {
	twitch := TwitchApi{}
	// read config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&twitch.Config)
	if err != nil {
		log.Fatal(err)
	}
	// set defaults
	twitch.BaseUrl = "https://api.twitch.tv/helix"
	return twitch
}

func (twitch *TwitchApi) GetClips() []TwitchClip {
	query := "?" + strings.Join([]string{
		"broadcaster_id=" + twitch.Config.Streamer.BroadcasterId,
		"started_at=" + time.Now().Add(-24*time.Hour).Truncate(24*time.Hour).UTC().Format(time.RFC3339),
		"first=100",
	}, "&")

	req, err := http.NewRequest(http.MethodGet, twitch.BaseUrl+"/clips"+query, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+twitch.Config.Me.BearerToken)
	req.Header.Set("Client-Id", twitch.Config.Me.ClientId)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	resBody := GetClipsResponse{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		log.Fatal(err)
	}

	slices.SortFunc(resBody.Data,
		func(a, b TwitchClip) int {
			return cmp.Compare(b.ViewCount, a.ViewCount)
		})

	filteredClips := func() []TwitchClip {
		var filtered []TwitchClip
		for _, clip := range resBody.Data[:25] {
			if clip.Duration <= 20 {
				filtered = append(filtered, clip)
			}
		}
		return filtered
	}()

	return filteredClips
}

func (clip *TwitchClip) DownloadClip() {
	clipUrl := clip.ThumbnailUrl[:strings.Index(clip.ThumbnailUrl, "-preview")] + ".mp4"

	res, err := http.Get(clipUrl)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	clipTitle, _ := filenamify.Filenamify(clip.Title, filenamify.Options{})
	clipTitle = strings.ReplaceAll(clipTitle, " ", "_")

	err = os.WriteFile("assets/.ignore/out/clip_"+clipTitle+".mp4", body, 0644)
	if err != nil {
		log.Fatal(err)
	}

}
