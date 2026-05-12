package security

import "golang.org/x/crypto/bcrypt"

const dummyBcryptHash = "$2a$12$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW"

func DummyPasswordHash() string {
	return dummyBcryptHash
}

func VerifyPassword(storedHash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(plain))
	if err == nil {
		return true
	}
	_ = bcrypt.CompareHashAndPassword([]byte(dummyBcryptHash), []byte(plain))
	return false
}
