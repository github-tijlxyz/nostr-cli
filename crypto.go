package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

const (
    saltSize = 16
    keySize  = 32
)

func decryptOld(ciphertext, password string) (string, error) {
    ciphertextBytes, err := hex.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    if len(ciphertextBytes) < saltSize+aes.BlockSize {
        return "", errors.New("Ciphertext is too short")
    }

    salt := ciphertextBytes[:saltSize]
    ciphertextBytes = ciphertextBytes[saltSize:]

    key := pbkdf2.Key([]byte(password), salt, 4096, keySize, sha1.New)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    iv := ciphertextBytes[:aes.BlockSize]
    ciphertextBytes = ciphertextBytes[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

    return string(ciphertextBytes), nil
}

