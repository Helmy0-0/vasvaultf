package dto

import "time"

type CreateWorkspaceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type WorkspaceResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uint      `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`
	Role        string    `json:"my_role"`
	CreatedAt   time.Time `json:"created_at"`
}

type WorkspaceMemberResponse struct {
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
}

type WorkspaceDetailResponse struct {
	ID          uint                      `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	OwnerID     uint                      `json:"owner_id"`
	Members     []WorkspaceMemberResponse `json:"members"`
}

type UpdateWorkspaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required"`
}