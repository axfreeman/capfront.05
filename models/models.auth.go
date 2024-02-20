// usermodels.go
// models that are employed in handling users, login, registration, authentication, etc etc

package models

type PostData struct {
	User string `json:"user"`
}

type UserLoginDetails struct {
	UserName string `json:"username" binding:"required" gorm:"primary_key"`
	Password string `json:"password" binding:"required"`
}
