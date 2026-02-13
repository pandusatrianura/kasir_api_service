package entity

type Role struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Users     []User `json:"users"`
}

type User struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Photo     string `json:"photo"`
	Phone     string `json:"phone"`
	Roles     []Role
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserRole struct {
	ID        uint `json:"id"`
	UserID    uint `json:"user_id"`
	RoleID    uint `json:"role_id"`
	User      User
	Role      Role
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID uint     `json:"user_id"`
	Email  string   `json:"email"`
	Role   []string `json:"role"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Roles string `json:"roles"`
	} `json:"user"`
}
