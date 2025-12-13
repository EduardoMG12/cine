package http

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/usecase/user"
)

type UserHandler struct {
	updateUserUC *user.UpdateUserUseCase
}

func NewUserHandler(updateUserUC *user.UpdateUserUseCase) *UserHandler {
	return &UserHandler{
		updateUserUC: updateUserUC,
	}
}

// UpdateUser godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserRequest true "User update data"
// @Success 200 {object} dto.APIResponse{data=dto.UserProfileResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/users/me [patch]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := h.updateUserUC.Execute(userID, &req)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "User profile updated successfully", result)
}
