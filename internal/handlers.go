package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	nativeDialog "github.com/sqweek/dialog"
)

// filterClients filtreler client listesini arama sorgusuna göre
func (s *AppState) filterClients(query string) {
	query = strings.ToLower(strings.TrimSpace(query))

	if query == "" {
		s.filteredClients = make([]Client, len(s.clients))
		copy(s.filteredClients, s.clients)
	} else {
		s.filteredClients = []Client{}
		for _, client := range s.clients {
			if strings.Contains(strings.ToLower(client.Company), query) ||
				strings.Contains(strings.ToLower(client.EBSVersion), query) ||
				strings.Contains(strings.ToLower(client.Notes), query) {
				s.filteredClients = append(s.filteredClients, client)
			}
		}
	}

	s.buildAccordion()
}

// openFile dosya açma dialogu gösterir ve client'ları yükler
func (s *AppState) openFile() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		path := reader.URI().Path()
		if err := s.loadClients(path); err != nil {
			dialog.ShowError(err, s.window)
			return
		}

		s.filterClients("")
		s.searchEntry.SetText("")
		dialog.ShowInformation(DialogTitleSuccess, DialogMsgFileLoaded, s.window)
	}, s.window)
}

// addClient yeni firma ekleme dialogu gösterir
func (s *AppState) addClient() {
	companyEntry := widget.NewEntry()
	companyEntry.SetPlaceHolder("Firma adını girin...")

	// EBS versiyon seçimi
	ebsVersionOptions := []string{"all", "r11", "r12", "12.1", "12.2", "Cloud"}
	ebsSelect := widget.NewSelect(ebsVersionOptions, func(value string) {})
	ebsSelect.SetSelected("12.2") // Varsayılan değer

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetPlaceHolder("Notlar...")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Firma Adı:", Widget: companyEntry},
			{Text: "EBS Versiyon:", Widget: ebsSelect},
			{Text: "Not:", Widget: notesEntry},
		},
	}

	// Custom dialog boyutu için
	formContainer := container.NewVBox(form)

	customDialog := dialog.NewCustomConfirm("Yeni Firma Ekle", "Ekle", "İptal", formContainer, func(ok bool) {
		if !ok || strings.TrimSpace(companyEntry.Text) == "" {
			return
		}

		newClient := Client{
			Company:    companyEntry.Text,
			EBSVersion: ebsSelect.Selected,
			Notes:      notesEntry.Text,
			Apps:       []AppInfo{},
		}

		s.clients = append(s.clients, newClient)
		s.filterClients(s.searchEntry.Text)

		if err := s.saveClients(); err != nil {
			dialog.ShowError(err, s.window)
			return
		}

		dialog.ShowInformation(DialogTitleSuccess, DialogMsgClientAdded, s.window)
	}, s.window)

	customDialog.Resize(fyne.NewSize(500, 300))
	customDialog.Show()
}

// editClient firma düzenleme dialogu gösterir
func (s *AppState) editClient(index int) {
	if index >= len(s.clients) {
		return
	}

	client := s.clients[index]

	companyEntry := widget.NewEntry()
	companyEntry.SetText(client.Company)

	ebsEntry := widget.NewEntry()
	ebsEntry.SetText(client.EBSVersion)

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetText(client.Notes)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Firma Adı:", Widget: companyEntry},
			{Text: "EBS Versiyon:", Widget: ebsEntry},
			{Text: "Not:", Widget: notesEntry},
		},
	}

	dialog.ShowForm("Firma Düzenle", "Kaydet", "İptal", form.Items, func(ok bool) {
		if !ok {
			return
		}

		s.clients[index].Company = companyEntry.Text
		s.clients[index].EBSVersion = ebsEntry.Text
		s.clients[index].Notes = notesEntry.Text

		s.filterClients(s.searchEntry.Text)

		if err := s.saveClients(); err != nil {
			dialog.ShowError(err, s.window)
			return
		}

		dialog.ShowInformation(DialogTitleSuccess, DialogMsgClientUpdated, s.window)
	}, s.window)
}

