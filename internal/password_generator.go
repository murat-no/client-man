package main

import (
	"crypto/rand"
	"math/big"
)

// PasswordGeneratorConfig şifre oluşturucu ayarları
type PasswordGeneratorConfig struct {
	Length         int
	UseUppercase   bool
	UseLowercase   bool
	UseNumbers     bool
	UseSpecialChar bool
}

// Karakter setleri
const (
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	numberChars    = "0123456789"
	specialChars   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// GeneratePassword verilen konfigürasyona göre güçlü şifre oluşturur
func GeneratePassword(config PasswordGeneratorConfig) (string, error) {
	// Karakter havuzunu oluştur
	var charPool string

	if config.UseUppercase {
		charPool += uppercaseChars
	}
	if config.UseLowercase {
		charPool += lowercaseChars
	}
	if config.UseNumbers {
		charPool += numberChars
	}
	if config.UseSpecialChar {
		charPool += specialChars
	}

	// Eğer hiç karakter seçilmemişse varsayılan olarak hepsini kullan
	if charPool == "" {
		charPool = uppercaseChars + lowercaseChars + numberChars + specialChars
	}

	// Şifreyi oluştur
	password := make([]byte, config.Length)
	charPoolLen := big.NewInt(int64(len(charPool)))

	for i := 0; i < config.Length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charPoolLen)
		if err != nil {
			return "", err
		}
		password[i] = charPool[randomIndex.Int64()]
	}

	return string(password), nil
}

// DefaultPasswordConfig varsayılan şifre yapılandırması
func DefaultPasswordConfig() PasswordGeneratorConfig {
	return PasswordGeneratorConfig{
		Length:         16,
		UseUppercase:   true,
		UseLowercase:   true,
		UseNumbers:     true,
		UseSpecialChar: true,
	}
}
