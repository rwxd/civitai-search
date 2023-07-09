CREATE TABLE IF NOT EXISTS images (
  id   INTEGER PRIMARY KEY,
  url  text NOT NULL UNIQUE,
  nsfw INTEGER NOT NULL,
  nsfwlevel text NOT NULL,
  prompt text NOT NULL,
  width INTEGER NOT NULL,
  height INTEGER NOT NULL,
  score INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tags (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  content  text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS  images_tags (
  image_id INTEGER NOT NULL,
  tag_id INTEGER NOT NULL,
  CONSTRAINT images_tags_imageforeign FOREIGN KEY (image_id) references images(id),
  CONSTRAINT images_tags_tag FOREIGN KEY (tag_id) references tags(id),
  CONSTRAINT images_tags_unique UNIQUE (image_id, tag_id)
);
