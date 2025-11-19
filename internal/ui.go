package main

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// createEditableLabel düzenlenebilir label ve kopyalama butonu oluşturur
func (s *AppState) createEditableLabel(text string, multiLine bool, clientIndex int, updateFunc func(*Client, string)) fyne.CanvasObject {
	// Eğer text hala encrypted ise (enc: prefix varsa), decrypt et
	if isEncrypted(text) {
		decrypted, err := decryptString(text)
		if err == nil {
			text = decrypted
		}
		// Hata olursa şifreli hali olarak devam et
	}

	editLabel := newEditableLabel(text, multiLine, func(newText string) {
		if clientIndex >= 0 && clientIndex < len(s.clients) {
			updateFunc(&s.clients[clientIndex], newText)
			if err := s.saveClients(); err != nil {
				dialog.ShowError(err, s.window)
			}
		}
	})

	// Kopyalama butonu - IconButton ile oluştur
	copyBtn := NewIconButtonSimple(
		theme.ContentCopyIcon(),
		"",
		fyne.NewSize(18, 18),
		"Kopyala - Metni panoya kopyala",
		func() {
			s.window.Clipboard().SetContent(editLabel.text)
		},
	)

	return container.NewBorder(nil, nil, copyBtn, nil, editLabel)
}

// createClickableURLLabel tıklanabilir URL label oluşturur
func (s *AppState) createClickableURLLabel(text string, clientIndex int, updateFunc func(*Client, string)) fyne.CanvasObject {
	// Eğer text hala encrypted ise (enc: prefix varsa), decrypt et
	if isEncrypted(text) {
		decrypted, err := decryptString(text)
		if err == nil {
			text = decrypted
		}
		// Hata olursa şifreli hali olarak devam et
	}

	urlLabel := newClickableURLLabel(text, func(newText string) {
		if clientIndex >= 0 && clientIndex < len(s.clients) {
			updateFunc(&s.clients[clientIndex], newText)
			if err := s.saveClients(); err != nil {
				dialog.ShowError(err, s.window)
			}
		}
	})

	// Kopyalama butonu - IconButton ile oluştur
	copyBtn := NewIconButtonSimple(
		theme.ContentCopyIcon(),
		"",
		fyne.NewSize(18, 18),
		"Kopyala - URL'yi panoya kopyala",
		func() {
			s.window.Clipboard().SetContent(urlLabel.text)
		},
	)

	return container.NewBorder(nil, nil, copyBtn, nil, urlLabel)
}

// createAppUsersWidget kullanıcı/şifre listesi oluşturur
func (s *AppState) createAppUsersWidget(appUsers []string, companyName string, appIdx int) *appUsersWidget {
	usersWidget := newAppUsersWidget(appUsers, func(newUsers []string) {
		// Gerçek client'ı firma adına göre bul
		for i := range s.clients {
			if s.clients[i].Company == companyName {
				if appIdx >= 0 && appIdx < len(s.clients[i].Apps) {
					s.clients[i].Apps[appIdx].AppUsers = newUsers
					if err := s.saveClients(); err != nil {
						dialog.ShowError(err, s.window)
					}
				}
				break
			}
		}
	})
	return usersWidget
}

// buildUI ana arayüzü oluşturur
func (s *AppState) buildUI() fyne.CanvasObject {
	// Search box
	s.searchEntry = widget.NewEntry()
	s.searchEntry.SetPlaceHolder("Firma ara...")
	s.searchEntry.OnChanged = func(text string) {
		s.filterClients(text)
	}

	// Liste container'ı oluştur
	s.listContainer = container.NewVBox()
	s.buildClientList()

	// Create hamburger button that shows menu items
	var hamburgerBtn *widget.Button
	hamburgerBtn = widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		// Show popup menu with options
		newFirmaItem := fyne.NewMenuItem("Yeni Firma", func() {
			s.addClient()
		})
		newFirmaItem.Icon = theme.ContentAddIcon()

		importItem := fyne.NewMenuItem("Müşteri Import", func() {
			s.importClientFromCustomer()
		})
		importItem.Icon = theme.DownloadIcon()

		menu := fyne.NewMenu("",
			newFirmaItem,
			importItem,
		)
		pos := fyne.NewPos(hamburgerBtn.Position().X, hamburgerBtn.Position().Y+hamburgerBtn.Size().Height)
		widget.NewPopUpMenu(menu, s.window.Canvas()).ShowAtPosition(pos)
	})

	// Search bar with hamburger menu on the right
	searchBar := container.NewBorder(
		nil, nil,
		nil,
		hamburgerBtn,
		s.searchEntry,
	)

	// Toolbar with file path
	toolbar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabel(fmt.Sprintf("📁 %s", filepath.Base(s.currentFile))),
	)

	// Main content
	content := container.NewBorder(
		container.NewVBox(searchBar, widget.NewSeparator()),
		toolbar,
		nil, nil,
		container.NewVScroll(s.listContainer),
	)

	return content
}

