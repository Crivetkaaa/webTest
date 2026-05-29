package utilit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

// ИСПРАВЛЕНО: Функция динамического получения ключа.
// Она сработает строго в момент вызова шифрования, когда .env уже загружен.
func getEncryptionKey() ([]byte, error) {
	key := []byte(os.Getenv("SECRET_KEY"))
	if len(key) != 32 {
		return nil, fmt.Errorf("SECRET_KEY должен быть строго 32 байта, текущая длина: %d", len(key))
	}
	return key, nil
}

func Encrypt(plainText string) (string, error) {
	if len(plainText) == 0 {
		return "", nil
	}

	// Запрашиваем ключ динамически
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(cipherText), nil
}

func Decrypt(hexText string) (string, error) {
	if len(hexText) == 0 {
		return "", nil
	}

	cipherText, err := hex.DecodeString(hexText)
	if err != nil {
		return "", err
	}

	// Запрашиваем ключ динамически
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("длина зашифрованного текста слишком мала")
	}

	nonce, actualCipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, actualCipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
