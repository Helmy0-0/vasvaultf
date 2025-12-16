# POST /api/v1/files

Method: POST (multipart/form-data)

URL: /api/v1/files

Auth: Bearer (required)

Form fields:

- `file` (file, required)
- `workspace_id` (optional)
- `category_ids` (optional, JSON array)

Response (200):

```json
{ "id":1, "file_name":"abc.pdf", "file_path":"/uploads/..","size":12345 }
```