// buildClientList özel liste ile firma listesini oluşturur
func (s *AppState) buildClientList() {
	s.listContainer.Objects = nil

	for i, client := range s.filteredClients {
		clientIndex := i
		item := s.createExpandableClientItem(client, clientIndex)
		s.listContainer.Objects = append(s.listContainer.Objects, item)
	}

	s.listContainer.Refresh()
}

// createExpandableClientItem genişletilebilir firma item'ı oluşturur
func (s *AppState) createExpandableClientItem(client Client, index int) fyne.CanvasObject {
	// Başlık metni - Renkli text badge'lerle
	companyLabel := widget.NewLabel(client.Company)
	companyLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Badge'leri container'a ekle
	badges := container.NewHBox()

	// VPN badge (yeşil text) - sadece VPN bilgisi varsa
	if client.VPN.App != "" || client.VPN.Host != "" || client.VPN.User != "" || client.VPN.Password != "" {
		vpnBadge := newBadge("VPN", color.RGBA{34, 197, 94, 255})
		badges.Add(vpnBadge)
	}

	// EBS Version badge (mavi text)
	ebsText := fallback(client.EBSVersion)
	if ebsText == "" {
		ebsText = "all"
	}
	ebsBadge := newBadge(ebsText, color.RGBA{59, 130, 246, 255})
	badges.Add(ebsBadge)

	var menuBtn fyne.CanvasObject

	// Hamburger menü butonu - daha küçük
	var menuOverlay *fyne.Container

	iconRes := theme.MoreVerticalIcon()
	icon := canvas.NewImageFromResource(iconRes)
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSize(12, 12)) // Daha küçük icon

	iconContainer := container.NewStack(icon)
	iconContainer.Resize(fyne.NewSize(16, 16)) // Daha küçük tıklanabilir alan

	menuBtn = NewIconButtonSimple(theme.MenuIcon(), "", fyne.NewSize(16, 16), "Menü - Dışa aktar, içe aktar, ayarlar", func() {
		// Eğer menü zaten açıksa kapat
		if menuOverlay != nil {
			s.window.Canvas().Overlays().Remove(menuOverlay)
			menuOverlay = nil
			return
		}

		// Menü öğeleri
		exportItem := newMenuItemWithIcon(theme.UploadIcon(), "Dışa Aktar", func() {
			if menuOverlay != nil {
				s.window.Canvas().Overlays().Remove(menuOverlay)
				menuOverlay = nil
			}
			s.exportClientForCustomer(index)
		})

		deleteItem := newMenuItemWithIcon(theme.DeleteIcon(), "Sil", func() {
			if menuOverlay != nil {
				s.window.Canvas().Overlays().Remove(menuOverlay)
				menuOverlay = nil
			}
			s.deleteClient(index)
		})

		// Menü içeriği
		menuItems := container.NewVBox(exportItem, deleteItem)

		// Border - theme'den açık gri çerçeve
		borderBg := canvas.NewRectangle(colorMenuBorder)

		// İç arka plan - theme'den koyu gri
		innerBg := canvas.NewRectangle(colorMenuBg)

		// Çerçeveli menü
		menuContent := container.NewStack(
			borderBg,
			container.NewPadded(
				container.NewStack(innerBg, menuItems),
			),
		)

		// Tıklanabilir arka plan - menünün dışına tıklayınca kapansın
		tapBg := newTappableBackground(func() {
			if menuOverlay != nil {
				s.window.Canvas().Overlays().Remove(menuOverlay)
				menuOverlay = nil
			}
		})

		// Buton pozisyonunu al
		btnPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(menuBtn)
		btnSize := menuBtn.Size()

		// Menü yüksekliğini dinamik hesapla: item sayısı * item yüksekliği + padding
		itemCount := 2 // exportItem + deleteItem
		itemHeight := float32(38)
		padding := float32(16) // NewPadded için toplam padding
		menuHeight := float32(itemCount)*itemHeight + padding

		// Menüyü konumlandır - hamburger menü yüksekliği kadar yukarı, genişliği kadar sağa
		menuX := btnPos.X - 180 + btnSize.Width             // Sağa kaydır
		menuY := btnPos.Y + btnSize.Height - btnSize.Height // Yukarı kaydır (aynı hizada)
		menuContent.Move(fyne.NewPos(menuX, menuY))
		menuContent.Resize(fyne.NewSize(180, menuHeight))

		// Overlay oluştur - önce arka plan, sonra menü (sıra önemli!)
		menuOverlay = container.NewWithoutLayout(tapBg, menuContent)

		// tapBg'yi tam ekran yap
		tapBg.Resize(s.window.Canvas().Size())

		s.window.Canvas().Overlays().Add(menuOverlay)
	})

	// Detay içeriği oluştur
	detailContent := s.createClientDetails(client, index)

	// Firma başlığı için custom header oluştur (badge'ler + hamburger menü)
	// accordionHeader yerine kendi header'ımızı oluşturalım

	// Sağ taraf - badges + menu button
	rightSide := container.NewHBox(badges, menuBtn)

	// Başlık satırı
	headerContent := container.NewBorder(nil, nil,
		companyLabel, // Sol
		rightSide,    // Sağ
		nil,
	)

	// Custom accordion header - sadece expand icon için
	// Title olarak boş string, content'i kendi header'ımızla değiştireceğiz
	dummyHeader := newAccordionHeader("", nil, nil, nil)

	// Header'ın içeriğini değiştir - title yerine kendi content'imizi koy
	// Bu biraz hack ama accordionHeader'ı multiple badge destekleyecek şekilde değiştirmekten daha basit

	// Expandable item oluştur
	expandableItem := newExpandableItem(dummyHeader, detailContent)

	// Header content'ini değiştir - expand icon'u koruyarak
	// accordionHeader'ın layout'unu taklit edelim

	// Önceki expand durumunu geri yükle
	if s.expandedCompanies[client.Company] {
		expandableItem.SetExpanded(true)
	}

	// Expand durumu değiştiğinde kaydet
	originalOnTap := dummyHeader.onTap
	dummyHeader.onTap = func() {
		if originalOnTap != nil {
			originalOnTap()
		}
		// Durumu kaydet
		s.expandedCompanies[client.Company] = expandableItem.IsExpanded()
	}

	// Ana container - arka plan ile
	bg := canvas.NewRectangle(colorDarkBlue)

	// expandableItem'ı kullanmak yerine manuel expand kontrolü yapalım
	// Çünkü firma header'ı çok özel (multiple badges + hamburger menu)

	detailContainer := container.NewVBox(detailContent)

	// Expand durumunu kontrol et
	if !s.expandedCompanies[client.Company] {
		detailContainer.Hide()
	}

	// Header'a tıklama event'i ekle
	tappableHeader := widget.NewButton("", func() {
		// Toggle expand
		s.expandedCompanies[client.Company] = !s.expandedCompanies[client.Company]
		if s.expandedCompanies[client.Company] {
			detailContainer.Show()
		} else {
			detailContainer.Hide()
		}
	})
	tappableHeader.Importance = widget.LowImportance

	// Button'ın görünümünü özelleştir - header content'i ile
	headerWithButton := container.NewStack(
		tappableHeader,
		container.NewPadded(headerContent),
	)

	itemContent := container.NewVBox(
		headerWithButton,
		detailContainer,
	)

	return container.NewStack(bg, container.NewPadded(itemContent))
}

