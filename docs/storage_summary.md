# GET /api/v1/storage/summary

Method: GET

URL: /api/v1/storage/summary

Auth: Bearer (required)

Description: Returns user's storage usage (max 5 GiB), used bytes, remaining bytes, and the latest uploaded file (if any).

Response (200):

```json
{
  "max_bytes": 5368709120,
  "used_bytes": 12345678,
  "remaining_bytes": 5356363442,
  "latest_file": [{ "id":1, "file_name":"abc.pdf", "size":12345678 }]
}
```
