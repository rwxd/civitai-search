package web

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rwxd/civitai-search/storage"
	log "github.com/sirupsen/logrus"
)

func StartServer(ip string, port string, dbPath string) {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	ctx := context.Background()

	// open db
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	// create tables
	if _, err := db.ExecContext(ctx, storage.Ddl); err != nil {
		panic(err)
	}

	queries := storage.New(db)

	r.GET("/", func(c *gin.Context) {
		log.Info("Root endpoint called")
		c.HTML(http.StatusOK, "index.html.tmpl", gin.H{})
	})

	r.GET("/search", func(c *gin.Context) {
		tags := c.QueryArray("tag")
		log.Info("Query for tags: ", tags)

		images, err := storage.GetImagesByTagsCustom(db, tags)
		if err != nil {
			log.Fatal(err)
		}

		data := GalleryData{
			Tags:   tags,
			Images: images,
		}

		c.HTML(http.StatusOK, "gallery.html.tmpl", data)
	})

	r.GET("/tags", func(c *gin.Context) {
		content := c.Query("content")
		log.Info("Searching for tag starting with: ", content)

		tags, err := queries.GetTagsStartingWith(ctx, content)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, tags)
	})

	log.Info("Starting server on ", ip, ":", port)
	r.Run(ip + ":" + port)
}

type GalleryData struct {
	Tags   []string
	Images []storage.Image
}
