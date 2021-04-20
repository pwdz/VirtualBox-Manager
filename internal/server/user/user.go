package user

import "github.com/dgrijalva/jwt-go"
// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type User struct {
	Username string `json:"username"` 
	Password string 
	Role     string 
	jwt.StandardClaims
}

func NewUser(username, password, role string) User {
	return User{
		Username: username,
		Password: password,
		Role:     role,
	}
}

func (u User) CheckPassword(pass string) bool {
	return u.Password == pass
}