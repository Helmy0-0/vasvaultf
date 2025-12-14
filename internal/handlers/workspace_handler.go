package handlers

import (
	"net/http"
	"vasvault/internal/dto"
	"vasvault/internal/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

type WorkspaceHandler struct {
	service services.WorkspaceService
}

func NewWorkspaceHandler(service services.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{service: service}
}

func (h *WorkspaceHandler) Create(c *gin.Context) {

	userIDCtx, exists := c.Get("userID") 
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDCtx.(uint)

	var req dto.CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspace, err := h.service.CreateWorkspace(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Workspace created",
		"data":    workspace,
	})
}

func (h *WorkspaceHandler) List(c *gin.Context) {
	userIDCtx, _ := c.Get("userID")
	userID := userIDCtx.(uint)
	
	search := c.Query("search")

	workspaces, err := h.service.GetMyWorkspaces(userID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch workspaces"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": workspaces})
}

func (h *WorkspaceHandler) Detail(c *gin.Context) {
    
    userIDCtx, _ := c.Get("userID")
    userID := userIDCtx.(uint)

    idParam := c.Param("id")
    workspaceID, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
        return
    }

    detail, err := h.service.GetWorkspaceDetail(userID, uint(workspaceID))
    if err != nil {
        
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": detail})
}

func (h *WorkspaceHandler) Update(c *gin.Context) {
    userIDCtx, _ := c.Get("userID")
    userID := userIDCtx.(uint)

    idParam := c.Param("id")
    workspaceID, _ := strconv.Atoi(idParam)

    var req dto.UpdateWorkspaceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ws, err := h.service.UpdateWorkspace(userID, uint(workspaceID), req)
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Workspace updated", "data": ws})
}

func (h *WorkspaceHandler) Delete(c *gin.Context) {
    userIDCtx, _ := c.Get("userID")
    userID := userIDCtx.(uint)

    idParam := c.Param("id")
    workspaceID, _ := strconv.Atoi(idParam)

    err := h.service.DeleteWorkspace(userID, uint(workspaceID))
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Workspace deleted successfully"})
}

func (h *WorkspaceHandler) AddMember(c *gin.Context) {
    userIDCtx, _ := c.Get("userID")
    workspaceID, _ := strconv.Atoi(c.Param("id"))

    var req dto.AddMemberRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.AddMember(userIDCtx.(uint), uint(workspaceID), req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) 
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Member added successfully"})
}

func (h *WorkspaceHandler) UpdateMemberRole(c *gin.Context) {
    userIDCtx, _ := c.Get("userID")
    workspaceID, _ := strconv.Atoi(c.Param("id"))
    targetUserID, _ := strconv.Atoi(c.Param("userId"))

    var req dto.UpdateMemberRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.UpdateMemberRole(userIDCtx.(uint), uint(workspaceID), uint(targetUserID), req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

func (h *WorkspaceHandler) RemoveMember(c *gin.Context) {
    userIDCtx, _ := c.Get("userID")
    workspaceID, _ := strconv.Atoi(c.Param("id"))
    targetUserID, _ := strconv.Atoi(c.Param("userId"))

    if err := h.service.RemoveMember(userIDCtx.(uint), uint(workspaceID), uint(targetUserID)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}