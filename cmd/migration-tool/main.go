package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const (
	encryptedPrefix = "enc:"
	keyringService  = "client-manager"
	keyringUser     = "encryption-key"
)

type VPNInfo struct {
	App           string `json:"app"`
	Host          string `json:"host"`
	User          string `json:"user"`
	Password      string `json:"password"`
	TwoFATokenApp string `json:"two_fa_token_app"`
	Not           string `json:"not"`
}

type DataInfo struct {
	JiraURI      string   `json:"jira_uri"`
	JiraUser     string   `json:"jira_user"`
	JiraPassword string   `json:"jira_password"`
	User         string   `json:"user"`
	PassReset    string   `json:"pass_reset"`
	RDC          []string `json:"rdc"`
	Hosts        []string `json:"hosts"`
	Not          string   `json:"not"`
}

type AppInfo struct {
	Type          string   `json:"type"`
	Name          string   `json:"name"`
	User          string   `json:"user"`
	Password      string   `json:"pass"`
	DBServerIP    string   `json:"db_server_ip"`
	TNS           string   `json:"tns"`
	AppServerIP   string   `json:"app_server_ip"`
	AppServerURI  string   `json:"app_server_uri"`
	AppServerUser string   `json:"app_server_user"`
	AppServerPass string   `json:"app_server_pass"`
	AppURI        string   `json:"app_uri"`
	AppUsers      []string `json:"app_users"`
	Not           string   `json:"not"`
}

type Client struct {
	Company    string    `json:"company"`
	EBSVersion string    `json:"ebs_version"`
	VPN        VPNInfo   `json:"vpn"`
	Data       DataInfo  `json:"data"`
	Apps       []AppInfo `json:"apps"`
	Notes      string    `json:"not"`
}

func isEncrypted(s string) bool {
	return len(s) >= len(encryptedPrefix) && s[:len(encryptedPrefix)] == encryptedPrefix
}

// getOldKeyFromKeyring retrieves the old encryption key from Windows keyring
func getOldKeyFromKeyring() ([]byte, error) {
	v, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		return nil, fmt.Errorf("keyring get failed: %w", err)
	}
	if v == "" {
		return nil, fmt.Errorf("key not found in keyring")
	}

	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}

	if len(data) != 32 {
		return nil, fmt.Errorf("invalid key length: %d (expected 32)", len(data))
	}

	return data, nil
}

