package dto

import "time"

type UploadFileRequest struct {
	WorkspaceId *uint  `json:"workspace_id" form:"workspace_id" binding:"omitempty"`
	CategoryIDs []uint `json:"category_ids" form:"category_ids[]" binding:"omitempty"`
}

type CategorySimple struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type FileResponse struct {
	ID          uint             `json:"id"`
	UserId      uint             `json:"user_id"`
	WorkspaceId *uint            `json:"workspace_id" binding:"omitempty"`
	FileName    string           `json:"file_name"`
	FilePath    string           `json:"file_path"`
	MimeType    string           `json:"mime_type"`
	Size        int64            `json:"size"`
	Categories  []CategorySimple `json:"categories,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

type AssignCategoriesRequest struct {
	CategoryIDs []uint `json:"category_ids" binding:"required"`
}

type StorageSummaryResponse struct {
	MaxBytes       int64          `json:"max_bytes"`
	UsedBytes      int64          `json:"used_bytes"`
	RemainingBytes int64          `json:"remaining_bytes"`
	LatestFiles    []FileResponse `json:"latest_files,omitempty"`
}
