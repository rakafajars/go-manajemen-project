package seed

import "golang.org/x/crypto/bcrypt"

func SeedAdmin() {
	password := "admin"
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
