package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// public key encrypt
//
func RsaEncrypt(data, keyBytes []byte) ([]byte, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// public key encrypt
// and base64
//
func RsaEncryptAndBase64(data, keyBytes []byte) (string, error) {
	text, err := RsaEncrypt(data, keyBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(text), nil
}
