-- name: GetImage :one
SELECT * FROM images
WHERE id = ? LIMIT 1;

-- name: ListImages :many
SELECT * FROM images
ORDER BY score;

-- name: GetImagesWithPrompt :many
SELECT *
FROM images
WHERE prompt LIKE '%' || $1 || '%'

-- name: CreateImage :one
INSERT INTO images (
  id, url, nsfw, nsfwlevel, prompt, width, height, score) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: ReplaceImage :one
REPLACE INTO images (
  id, url, nsfw, nsfwlevel, prompt, width, height, score) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteImage :exec
DELETE FROM images
WHERE id = ?;

-- name: GetTag :one
SELECT * FROM tags
WHERE content = ? LIMIT 1;

-- name: GetTagsStartingWith :many
SELECT *
FROM tags
WHERE content LIKE ? || '%';

-- name: ListTags :many
SELECT * FROM tags
ORDER BY content;

-- name: CreateTag :one
INSERT INTO tags (
	content
) VALUES (
  ?
)
RETURNING *;

-- name: ReplaceTag :one
REPLACE INTO tags (
	content
) VALUES (
  ?
)
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE content = ?;


-- name: CreateImageTag :one
INSERT INTO images_tags (
	image_id, tag_id
) VALUES (
  ?, ?
)
RETURNING *;


-- name: ReplaceImageTag :one
REPLACE INTO images_tags (
	image_id, tag_id
) VALUES (
  ?, ?
)
RETURNING *;

-- name: GetImageTag :one
SELECT * FROM images_tags
WHERE image_id = ? AND tag_id = ? LIMIT 1;

