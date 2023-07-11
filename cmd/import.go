package cmd

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rwxd/civitai-search/civitai"
	"github.com/rwxd/civitai-search/storage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, _ []string) {
		timeFrame, err := cmd.Flags().GetString("time-frame")
		if err != nil {
			panic(err)
		}

		sort, err := cmd.Flags().GetString("sort")
		if err != nil {
			panic(err)
		}

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

		err = civitai.LoadCivitaiImagesToDB(ctx, queries, 100000, 0, timeFrame, sort)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().String("time-frame", "Day", "Time frame to download new images")
	importCmd.Flags().String("sort", "Newest", "Sort images by")
	importCmd.Flags().String("db", "./db.sqlite", "Path to the sqlite database")
}
