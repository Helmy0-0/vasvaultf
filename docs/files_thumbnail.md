# GET /api/v1/files/:id/thumbnail

Method: GET

URL: `/api/v1/files/:id/thumbnail`

Auth: Bearer (required)

Description:

Returns a small cached thumbnail for the specified file. Thumbnails are generated on first request and cached under `./uploads/thumbs/` as `<filename>.thumb.jpg`.

Behavior:
- Only supported for files with an image `Content-Type` (e.g. `image/jpeg`, `image/png`). Non-image requests return HTTP 400.
- Thumbnails are generated at 200x200 px (center-cropped) and saved as JPEG with 80% quality.
- If a cached thumbnail exists, it is served directly to avoid repeated image processing.

Examples:

curl (using bearer token):

```bash
curl -H "Authorization: Bearer <token>" \
  -o thumb.jpg \
  http://localhost:8080/api/v1/files/123/thumbnail
```

Client usage notes:
- For summary lists, prefer requesting thumbnails (`/thumbnail`) rather than full files to reduce bandwidth.
- If you expect many thumbnails on list pages, pre-generate thumbnails at upload time or use lazy-loading in the UI.
