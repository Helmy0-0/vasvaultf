package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"vasvault/internal/dto"
	"vasvault/internal/services"
	"vasvault/pkg/utils"

	"github.com/disintegration/imaging"
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

// Download - GET /files/:id/download
func (h *FileHandler) Download(c *gin.Context) {
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

	// Serve file directly (Gin will set content-type)
	c.File(response.FilePath)
}

// Thumbnail - GET /files/:id/thumbnail
// Generates a cached thumbnail (200x200) and serves it.
func (h *FileHandler) Thumbnail(c *gin.Context) {
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

	// only handle image types
	if !strings.HasPrefix(response.MimeType, "image/") {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "thumbnail only supported for images")
		return
	}

	thumbDir := filepath.Join("./uploads", "thumbs")
	if err := os.MkdirAll(thumbDir, os.ModePerm); err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, "failed to create thumb dir")
		return
	}

	base := filepath.Base(response.FilePath)
	thumbPath := filepath.Join(thumbDir, base+".thumb.jpg")

	// serve cached thumbnail if exists
	if _, err := os.Stat(thumbPath); err == nil {
		c.File(thumbPath)
		return
	}

	// generate thumbnail
	img, err := imaging.Open(response.FilePath)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, "failed to open image")
		return
	}
	thumb := imaging.Thumbnail(img, 200, 200, imaging.Lanczos)
	if err := imaging.Save(thumb, thumbPath, imaging.JPEGQuality(80)); err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, "failed to save thumbnail")
		return
	}

	c.File(thumbPath)
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

func (h *FileHandler) ListByWorkspace(c *gin.Context) {
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

	idParam := c.Param("id")
	workspaceID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "invalid workspace id")
		return
	}

	resp, err := h.FileService.ListFilesByWorkspace(userID, uint(workspaceID))
	if err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, resp, "ok")
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

func (h *FileHandler) StorageSummary(c *gin.Context) {
	userID := c.GetUint("userID")

	resp, err := h.FileService.GetStorageSummary(userID)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, resp, "ok")
}
