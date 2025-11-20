package main

//goversioninfo -64 -o resource.syso versioninfo.json
//go build -ldflags "-H windowsgui" -o client-manager.exe .\internal\.
//go build -trimpath -ldflags="-s -w -H windowsgui" -o client-manager.exe .\internal\.
//go run .\internal\. 2>&1

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	fynetooltip "github.com/dweymouth/fyne-tooltip"
)

func main() {
	state := &AppState{
		expandedCompanies: make(map[string]bool),
		expandedApps:      make(map[string]map[int]bool),
		activeTabIndex:    make(map[string]int),
	}
	state.myApp = app.NewWithID(AppID)
	state.myApp.Settings().SetTheme(&blueTheme{Theme: theme.DefaultTheme()})

	state.window = state.myApp.NewWindow(AppName)
	state.window.Resize(fyne.NewSize(DefaultWindowWidth, DefaultWindowHeight))
	state.window.CenterOnScreen() // Pencereyi ekran ortasında aç

	// Icon ayarla - bundled resource kullan
	state.window.SetIcon(resourceAppiconPng)

	state.currentFile = DefaultJSONFile

	if err := state.loadClients(state.currentFile); err != nil {
		// Dosya yüklenemezse sadece uyarı göster, dosyayı bozma
		if os.IsNotExist(err) {
			dialog.ShowInformation("Bilgi", DefaultJSONFile+" dosyası bulunamadı. Yeni dosya oluşturulacak.", state.window)
		} else {
			dialog.ShowError(fmt.Errorf("JSON dosyası okunamadı: %w Dosya yedeklendi ve boş başlatıldı", err), state.window)
			// Bozuk dosyayı yedekle
			backupFile := state.currentFile + ".backup"
			os.Rename(state.currentFile, backupFile)
		}
		state.clients = []Client{}
		state.filteredClients = []Client{}
	}

	content := state.buildUI()
	// Tooltip layer'ını ekle
	contentWithTooltips := fynetooltip.AddWindowToolTipLayer(content, state.window.Canvas())
	state.window.SetContent(contentWithTooltips)

	state.window.ShowAndRun()
}
