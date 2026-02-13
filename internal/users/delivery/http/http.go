package http

import (
	"net/http"

	"github.com/pandusatrianura/kasir_api_service/api/middleware"
	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/internal/users/entity"
	"github.com/pandusatrianura/kasir_api_service/internal/users/service"
	"github.com/pandusatrianura/kasir_api_service/pkg/convert"
	"github.com/pandusatrianura/kasir_api_service/pkg/response"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UserHandler struct {
	service service.IUserService
}

func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{service: service}
}

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Tags users
// @Accept json
// @Produce json
// @Param login body entity.LoginRequest true "Login Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, constants.ErrorCode, constants.ErrInvalidMethod, nil)
		return
	}

	var (
		requestLogin entity.LoginRequest
		user         *entity.User
		err          error
	)
	if err := response.ParseJSON(r, &requestLogin); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidCategoryRequest, err)
		return
	}

	if user, err = h.service.GetUserByEmail(requestLogin.Email); err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Get User By Email failed", err)
		return
	}

	if user == nil {
		log.Errorf("[AuthController.Login] Login - 4: User not found")
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "User not found", nil)
		return
	}

	isSame := convert.CheckPasswordHash(requestLogin.Password, user.Password)
	if !isSame {
		log.Errorf("[AuthController.Login] Login - 5: Invalid email or password")
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Invalid email or password", nil)
		return
	}

	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	loginResp := entity.LoginResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   roles,
	}

	rolesStr := ""
	if len(loginResp.Role) > 0 {
		rolesStr = loginResp.Role[0]
		for i := 1; i < len(loginResp.Role); i++ {
			rolesStr += "," + loginResp.Role[i]
		}
	}

	jwtConfig := middleware.NewJWTConfig(viper.GetString("JWT_SECRET_KEY"), viper.GetString("JWT_ISSUER"), viper.GetString("JWT_DURATION"))
	token, err := middleware.GenerateJWT(user.ID, user.Email, rolesStr, jwtConfig)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Generate JWT failed", err)
		return
	}

	resp := entity.AuthResponse{
		Token: token,
		User: struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
			Roles string `json:"roles"`
		}{
			ID:    loginResp.UserID,
			Email: loginResp.Email,
			Roles: rolesStr,
		},
	}

	response.Success(w, http.StatusOK, constants.SuccessCode, "Login successfully", resp)
}

// Logout godoc
// @Summary Logout a user
// @Description Logout a user by clearing the authentication cookie or token
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/logout [post]
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, constants.ErrorCode, constants.ErrInvalidMethod, nil)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	response.Success(w, http.StatusOK, constants.SuccessCode, "Logout successfully", nil)
}
