package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rwxd/civitai-search/civitai"
	"github.com/rwxd/civitai-search/storage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// open db
		db, err := sql.Open("sqlite3", "db.sqlite")
		if err != nil {
			panic(err)
		}

		// create tables
		if _, err := db.ExecContext(ctx, storage.Ddl); err != nil {
			panic(err)
		}

		queries := storage.New(db)
		images, err := queries.ListImages(ctx)
		if err != nil {
			panic(err)
		}
		log.Infof("%d images already in the database \n", len(images))

		GetCivitaiImages(ctx, queries, 1000)
	},
}

func GetCivitaiImages(ctx context.Context, queries *storage.Queries, number int) ([]storage.Image, error) {
	cursor := 1
	output := []storage.Image{}
	for i := 1; i < number; i++ {
		images, err := civitai.GetCivitaiImages(cursor)
		if err != nil {
			log.Warn(err)
			if strings.Contains(err.Error(), "response") {
				time.Sleep(10 * time.Second)
			}
			continue
		}
		cursor = images.Metadata.NextCursor

		for _, image := range images.Items {
			if image.Meta.Prompt == "" {
				log.Debug("Image has empty prompt, skipping")
				continue
			}

			imageTags := sanitizePrompt(image.Meta.Prompt)

			if len(imageTags) > 0 {
				nsfw := 0
				if image.Nsfw {
					nsfw = 1
				}
				score := image.Stats.LikeCount + image.Stats.HeartCount + (image.Stats.CommentCount * 2) + image.Stats.CryCount + image.Stats.LaughCount - image.Stats.DislikeCount
				imageParams := storage.CreateImageParams{
					ID:        int64(image.ID),
					Url:       image.URL,
					Nsfw:      int64(nsfw),
					Nsfwlevel: image.NsfwLevel,
					Prompt:    image.Meta.Prompt,
					Width:     int64(image.Width),
					Height:    int64(image.Height),
					Score:     int64(score),
				}

				dbImage, err := queries.GetImage(ctx, imageParams.ID)
				if err != nil {
					log.Infof("Creating image %s", imageParams.Url)
					dbImage, err = queries.CreateImage(ctx, imageParams)
					if err != nil {
						log.Error(err)
						continue
					}
				}
				for _, tag := range imageTags {
					tag = strings.ToLower(tag)
					dbTag, err := queries.GetTag(ctx, tag)
					if err != nil {
						log.Debugf("Creating tag \"%s\"\n", tag)
						dbTag, err = queries.CreateTag(ctx, tag)
						if err != nil {
							log.Error(err)
							continue
						}
					}

					_, err = queries.GetImageTag(ctx, storage.GetImageTagParams{ImageID: dbImage.ID, TagID: dbTag.ID})
					if err != nil {
						log.Debugf("Creating tag - image relation")
						_, err := queries.CreateImageTag(ctx, storage.CreateImageTagParams{ImageID: dbImage.ID, TagID: dbTag.ID})
						if err != nil {
							log.Error(err)
							continue
						}
					}
				}
			}
		}
	}
	return output, nil
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

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
