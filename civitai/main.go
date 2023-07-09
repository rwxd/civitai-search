package civitai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

type CivitaiAnswer struct {
	Items    []CivitaiAnswerItem   `json:"items,omitempty"`
	Metadata CivitaiAnswerMetadata `json:"metadata,omitempty"`
}

type CivitaiAnswerItem struct {
	ID        int        `json:"id,omitempty"`
	URL       string     `json:"url,omitempty"`
	Hash      string     `json:"hash,omitempty"`
	Width     int        `json:"width,omitempty"`
	Height    int        `json:"height,omitempty"`
	NsfwLevel string     `json:"nsfwLevel,omitempty"`
	Nsfw      bool       `json:"nsfw,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	PostID    int        `json:"postId,omitempty"`
	Stats     ImageStats `json:"stats,omitempty"`
	Meta      ImageMeta  `json:"meta,omitempty"`
}

type ImageStats struct {
	CryCount     int `json:"cryCount,omitempty"`
	LaughCount   int `json:"laughCount,omitempty"`
	LikeCount    int `json:"likeCount,omitempty"`
	DislikeCount int `json:"dislikeCount,omitempty"`
	HeartCount   int `json:"heartCount,omitempty"`
	CommentCount int `json:"commentCount,omitempty"`
}

type ImageMeta struct {
	Ensd     string  `json:"ENSD,omitempty"`
	Size     string  `json:"Size,omitempty"`
	Seed     int64   `json:"seed,omitempty"`
	Model    string  `json:"Model,omitempty"`
	Steps    int     `json:"steps,omitempty"`
	Prompt   string  `json:"prompt,omitempty"`
	Sampler  string  `json:"sampler,omitempty"`
	CfgScale float32 `json:"cfgScale,omitempty"`
}

type CivitaiAnswerMetadata struct {
	NextCursor int    `json:"nextCursor,omitempty"`
	NextPage   string `json:"nextPage,omitempty"`
}

func GetCivitaiImages(cursor int) (*CivitaiAnswer, error) {
	params := url.Values{}
	params.Add("limit", "200")
	params.Add("cursor", fmt.Sprintf("%d", cursor))
	url := "https://civitai.com/api/v1/images?" + params.Encode()

	log.Info("Requesting url ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &CivitaiAnswer{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &CivitaiAnswer{}, err
	}

	if resp.StatusCode != 200 {
		log.Warn(string(body))
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	log.Debug("Unmarshaling answer")
	var answer CivitaiAnswer
	err = json.Unmarshal(body, &answer)
	if err != nil {
		return &answer, err
	}

	return &answer, nil
}