// deleteClient firma silme onay dialogu gösterir
func (s *AppState) deleteClient(index int) {
	if index >= len(s.clients) {
		return
	}

	client := s.clients[index]

	dialog.ShowConfirm(
		DialogTitleDeleteConfirm,
		fmt.Sprintf("'%s' firmasını silmek istediğinize emin misiniz?", client.Company),
		func(ok bool) {
			if !ok {
				return
			}

			s.clients = append(s.clients[:index], s.clients[index+1:]...)
			s.filterClients(s.searchEntry.Text)

			if err := s.saveClients(); err != nil {
				dialog.ShowError(err, s.window)
				return
			}

			dialog.ShowInformation(DialogTitleSuccess, DialogMsgClientDeleted, s.window)
		},
		s.window,
	)
}

// exportClientForCustomer müşteri için VPN bilgileri olmadan export eder (şifreli)
func (s *AppState) exportClientForCustomer(index int) {
	if index >= len(s.clients) {
		return
	}

	client := s.clients[index]

	// VPN bilgilerini temizle
	clientCopy := client
	clientCopy.VPN = VPNInfo{}

	// Şifreleme yap (export dosyasında da şifre tutulsun)
	if err := encryptClientsInPlace([]Client{clientCopy}); err != nil {
		dialog.ShowError(fmt.Errorf("şifreleme hatası: %w", err), s.window)
		return
	}

	// Native Windows dialog kullan
	filename, err := nativeDialog.File().
		Title(DialogTitleSaveData).
		Filter("JSON Dosyası", "json").
		SetStartFile(client.Company + "_export.json").
		Save()

	if err != nil {
		// Kullanıcı iptal etti
		return
	}

	// JSON'a çevir ve kaydet
	data, err := json.MarshalIndent([]Client{clientCopy}, "", "  ")
	if err != nil {
		dialog.ShowError(err, s.window)
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		dialog.ShowError(err, s.window)
		return
	}

	dialog.ShowInformation(DialogTitleSuccess, DialogMsgDataExported, s.window)
}

// importClientFromCustomer müşteriden gelen JSON'u import eder (VPN bilgisi eklemeden)
func (s *AppState) importClientFromCustomer() {
	// Native Windows dialog kullan
	filename, err := nativeDialog.File().
		Title(DialogTitleOpenData).
		Filter("JSON Dosyası", "json").
		Load()

	if err != nil {
		// Kullanıcı iptal etti
		return
	}

	// Dosyayı oku
	data, err := os.ReadFile(filename)
	if err != nil {
		dialog.ShowError(fmt.Errorf(DialogMsgFileReadError+": %v", err), s.window)
		return
	}

	// Try to unmarshal into an array first (export may contain []Client)
	var importedClients []Client
	if err := json.Unmarshal(data, &importedClients); err != nil {
		// Fallback: try single Client object
		var single Client
		if err2 := json.Unmarshal(data, &single); err2 != nil {
			dialog.ShowError(fmt.Errorf(DialogMsgJSONReadError+": %v", err), s.window)
			return
		}
		importedClients = []Client{single}
	}

	// Şifreli alanları decrypt et (eğer şifreliyse)
	if err := decryptClientsInPlace(importedClients); err != nil {
		dialog.ShowError(fmt.Errorf("decrypt hatası: %v", err), s.window)
		return
	}

	// Process each imported client
	for _, imp := range importedClients {
		client := imp // capture

		// Firma adını kontrol et
		if strings.TrimSpace(client.Company) == "" {
			dialog.ShowError(errors.New(DialogMsgInvalidClientInfo), s.window)
			continue
		}

		// Mevcut firmayı kontrol et
		found := false
		foundIndex := -1
		for i, c := range s.clients {
			if c.Company == client.Company {
				found = true
				foundIndex = i
				break
			}
		}

		if found {
			// capture index and client for closure
			idx := foundIndex
			localClient := client
			dialog.ShowConfirm(
				"Firma Zaten Var",
				fmt.Sprintf("'%s' firması zaten mevcut. Üzerine yazmak ister misiniz?\n\nUyarı: Mevcut VPN bilgileri korunacak!", localClient.Company),
				func(ok bool) {
					if ok {
						// Mevcut VPN bilgilerini koru
						localClient.VPN = s.clients[idx].VPN
						s.clients[idx] = localClient
						s.filterClients(s.searchEntry.Text)
						if err := s.saveClients(); err != nil {
							dialog.ShowError(err, s.window)
							return
						}
						dialog.ShowInformation(DialogTitleSuccess, "Firma güncellendi!", s.window)
					}
				},
				s.window,
			)
		} else {
			// Yeni firma olarak ekle
			s.clients = append(s.clients, client)
			s.filterClients(s.searchEntry.Text)
			if err := s.saveClients(); err != nil {
				dialog.ShowError(err, s.window)
				continue
			}
			dialog.ShowInformation(DialogTitleSuccess, DialogMsgDataImported, s.window)
		}
	}
}

