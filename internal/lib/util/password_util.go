package util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordUtil struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// GeneratePassword is used to generate a new password hash for storing and
// comparing at a later date.
func GenPasswordHash(password string) (string, error) {
	this := &PasswordUtil{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}
	// Generate a Salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, this.time, this.memory, this.threads, this.keyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, this.memory, this.time, this.threads, b64Salt, b64Hash)
	return full, nil
}

// ComparePassword is used to compare a user-inputted password to a hash to see
// if the password matches or not.
func ComparePasswords(password, hashed string) (bool, error) {
	this := &PasswordUtil{}
	parts := strings.Split(hashed, "$")

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &this.memory, &this.time, &this.threads)
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	this.keyLen = uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(password), salt, this.time, this.memory, this.threads, this.keyLen)
	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1), nil
}
