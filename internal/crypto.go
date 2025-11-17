package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

const (
	encryptedPrefix = "enc:"
	// Hardcoded salt for key derivation
	encryptionSalt = "client-manager-secret-salt-2024-v1"
)

// isEncrypted checks if a string is already encrypted (has enc: prefix)
func isEncrypted(s string) bool {
	return len(s) >= len(encryptedPrefix) && s[:len(encryptedPrefix)] == encryptedPrefix
}

// deriveKey derives a fixed 32-byte key from the hardcoded salt
func deriveKey() []byte {
	hash := sha256.Sum256([]byte(encryptionSalt))
	return hash[:32]
}

// encryptString encrypts plaintext and returns a string with prefix enc:
func encryptString(plain string) (string, error) {
	if plain == "" {
		return "", nil
	}
	if len(plain) >= len(encryptedPrefix) && plain[:len(encryptedPrefix)] == encryptedPrefix {
		// already encrypted
		return plain, nil
	}

	key := deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ct := gcm.Seal(nil, nonce, []byte(plain), nil)
	out := append(nonce, ct...)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(out), nil
}

// decryptString reverses encryptString if string is prefixed with enc:
func decryptString(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	if len(s) < len(encryptedPrefix) || s[:len(encryptedPrefix)] != encryptedPrefix {
		// not encrypted
		return s, nil
	}
	b64 := s[len(encryptedPrefix):]
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}

	key := deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", errors.New("ciphertext too short")
	}
	nonce := data[:ns]
	ct := data[ns:]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

// encryptClientsInPlace encrypts password fields in clients slice in-place before saving.
func encryptClientsInPlace(clients []Client) error {
	for i := range clients {
		// VPN
		if v := clients[i].VPN.Password; v != "" {
			enc, err := encryptString(v)
			if err != nil {
				return err
			}
			clients[i].VPN.Password = enc
		}
		// Client Data
		if v := clients[i].Data.JiraPassword; v != "" {
			enc, err := encryptString(v)
			if err != nil {
				return err
			}
			clients[i].Data.JiraPassword = enc
		}
		// Apps
		for j := range clients[i].Apps {
			if v := clients[i].Apps[j].Password; v != "" {
				enc, err := encryptString(v)
				if err != nil {
					return err
				}
				clients[i].Apps[j].Password = enc
			}
			if v := clients[i].Apps[j].AppServerPass; v != "" {
				enc, err := encryptString(v)
				if err != nil {
					return err
				}
				clients[i].Apps[j].AppServerPass = enc
			}
			if v := clients[i].Apps[j].AppServerUser; v != "" {
				// do not encrypt usernames
				_ = v
			}
		}
	}
	return nil
}

// decryptClientsInPlace decrypts password fields in clients slice in-place after loading.
func decryptClientsInPlace(clients []Client) error {
	for i := range clients {
		if v := clients[i].VPN.Password; v != "" {
			dec, err := decryptString(v)
			if err != nil {
				return err
			}
			clients[i].VPN.Password = dec
		}
		if v := clients[i].Data.JiraPassword; v != "" {
			dec, err := decryptString(v)
			if err != nil {
				return err
			}
			clients[i].Data.JiraPassword = dec
		}
		for j := range clients[i].Apps {
			if v := clients[i].Apps[j].Password; v != "" {
				dec, err := decryptString(v)
				if err != nil {
					return err
				}
				clients[i].Apps[j].Password = dec
			}
			if v := clients[i].Apps[j].AppServerPass; v != "" {
				dec, err := decryptString(v)
				if err != nil {
					return err
				}
				clients[i].Apps[j].AppServerPass = dec
			}
		}
	}
	return nil
}
