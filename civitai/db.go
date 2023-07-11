package civitai

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rwxd/civitai-search/storage"
	log "github.com/sirupsen/logrus"
)

func LoadCivitaiImagesToDB(ctx context.Context, queries *storage.Queries, number int, cursor int, requestTimeFrame string, requestSort string) error {
	for i := 1; i < number; i++ {
		images, err := GetCivitaiImages(cursor, requestTimeFrame, requestSort)
		if err != nil {
			log.Warnf("Failed to get Civitai images: %v", err)
			if strings.Contains(err.Error(), "response") {
				time.Sleep(10 * time.Second)
			}
			continue
		}
		if images.Metadata.NextCursor < cursor {
			log.Info("Provided next cursor is less than the current")
			return fmt.Errorf("cursor reset, all images scraped")
		}
		cursor = images.Metadata.NextCursor

		for _, image := range images.Items {
			if image.Meta.Prompt == "" {
				log.Debug("Skipping image due to empty prompt")
				continue
			}

			imageTags := sanitizePrompt(image.Meta.Prompt)

			if len(imageTags) > 0 {
				nsfw := 0
				if image.Nsfw {
					nsfw = 1
				}

				score := image.Stats.LikeCount + image.Stats.HeartCount + (image.Stats.CommentCount * 2) + image.Stats.CryCount + image.Stats.LaughCount - image.Stats.DislikeCount

				imageParams := storage.ReplaceImageParams{
					ID:        int64(image.ID),
					Url:       image.URL,
					Nsfw:      int64(nsfw),
					Nsfwlevel: image.NsfwLevel,
					Prompt:    image.Meta.Prompt,
					Width:     int64(image.Width),
					Height:    int64(image.Height),
					Score:     int64(score),
				}

				log.Infof("Replacing image with URL: %s", imageParams.Url)
				dbImage, err := queries.ReplaceImage(ctx, imageParams)
				if err != nil {
					log.Errorf("Failed to replace image: %v", err)
					continue
				}

				for _, tag := range imageTags {
					tag = strings.ToLower(tag)
					dbTag, err := queries.GetTag(ctx, tag)
					if err != nil {
						log.Debugf("Creating tag \"%s\"", tag)
						dbTag, err = queries.CreateTag(ctx, tag)
						if err != nil {
							log.Errorf("Failed to create tag: %v", err)
							continue
						}
					}

					_, err = queries.GetImageTag(ctx, storage.GetImageTagParams{ImageID: dbImage.ID, TagID: dbTag.ID})
					if err != nil {
						log.Debugf("Creating tag - image relation")
						_, err := queries.CreateImageTag(ctx, storage.CreateImageTagParams{ImageID: dbImage.ID, TagID: dbTag.ID})
						if err != nil {
							log.Errorf("Failed to create tag-image relation: %v", err)
							continue
						}
					}
				}
			}
		}
	}
	return nil
}

func sanitizePrompt(prompt string) []string {
	output := []string{}
	tags_by_comma := strings.Split(prompt, ",")
	for _, tag := range tags_by_comma {
		for _, replace := range []string{"(", ")", "^", "[", "]", "|"} {
			tag = strings.ReplaceAll(tag, replace, "")
		}

		lora, err := extractLora(tag)
		if err == nil {
			tag = lora
		}

		removedColon, err := extractColon(tag)
		if err == nil {
			tag = removedColon
		}

		tag = strings.TrimSpace(tag)
		if len(tag) > 0 {
			output = append(output, tag)
		}
	}
	return output
}

func extractLora(text string) (string, error) {
	re := regexp.MustCompile(`<(\w+):([^:>]*)(?::([^:>]*))?(?::([^:>]*))?>`)
	match := re.FindStringSubmatch(text)
	if len(match) > 0 {
		return match[0], nil
	}
	return "", fmt.Errorf("did not match")
}

func extractColon(text string) (string, error) {
	re := regexp.MustCompile(`(\w+)\:`)
	match := re.FindStringSubmatch(text)
	if len(match) > 0 {
		return match[1], nil
	}
	return "", fmt.Errorf("did not match")
}
