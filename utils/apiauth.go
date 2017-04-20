package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	//	"laatoosdk/log"
)

func EncryptWithKey(publickey *rsa.PublicKey, message []byte) ([]byte, error) {
	encryptedmsg, err := rsa.EncryptOAEP(md5.New(), rand.Reader, publickey, message, []byte(""))
	if err != nil {
		return nil, err
	}
	return encryptedmsg, nil
}

func DecryptWithKey(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	out, err := rsa.DecryptOAEP(md5.New(), rand.Reader, privateKey, message, []byte(""))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func LoadPublicKey(path string) (key *rsa.PublicKey, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("ssh: no key found")
	}
	switch block.Type {
	case "PUBLIC KEY":
		keyInt, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return keyInt.(*rsa.PublicKey), nil
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
}

// loadPrivateKey loads an parses a PEM encoded private key file.
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("ssh: no key found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
}
