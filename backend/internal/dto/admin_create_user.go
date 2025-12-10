package dto

// AdminCreateUserRequest is used by super_admin to create an admin user.
type AdminCreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}
