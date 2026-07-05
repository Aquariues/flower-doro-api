# Flower Images

Flower image files should live in S3 or an S3-compatible object store. The API stores only URLs and metadata in Postgres.

## Data Model

The `flowers` table stores:

- `asset_name`: local/app fallback name, for example `daisy`.
- `image_url`: full-size image URL.
- `thumbnail_url`: smaller image URL for book grids or compact garden views.

The app should prefer:

1. `thumbnail_url` for small list/book/garden previews.
2. `image_url` for larger flower detail views.
3. bundled local asset from `asset_name` when URLs are empty or offline.

## Recommended S3 Layout

```text
s3://flowerdoro-assets/flowers/daisy.png
s3://flowerdoro-assets/flowers/daisy-thumb.png
s3://flowerdoro-assets/flowers/rose.png
s3://flowerdoro-assets/flowers/rose-thumb.png
```

Public CDN/object URLs can then be saved in GoAdmin or through the flower API.

## API Example

```json
{
  "kind": "daisy",
  "asset_name": "daisy",
  "image_url": "https://cdn.example.com/flowers/daisy.png",
  "thumbnail_url": "https://cdn.example.com/flowers/daisy-thumb.png"
}
```

