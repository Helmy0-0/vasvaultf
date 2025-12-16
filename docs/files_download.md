# GET /api/v1/files/:id/download

Method: GET

URL: `/api/v1/files/:id/download`

Auth: Bearer (required)

Description:

Returns the raw file contents for the given file id. The endpoint performs the same authorization checks as other protected file endpoints (the user must own the file or have permission to access it).

Behavior:
- Responds with the file binary. Gin's `c.File()` is used, so the `Content-Type` header is set based on the file and the content is streamed.
- For large files the response is streamed directly from disk; the server does not load the whole file into memory.

Examples:

curl (using bearer token):

```bash
curl -H "Authorization: Bearer <token>" \
  -o myfile.jpg \
  http://localhost:8080/api/v1/files/123/download
```

Notes:
- If you want browsers to open the file inline or force download, adjust client behavior. The current implementation uses `c.File` which typically displays inline for recognized types.
- If you need `Content-Disposition` control (attachment vs inline), we can update the handler to set it explicitly.
