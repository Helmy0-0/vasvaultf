package handlers

import (
	"net/http"
	"strconv"
	"vasvault/internal/dto"
	"vasvault/internal/services"
	"vasvault/pkg/utils"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	FileService services.FileServiceInterface
}

func NewFileHandler(fileService services.FileServiceInterface) *FileHandler {
	return &FileHandler{
		FileService: fileService,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		utils.RespondJSON(c, http.StatusUnauthorized, nil, "user not found in context")
		return
	}

	userID, ok := uid.(uint)
	if !ok {
		if fid, ok := uid.(float64); ok {
			userID = uint(fid)
		} else {
			utils.RespondJSON(c, http.StatusInternalServerError, nil, "invalid user id")
			return
		}
	}

	var request dto.UploadFileRequest
	_ = c.ShouldBind(&request)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "file is required")
		return
	}
	defer file.Close()

	workspaceIDStr := c.PostForm("workspace_id")
    if workspaceIDStr != "" {
  
        parsedID, err := strconv.ParseUint(workspaceIDStr, 10, 64)
        if err == nil {
            
            workspaceID := uint(parsedID)
            request.WorkspaceID = &workspaceID 
        }
    }

	response, err := h.FileService.UploadFile(userID, file, header, request)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, err.Error())
	}

	utils.RespondJSON(c, http.StatusOK, response, "file uploaded successfully")
}

func (h *FileHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	fileID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "invalid file id")
		return
	}

	response, err := h.FileService.GetFileByID(uint(fileID))
	if err != nil {
		utils.RespondJSON(c, http.StatusNotFound, nil, "file not found")
		return
	}

	utils.RespondJSON(c, http.StatusOK, response, "ok")
}

func (h *FileHandler) ListMyFiles(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		utils.RespondJSON(c, http.StatusUnauthorized, nil, "user not found in context")
		return
	}

	userID, ok := uid.(uint)
	if !ok {
		if fid, ok := uid.(float64); ok {
			userID = uint(fid)
		} else {
			utils.RespondJSON(c, http.StatusInternalServerError, nil, "invalid user id")
			return
		}
	}

	categoryIDParam := c.Query("categoryId")
	var categoryID *uint

	if categoryIDParam != "" {
		parseID, err := strconv.ParseUint(categoryIDParam, 10, 64)
		if err != nil {
			utils.RespondJSON(c, http.StatusBadRequest, nil, "invalid category id")
			return
		}
		id := uint(parseID)
		categoryID = &id
	}

	response, err := h.FileService.ListUserFilesWithOptionalCategory(userID, categoryID)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	utils.RespondJSON(c, http.StatusOK, response, "ok")
}

func (h *FileHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	fileID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "invalid file id")
		return
	}

	if err := h.FileService.DeleteFile(uint(fileID)); err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	utils.RespondJSON(c, http.StatusOK, nil, "file deleted successfully")
}

// AssignCategories - POST /files/:id/categories/assign
func (h *FileHandler) AssignCategories(c *gin.Context) {
	userID := c.GetUint("userID")
	idParam := c.Param("id")
	fileID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "invalid file id")
		return
	}

	var req dto.AssignCategoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := h.FileService.AssignCategories(userID, uint(fileID), req.CategoryIDs); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, nil, "categories assigned successfully")
}

// RemoveCategories - POST /files/:id/categories/remove
func (h *FileHandler) RemoveCategories(c *gin.Context) {
	userID := c.GetUint("userID")
	idParam := c.Param("id")
	fileID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "invalid file id")
		return
	}

	var req dto.AssignCategoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := h.FileService.RemoveCategories(userID, uint(fileID), req.CategoryIDs); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, nil, "categories removed successfully")
}

// UpdateCategories - PUT /files/:id/categories
func (h *FileHandler) UpdateCategories(c *gin.Context) {
	userID := c.GetUint("userID")
	idParam := c.Param("id")
	fileID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "invalid file id")
		return
	}

	var req dto.AssignCategoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := h.FileService.UpdateCategories(userID, uint(fileID), req.CategoryIDs); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, nil, "categories updated successfully")
}

// StorageSummary - GET /storage/summary
func (h *FileHandler) StorageSummary(c *gin.Context) {
	userID := c.GetUint("userID")

	resp, err := h.FileService.GetStorageSummary(userID)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, resp, "ok")
}
