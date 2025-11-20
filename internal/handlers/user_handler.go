package handlers

import (
	"errors"
	"net/http"
	"vasvault/internal/dto"
	"vasvault/internal/services"
	"vasvault/pkg/utils"
	apperrors "vasvault/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Me(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		utils.RespondJSON(c, http.StatusUnauthorized, nil, "user not found in context")
		return
	}

	id, ok := uid.(uint)
	if !ok {
		if fid, ok := uid.(float64); ok {
			id = uint(fid)
		} else {
			utils.RespondJSON(c, http.StatusInternalServerError, nil, "invalid user id in context")
			return
		}
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		utils.RespondJSON(c, http.StatusNotFound, nil, "user not found")
		return
	}

	response := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	utils.RespondJSON(c, http.StatusOK, response, "ok")
}

func (h *UserHandler) Register(c *gin.Context) {
	var userRequest dto.RegisterRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "Validation error")
		return
	}

	resp, err := h.userService.Register(userRequest)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrEmailExists):
			utils.RespondJSON(c, http.StatusConflict, nil, "email already exists")
			return
		case errors.Is(err, apperrors.ErrUsernameExists):
			utils.RespondJSON(c, http.StatusConflict, nil, "username already exists")
			return
		default:
			utils.RespondJSON(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
	}

	token, err := utils.GenerateTokenPair(resp.Username, resp.ID)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to generate tokens")
		return
	}

	response := dto.AuthResponse{
		User: *resp,
		Token: dto.TokenResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}

	utils.RespondJSON(c, http.StatusOK, response, "User created successfully")
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "Validation error")
		return
	}

	userResp, err := h.userService.Login(loginRequest)
	if err != nil {
		utils.RespondJSON(c, http.StatusUnauthorized, nil, "Invalid email or password")
		return
	}
	token, err := utils.GenerateTokenPair(userResp.Username, userResp.ID)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to generate tokens")
		return
	}

	response := dto.AuthResponse{
		User: dto.UserResponse{
			ID:       userResp.ID,
			Email:    userResp.Email,
			Username: userResp.Username,
		},
		Token: dto.TokenResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}

	utils.RespondJSON(c, http.StatusOK, response, "Login successful")
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		utils.RespondJSON(c, http.StatusUnauthorized, nil, "user not found in context")
		return
	}

	id, ok := uid.(uint)
	if !ok {
		if fid, ok := uid.(float64); ok {
			id = uint(fid)
		} else {
			utils.RespondJSON(c, http.StatusInternalServerError, nil, "invalid user id in context")
			return
		}
	}

	var updateRequest dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "Validation error")
		return
	}

	resp, err := h.userService.UpdateUser(id, updateRequest)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, resp, "Profile updated successfully")
}

func (h *UserHandler) Refresh(c *gin.Context) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, nil, "Validation error")
		return
	}

	userResp, err := h.userService.Refresh(refreshRequest.RefreshToken)
	if err != nil {
		utils.RespondJSON(c, http.StatusUnauthorized, nil, "Invalid refresh token")
		return
	}

	createdUser, err := h.userService.GetUser(userResp.ID)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to fetch user")
		return
	}

	token, err := utils.GenerateTokenPair(createdUser.Username, createdUser.ID)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to generate tokens")
		return
	}
	response := dto.AuthResponse{
		User: dto.UserResponse{
			ID:       createdUser.ID,
			Email:    createdUser.Email,
			Username: createdUser.Username,
		},
		Token: dto.TokenResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}

	utils.RespondJSON(c, http.StatusOK, response, "Token refreshed successfully")
}