// exportAllClientsForCustomer tüm firmaları VPN bilgisi olmadan export eder (şifreli)
func (s *AppState) exportAllClientsForCustomer() {
	if len(s.clients) == 0 {
		dialog.ShowInformation(DialogTitleInfo, DialogMsgNoClientsToExport, s.window)
		return
	}

	// Tüm client'ları kopyala ve VPN bilgilerini temizle
	clientsCopy := make([]Client, len(s.clients))
	for i, client := range s.clients {
		clientsCopy[i] = client
		clientsCopy[i].VPN = VPNInfo{}
	}

	// Şifreleme yap (export dosyasında da şifreler tutulsun)
	if err := encryptClientsInPlace(clientsCopy); err != nil {
		dialog.ShowError(fmt.Errorf("şifreleme hatası: %w", err), s.window)
		return
	}

	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		data, err := json.MarshalIndent(clientsCopy, "", "  ")
		if err != nil {
			dialog.ShowError(err, s.window)
			return
		}

		if _, err := writer.Write(data); err != nil {
			dialog.ShowError(err, s.window)
			return
		}

		dialog.ShowInformation(DialogTitleSuccess, fmt.Sprintf(DialogMsgExportForClient, len(clientsCopy)), s.window)
	}, s.window)

	saveDialog.SetFileName("all_clients_export.json")
	saveDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	saveDialog.Show()
}

// addApp boş yeni ortam ekler
func (s *AppState) addApp(filteredIndex int) {
	if len(s.clients) == 0 {
		dialog.ShowInformation(DialogTitleInfo, DialogMsgAddClientFirst, s.window)
		return
	}

	if filteredIndex < 0 || filteredIndex >= len(s.filteredClients) {
		dialog.ShowError(errors.New(DialogMsgInvalidClientSelection), s.window)
		return
	}

	// filteredClients'tan firma adını al
	companyName := s.filteredClients[filteredIndex].Company

	// Gerçek clients dizisinde firma adına göre bul
	realClientIndex := -1
	for i, c := range s.clients {
		if c.Company == companyName {
			realClientIndex = i
			break
		}
	}

	if realClientIndex == -1 {
		dialog.ShowError(errors.New(DialogMsgClientNotFound), s.window)
		return
	}

	// Seçili client'a boş ortam ekle
	newApp := AppInfo{
		Type:         "TEST",
		Name:         "Yeni Ortam",
		AppServerURI: "",
		AppURI:       "",
		AppUsers:     []string{}, // Boş slice ile başlat
	}

	// Seçili client'a ekle
	s.clients[realClientIndex].Apps = append(s.clients[realClientIndex].Apps, newApp)

	// Sadece kaydet ve UI'ı yenile
	if err := s.saveClients(); err != nil {
		dialog.ShowError(err, s.window)
		return
	}

	// Filtreyi yeniden uygula
	s.filterClients(s.searchEntry.Text)
}

