package src

import (
	"cmp"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

type TwitchClipSortingStrategy uint

const (
	TWITCH_CLIP_SORT_BY_CREATED_AT TwitchClipSortingStrategy = 0
	TWITCH_CLIP_SORT_BY_VIEWS                                = 1
	TWITCH_CLIP_SORT_BY_DURATION                             = 2
)

type TwitchClipSortingDirection uint

const (
	DESC TwitchClipSortingDirection = 0
	ASC                             = 1
)

type TwitchClipFilterOptions struct {
	Limit         int
	MaxDuration   float32
	CreatedAfter  time.Time
	CreatedBefore time.Time
	GameId        string
	Language      string
	IncludesTitle string
}

type TwitchClipArray []TwitchClip

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
	Data       TwitchClipArray `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

type TwitchApi struct {
	BaseUrl string
	Config  TwitchConfig
}

func NewTwitchApi(config TwitchConfig) TwitchApi {
	twitch := TwitchApi{}
	twitch.BaseUrl = "https://api.twitch.tv/helix"
	twitch.Config = config
	return twitch
}

func (twitch *TwitchApi) GetClips() TwitchClipArray {
	query := "?" + strings.Join([]string{
		"broadcaster_id=" + twitch.Config.BroadcasterId,
		"started_at=" + time.Now().Add(-24*time.Hour).Truncate(24*time.Hour).UTC().Format(time.RFC3339),
		"first=100",
	}, "&")

	req, err := http.NewRequest(http.MethodGet, twitch.BaseUrl+"/clips"+query, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+twitch.Config.BearerToken)
	req.Header.Set("Client-Id", twitch.Config.ClientId)

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

	return resBody.Data
}

func (clips TwitchClipArray) SortClips(sortingStrategy TwitchClipSortingStrategy, sortDirection ...TwitchClipSortingDirection) TwitchClipArray {
	if len(sortDirection) == 0 {
		sortDirection = []TwitchClipSortingDirection{DESC}
	}
	sorted := clips
	slices.SortFunc(sorted, func(a, b TwitchClip) int {
		result := 0
		switch sortingStrategy {
		case TWITCH_CLIP_SORT_BY_CREATED_AT:
			return int(a.CreatedAt.Sub(b.CreatedAt).Nanoseconds())
		case TWITCH_CLIP_SORT_BY_VIEWS:
			return cmp.Compare(b.ViewCount, a.ViewCount)
		case TWITCH_CLIP_SORT_BY_DURATION:
			return cmp.Compare(b.Duration, a.Duration)
		}
		if sortDirection[0] == ASC {
			result *= -1
		}
		return result
	})
	return sorted
}

func (clips TwitchClipArray) FilterClips(options TwitchClipFilterOptions) TwitchClipArray {
	var filtered TwitchClipArray
	for _, clip := range clips {
		if options.MaxDuration != 0 && options.MaxDuration <= clip.Duration {
			continue
		}
		if !options.CreatedAfter.IsZero() && options.CreatedAfter.After(clip.CreatedAt) {
			continue
		}
		if !options.CreatedBefore.IsZero() && options.CreatedBefore.Before(clip.CreatedAt) {
			continue
		}
		if options.GameId != "" && options.GameId != clip.GameId {
			continue
		}
		if options.Language != "" && options.Language != clip.Language {
			continue
		}
		if !strings.Contains(clip.Title, options.IncludesTitle) {
			continue
		}
		filtered = append(filtered, clip)

		if options.Limit <= len(filtered) {
			break
		}
	}
	return filtered
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

	clipTitle := Pathify(clip.Title)
	SaveToDisk(body, "@videos/clips/"+clipTitle+".mp4")
}

func (clips TwitchClipArray) SaveJson() {
	SaveToDisk([]byte(clips.String()), "@videos/clips/@random.json")
}

func (clips TwitchClipArray) String() string {
	jayson, _ := json.MarshalIndent(clips, "", "  ")
	return string(jayson)
}

func (clip *TwitchClip) String() string {
	jayson, _ := json.MarshalIndent(clip, "", "  ")
	return string(jayson)
}
