package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha1"
    "encoding/hex"
    "golang.org/x/crypto/pbkdf2"
    "io"
    "errors"
)

const (
    saltSize = 16
    keySize  = 32
)

func encrypt(plaintext, password string) (string, error) {
    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return "", err
    }

    key := pbkdf2.Key([]byte(password), salt, 4096, keySize, sha1.New)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

    result := append(salt, ciphertext...)

    return hex.EncodeToString(result), nil
}

func decrypt(ciphertext, password string) (string, error) {
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

