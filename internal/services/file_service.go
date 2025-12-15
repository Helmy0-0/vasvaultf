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
	ListUserFilesWithOptionalCategory(userID uint, categoryID *uint) ([]dto.FileResponse, error)
	DeleteFile(fileID uint) error
	AssignCategories(userID, fileID uint, categoryIDs []uint) error
	RemoveCategories(userID, fileID uint, categoryIDs []uint) error
	UpdateCategories(userID, fileID uint, categoryIDs []uint) error
	GetStorageSummary(userID uint) (*dto.StorageSummaryResponse, error)
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
		WorkspaceID: request.WorkspaceId,
		UploadedAt:  time.Now(),
	}

	if err := s.repository.Create(model); err != nil {
		return nil, fmt.Errorf("failed to store file metadata: %w", err)
	}

	// Assign categories jika ada
	if len(request.CategoryIDs) > 0 {
		if err := s.repository.AssignCategories(model.ID, request.CategoryIDs); err != nil {
			return nil, fmt.Errorf("failed to assign categories: %w", err)
		}
	}

	// Reload file dengan categories untuk response
	fileWithCategories, err := s.repository.FindByIDWithCategories(model.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load file with categories: %w", err)
	}

	var categories []dto.CategorySimple
	for _, cat := range fileWithCategories.Categories {
		categories = append(categories, dto.CategorySimple{
			ID:    cat.ID,
			Name:  cat.Name,
			Color: cat.Color,
		})
	}

	response := dto.FileResponse{
		ID:          model.ID,
		UserId:      model.UserID,
		WorkspaceId: model.WorkspaceID,
		FileName:    model.Filename,
		FilePath:    model.Filepath,
		MimeType:   model.Mimetype,
		Size:       model.Size,
		Categories: categories,
		CreatedAt:  model.UploadedAt,
	}

	return &response, nil
}

func (s *FileService) GetFileByID(fileID uint) (*dto.FileResponse, error) {
	file, err := s.repository.FindByIDWithCategories(fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	var categories []dto.CategorySimple
	for _, cat := range file.Categories {
		categories = append(categories, dto.CategorySimple{
			ID:    cat.ID,
			Name:  cat.Name,
			Color: cat.Color,
		})
	}

	response := dto.FileResponse{
		ID:          file.ID,
		UserId:      file.UserID,
		WorkspaceId: file.WorkspaceID,
		FileName:    file.Filename,
		FilePath:    file.Filepath,
		MimeType:   file.Mimetype,
		Size:       file.Size,
		Categories: categories,
		CreatedAt:  file.UploadedAt,
	}
	return &response, nil
}

func (s *FileService) ListUserFiles(userID uint) ([]dto.FileResponse, error) {
	files, err := s.repository.ListUserFilesWithCategories(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user file: %w", err)
	}
	var responses []dto.FileResponse
	for _, f := range files {
		var categories []dto.CategorySimple
		for _, cat := range f.Categories {
			categories = append(categories, dto.CategorySimple{
				ID:    cat.ID,
				Name:  cat.Name,
				Color: cat.Color,
			})
		}

		responses = append(responses, dto.FileResponse{
			ID:          f.ID,
			UserId:      f.UserID,
			WorkspaceId: f.WorkspaceID,
			FileName:    f.Filename,
			FilePath:    f.Filepath,
			MimeType:   f.Mimetype,
			Size:       f.Size,
			Categories: categories,
			CreatedAt:  f.UploadedAt,
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

// AssignCategories menambahkan kategori ke file (tidak menghapus kategori yang sudah ada)
func (s *FileService) AssignCategories(userID, fileID uint, categoryIDs []uint) error {
	// Validasi file milik user
	file, err := s.repository.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found")
	}
	if file.UserID != userID {
		return fmt.Errorf("unauthorized: file does not belong to user")
	}

	if err := s.repository.AssignCategories(fileID, categoryIDs); err != nil {
		return fmt.Errorf("failed to assign categories: %w", err)
	}
	return nil
}

// RemoveCategories menghapus kategori tertentu dari file
func (s *FileService) RemoveCategories(userID, fileID uint, categoryIDs []uint) error {
	// Validasi file milik user
	file, err := s.repository.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found")
	}
	if file.UserID != userID {
		return fmt.Errorf("unauthorized: file does not belong to user")
	}

	if err := s.repository.RemoveCategories(fileID, categoryIDs); err != nil {
		return fmt.Errorf("failed to remove categories: %w", err)
	}
	return nil
}

// UpdateCategories mengganti semua kategori file dengan yang baru
func (s *FileService) UpdateCategories(userID, fileID uint, categoryIDs []uint) error {
	// Validasi file milik user
	file, err := s.repository.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found")
	}
	if file.UserID != userID {
		return fmt.Errorf("unauthorized: file does not belong to user")
	}

	// Clear semua kategori lama
	if err := s.repository.ClearAllCategories(fileID); err != nil {
		return fmt.Errorf("failed to clear categories: %w", err)
	}

	// Assign kategori baru
	if len(categoryIDs) > 0 {
		if err := s.repository.AssignCategories(fileID, categoryIDs); err != nil {
			return fmt.Errorf("failed to assign new categories: %w", err)
		}
	}

	return nil
}

func (s *FileService) ListUserFilesWithOptionalCategory(userID uint, categoryID *uint) ([]dto.FileResponse, error) {
	files, err := s.repository.ListUserFilesWithOptionalCategory(userID, categoryID)
	if err != nil {
		return nil, err
	}

	var response []dto.FileResponse
	for _, file := range files {
		var categories []dto.CategorySimple
		for _, cat := range file.Categories {
			categories = append(categories, dto.CategorySimple{
				ID:    cat.ID,
				Name:  cat.Name,
				Color: cat.Color,
			})
		}

		response = append(response, dto.FileResponse{
			ID:          file.ID,
			UserId:      file.UserID,
			WorkspaceId: file.WorkspaceID,
			FileName:    file.Filename,
			FilePath:    file.Filepath,
			MimeType:   file.Mimetype,
			Size:       file.Size,
			Categories: categories,
			CreatedAt:  file.UploadedAt,
		})
	}
	
	return response, nil
}

func (s *FileService) GetStorageSummary(userID uint) (*dto.StorageSummaryResponse, error) {
	// Max storage: 5 GiB
	const maxBytes int64 = 5 * 1024 * 1024 * 1024

	used, err := s.repository.TotalUserStorage(userID)
	if err != nil {
		return nil, err
	}
	files, err := s.repository.GetLatestFilesForUser(userID, 10)
	if err != nil {
		files = []models.File{}
	}

	var latestDtos []dto.FileResponse
	for _, latest := range files {
		var categories []dto.CategorySimple
		for _, cat := range latest.Categories {
			categories = append(categories, dto.CategorySimple{ID: cat.ID, Name: cat.Name, Color: cat.Color})
		}
		f := dto.FileResponse{
			ID:         latest.ID,
			UserId:     latest.UserID,
			WorkspaceId:latest.WorkspaceID,
			FileName:   latest.Filename,
			FilePath:   latest.Filepath,
			MimeType:   latest.Mimetype,
			Size:       latest.Size,
			Categories: categories,
			CreatedAt:  latest.UploadedAt,
		}
		latestDtos = append(latestDtos, f)
	}

	remaining := maxBytes - used
	if remaining < 0 {
		remaining = 0
	}

	return &dto.StorageSummaryResponse{
		MaxBytes:       maxBytes,
		UsedBytes:      used,
		RemainingBytes: remaining,
		LatestFiles:    latestDtos,
	}, nil
}