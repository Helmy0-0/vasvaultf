package repositories

import (
	"vasvault/internal/models"

	"gorm.io/gorm"
)

type FileRepositoryInterface interface {
	Create(file *models.File) error
	FindByID(id uint) (*models.File, error)
	FindByIDWithCategories(id uint) (*models.File, error)
	ListUserFiles(userID uint) ([]models.File, error)
	ListUserFilesWithCategories(userID uint) ([]models.File, error)
	ListUserFilesWithOptionalCategory(userID uint, categoryID *uint) ([]models.File, error)
	Delete(fileID uint) error
	AssignCategories(fileID uint, categoryIDs []uint) error
	RemoveCategories(fileID uint, categoryIDs []uint) error
	ClearAllCategories(fileID uint) error
	TotalUserStorage(userID uint) (int64, error)
	GetLatestFileForUser(userID uint) (*models.File, error)
	GetLatestFilesForUser(userID uint, limit int) ([]models.File, error)
}

type FileRepository struct {
	db *gorm.DB
}

func (r *FileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
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

func (r *FileRepository) FindByIDWithCategories(id uint) (*models.File, error) {
	var file models.File
	if err := r.db.Preload("Categories").Where("id = ?", id).First(&file).Error; err != nil {
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

func (r *FileRepository) ListUserFilesWithCategories(userID uint) ([]models.File, error) {
	var files []models.File
	if err := r.db.Preload("Categories").Where("user_id = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) Delete(fileID uint) error {
	return r.db.Delete(&models.File{}, fileID).Error
}

// AssignCategories menambahkan kategori ke file
func (r *FileRepository) AssignCategories(fileID uint, categoryIDs []uint) error {
	var file models.File
	if err := r.db.First(&file, fileID).Error; err != nil {
		return err
	}

	var categories []models.Category
	if err := r.db.Find(&categories, categoryIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&file).Association("Categories").Append(&categories)
}

// RemoveCategories menghapus kategori dari file
func (r *FileRepository) RemoveCategories(fileID uint, categoryIDs []uint) error {
	var file models.File
	if err := r.db.First(&file, fileID).Error; err != nil {
		return err
	}

	var categories []models.Category
	if err := r.db.Find(&categories, categoryIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&file).Association("Categories").Delete(&categories)
}

// ClearAllCategories menghapus semua kategori dari file
func (r *FileRepository) ClearAllCategories(fileID uint) error {
	var file models.File
	if err := r.db.First(&file, fileID).Error; err != nil {
		return err
	}

	return r.db.Model(&file).Association("Categories").Clear()
}

func (r *FileRepository) ListUserFilesWithOptionalCategory(userID uint, categoryID *uint) ([]models.File, error) {
	var files []models.File
	query := r.db.Preload("Categories").Where("user_id = ?", userID)

	if categoryID != nil {
		query = query.Joins("JOIN file_categories ON file_categories.file_id = files.id").Where("file_categories.category_id = ?", *categoryID)
	}
	if err := query.Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) TotalUserStorage(userID uint) (int64, error) {
	var total int64
	if err := r.db.Model(&models.File{}).Select("COALESCE(SUM(size),0)").Where("user_id = ?", userID).Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *FileRepository) GetLatestFileForUser(userID uint) (*models.File, error) {
	var file models.File
	if err := r.db.Preload("Categories").Where("user_id = ?", userID).Order("uploaded_at desc").First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) GetLatestFilesForUser(userID uint, limit int) ([]models.File, error) {
	var files []models.File
	if err := r.db.Preload("Categories").Where("user_id = ?", userID).Order("uploaded_at desc").Limit(limit).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
