package dto

import "time"

type UploadFileRequest struct {
	FolderId *uint `json:"folder_id" binding:"omitempty"`
}

type FileResponse struct {
	ID 			uint 		`json:"id"`
	UserId		uint 		`json:"user_id"`
	FolderId	*uint		`json:"folder_id" binding:"omitempty"`
	FileName 	string		`json:"file_name"`
	FilePath	string		`json:"file_path"`
	MimeType	string		`json:"mime_type"`
	Size		int64  		`json:"size"`
	CreatedAt 	time.Time 	`json:"created_at"`
}