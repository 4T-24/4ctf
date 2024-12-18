package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

func TestPbkdf2ReturnFalse(t *testing.T) {
	pass := NewPassword(sha256.New, 64, 64, 15000)
	hashed := pass.HashPassword("12345")
	cipherText := hashed.CipherText
	salt := hashed.Salt

	isValid := pass.VerifyPassword("1234", cipherText, salt)

	if isValid {
		t.Error("Verify Password was expected to return false : but result is ", isValid)
	}
}

func TestPbkdf2ReturnTrue(t *testing.T) {
	pass := NewPassword(sha256.New, 7, 64, 15000)
	hashed := pass.HashPassword("12345")
	cipherText := hashed.CipherText
	salt := hashed.Salt

	isValid := pass.VerifyPassword("12345", cipherText, salt)
	if !isValid {
		t.Error("Verify Password was expected to return true : but result is ", isValid)
	}
}

func TestPbkdf2DefaultReturnFalse(t *testing.T) {
	var password []byte
	rand.Read(password)
	randomPassword := string(password)

	hashed := HashPassword(randomPassword)
	isValid := VerifyPassword("1234", hashed)

	if isValid {
		t.Error("Verify Password was expected to return false : but result is ", isValid)
	}
}

func TestPbkdf2DefaultReturnTrue(t *testing.T) {
	var password []byte
	rand.Read(password)
	randomPassword := string(password)

	hashed := HashPassword(randomPassword)
	isValid := VerifyPassword(randomPassword, hashed)

	if !isValid {
		t.Error("Verify Password was expected to return true : but result is ", isValid)
	}
}

func BenchmarkPBKDF2HashOneThousandIterations(b *testing.B) {
	pass := NewPassword(sha256.New, 64, 64, 1000)
	for i := 0; i < b.N; i++ {
		pass.HashPassword("12345")
	}
}

func BenchmarkPBKDF2HashFifteenThousandIterations(b *testing.B) {
	pass := NewPassword(sha256.New, 64, 64, 15000)
	for i := 0; i < b.N; i++ {
		pass.HashPassword("12345")
	}
}

func BenchmarkPBKDF2Verify(b *testing.B) {
	pass := NewPassword(sha256.New, 64, 64, 1000)
	hashed := pass.HashPassword("12345")
	cipherText := hashed.CipherText
	salt := hashed.Salt
	for i := 0; i < b.N; i++ {
		pass.VerifyPassword("12345", cipherText, salt)
	}
}
