// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: query.sql

package storage

import (
	"context"
)

const createImage = `-- name: CreateImage :one
INSERT INTO images (
  id, url, nsfw, nsfwlevel, prompt, width, height, score) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING id, url, nsfw, nsfwlevel, prompt, width, height, score
`

type CreateImageParams struct {
	ID        int64
	Url       string
	Nsfw      int64
	Nsfwlevel string
	Prompt    string
	Width     int64
	Height    int64
	Score     int64
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) (Image, error) {
	row := q.db.QueryRowContext(ctx, createImage,
		arg.ID,
		arg.Url,
		arg.Nsfw,
		arg.Nsfwlevel,
		arg.Prompt,
		arg.Width,
		arg.Height,
		arg.Score,
	)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Nsfw,
		&i.Nsfwlevel,
		&i.Prompt,
		&i.Width,
		&i.Height,
		&i.Score,
	)
	return i, err
}

const createImageTag = `-- name: CreateImageTag :one
INSERT INTO images_tags (
	image_id, tag_id
) VALUES (
  ?, ?
)
RETURNING image_id, tag_id
`

type CreateImageTagParams struct {
	ImageID int64
	TagID   int64
}

func (q *Queries) CreateImageTag(ctx context.Context, arg CreateImageTagParams) (ImagesTag, error) {
	row := q.db.QueryRowContext(ctx, createImageTag, arg.ImageID, arg.TagID)
	var i ImagesTag
	err := row.Scan(&i.ImageID, &i.TagID)
	return i, err
}

const createTag = `-- name: CreateTag :one
INSERT INTO tags (
	content
) VALUES (
  ?
)
RETURNING id, content
`

func (q *Queries) CreateTag(ctx context.Context, content string) (Tag, error) {
	row := q.db.QueryRowContext(ctx, createTag, content)
	var i Tag
	err := row.Scan(&i.ID, &i.Content)
	return i, err
}

const deleteImage = `-- name: DeleteImage :exec
DELETE FROM images
WHERE id = ?
`

func (q *Queries) DeleteImage(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteImage, id)
	return err
}

const deleteTag = `-- name: DeleteTag :exec
DELETE FROM tags
WHERE content = ?
`

func (q *Queries) DeleteTag(ctx context.Context, content string) error {
	_, err := q.db.ExecContext(ctx, deleteTag, content)
	return err
}

const getImage = `-- name: GetImage :one
SELECT id, url, nsfw, nsfwlevel, prompt, width, height, score FROM images
WHERE id = ? LIMIT 1
`

func (q *Queries) GetImage(ctx context.Context, id int64) (Image, error) {
	row := q.db.QueryRowContext(ctx, getImage, id)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Nsfw,
		&i.Nsfwlevel,
		&i.Prompt,
		&i.Width,
		&i.Height,
		&i.Score,
	)
	return i, err
}

const getImageTag = `-- name: GetImageTag :one
SELECT image_id, tag_id FROM images_tags
WHERE image_id = ? AND tag_id = ? LIMIT 1
`

type GetImageTagParams struct {
	ImageID int64
	TagID   int64
}

func (q *Queries) GetImageTag(ctx context.Context, arg GetImageTagParams) (ImagesTag, error) {
	row := q.db.QueryRowContext(ctx, getImageTag, arg.ImageID, arg.TagID)
	var i ImagesTag
	err := row.Scan(&i.ImageID, &i.TagID)
	return i, err
}

const getTag = `-- name: GetTag :one
SELECT id, content FROM tags
WHERE content = ? LIMIT 1
`

func (q *Queries) GetTag(ctx context.Context, content string) (Tag, error) {
	row := q.db.QueryRowContext(ctx, getTag, content)
	var i Tag
	err := row.Scan(&i.ID, &i.Content)
	return i, err
}

const getTagsStartingWith = `-- name: GetTagsStartingWith :many
SELECT id, content
FROM tags
WHERE content LIKE ? || '%'
`

func (q *Queries) GetTagsStartingWith(ctx context.Context, content string) ([]Tag, error) {
	rows, err := q.db.QueryContext(ctx, getTagsStartingWith, content)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tag
	for rows.Next() {
		var i Tag
		if err := rows.Scan(&i.ID, &i.Content); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listImages = `-- name: ListImages :many
SELECT id, url, nsfw, nsfwlevel, prompt, width, height, score FROM images
ORDER BY score
`

func (q *Queries) ListImages(ctx context.Context) ([]Image, error) {
	rows, err := q.db.QueryContext(ctx, listImages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ID,
			&i.Url,
			&i.Nsfw,
			&i.Nsfwlevel,
			&i.Prompt,
			&i.Width,
			&i.Height,
			&i.Score,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTags = `-- name: ListTags :many
SELECT id, content FROM tags
ORDER BY content
`

func (q *Queries) ListTags(ctx context.Context) ([]Tag, error) {
	rows, err := q.db.QueryContext(ctx, listTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tag
	for rows.Next() {
		var i Tag
		if err := rows.Scan(&i.ID, &i.Content); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const replaceImage = `-- name: ReplaceImage :one
REPLACE INTO images (
  id, url, nsfw, nsfwlevel, prompt, width, height, score) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING id, url, nsfw, nsfwlevel, prompt, width, height, score
`

type ReplaceImageParams struct {
	ID        int64
	Url       string
	Nsfw      int64
	Nsfwlevel string
	Prompt    string
	Width     int64
	Height    int64
	Score     int64
}

func (q *Queries) ReplaceImage(ctx context.Context, arg ReplaceImageParams) (Image, error) {
	row := q.db.QueryRowContext(ctx, replaceImage,
		arg.ID,
		arg.Url,
		arg.Nsfw,
		arg.Nsfwlevel,
		arg.Prompt,
		arg.Width,
		arg.Height,
		arg.Score,
	)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Nsfw,
		&i.Nsfwlevel,
		&i.Prompt,
		&i.Width,
		&i.Height,
		&i.Score,
	)
	return i, err
}

const replaceImageTag = `-- name: ReplaceImageTag :one
REPLACE INTO images_tags (
	image_id, tag_id
) VALUES (
  ?, ?
)
RETURNING image_id, tag_id
`

type ReplaceImageTagParams struct {
	ImageID int64
	TagID   int64
}

func (q *Queries) ReplaceImageTag(ctx context.Context, arg ReplaceImageTagParams) (ImagesTag, error) {
	row := q.db.QueryRowContext(ctx, replaceImageTag, arg.ImageID, arg.TagID)
	var i ImagesTag
	err := row.Scan(&i.ImageID, &i.TagID)
	return i, err
}

const replaceTag = `-- name: ReplaceTag :one
REPLACE INTO tags (
	content
) VALUES (
  ?
)
RETURNING id, content
`

func (q *Queries) ReplaceTag(ctx context.Context, content string) (Tag, error) {
	row := q.db.QueryRowContext(ctx, replaceTag, content)
	var i Tag
	err := row.Scan(&i.ID, &i.Content)
	return i, err
}
