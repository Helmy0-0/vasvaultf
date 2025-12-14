package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"vasvault/internal/dto"
	"vasvault/internal/models"
	"vasvault/internal/repositories"

	"github.com/google/uuid"
)

type FileServiceInterface interface {
	UploadFile(userID uint, file multipart.File, header *multipart.FileHeader, request dto.UploadFileRequest) (*dto.FileResponse, error)
	GetFileByID(fileID uint) (*dto.FileResponse, error)
	ListUserFiles(userID uint) ([]dto.FileResponse, error)
	DeleteFile(fileID uint) error
}

type FileService struct {
	repository repositories.FileRepositoryInterface
	basePath   string
}

func NewFileService(repo repositories.FileRepositoryInterface, basePath string) FileServiceInterface {
	return &FileService{
		repository: repo,
		basePath:   basePath,
	}
}

func (s *FileService) UploadFile(userID uint, file multipart.File, header *multipart.FileHeader, request dto.UploadFileRequest) (*dto.FileResponse, error) {
	if _, err := os.Stat(s.basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(s.basePath, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	ext := filepath.Ext(header.Filename)
	newName := uuid.New().String() + ext
	fullPath := filepath.Join(s.basePath, newName)

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	model := &models.File{
		Filename:    newName,
		Filepath:    fullPath,
		Mimetype:    header.Header.Get("Content-Type"),
		Size:        header.Size,
		UserID:      userID,
		WorkspaceID: request.FolderId,
		UploadedAt:  time.Now(),
	}

	if err := s.repository.Create(model); err != nil {
		return nil, fmt.Errorf("failed to store file metadata: %w", err)
	}

	response := dto.FileResponse{
		ID:        model.ID,
		UserId:    model.UserID,
		FolderId:  model.WorkspaceID,
		FileName:  model.Filename,
		FilePath:  model.Filepath,
		MimeType:  model.Mimetype,
		Size:      model.Size,
		CreatedAt: model.UploadedAt,
	}

	return &response, nil
}

func (s *FileService) GetFileByID(fileID uint) (*dto.FileResponse, error) {
	file, err := s.repository.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}
	response := dto.FileResponse{
		ID:        file.ID,
		UserId:    file.UserID,
		FolderId:  file.WorkspaceID,
		FileName:  file.Filename,
		FilePath:  file.Filepath,
		MimeType:  file.Mimetype,
		Size:      file.Size,
		CreatedAt: file.UploadedAt,
	}
	return &response, nil
}

func (s *FileService) ListUserFiles(userID uint) ([]dto.FileResponse, error) {
	files, err := s.repository.ListUserFiles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user file: %w", err)
	}
	var responses []dto.FileResponse
	for _, f := range files {
		responses = append(responses, dto.FileResponse{
			ID:        f.ID,
			UserId:    f.UserID,
			FolderId:  f.WorkspaceID,
			FileName:  f.Filename,
			FilePath:  f.Filepath,
			MimeType:  f.Mimetype,
			Size:      f.Size,
			CreatedAt: f.UploadedAt,
		})
	}
	return responses, err
}

func (s *FileService) DeleteFile(fileID uint) error {
	file, err := s.repository.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}
	if err := os.Remove(file.Filepath); err != nil {
		return fmt.Errorf("failed to delete file from fisk: %w", err)
	}
	if err := s.repository.Delete(fileID); err != nil {
		return fmt.Errorf("failed to delete file metadata: %w", err)
	}
	return nil
}
