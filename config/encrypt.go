package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

func Encrypt(sender, pass string) (string, error) {
	encMsg, err := aesencrypt(sender, pass)
	if err != nil {
		return "", err
	}
	return encMsg, nil
}
func Decrypt(sender, pass string) (string, error) {
	password, err := aesdecrypt(sender, pass)
	if err != nil {
		return "", err
	}
	return password, nil
}

func aesencrypt(key string, msg string) (ret string, err error) {
	c, err := aes.NewCipher(createHash("gibberish" + key + "gibb"))
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	text := gcm.Seal(nonce, nonce, []byte(msg), nil)
	ret = hex.EncodeToString(text)
	return
}
func aesdecrypt(key, msg string) (ret string, err error) {
	data, err := hex.DecodeString(msg)
	if err != nil {
		err = errors.New("Cannot decode the password")
		return
	}
	c, err := aes.NewCipher(createHash("gibberish" + key + "gibb"))
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return
	}

	nonceSize := gcm.NonceSize()
	nonce, text := data[:nonceSize], data[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, text, nil)
	if err != nil {
		return
	}
	ret = string(plainText)
	return
}
func createHash(key string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)
}
