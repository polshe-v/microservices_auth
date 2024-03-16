package utils

import "golang.org/x/crypto/bcrypt"

// VerifyPassword compares stored and provided password hashes.
func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}
