package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// AppState holds the application state
type AppState struct {
	myApp             fyne.App
	window            fyne.Window
	clients           []Client
	filteredClients   []Client
	listContainer     *fyne.Container
	searchEntry       *widget.Entry
	currentFile       string
	expandedCompanies map[string]bool         // Firma adı -> açık/kapalı durumu
	expandedApps      map[string]map[int]bool // Firma adı -> (App index -> açık/kapalı)
	activeTabIndex    map[string]int          // Firma adı -> aktif tab index
}

// FileManager handles file I/O operations
type FileManager struct {
	filePath string
}

// NewFileManager creates a new FileManager instance
func NewFileManager(filePath string) *FileManager {
	return &FileManager{
		filePath: filePath,
	}
}

// SetFilePath updates the current file path
func (fm *FileManager) SetFilePath(filePath string) {
	fm.filePath = filePath
}

// GetFilePath returns the current file path
func (fm *FileManager) GetFilePath() string {
	return fm.filePath
}

// LoadClients reads and parses client data from a JSON file
func (s *AppState) loadClients(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &s.clients); err != nil {
		return err
	}

	// Check if any password fields are unencrypted and encrypt them if needed
	hasUnencrypted := false
	for i := range s.clients {
		if (s.clients[i].VPN.Password != "" && !isEncrypted(s.clients[i].VPN.Password)) ||
			(s.clients[i].Data.JiraPassword != "" && !isEncrypted(s.clients[i].Data.JiraPassword)) {
			hasUnencrypted = true
			break
		}
		for j := range s.clients[i].Apps {
			if (s.clients[i].Apps[j].Password != "" && !isEncrypted(s.clients[i].Apps[j].Password)) ||
				(s.clients[i].Apps[j].AppServerPass != "" && !isEncrypted(s.clients[i].Apps[j].AppServerPass)) {
				hasUnencrypted = true
				break
			}
		}
		if hasUnencrypted {
			break
		}
	}

	// If unencrypted passwords found, encrypt them and save immediately
	if hasUnencrypted {
		// Yedek dosya oluştur - orijinal plaintext JSON'u koru
		backupPath := path + ".backup"
		if err := os.Rename(path, backupPath); err != nil {
			return fmt.Errorf("backup oluşturulamadı: %w", err)
		}

		clientsCopy := make([]Client, len(s.clients))
		copy(clientsCopy, s.clients)
		if err := encryptClientsInPlace(clientsCopy); err != nil {
			// Hata durumunda yedegi geri yükle
			os.Rename(backupPath, path)
			return err
		}
		// Update s.clients with encrypted values
		s.clients = clientsCopy

		// Write encrypted data back to file
		encData, err := json.MarshalIndent(s.clients, "", "  ")
		if err != nil {
			// Hata durumunda yedegi geri yükle
			os.Rename(backupPath, path)
			return err
		}
		if err := os.WriteFile(path, encData, 0644); err != nil {
			// Hata durumunda yedegi geri yükle
			os.Rename(backupPath, path)
			return err
		}
	}

	// Decrypt any encrypted password fields after loading (or after migration)
	if err := decryptClientsInPlace(s.clients); err != nil {
		return err
	}

	s.filteredClients = make([]Client, len(s.clients))
	copy(s.filteredClients, s.clients)
	s.currentFile = path
	s.window.SetTitle(fmt.Sprintf("Client Info Manager — %s", filepath.Base(path)))

	return nil
}

// SaveClients writes client data to JSON file
func (s *AppState) saveClients() error {
	// Make a copy of clients and encrypt password fields before writing
	clientsCopy := make([]Client, len(s.clients))
	copy(clientsCopy, s.clients)
	if err := encryptClientsInPlace(clientsCopy); err != nil {
		return err
	}

	data, err := json.MarshalIndent(clientsCopy, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.currentFile, data, 0644)
}
