package services

import (
	"vasvault/internal/dto"
	"vasvault/internal/models"
	"vasvault/internal/repositories"
	"errors"
)

type WorkspaceService interface {
	CreateWorkspace(userID uint, req dto.CreateWorkspaceRequest) (*models.Workspace, error)
	GetMyWorkspaces(userID uint, search string) ([]dto.WorkspaceResponse, error)
	GetWorkspaceDetail(userID uint, workspaceID uint) (*dto.WorkspaceDetailResponse, error)
	UpdateWorkspace(userID uint, workspaceID uint, req dto.UpdateWorkspaceRequest) (*models.Workspace, error)
    DeleteWorkspace(userID uint, workspaceID uint) error
	AddMember(requesterID uint, workspaceID uint, req dto.AddMemberRequest) error
    UpdateMemberRole(requesterID uint, workspaceID uint, targetUserID uint, req dto.UpdateMemberRoleRequest) error
    RemoveMember(requesterID uint, workspaceID uint, targetUserID uint) error
}

type workspaceService struct {
	repo repositories.WorkspaceRepository
	userRepo repositories.UserRepositoryInterface
}

func NewWorkspaceService(repo repositories.WorkspaceRepository, userRepo repositories.UserRepositoryInterface) WorkspaceService {
	return &workspaceService{
        repo: repo, 
        userRepo: userRepo, 
    }
}
func (s *workspaceService) CreateWorkspace(userID uint, req dto.CreateWorkspaceRequest) (*models.Workspace, error) {
	
	workspace := &models.Workspace{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
	}

	member := &models.WorkspaceMember{
		UserID: userID,
		Role:   models.RoleOwner,
	}

	err := s.repo.CreateWithMember(workspace, member)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func (s *workspaceService) GetMyWorkspaces(userID uint, search string) ([]dto.WorkspaceResponse, error) {
	memberships, err := s.repo.FindByUserID(userID, search)
	if err != nil {
		return nil, err
	}

	var responses []dto.WorkspaceResponse
	for _, m := range memberships {
		responses = append(responses, dto.WorkspaceResponse{
			ID:          m.Workspace.ID,
			Name:        m.Workspace.Name,
			Description: m.Workspace.Description,
			Role:        m.Role,
			OwnerName:   m.Workspace.Owner.Username, 
		})
	}

	return responses, nil
}

func (s *workspaceService) GetWorkspaceDetail(userID uint, workspaceID uint) (*dto.WorkspaceDetailResponse, error) {
    workspace, err := s.repo.FindByID(workspaceID)
    if err != nil {
        return nil, err
    }

    isMember := false
    for _, m := range workspace.Memberships {
        if m.UserID == userID {
            isMember = true
            break
        }
    }

    if !isMember {
 
        return nil, errors.New("you are not a member of this workspace") 
    }

  
    var memberResponses []dto.WorkspaceMemberResponse
    for _, m := range workspace.Memberships {
        memberResponses = append(memberResponses, dto.WorkspaceMemberResponse{
            UserID:   m.UserID,
            Name:     m.User.Username, 
            Email:    m.User.Email,
            Role:     m.Role,
            JoinedAt: m.JoinedAt.Format("2006-01-02"),
        })
    }

    response := &dto.WorkspaceDetailResponse{
        ID:          workspace.ID,
        Name:        workspace.Name,
        Description: workspace.Description,
        OwnerID:     workspace.OwnerID,
        Members:     memberResponses,
    }

    return response, nil
}

func (s *workspaceService) UpdateWorkspace(userID uint, workspaceID uint, req dto.UpdateWorkspaceRequest) (*models.Workspace, error) {

    workspace, err := s.repo.FindByID(workspaceID)
    if err != nil {
        return nil, err
    }

    isAuthorized := false
    for _, m := range workspace.Memberships {
        if m.UserID == userID {
        
            if m.Role == "owner" || m.Role == "admin" {
                isAuthorized = true
            }
            break
        }
    }

    if !isAuthorized {
        return nil, errors.New("unauthorized: only owner or admin can update workspace")
    }

    if req.Name != "" {
        workspace.Name = req.Name
    }
    workspace.Description = req.Description

    if err := s.repo.Update(workspace); err != nil {
        return nil, err
    }

    return workspace, nil
}

func (s *workspaceService) DeleteWorkspace(userID uint, workspaceID uint) error {
    
    workspace, err := s.repo.FindByID(workspaceID)
    if err != nil {
        return err
    }

    if workspace.OwnerID != userID {
        return errors.New("unauthorized: only owner can delete workspace")
    }

    return s.repo.Delete(workspaceID)
}

func (s *workspaceService) AddMember(requesterID uint, workspaceID uint, req dto.AddMemberRequest) error {

    requester, err := s.repo.FindMember(workspaceID, requesterID)
    if err != nil {
        return errors.New("access denied: you are not a member of this workspace")
    }
    if requester.Role != "owner" && requester.Role != "admin" {
        return errors.New("unauthorized: only owner or admin can add members")
    }

    targetUser, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return errors.New("user with this email not found")
    }

    _, err = s.repo.FindMember(workspaceID, targetUser.ID)
    if err == nil {
        return errors.New("user is already a member of this workspace")
    }


    newMember := models.WorkspaceMember{
        WorkspaceID: workspaceID,
        UserID:      targetUser.ID,
        Role:        "viewer", 
    }

    return s.repo.AddMember(&newMember)
}

func (s *workspaceService) UpdateMemberRole(requesterID uint, workspaceID uint, targetUserID uint, req dto.UpdateMemberRoleRequest) error {

    requester, err := s.repo.FindMember(workspaceID, requesterID)
    if err != nil || (requester.Role != "owner" && requester.Role != "admin") {
        return errors.New("unauthorized")
    }

    targetMember, err := s.repo.FindMember(workspaceID, targetUserID)
    if err != nil {
        return errors.New("member not found")
    }

    if targetMember.Role == "owner" {
        return errors.New("cannot change role of the owner")
    }

    targetMember.Role = req.Role
    return s.repo.UpdateMember(targetMember)
}

func (s *workspaceService) RemoveMember(requesterID uint, workspaceID uint, targetUserID uint) error {
 
    requester, err := s.repo.FindMember(workspaceID, requesterID)
    if err != nil || (requester.Role != "owner" && requester.Role != "admin") {
        return errors.New("unauthorized")
    }

    if requesterID == targetUserID { return errors.New("cannot kick yourself, please leave instead") }

    targetMember, err := s.repo.FindMember(workspaceID, targetUserID)
    if err != nil {
        return errors.New("member not found")
    }
    if targetMember.Role == "owner" {
        return errors.New("cannot remove workspace owner")
    }

    return s.repo.RemoveMember(workspaceID, targetUserID)
}