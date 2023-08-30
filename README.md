# Civitai Prompt Search

## Database Code Generation

```bash
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

```bash
cd storage
sqlc generate
```

## Indexing the Database

```bash
podman exec -ti civitai-search_civitai-search_1 ./main import --db /db/database.sqlite --sort "Most Reactions" --time-frame "AllTime"
```
