package storage

import (
	"database/sql"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Extra functions to help with the database

// GetImagesByTags queries the database for images that matches all the tags
func GetImagesByTagsCustom(db *sql.DB, tags []string) ([]Image, error) {
	log.Debug("Querying database for tags: ", tags)

	// Generate a string with '?' placeholders for each tag. Each '?' is separated by a comma.
	// The placeholders will be used in the SQL query to substitute the tag values.
	placeholders := strings.Repeat("?,", len(tags)-1) + "?"

	// Use the placeholders to construct the SQL query.
	// This query selects images that are tagged with all the tags in the provided list.
	// fmt.Sprintf is used to insert the placeholders into the query string at the correct position (%s).
	// It is safe to use fmt.Sprintf here because the placeholders string does not contain any user-provided data.
	// The '?' in the HAVING clause is a placeholder for the count of tags, to ensure that we only select images that are tagged with all the provided tags.
	query := fmt.Sprintf(`SELECT images.*
	FROM images
	JOIN images_tags ON images.id = images_tags.image_id
	JOIN tags ON images_tags.tag_id = tags.id
	WHERE tags.content IN (%s)
	GROUP BY images.id
	HAVING COUNT(DISTINCT tags.id) = ?
	ORDER BY images.score DESC`, placeholders)

	log.Debug("Generated query: ", query)

	// Create a slice to store the tags and the count of tags.
	args := make([]interface{}, len(tags)+1)

	for i, tag := range tags {
		args[i] = tag
	}

	args[len(tags)] = len(tags)

	rows, err := db.Query(query, args...)
	if err != nil {
		return []Image{}, err
	}
	defer rows.Close()

	images := []Image{}

	for rows.Next() {
		var image Image
		err = rows.Scan(&image.ID, &image.Url, &image.Nsfw, &image.Nsfwlevel, &image.Prompt, &image.Width, &image.Height, &image.Score)
		if err != nil {
			return images, err
		}
		images = append(images, image)
	}
	return images, nil
}

func GetImagesByPromptTagsCustom(db *sql.DB, tags []string) ([]Image, error) {
	log.Debug("Querying database for tags: ", tags)

	// Create a slice to store the arguments for the SQL query.
	args := make([]interface{}, len(tags))

	// Generate a SQL condition for each tag and store the tag in the args slice.
	conditions := make([]string, len(tags))
	for i, tag := range tags {
		conditions[i] = "prompt LIKE ?"
		args[i] = "%" + tag + "%"
	}

	// Join the conditions together using OR.
	conditionsStr := strings.Join(conditions, " AND ")

	// Construct the SQL query.
	query := "SELECT * FROM images WHERE " + conditionsStr

	log.Debug("Generated query: ", query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return []Image{}, err
	}
	defer rows.Close()

	images := []Image{}

	for rows.Next() {
		var image Image
		err = rows.Scan(&image.ID, &image.Url, &image.Nsfw, &image.Nsfwlevel, &image.Prompt, &image.Width, &image.Height, &image.Score)
		if err != nil {
			return images, err
		}
		images = append(images, image)
	}
	return images, nil
}
