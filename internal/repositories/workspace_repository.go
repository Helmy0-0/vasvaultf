package repositories

import (
	"vasvault/internal/models"

	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	CreateWithMember(workspace *models.Workspace, member *models.WorkspaceMember) error
	FindByUserID(userID uint, search string) ([]models.WorkspaceMember, error)
	FindByID(workspaceID uint) (*models.Workspace, error)
	Update(workspace *models.Workspace) error
	Delete(id uint) error
	AddMember(member *models.WorkspaceMember) error
	UpdateMember(member *models.WorkspaceMember) error
	RemoveMember(workspaceID uint, userID uint) error
	FindMember(workspaceID uint, userID uint) (*models.WorkspaceMember, error)
}

type workspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) CreateWithMember(workspace *models.Workspace, member *models.WorkspaceMember) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(workspace).Error; err != nil {
			return err
		}

		member.WorkspaceID = workspace.ID

		if err := tx.Create(member).Error; err != nil {
			return err
		}

		if err := tx.Preload("Owner").First(workspace, workspace.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *workspaceRepository) FindByUserID(userID uint, search string) ([]models.WorkspaceMember, error) {
	var memberships []models.WorkspaceMember

	query := r.db.Preload("Workspace").
		Preload("Workspace.Owner").
		Where("user_id = ?", userID)

	if search != "" {
		query = query.Joins("JOIN workspaces ON workspaces.id = workspace_members.workspace_id").
			Where("workspaces.name ILIKE ?", "%"+search+"%")
	}

	err := query.Find(&memberships).Error
	return memberships, err
}

func (r *workspaceRepository) FindByID(id uint) (*models.Workspace, error) {
	var workspace models.Workspace

	err := r.db.Preload("Memberships").
		Preload("Memberships.User").
		First(&workspace, id).Error

	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *workspaceRepository) Update(workspace *models.Workspace) error {
	return r.db.Model(workspace).Select("Name", "Description").Updates(workspace).Error
}

func (r *workspaceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Workspace{}, id).Error
}

func (r *workspaceRepository) AddMember(member *models.WorkspaceMember) error {
	return r.db.Create(member).Error
}

func (r *workspaceRepository) UpdateMember(member *models.WorkspaceMember) error {
	return r.db.Save(member).Error
}

func (r *workspaceRepository) RemoveMember(workspaceID uint, userID uint) error {
	return r.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).Delete(&models.WorkspaceMember{}).Error
}

func (r *workspaceRepository) FindMember(workspaceID uint, userID uint) (*models.WorkspaceMember, error) {
	var member models.WorkspaceMember
	err := r.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&member).Error
	return &member, err
}
