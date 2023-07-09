package storage

import (
	_ "embed"
)

//go:embed schema.sql
var Ddl string