// deleteApp ortamı siler
func (s *AppState) deleteApp(filteredIndex, appIndex int) {
	if filteredIndex < 0 || filteredIndex >= len(s.filteredClients) {
		return
	}

	// filteredClients'tan firma adını al
	companyName := s.filteredClients[filteredIndex].Company

	// Gerçek clients dizisinde firma adına göre bul
	realClientIndex := -1
	for i, c := range s.clients {
		if c.Company == companyName {
			realClientIndex = i
			break
		}
	}

	if realClientIndex == -1 {
		return
	}

	client := &s.clients[realClientIndex]
	if appIndex < 0 || appIndex >= len(client.Apps) {
		return
	}

	appName := client.Apps[appIndex].Name

	// Onay dialogu göster
	dialog.ShowConfirm(DialogTitleDeleteEnv,
		fmt.Sprintf("'%s' ortamını silmek istediğinizden emin misiniz?", appName),
		func(confirmed bool) {
			if !confirmed {
				return
			}

			// Ortamı listeden çıkar
			client.Apps = append(client.Apps[:appIndex], client.Apps[appIndex+1:]...)

			// Kaydet ve UI'ı yenile
			if err := s.saveClients(); err != nil {
				dialog.ShowError(err, s.window)
				return
			}

			// Filtreyi yeniden uygula
			s.filterClients(s.searchEntry.Text)
		}, s.window)
}

// openSSHShell SSH shell'i açar
func (s *AppState) openSSHShell(app AppInfo) {
	// Validasyon: IP ve User gerekli
	if strings.TrimSpace(app.AppServerIP) == "" || strings.TrimSpace(app.AppServerUser) == "" {
		dialog.ShowInformation(DialogTitleSSH, DialogMsgSSHConfig, s.window)
		return
	}

	// SSH komutunu oluştur
	sshCmd := fmt.Sprintf("ssh %s@%s", app.AppServerUser, app.AppServerIP)

	// Platform'a göre terminal aç
	if runtime.GOOS == "windows" {
		// Windows: cmd ile SSH aç, yeni window'da, exit yazınca window kapatılır
		// "start" komutu yeni window açar
		// "/c" flag'ı: komut bitince window'u kapat
		psCmd := fmt.Sprintf("start cmd /c ssh %s@%s", app.AppServerUser, app.AppServerIP)

		cmd := exec.Command("cmd", "/c", psCmd)
		if err := cmd.Start(); err != nil {
			// Hata: şifreyi panoya kopyala
			s.window.Clipboard().SetContent(app.AppServerPass)
			dialog.ShowInformation(DialogTitleSSH, DialogMsgSSHPasswordCopy, s.window)
		}
	} else if runtime.GOOS == "darwin" {
		// macOS: Terminal.app ile aç
		script := fmt.Sprintf("tell app \"Terminal\" to do script \"%s; exit\"", sshCmd)
		cmd := exec.Command("osascript", "-e", script)
		if err := cmd.Start(); err != nil {
			// Hata: şifreyi panoya kopyala
			s.window.Clipboard().SetContent(app.AppServerPass)
			dialog.ShowInformation(DialogTitleSSH, DialogMsgSSHPasswordCopy, s.window)
		}
	} else {
		// Linux: xterm veya gnome-terminal ile aç
		termCmd := fmt.Sprintf("%s; exit", sshCmd)

		// xterm dene
		cmd := exec.Command("xterm", "-hold", "-e", "bash", "-c", termCmd)
		if err := cmd.Start(); err != nil {
			// xterm başarısız, gnome-terminal dene
			cmd = exec.Command("gnome-terminal", "--", "bash", "-c", termCmd)
			if err := cmd.Start(); err != nil {
				// Her iki terminal de başarısız: şifreyi panoya kopyala
				s.window.Clipboard().SetContent(app.AppServerPass)
				dialog.ShowInformation(DialogTitleSSH, DialogMsgSSHPasswordCopy, s.window)
			}
		}
	}
}
