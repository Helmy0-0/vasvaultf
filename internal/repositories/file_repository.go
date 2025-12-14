package repositories

import (
	"gorm.io/gorm"
	"vasvault/internal/models"
)

type FileRepositoryInterface interface {
	Create(file *models.File) error
	FindByID(id uint) (*models.File, error)
	ListUserFiles(userID uint) ([]models.File, error)
	Delete(fileID uint) error
}

type FileRepository struct {
	db *gorm.DB
}

// Create implements FileRepositoryInterface.
func (r *FileRepository) Create(file *models.File) error {
	return 	r.db.Create(file).Error
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Upload(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *FileRepository) FindByID(id uint) (*models.File, error) {
	var file models.File
	if err := r.db.Where("id = ?", id).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) ListUserFiles(userID uint) ([]models.File, error) {
	var files []models.File
	if err := r.db.Where("user_id = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) Delete(fileID uint) error {
	return r.db.Delete(&models.File{}, fileID).Error
}