// buildAccordion accordion'u filtreli client'larla yeniden oluşturur
func (s *AppState) buildAccordion() {
	s.buildClientList()
}

// wrapWithBlueBackground koyu mavi arka plan ile wrap eden yardımcı fonksiyon
func wrapWithBlueBackground(content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(colorDarkBlue) // Koyu mavi - tab içerikleri için

	// İçeriğin maksimum genişliğini sınırlandır
	maxWidthContent := newMaxWidthContainer(content, TabContentMaxWidth)

	paddedContent := container.NewPadded(maxWidthContent)
	return container.NewStack(bg, paddedContent)
}

// createClientDetails firma detaylarını (tabs) oluşturur
func (s *AppState) createClientDetails(client Client, index int) fyne.CanvasObject {
	// Tabs container
	tabs := container.NewAppTabs()

	// Firma Tab
	firmaContent := widget.NewForm(
		s.createCustomTextBoxItem("Firma Adı", client.Company, false, false, false, index, func(c *Client, v string) { c.Company = v }),
		s.createCustomTextBoxItem("EBS Versiyon", fallback(client.EBSVersion), false, false, false, index, func(c *Client, v string) { c.EBSVersion = v }),
		s.createCustomTextBoxItem("Not", fallback(client.Notes), false, true, false, index, func(c *Client, v string) { c.Notes = v }),
	)
	tabs.Append(container.NewTabItemWithIcon(TabNameCompany, theme.InfoIcon(), wrapWithBlueBackground(firmaContent)))

	// VPN Tab
	vpnForm := widget.NewForm(
		s.createCustomTextBoxItem("Uygulama", fallback(client.VPN.App), false, false, false, index, func(c *Client, v string) { c.VPN.App = v }),
		s.createCustomTextBoxItem("Host", fallback(client.VPN.Host), false, false, false, index, func(c *Client, v string) { c.VPN.Host = v }),
		s.createCustomTextBoxItem("Kullanıcı", fallback(client.VPN.User), false, false, false, index, func(c *Client, v string) { c.VPN.User = v }),
		s.createCustomTextBoxItem("Parola", fallback(client.VPN.Password), true, false, false, index, func(c *Client, v string) { c.VPN.Password = v }),
		s.createCustomTextBoxItem("2FA", fallback(client.VPN.TwoFATokenApp), false, false, false, index, func(c *Client, v string) { c.VPN.TwoFATokenApp = v }),
		s.createCustomTextBoxItem("Not", fallback(client.VPN.Notes), false, true, false, index, func(c *Client, v string) { c.VPN.Notes = v }),
	)
	tabs.Append(container.NewTabItem(TabNameVPN, wrapWithBlueBackground(vpnForm)))

	// Data Accordion
	dataContent := widget.NewForm(
		s.createCustomTextBoxItem("Jira URI", fallback(client.Data.JiraURI), false, false, true, index, func(c *Client, v string) { c.Data.JiraURI = v }),
		s.createCustomTextBoxItem("Jira User", fallback(client.Data.JiraUser), false, false, false, index, func(c *Client, v string) { c.Data.JiraUser = v }),
		s.createCustomTextBoxItem("Jira Pass", fallback(client.Data.JiraPassword), true, false, false, index, func(c *Client, v string) { c.Data.JiraPassword = v }),
		s.createCustomTextBoxItem("Kullanıcı", fallback(client.Data.User), false, false, false, index, func(c *Client, v string) { c.Data.User = v }),
		s.createCustomTextBoxItem("Pass Reset", fallback(client.Data.PasswordReset), false, false, false, index, func(c *Client, v string) { c.Data.PasswordReset = v }),
	)

	// RDC - Custom Expandable Item
	rdcContainer := container.NewVBox()
	if len(client.Data.RDC) > 0 {
		rdcTextBox := NewCustomTextBox(strings.Join(client.Data.RDC, "\n"), false, true, false, func(v string) {
			client.Data.RDC = strings.Split(strings.TrimSpace(v), "\n")
			if err := s.saveClients(); err != nil {
				dialog.ShowError(err, s.window)
			}
		}, func() fyne.Window {
			return s.window
		})

		rdcBadge := newBadge(fmt.Sprintf("%d", len(client.Data.RDC)), color.RGBA{59, 130, 246, 255})
		rdcHeader := newAccordionHeader("RDC", rdcBadge, []fyne.CanvasObject{}, nil)
		rdcItem := newExpandableItem(rdcHeader, rdcTextBox)
		rdcContainer.Add(rdcItem)
	}

	// Hosts - Custom Expandable Item
	hostsContainer := container.NewVBox()
	if len(client.Data.Hosts) > 0 {
		hostsTextBox := NewCustomTextBox(strings.Join(client.Data.Hosts, "\n"), false, true, false, func(v string) {
			client.Data.Hosts = strings.Split(strings.TrimSpace(v), "\n")
			if err := s.saveClients(); err != nil {
				dialog.ShowError(err, s.window)
			}
		}, func() fyne.Window {
			return s.window
		})

		hostsBadge := newBadge(fmt.Sprintf("%d", len(client.Data.Hosts)), color.RGBA{59, 130, 246, 255})
		hostsHeader := newAccordionHeader("Hosts", hostsBadge, []fyne.CanvasObject{}, nil)
		hostsItem := newExpandableItem(hostsHeader, hostsTextBox)
		hostsContainer.Add(hostsItem)
	}

	sistemContent := container.NewVBox(
		dataContent,
		widget.NewSeparator(),
		rdcContainer,
		hostsContainer,
	)
	tabs.Append(container.NewTabItem(TabNameSystem, wrapWithBlueBackground(sistemContent)))

	// Apps - Custom Expandable Items
	if len(client.Apps) > 0 {
		appsContainer := container.NewVBox()
		for appIdx, app := range client.Apps {
			idx := appIdx // Capture the index value to avoid closure issues
			appTypeOptions := []string{"DEV", "TEST", "UAT", "PREP", "PROD"}
			appForm := widget.NewForm(
				s.createCustomComboBoxItem("Tip", fallback(app.Type), appTypeOptions, index, func(c *Client, v string) { c.Apps[idx].Type = v }),
				s.createCustomTextBoxItem("İsim", fallback(app.Name), false, false, false, index, func(c *Client, v string) { c.Apps[idx].Name = v }),
				s.createCustomTextBoxItem("User", fallback(app.User), false, false, false, index, func(c *Client, v string) { c.Apps[idx].User = v }),
				s.createCustomTextBoxItem("Pass", fallback(app.Password), true, false, false, index, func(c *Client, v string) { c.Apps[idx].Password = v }),
				s.createCustomTextBoxItem("DB IP", fallback(app.DBServerIP), false, false, false, index, func(c *Client, v string) { c.Apps[idx].DBServerIP = v }),
				s.createCustomTextBoxItem("TNS", fallback(app.TNS), false, false, false, index, func(c *Client, v string) { c.Apps[idx].TNS = v }),
				s.createCustomTextBoxItem("App IP", fallback(app.AppServerIP), false, false, false, index, func(c *Client, v string) { c.Apps[idx].AppServerIP = v }),
				s.createCustomTextBoxItem("App URI", fallback(app.AppServerURI), false, false, true, index, func(c *Client, v string) { c.Apps[idx].AppServerURI = v }),
				s.createCustomTextBoxItem("App User", fallback(app.AppServerUser), false, false, false, index, func(c *Client, v string) { c.Apps[idx].AppServerUser = v }),
				s.createCustomTextBoxItem("App Pass", fallback(app.AppServerPass), true, false, false, index, func(c *Client, v string) { c.Apps[idx].AppServerPass = v }),
				s.createCustomTextBoxItem("SSH Params", fallback(app.SSHParams), false, false, false, index, func(c *Client, v string) { c.Apps[idx].SSHParams = v }),
				s.createCustomTextBoxItem("URI", fallback(app.AppURI), false, false, true, index, func(c *Client, v string) { c.Apps[idx].AppURI = v }),
			)

			// App Users - Custom Expandable Item
			usersWidget := s.createAppUsersWidget(app.AppUsers, client.Company, appIdx)
			// Düzenle butonu - IconButton ile oluştur
			currentUsersWidget := usersWidget
			editBtn := NewIconButtonSimple(theme.DocumentCreateIcon(), "Düzenle", fyne.NewSize(18, 18), "Düzenle - Kullanıcı adı ve şifreleri düzenle", func() {
				// Doğrudan startEdit çağır
				currentUsersWidget.startEdit()
			})

			// AppUsers için badge ve header
			usersBadge := newBadge(fmt.Sprintf("%d", len(app.AppUsers)), color.RGBA{59, 130, 246, 255})
			usersHeader := newAccordionHeader("", usersBadge, []fyne.CanvasObject{editBtn}, nil)
			usersItem := newExpandableItem(usersHeader, usersWidget)

			appForm.Append(FormLabelAppUsers, usersItem)

			// Silme butonu - IconButton ile oluştur
			deleteIcon := NewIconButtonSimple(
				theme.DeleteIcon(),
				"",
				fyne.NewSize(18, 18),
				"Sil - Bu ortamı ve tüm verilerini kalıcı olarak sil",
				func() {
					s.deleteApp(index, idx)
				},
			)

			// Ortam başlık metni
			appTitleText := fmt.Sprintf("%s - %s", fallback(app.Type), fallback(app.Name))

			// Badge yok şimdilik, istenirse eklenebilir
			// SSH Shell butonu - IP ve User varsa ekle
			headerButtons := []fyne.CanvasObject{}
			if fallback(app.AppServerIP) != "—" && fallback(app.AppServerUser) != "—" {
				sshBtn := NewIconButtonSimple(
					theme.ComputerIcon(),
					"",
					fyne.NewSize(18, 18),
					"SSH - Sunucuya SSH bağlantısı aç",
					func() {
						// AppIndex'den doğru app'i al
						if idx < len(s.clients[index].Apps) {
							s.openSSHShell(s.clients[index].Apps[idx])
						}
					},
				)
				headerButtons = append(headerButtons, sshBtn)
			}

			// Silme butonu
			headerButtons = append(headerButtons, deleteIcon)

			// Custom accordion header oluştur
			header := newAccordionHeader(
				appTitleText,
				nil,           // Badge yok
				headerButtons, // SSH butonu + Silme butonu
				nil,           // onTap daha sonra expandableItem tarafından set edilecek
			)

			// Ortam tipine göre border rengi al
			borderColor := getAppTypeBorderColor(app.Type)

			// Form'u önce border ile sarmala, sonra background ekle
			borderedForm := newBorderedContainer(appForm, borderColor, 5)
			contentWithBg := wrapWithBlueBackground(borderedForm)

			// Expandable item oluştur
			expandableApp := newExpandableItem(header, contentWithBg)

			// Önceki expand durumunu geri yükle
			if s.expandedApps[client.Company] == nil {
				s.expandedApps[client.Company] = make(map[int]bool)
			}
			if s.expandedApps[client.Company][appIdx] {
				expandableApp.SetExpanded(true)
			}

			// Expand durumu değiştiğinde kaydet
			originalOnTap := header.onTap
			header.onTap = func() {
				if originalOnTap != nil {
					originalOnTap()
				}
				// Durumu kaydet
				if s.expandedApps[client.Company] == nil {
					s.expandedApps[client.Company] = make(map[int]bool)
				}
				s.expandedApps[client.Company][idx] = expandableApp.IsExpanded()
			}

			// Container'a ekle - MaxLayout ile tam genişlik
			appsContainer.Add(container.NewMax(expandableApp))
			appsContainer.Add(widget.NewSeparator())
		}

		// Yeni ortam ekleme düğmesi - IconButton ile
		addAppBtn := NewIconButtonSimple(theme.ContentAddIcon(), "Yeni Ortam", fyne.NewSize(24, 24), "Yeni Ortam - Firma altına yeni ortam (dev, test, prod vb.) ekle", func() {
			s.addApp(index)
		})

		// Container ve butonı container'a koy
		appsWithButton := container.NewBorder(
			addAppBtn, // Top'ta buton
			nil, nil, nil,
			appsContainer, // Content
		)

		tabs.Append(container.NewTabItem(TabNameEnvironments, wrapWithBlueBackground(appsWithButton)))
	}

	// Önceki aktif tab'ı geri yükle
	if savedTabIndex, ok := s.activeTabIndex[client.Company]; ok {
		if savedTabIndex >= 0 && savedTabIndex < len(tabs.Items) {
			tabs.SelectIndex(savedTabIndex)
		}
	}

	// Tab değiştiğinde kaydet
	tabs.OnSelected = func(item *container.TabItem) {
		// Mevcut tab index'ini bul
		for i, tabItem := range tabs.Items {
			if tabItem == item {
				s.activeTabIndex[client.Company] = i
				break
			}
		}
	}

	// Tabs'ı kahverengi arka plan ile wrap et
	brownBg := canvas.NewRectangle(colorDarkcyan)
	tabsWithBg := container.NewStack(brownBg, tabs)

	return tabsWithBg
}

// fallback boş string için varsayılan değer döndürür
func fallback(value string) string {
	if strings.TrimSpace(value) == "" {
		return "—"
	}
	return value
}
