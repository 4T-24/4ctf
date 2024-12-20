package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

const (
	MinSaltSize = 10
)

var (
	DefaultPassword = NewPassword(sha512.New, 16, 32, 210000)
)

func HashPassword(password string) string {
	return DefaultPassword.HashPassword(password).String()
}

func VerifyPassword(password, stored string) bool {
	parts := ParseHashResult(stored)
	return DefaultPassword.VerifyPassword(password, parts.CipherText, parts.Salt)
}

type Password struct {
	Diggest    func() hash.Hash
	SaltSize   int
	KeyLen     int
	Iterations int
}

type HashResult struct {
	CipherText string
	Salt       string
}

func ParseHashResult(s string) HashResult {
	parts := bytes.Split([]byte(s), []byte(":"))
	return HashResult{CipherText: string(parts[0]), Salt: string(parts[1])}
}

func (h HashResult) String() string {
	return h.CipherText + ":" + h.Salt
}
func NewPassword(diggest func() hash.Hash, saltSize int, keyLen int, iter int) *Password {
	if saltSize < MinSaltSize {
		saltSize = MinSaltSize
	}

	return &Password{
		Diggest:    diggest,
		SaltSize:   saltSize,
		KeyLen:     keyLen,
		Iterations: iter,
	}
}

func (p *Password) genSalt() string {
	saltBytes := make([]byte, p.SaltSize)
	_, err := rand.Read(saltBytes)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(saltBytes)
}

func (p *Password) HashPassword(password string) HashResult {
	saltString := p.genSalt()
	salt := bytes.NewBufferString(saltString).Bytes()
	df := pbkdf2.Key([]byte(password), salt, p.Iterations, p.KeyLen, p.Diggest)
	cipherText := base64.StdEncoding.EncodeToString(df)
	return HashResult{CipherText: cipherText, Salt: saltString}
}

func (p *Password) VerifyPassword(password, cipherText, salt string) bool {
	saltBytes := bytes.NewBufferString(salt).Bytes()
	df := pbkdf2.Key([]byte(password), saltBytes, p.Iterations, p.KeyLen, p.Diggest)

	return equal(cipherText, df)
}

// check per bit by applying bitwise XOR
// first, decode the base64 string to bytes
// for example
// 114  1110010
// 114  1110010
// ----------------- xor
//
//	0000000
func equal(cipherText string, newCipherText []byte) bool {
	x, _ := base64.StdEncoding.DecodeString(cipherText)
	diff := uint64(len(x)) ^ uint64(len(newCipherText))

	for i := 0; i < len(x) && i < len(newCipherText); i++ {
		diff |= uint64(x[i]) ^ uint64(newCipherText[i])
	}

	return diff == 0
}
