package main

import (
	"crypto/sha256"
	"fmt"
)

func encrypt(passpharse string) string {
	passwordEncryptor := sha256.New()
	passwordEncryptor.Write([]byte(passpharse))
	return fmt.Sprintf("%x", passwordEncryptor.Sum(nil))
}