// decryptStringWithOldKey decrypts a string using the old keyring-based key
func decryptStringWithOldKey(s string, key []byte) (string, error) {
	if s == "" {
		return "", nil
	}
	if len(s) < len(encryptedPrefix) || s[:len(encryptedPrefix)] != encryptedPrefix {
		return s, nil
	}

	b64 := s[len(encryptedPrefix):]
	data, err := base64.StdEncoding.DecodeString(b64)
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

	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce := data[:ns]
	ct := data[ns:]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}

	return string(pt), nil
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║  Client Manager - Encryption Migration Tool               ║")
	fmt.Println("║  Eski keyring şifrelerini yeni encryption'a migrate eder  ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// client_info.json dosyasını bul
	jsonFile := "client_info.json"
	if _, err := os.Stat(jsonFile); err != nil {
		// Eğer current directory'de yoksa, parent directory'de ara
		if os.IsNotExist(err) {
			jsonFile = filepath.Join("..", jsonFile)
			if _, err := os.Stat(jsonFile); err != nil {
				fmt.Printf("❌ Hata: client_info.json bulunamadı\n")
				fmt.Println("Tool'u client_info.json ile aynı dizinde çalıştırın.")
				return
			}
		}
	}

	fmt.Printf("📁 Dosya: %s\n", jsonFile)
	fmt.Println()

	// JSON dosyasını oku
	fmt.Println("📂 Dosya okunuyor...")
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Printf("❌ Hata: Dosya okunamadı: %v\n", err)
		return
	}

	var clients []Client
	if err := json.Unmarshal(data, &clients); err != nil {
		fmt.Printf("❌ Hata: JSON parse edilemedi: %v\n", err)
		return
	}

	fmt.Printf("✓ %d firma yüklendi\n", len(clients))
	fmt.Println()

	// Eski keyring'ten key'i al
	fmt.Println("🔑 Eski keyring anahtarı okunuyor...")
	oldKey, err := getOldKeyFromKeyring()
	if err != nil {
		fmt.Printf("❌ Hata: %v\n", err)
		fmt.Println()
		fmt.Println("⚠️  Keyring'te anahtar bulunamadı.")
		fmt.Println("Muhtemelen eski versiyon yüklenmedikçe bu araç çalıştırılamaz.")
		return
	}
	fmt.Println("✓ Eski anahtar bulundu")
	fmt.Println()

	// Şifreleri decrypt et (eski key ile)
	fmt.Println("🔓 Şifreler decrypt ediliyor (eski key ile)...")
	decryptedCount := 0

	for i := range clients {
		// VPN Password
		if clients[i].VPN.Password != "" && isEncrypted(clients[i].VPN.Password) {
			if dec, err := decryptStringWithOldKey(clients[i].VPN.Password, oldKey); err == nil {
				clients[i].VPN.Password = dec
				decryptedCount++
				fmt.Printf("  ✓ %s VPN şifresi\n", clients[i].Company)
			}
		}

		// Jira Password
		if clients[i].Data.JiraPassword != "" && isEncrypted(clients[i].Data.JiraPassword) {
			if dec, err := decryptStringWithOldKey(clients[i].Data.JiraPassword, oldKey); err == nil {
				clients[i].Data.JiraPassword = dec
				decryptedCount++
			}
		}

		// App Passwords
		for j := range clients[i].Apps {
			if clients[i].Apps[j].Password != "" && isEncrypted(clients[i].Apps[j].Password) {
				if dec, err := decryptStringWithOldKey(clients[i].Apps[j].Password, oldKey); err == nil {
					clients[i].Apps[j].Password = dec
					decryptedCount++
				}
			}
			if clients[i].Apps[j].AppServerPass != "" && isEncrypted(clients[i].Apps[j].AppServerPass) {
				if dec, err := decryptStringWithOldKey(clients[i].Apps[j].AppServerPass, oldKey); err == nil {
					clients[i].Apps[j].AppServerPass = dec
					decryptedCount++
				}
			}
		}
	}

	fmt.Printf("✓ %d şifre decrypt edildi\n", decryptedCount)
	fmt.Println()

	// Backup oluştur
	backupFile := jsonFile + ".pre_migration_backup"
	fmt.Println("💾 Backup oluşturuluyor...")
	if err := os.WriteFile(backupFile, data, 0600); err != nil {
		fmt.Printf("⚠️  Backup oluşturulamadı: %v\n", err)
	} else {
		fmt.Printf("✓ Backup: %s\n", backupFile)
	}
	fmt.Println()

	// Yeni JSON'ı yaz (plaintext olarak - migration tool başarılı olduğu için)
	fmt.Println("💾 Güncellenmiş veriler kaydediliyor...")
	newData, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		fmt.Printf("❌ Hata: JSON marshal edilemedi: %v\n", err)
		return
	}

	if err := os.WriteFile(jsonFile, newData, 0600); err != nil {
		fmt.Printf("❌ Hata: Dosya yazılamadı: %v\n", err)
		return
	}

	fmt.Printf("✓ %s güncellendi\n", jsonFile)
	fmt.Println()

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║ ✅ Migration başarılı!                                    ║")
	fmt.Println("║                                                            ║")
	fmt.Println("║ Şimdi client-manager.exe uygulamasını açabilirsiniz.      ║")
	fmt.Println("║ Şifreler otomatik olarak yeni anahtar ile encrypt edilecek║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
}
