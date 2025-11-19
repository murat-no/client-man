package main

import (
	"image/color"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	ttwidget "github.com/dweymouth/fyne-tooltip/widget"
)

// tappableIcon tıklanabilir küçük icon
type tappableIcon struct {
	widget.BaseWidget
	icon     *canvas.Image
	size     fyne.Size
	onTapped func()
	hovered  bool
}

func newTappableIcon(icon *canvas.Image, onTapped func()) *tappableIcon {
	t := &tappableIcon{
		icon:     icon,
		size:     fyne.NewSize(20, 20), // Varsayılan boyut
		onTapped: onTapped,
	}
	t.ExtendBaseWidget(t)
	return t
}

func newTappableIconWithSize(icon *canvas.Image, size fyne.Size, onTapped func()) *tappableIcon {
	t := &tappableIcon{
		icon:     icon,
		size:     size,
		onTapped: onTapped,
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *tappableIcon) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.Transparent)
	bg.SetMinSize(t.size)

	content := container.NewStack(bg, t.icon)
	return widget.NewSimpleRenderer(content)
}

func (t *tappableIcon) Tapped(_ *fyne.PointEvent) {
	if t.onTapped != nil {
		t.onTapped()
	}
}

func (t *tappableIcon) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (t *tappableIcon) MouseIn(*desktop.MouseEvent) {
	t.hovered = true
	t.Refresh()
}

func (t *tappableIcon) MouseOut() {
	t.hovered = false
	t.Refresh()
}

// ============================================================================
// IconButton - Custom button widget with icon, text, tooltip, and hover support
// Uses fyne-tooltip for professional tooltip rendering
// ============================================================================

type IconButton struct {
	ttwidget.ToolTipWidget
	icon       fyne.Resource
	text       string // Buton metni (opsiyonel)
	size       fyne.Size
	onTapped   func()
	onHoverIn  func()
	onHoverOut func()
	hovered    bool

	// Internal
	icon_widget *widget.Icon
	textLabel   *widget.Label
	bg          *canvas.Rectangle
}

// NewIconButton creates a new icon button with optional text, tooltip and hover callbacks
func NewIconButton(icon fyne.Resource, text string, size fyne.Size, tooltip string, onTapped func(), onHoverIn, onHoverOut func()) *IconButton {
	b := &IconButton{
		icon:       icon,
		text:       text,
		size:       size,
		onTapped:   onTapped,
		onHoverIn:  onHoverIn,
		onHoverOut: onHoverOut,
		bg:         canvas.NewRectangle(color.Transparent),
	}
	b.ExtendBaseWidget(b)
	if tooltip != "" {
		b.SetToolTip(tooltip)
	}
	return b
}

// NewIconButtonSimple creates a simple icon button without hover callbacks
func NewIconButtonSimple(icon fyne.Resource, text string, size fyne.Size, tooltip string, onTapped func()) *IconButton {
	return NewIconButton(icon, text, size, tooltip, onTapped, nil, nil)
}

func (b *IconButton) CreateRenderer() fyne.WidgetRenderer {
	return newIconButtonRenderer(b)
}

type iconButtonRenderer struct {
	button  *IconButton
	bg      *canvas.Rectangle
	icon    *widget.Icon
	textLbl *widget.Label
}

func (r *iconButtonRenderer) Layout(space fyne.Size) {
	r.bg.Resize(space)
	if r.textLbl != nil {
		// Icon ve text
		iconSize := r.button.size
		r.icon.Resize(iconSize)
		r.icon.Move(fyne.NewPos(2, (space.Height-iconSize.Height)/2))

		textSize := r.textLbl.MinSize()
		r.textLbl.Resize(textSize)
		r.textLbl.Move(fyne.NewPos(iconSize.Width+4, (space.Height-textSize.Height)/2))
	} else {
		// Sadece icon
		iconSize := r.button.size
		r.icon.Resize(iconSize)
		r.icon.Move(fyne.NewPos((space.Width-iconSize.Width)/2, (space.Height-iconSize.Height)/2))
	}
}

func (r *iconButtonRenderer) MinSize() fyne.Size {
	if r.textLbl != nil {
		textSize := r.textLbl.MinSize()
		return fyne.NewSize(r.button.size.Width+4+textSize.Width, r.button.size.Height)
	}
	return r.button.size
}

func (r *iconButtonRenderer) Refresh() {
	r.icon.SetResource(r.button.icon)
	if r.textLbl != nil {
		r.textLbl.SetText(r.button.text)
	}
	if r.button.hovered {
		r.bg.FillColor = color.NRGBA{R: 200, G: 200, B: 200, A: 50}
	} else {
		r.bg.FillColor = color.Transparent
	}
	r.bg.Refresh()
}

func (r *iconButtonRenderer) Objects() []fyne.CanvasObject {
	if r.textLbl != nil {
		return []fyne.CanvasObject{r.bg, r.icon, r.textLbl}
	}
	return []fyne.CanvasObject{r.bg, r.icon}
}

func (r *iconButtonRenderer) Destroy() {}

func newIconButtonRenderer(b *IconButton) *iconButtonRenderer {
	r := &iconButtonRenderer{
		button: b,
		bg:     canvas.NewRectangle(color.Transparent),
		icon:   widget.NewIcon(b.icon),
	}
	if b.text != "" {
		r.textLbl = widget.NewLabel(b.text)
	}
	return r
}

func (b *IconButton) Tapped(_ *fyne.PointEvent) {
	if b.onTapped != nil {
		b.onTapped()
	}
}

func (b *IconButton) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (b *IconButton) MouseIn(e *desktop.MouseEvent) {
	b.ToolTipWidget.MouseIn(e)
	b.hovered = true
	if b.onHoverIn != nil {
		b.onHoverIn()
	}
	b.Refresh()
}

func (b *IconButton) MouseOut() {
	b.ToolTipWidget.MouseOut()
	b.hovered = false
	if b.onHoverOut != nil {
		b.onHoverOut()
	}
	b.Refresh()
}

func (b *IconButton) MouseMoved(e *desktop.MouseEvent) {
	b.ToolTipWidget.MouseMoved(e)
}

// SetIcon updates the icon resource and refreshes the widget so UI reflects the change immediately.
func (b *IconButton) SetIcon(res fyne.Resource) {
	b.icon = res
	b.Refresh()
}

// tappableBackground tıklanabilir şeffaf arka plan (overlay kapatmak için)
type tappableBackground struct {
	widget.BaseWidget
	onTapped func()
	bg       *canvas.Rectangle
}

func newTappableBackground(onTapped func()) *tappableBackground {
	t := &tappableBackground{
		onTapped: onTapped,
		bg:       canvas.NewRectangle(color.Transparent),
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *tappableBackground) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.bg)
}

func (t *tappableBackground) Tapped(_ *fyne.PointEvent) {
	if t.onTapped != nil {
		t.onTapped()
	}
}

func (t *tappableBackground) Resize(size fyne.Size) {
	t.BaseWidget.Resize(size)
	t.bg.Resize(size)
}

// menuItem menü benzeri tıklanabilir item
type menuItem struct {
	widget.BaseWidget
	text      string
	iconRes   fyne.Resource
	onTapped  func()
	bg        *canvas.Rectangle
	iconImage *canvas.Image
	textLabel *widget.Label
	hovered   bool
}

func newMenuItemWithIcon(iconRes fyne.Resource, text string, onTapped func()) *menuItem {
	item := &menuItem{
		text:     text,
		iconRes:  iconRes,
		onTapped: onTapped,
		bg:       canvas.NewRectangle(color.Transparent),
	}

	// Icon için image
	item.iconImage = canvas.NewImageFromResource(iconRes)
	item.iconImage.FillMode = canvas.ImageFillContain
	item.iconImage.SetMinSize(fyne.NewSize(16, 16))

	// Text label
	item.textLabel = widget.NewLabel(text)
	item.textLabel.Alignment = fyne.TextAlignLeading

	item.ExtendBaseWidget(item)
	return item
}

// Backward compatibility için eski newMenuItem fonksiyonu
func newMenuItem(icon, text string, onTapped func()) *menuItem {
	item := &menuItem{
		text:     text,
		onTapped: onTapped,
		bg:       canvas.NewRectangle(color.Transparent),
	}

	// Icon text olarak
	iconText := canvas.NewText(icon, theme.ForegroundColor())
	iconText.TextSize = 16
	iconText.Alignment = fyne.TextAlignCenter

	// Text label
	item.textLabel = widget.NewLabel(text)
	item.textLabel.Alignment = fyne.TextAlignLeading

	item.ExtendBaseWidget(item)
	return item
}

func (m *menuItem) CreateRenderer() fyne.WidgetRenderer {
	m.bg.SetMinSize(fyne.NewSize(180, 32))

	var iconContainer *fyne.Container
	if m.iconImage != nil {
		// Icon resource varsa image kullan
		iconContainer = container.NewCenter(m.iconImage)
		iconContainer.Resize(fyne.NewSize(24, 24))
	} else {
		// Yoksa boş container
		iconContainer = container.NewCenter()
		iconContainer.Resize(fyne.NewSize(24, 24))
	}

	contentBox := container.NewBorder(nil, nil, iconContainer, nil, m.textLabel)
	paddedContent := container.NewPadded(contentBox)
	content := container.NewStack(m.bg, paddedContent)

	return widget.NewSimpleRenderer(content)
}

func (m *menuItem) Tapped(_ *fyne.PointEvent) {
	if m.onTapped != nil {
		m.onTapped()
	}
}

func (m *menuItem) TappedSecondary(_ *fyne.PointEvent) {}

func (m *menuItem) MouseIn(_ *desktop.MouseEvent) {
	m.hovered = true
	m.bg.FillColor = colorMenuHover // Theme'den hover rengi
	m.bg.Refresh()
}

func (m *menuItem) MouseOut() {
	m.hovered = false
	m.bg.FillColor = color.Transparent // Normal durumda şeffaf
	m.bg.Refresh()
}

func (m *menuItem) MouseMoved(_ *desktop.MouseEvent) {}

// parseURL URL parse eder, http:// yoksa ekler
func parseURL(urlStr string) *url.URL {
	urlStr = strings.TrimSpace(urlStr)
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}
	parsed, _ := url.Parse(urlStr)
	return parsed
}

// Özel tıklanabilir kopyalama düğmesi
type copyButton struct {
	widget.BaseWidget
	onTapped func()
	bg       *canvas.Rectangle
	icon     *widget.Icon
	hovered  bool
}

func newCopyButton(onTapped func()) *copyButton {
	btn := &copyButton{
		onTapped: onTapped,
		bg:       canvas.NewRectangle(colorDarkBlue),
		icon:     widget.NewIcon(theme.ContentCopyIcon()),
	}
	btn.ExtendBaseWidget(btn)
	return btn
}

func (b *copyButton) CreateRenderer() fyne.WidgetRenderer {
	b.bg.SetMinSize(fyne.NewSize(28, 28))
	iconContainer := container.NewCenter(b.icon)
	content := container.NewStack(b.bg, iconContainer)

	return widget.NewSimpleRenderer(content)
}

func (b *copyButton) Tapped(_ *fyne.PointEvent) {
	if b.onTapped != nil {
		b.onTapped()
	}
}

func (b *copyButton) TappedSecondary(_ *fyne.PointEvent) {}

func (b *copyButton) MouseIn(_ *desktop.MouseEvent) {
	b.hovered = true
	// Hover'da biraz daha açık mavi
	b.bg.FillColor = colorLightBlue
	b.bg.Refresh()
}

func (b *copyButton) MouseOut() {
	b.hovered = false
	// Normal koyu mavi
	b.bg.FillColor = colorDarkBlue
	b.bg.Refresh()
}

func (b *copyButton) MouseMoved(_ *desktop.MouseEvent) {}

func (b *copyButton) setIcon(res fyne.Resource) {
	b.icon.SetResource(res)
	b.icon.Refresh()
}

// clickableURLLabel tıklanabilir hyperlink stili label
type clickableURLLabel struct {
	widget.BaseWidget
	text      string
	label     *canvas.Text
	onSave    func(string)
	editing   bool
	entry     *widget.Entry
	container *fyne.Container
}

func newClickableURLLabel(text string, onSave func(string)) *clickableURLLabel {
	ul := &clickableURLLabel{
		text:   text,
		onSave: onSave,
	}

	// Canvas.Text ile mavi altı çizgili text oluştur
	ul.label = canvas.NewText(text, color.RGBA{R: 56, G: 118, B: 233, A: 255}) // Mavi
	ul.label.TextStyle = fyne.TextStyle{}                                      // Fyne'da underline yok, ama renk yeterli
	ul.label.TextSize = 14

	ul.ExtendBaseWidget(ul)
	return ul
}

func (ul *clickableURLLabel) CreateRenderer() fyne.WidgetRenderer {
	ul.container = container.NewStack(ul.label)
	return widget.NewSimpleRenderer(ul.container)
}

func (ul *clickableURLLabel) Tapped(ev *fyne.PointEvent) {
	if ul.editing {
		return
	}

	// Shift basılıysa edit moda geç
	if ev != nil && (ev.AbsolutePosition.X < 0) { // Shift kontrolü için workaround
		ul.startEdit()
		return
	}

	urlStr := strings.TrimSpace(ul.text)
	if urlStr == "" || urlStr == "—" {
		return
	}

	// URL'yi tarayıcıda aç
	if parsedURL := parseURL(urlStr); parsedURL != nil {
		_ = fyne.CurrentApp().OpenURL(parsedURL)
	}
}

func (ul *clickableURLLabel) TappedSecondary(_ *fyne.PointEvent) {
	// Sağ tık ile de edit
	ul.startEdit()
}

func (ul *clickableURLLabel) DoubleTapped(_ *fyne.PointEvent) {
	// Çift tıklama ile edit
	ul.startEdit()
}

func (ul *clickableURLLabel) startEdit() {
	if ul.editing || ul.container == nil {
		return
	}

	ul.editing = true
	ul.entry = widget.NewEntry()
	ul.entry.SetText(ul.text)

	// Kaydet butonu (✓)
	saveBtn := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		ul.saveEdit(ul.entry.Text)
	})
	saveBtn.Importance = widget.HighImportance

	// İptal butonu (✕)
	cancelBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		ul.cancelEdit()
	})
	cancelBtn.Importance = widget.LowImportance

	ul.entry.OnSubmitted = func(newText string) {
		ul.saveEdit(newText)
	}

	if ul.container != nil {
		// Entry + butonlar
		editContainer := container.NewBorder(nil, nil, nil,
			container.NewHBox(saveBtn, cancelBtn),
			ul.entry,
		)
		ul.container.Objects = []fyne.CanvasObject{editContainer}
		ul.container.Refresh()
		if canvas := fyne.CurrentApp().Driver().CanvasForObject(ul.entry); canvas != nil {
			canvas.Focus(ul.entry)
		}
	}
}

func (ul *clickableURLLabel) saveEdit(newText string) {
	ul.text = newText
	ul.updateLabel()
	if ul.container != nil {
		ul.container.Objects = []fyne.CanvasObject{ul.label}
		ul.container.Refresh()
	}
	ul.editing = false
	if ul.onSave != nil {
		ul.onSave(newText)
	}
}

func (ul *clickableURLLabel) cancelEdit() {
	if ul.container != nil {
		ul.container.Objects = []fyne.CanvasObject{ul.label}
		ul.container.Refresh()
	}
	ul.editing = false
}

func (ul *clickableURLLabel) updateLabel() {
	ul.label.Text = ul.text
	ul.label.Refresh()
}

func (ul *clickableURLLabel) Cursor() desktop.Cursor {
	if !ul.editing {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

func (ul *clickableURLLabel) MouseIn(*desktop.MouseEvent) {}
func (ul *clickableURLLabel) MouseOut()                   {}

// appUsersWidget kullanıcı/şifre listesi widget'ı
type appUsersWidget struct {
	widget.BaseWidget
	users   []string
	onSave  func([]string)
	list    *fyne.Container
	editing bool
	entry   *widget.Entry
}

func newAppUsersWidget(users []string, onSave func([]string)) *appUsersWidget {
	w := &appUsersWidget{
		users:  users,
		onSave: onSave,
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *appUsersWidget) CreateRenderer() fyne.WidgetRenderer {
	w.list = w.buildUserList()
	return widget.NewSimpleRenderer(w.list)
}

func (w *appUsersWidget) buildUserList() *fyne.Container {
	items := []fyne.CanvasObject{}

	for _, userLine := range w.users {
		userLine = strings.TrimSpace(userLine)
		if userLine == "" {
			continue
		}

		// "kullanıcı/şifre" formatını parse et
		parts := strings.SplitN(userLine, "/", 2)
		username := parts[0]
		password := ""
		if len(parts) > 1 {
			password = parts[1]
		}

		// Kullanıcı label
		userLabel := widget.NewLabel("👤 " + username)
		userLabel.TextStyle = fyne.TextStyle{Bold: true}

		// Şifre label (gizli)
		passLabel := widget.NewLabel("🔒 " + strings.Repeat("•", len(password)))

		// Göster/Gizle butonu
		var showBtn *widget.Button
		isVisible := false
		showBtn = widget.NewButtonWithIcon("", theme.VisibilityIcon(), func() {
			isVisible = !isVisible
			if isVisible {
				passLabel.SetText("🔒 " + password)
				showBtn.SetIcon(theme.VisibilityOffIcon())
			} else {
				passLabel.SetText("🔒 " + strings.Repeat("•", len(password)))
				showBtn.SetIcon(theme.VisibilityIcon())
			}
		})
		showBtn.Importance = widget.LowImportance

		// Kopyalama butonu
		copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
			fyne.CurrentApp().Driver().AllWindows()[0].Clipboard().SetContent(password)
		})
		copyBtn.Importance = widget.LowImportance

		// Satır layout
		userRow := container.NewBorder(nil, nil,
			userLabel,
			container.NewHBox(showBtn, copyBtn),
			passLabel,
		)

		items = append(items, userRow)
	}

	// Eğer liste boşsa bilgi göster
	if len(items) == 0 {
		items = append(items, widget.NewLabel("Kullanıcı yok"))
	}

	return container.NewVBox(items...)
}

func (w *appUsersWidget) startEdit() {
	if w.editing {
		return
	}

	w.editing = true
	w.entry = widget.NewMultiLineEntry()

	// users nil veya boş olabilir, güvenli join
	if w.users == nil {
		w.users = []string{}
	}
	w.entry.SetText(strings.Join(w.users, "\n"))
	w.entry.Wrapping = fyne.TextWrapWord

	// Kaydet butonu
	saveBtn := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		text := strings.TrimSpace(w.entry.Text)
		var newUsers []string

		if text != "" {
			newUsers = strings.Split(text, "\n")
		} else {
			newUsers = []string{}
		}

		w.users = newUsers

		// Liste container'ı yeniden oluştur
		newList := w.buildUserList()
		if w.list != nil {
			w.list.Objects = newList.Objects
			w.list.Refresh()
		}
		w.editing = false

		if w.onSave != nil {
			w.onSave(newUsers)
		}
	})
	saveBtn.Importance = widget.HighImportance

	// İptal butonu
	cancelBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		// Liste container'ı yeniden oluştur
		newList := w.buildUserList()
		if w.list != nil {
			w.list.Objects = newList.Objects
			w.list.Refresh()
		}
		w.editing = false
	})
	cancelBtn.Importance = widget.LowImportance

	editContainer := container.NewBorder(
		nil,
		container.NewHBox(saveBtn, cancelBtn),
		nil, nil,
		w.entry,
	)

	if w.list != nil {
		w.list.Objects = []fyne.CanvasObject{editContainer}
		w.list.Refresh()
	}
}

// Düzenlenebilir label widget'ı
type editableLabel struct {
	widget.BaseWidget
	text      string
	label     *widget.Label
	entry     *widget.Entry
	editing   bool
	multiLine bool
	onSave    func(string)
	container *fyne.Container
}

func newEditableLabel(text string, multiLine bool, onSave func(string)) *editableLabel {
	el := &editableLabel{
		text:      text,
		label:     widget.NewLabel(text),
		multiLine: multiLine,
		onSave:    onSave,
	}
	el.label.Wrapping = fyne.TextWrapWord
	el.ExtendBaseWidget(el)
	return el
}

func (el *editableLabel) CreateRenderer() fyne.WidgetRenderer {
	el.container = container.NewStack(el.label)
	return widget.NewSimpleRenderer(el.container)
}

func (el *editableLabel) DoubleTapped(_ *fyne.PointEvent) {
	if el.editing {
		// Zaten edit modundaysa tüm metni seç: önce entry'ye focus ver, sonra SelectAll shortcut'ını gönder
		if el.entry != nil {
			if canvas := fyne.CurrentApp().Driver().CanvasForObject(el.entry); canvas != nil {
				canvas.Focus(el.entry)
			}
			// Ctrl+A (SelectAll) shortcut'ını uygula
			selectAll := &fyne.ShortcutSelectAll{}
			el.entry.TypedShortcut(selectAll)
		}
		return
	}

	el.editing = true
	if el.multiLine {
		el.entry = widget.NewMultiLineEntry()
		el.entry.Wrapping = fyne.TextWrapWord
	} else {
		el.entry = widget.NewEntry()
	}
	el.entry.SetText(el.text)

	// OK düğmesi
	okBtn := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		el.text = el.entry.Text
		if el.onSave != nil {
			el.onSave(el.text)
		}
		el.label.SetText(el.text)
		el.editing = false
		el.container.Objects = []fyne.CanvasObject{el.label}
		el.container.Refresh()
	})
	okBtn.Importance = widget.HighImportance

	// Cancel düğmesi
	cancelBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		el.editing = false
		el.container.Objects = []fyne.CanvasObject{el.label}
		el.container.Refresh()
	})
	cancelBtn.Importance = widget.LowImportance

	buttons := container.NewHBox(okBtn, cancelBtn)
	editBox := container.NewBorder(nil, nil, nil, buttons, el.entry)

	el.container.Objects = []fyne.CanvasObject{editBox}
	el.container.Refresh()
	// Fokusu entry'e ver ki Ctrl+A ve diğer klavye kısayolları çalışsın
	if canvas := fyne.CurrentApp().Driver().CanvasForObject(el.entry); canvas != nil {
		canvas.Focus(el.entry)
	}
}

func (el *editableLabel) Tapped(_ *fyne.PointEvent)          {}
func (el *editableLabel) TappedSecondary(_ *fyne.PointEvent) {}

func (el *editableLabel) setText(text string) {
	el.text = text
	el.label.SetText(text)
	el.label.Refresh()
}

// Panoya kopyalama butonu ile düzenlenebilir label oluştur (tek satır)
func createLabelWithCopy(text string, window fyne.Window, onSave func(string)) fyne.CanvasObject {
	// Düzenlenebilir label (tek satır)
	editLabel := newEditableLabel(text, false, onSave)

	// Kopyalama butonu - IconButton ile oluştur
	copyBtn := NewIconButtonSimple(
		theme.ContentCopyIcon(),
		"",
		fyne.NewSize(18, 18),
		"Kopyala - Metni panoya kopyala",
		func() {
			window.Clipboard().SetContent(editLabel.text)
		},
	)

	return container.NewBorder(nil, nil, copyBtn, nil, editLabel)
}

// Şifre alanı için özel widget (gizli/görünür toggle ve düzenlenebilir)
func createPasswordLabelWithCopy(text string, window fyne.Window) fyne.CanvasObject {
	// Gizli/görünür durumu
	hidden := true

	// Düzenlenebilir label (şifre için) - bu fonksiyon artık kullanılmıyor
	editLabel := newEditableLabel(strings.Repeat("•", len(text)), false, func(newText string) {
		text = newText
		dialog.ShowInformation("Bilgi", "Şifre değiştirildi", window)
	})

	// Kopyalama butonu
	copyBtn := NewIconButtonSimple(
		theme.ContentCopyIcon(),
		"",
		fyne.NewSize(18, 18),
		"Kopyala - Şifreyi panoya kopyala",
		func() {
			window.Clipboard().SetContent(text)
		},
	)

	// Göz ikonu butonu - toggle görünürlük
	eyeBtn := NewIconButton(
		theme.VisibilityIcon(),
		"",
		fyne.NewSize(18, 18),
		"Göster/Gizle - Şifreyi göster veya gizle",
		func() {
			hidden = !hidden
			if hidden {
				editLabel.label.SetText(strings.Repeat("•", len(text)))
			} else {
				editLabel.label.SetText(text)
			}
			editLabel.label.Refresh()
		},
		nil,
		nil,
	)

	return container.NewBorder(nil, nil, copyBtn, eyeBtn, editLabel)
}

func (s *AppState) createEditablePasswordLabel(text string, clientIndex int, updateFunc func(*Client, string)) fyne.CanvasObject {
	hidden := true

	// Eğer şifre hala encrypted ise (enc: prefix varsa), decrypt et
	if isEncrypted(text) {
		decrypted, err := decryptString(text)
		if err == nil {
			text = decrypted
		}
		// Hata olursa şifreli hali olarak devam et
	}

	editLabel := newEditableLabel(strings.Repeat("•", len(text)), false, func(newText string) {
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
		"Kopyala - Şifreyi panoya kopyala",
		func() {
			s.window.Clipboard().SetContent(text)
		},
	)

	// Göz ikonu butonu - toggle görünürlük
	eyeBtn := NewIconButton(
		theme.VisibilityIcon(),
		"",
		fyne.NewSize(18, 18),
		"Göster/Gizle - Şifreyi göster veya gizle",
		func() {
			hidden = !hidden
			if hidden {
				editLabel.label.SetText(strings.Repeat("•", len(text)))
			} else {
				editLabel.label.SetText(text)
			}
			editLabel.label.Refresh()
		},
		nil,
		nil,
	)

	return container.NewBorder(nil, nil, copyBtn, eyeBtn, editLabel)
}

// Çift tıklanabilir select wrapper
type tappableSelect struct {
	widget.BaseWidget
	selectWidget     *widget.Select
	label            *widget.Label
	container        *fyne.Container
	editing          bool
	originalOnChange func(string)
}

func newTappableSelect(options []string, selected string, onChange func(string)) *tappableSelect {
	ts := &tappableSelect{
		label:            widget.NewLabel(selected),
		originalOnChange: onChange,
		editing:          false,
	}
	ts.label.Wrapping = fyne.TextWrapWord

	// Container'ı önce oluştur - başlangıçta label ile
	ts.container = container.NewStack(ts.label)

	ts.selectWidget = widget.NewSelect(options, func(sel string) {
		ts.label.SetText(sel)
		ts.editing = false
		// Container'ı label'a geri döndür
		ts.container.Objects = []fyne.CanvasObject{ts.label}
		ts.container.Refresh()
		if ts.originalOnChange != nil {
			ts.originalOnChange(sel)
		}
	})
	ts.selectWidget.SetSelected(selected)

	ts.ExtendBaseWidget(ts)
	return ts
}

func (ts *tappableSelect) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(ts.container)
}

func (ts *tappableSelect) DoubleTapped(_ *fyne.PointEvent) {
	if !ts.editing {
		ts.editing = true
		// Container'ı select widget'a çevir
		ts.container.Objects = []fyne.CanvasObject{ts.selectWidget}
		ts.container.Refresh()
	}
}

func (ts *tappableSelect) Tapped(_ *fyne.PointEvent) {}

func (ts *tappableSelect) TappedSecondary(_ *fyne.PointEvent) {}

// Select dropdown ile düzenlenebilir alan
func (s *AppState) createEditableSelect(text string, options []string, clientIndex int, updateFunc func(*Client, string)) fyne.CanvasObject {
	// Eğer text hala encrypted ise (enc: prefix varsa), decrypt et
	if isEncrypted(text) {
		decrypted, err := decryptString(text)
		if err == nil {
			text = decrypted
		}
		// Hata olursa şifreli hali olarak devam et
	}

	tappable := newTappableSelect(options, text, func(selected string) {
		if clientIndex >= 0 && clientIndex < len(s.clients) {
			updateFunc(&s.clients[clientIndex], selected)
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
		"Kopyala - Değeri panoya kopyala",
		func() {
			s.window.Clipboard().SetContent(tappable.label.Text)
		},
	)

	return container.NewBorder(nil, nil, copyBtn, nil, tappable)
}

// badge renkli badge widget'ı
type badge struct {
	widget.BaseWidget
	text  string
	color color.Color
}

func newBadge(text string, bgColor color.Color) *badge {
	b := &badge{text: text, color: bgColor}
	b.ExtendBaseWidget(b)
	return b
}

func (b *badge) CreateRenderer() fyne.WidgetRenderer {
	// Arka plan yok, sadece renkli text
	label := canvas.NewText(b.text, b.color) // Arka plan rengi yerine text rengini kullan
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle = fyne.TextStyle{Bold: true} // Text için bold daha iyi
	label.TextSize = 10                          // Biraz büyütelim

	objects := []fyne.CanvasObject{label} // Sadece label, bg yok

	return &badgeRenderer{
		badge:   b,
		bg:      nil, // Arka plan yok
		label:   label,
		objects: objects,
	}
}

func (b *badge) MinSize() fyne.Size {
	return fyne.NewSize(30, 18) // Daha kompakt
}

type badgeRenderer struct {
	badge   *badge
	bg      *canvas.Rectangle
	label   *canvas.Text
	objects []fyne.CanvasObject
}

func (r *badgeRenderer) Layout(size fyne.Size) {
	if r.bg != nil {
		r.bg.Resize(size)
	}
	r.label.Resize(size)
	r.label.Move(fyne.NewPos(0, 0)) // Merkeze hizala
}

func (r *badgeRenderer) MinSize() fyne.Size {
	return r.badge.MinSize()
}

func (r *badgeRenderer) Refresh() {
	if r.bg != nil {
		r.bg.FillColor = r.badge.color
		r.bg.Refresh()
	}
	r.label.Text = r.badge.text
	r.label.Color = r.badge.color // Text rengini güncelle
	r.label.Refresh()
}

func (r *badgeRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *badgeRenderer) Destroy() {}

// accordionHeader custom accordion başlık widget'ı - sol metin, sağ badge ve butonlar
type accordionHeader struct {
	widget.BaseWidget
	title      string
	badge      *badge
	buttons    []fyne.CanvasObject
	expanded   bool
	onTap      func()
	container  *fyne.Container
	expandIcon *widget.Button
}

func newAccordionHeader(title string, badge *badge, buttons []fyne.CanvasObject, onTap func()) *accordionHeader {
	h := &accordionHeader{
		title:   title,
		badge:   badge,
		buttons: buttons,
		onTap:   onTap,
	}
	h.ExtendBaseWidget(h)
	return h
}

func (h *accordionHeader) CreateRenderer() fyne.WidgetRenderer {
	// Başlık label
	titleLabel := widget.NewLabel(h.title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Expand/collapse ikonu
	h.expandIcon = widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
		h.expanded = !h.expanded
		if h.expanded {
			h.expandIcon.SetIcon(theme.MenuDropUpIcon())
		} else {
			h.expandIcon.SetIcon(theme.MenuDropDownIcon())
		}
		if h.onTap != nil {
			h.onTap()
		}
	})
	h.expandIcon.Importance = widget.LowImportance

	// Sağ taraf: badge + butonlar + expand icon
	rightItems := []fyne.CanvasObject{}
	if h.badge != nil {
		rightItems = append(rightItems, h.badge)
	}
	rightItems = append(rightItems, h.buttons...)
	rightItems = append(rightItems, h.expandIcon)

	rightContainer := container.NewHBox(rightItems...)

	// Ana container
	h.container = container.NewBorder(
		nil, nil,
		nil, rightContainer, // Sağda: badge + butonlar + expand
		titleLabel, // Solda: başlık
	)

	return widget.NewSimpleRenderer(h.container)
}

func (h *accordionHeader) Tapped(_ *fyne.PointEvent) {
	h.expanded = !h.expanded
	if h.expanded {
		h.expandIcon.SetIcon(theme.MenuDropUpIcon())
	} else {
		h.expandIcon.SetIcon(theme.MenuDropDownIcon())
	}
	if h.onTap != nil {
		h.onTap()
	}
}

// expandableItem custom accordion item - header + content
type expandableItem struct {
	widget.BaseWidget
	header      *accordionHeader
	content     fyne.CanvasObject
	expanded    bool
	borderColor color.Color
	borderWidth float32
}

func newExpandableItem(header *accordionHeader, content fyne.CanvasObject) *expandableItem {
	return newExpandableItemWithBorder(header, content, color.Transparent, 0)
}

func newExpandableItemWithBorder(header *accordionHeader, content fyne.CanvasObject, borderColor color.Color, borderWidth float32) *expandableItem {
	item := &expandableItem{
		header:      header,
		content:     content,
		expanded:    false,
		borderColor: borderColor,
		borderWidth: borderWidth,
	}

	// Header'a tap event'i bağla
	header.onTap = func() {
		item.expanded = !item.expanded
		item.Refresh()
	}

	item.ExtendBaseWidget(item)
	return item
}

// SetExpanded expand durumunu dışarıdan ayarlar
func (item *expandableItem) SetExpanded(expanded bool) {
	item.expanded = expanded
	item.Refresh()
}

// IsExpanded mevcut expand durumunu döner
func (item *expandableItem) IsExpanded() bool {
	return item.expanded
}

func (item *expandableItem) CreateRenderer() fyne.WidgetRenderer {
	container := container.NewVBox(item.header)

	if item.expanded {
		container.Add(item.content)
	}

	// Border için rectangle oluştur
	var border *canvas.Rectangle
	if item.borderWidth > 0 {
		border = canvas.NewRectangle(item.borderColor)
	}

	return &expandableItemRenderer{
		item:      item,
		container: container,
		border:    border,
	}
}

type expandableItemRenderer struct {
	item      *expandableItem
	container *fyne.Container
	border    *canvas.Rectangle
}

func (r *expandableItemRenderer) Layout(size fyne.Size) {
	if r.border != nil {
		// Border'ı tam genişliğe yay
		r.border.Resize(size)
		r.border.Move(fyne.NewPos(0, 0))

		// Container'ı border içine yerleştir
		inset := r.item.borderWidth
		r.container.Resize(fyne.NewSize(size.Width-2*inset, size.Height-2*inset))
		r.container.Move(fyne.NewPos(inset, inset))
	} else {
		r.container.Resize(size)
	}
}

func (r *expandableItemRenderer) MinSize() fyne.Size {
	minSize := r.container.MinSize()
	if r.border != nil {
		// Border kalınlığını MinSize'a ekle
		inset := r.item.borderWidth * 2
		return fyne.NewSize(minSize.Width+inset, minSize.Height+inset)
	}
	return minSize
}

func (r *expandableItemRenderer) Refresh() {
	r.container.Objects = []fyne.CanvasObject{r.item.header}

	if r.item.expanded {
		r.container.Objects = append(r.container.Objects, r.item.content)
	}

	r.container.Refresh()
}

func (r *expandableItemRenderer) Objects() []fyne.CanvasObject {
	if r.border != nil {
		return []fyne.CanvasObject{r.border, r.container}
	}
	return []fyne.CanvasObject{r.container}
}

func (r *expandableItemRenderer) Destroy() {}

// borderedContainer - İçeriği renkli çerçeve ile sarmalayan container
type borderedContainer struct {
	widget.BaseWidget
	content     fyne.CanvasObject
	borderColor color.Color
	borderWidth float32
}

func newBorderedContainer(content fyne.CanvasObject, borderColor color.Color, borderWidth float32) *borderedContainer {
	bc := &borderedContainer{
		content:     content,
		borderColor: borderColor,
		borderWidth: borderWidth,
	}
	bc.ExtendBaseWidget(bc)
	return bc
}

func (bc *borderedContainer) CreateRenderer() fyne.WidgetRenderer {
	border := canvas.NewRectangle(bc.borderColor)

	return &borderedContainerRenderer{
		container: bc,
		border:    border,
	}
}

type borderedContainerRenderer struct {
	container *borderedContainer
	border    *canvas.Rectangle
}

func (r *borderedContainerRenderer) Layout(size fyne.Size) {
	// Border sadece sol kenarda dikey çizgi olarak
	r.border.Resize(fyne.NewSize(r.container.borderWidth, size.Height))
	r.border.Move(fyne.NewPos(0, 0))

	// Content border'dan sonra başlasın
	contentX := r.container.borderWidth + 2 // Border genişliği + küçük boşluk
	r.container.content.Resize(fyne.NewSize(size.Width-contentX, size.Height))
	r.container.content.Move(fyne.NewPos(contentX, 0))
}

func (r *borderedContainerRenderer) MinSize() fyne.Size {
	contentMin := r.container.content.MinSize()
	// Sol border genişliği + küçük boşluk ekle
	extraWidth := r.container.borderWidth + 2
	return fyne.NewSize(contentMin.Width+extraWidth, contentMin.Height)
}

func (r *borderedContainerRenderer) Refresh() {
	r.border.FillColor = r.container.borderColor
	r.border.Refresh()
	canvas.Refresh(r.container.content)
}

func (r *borderedContainerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.border, r.container.content}
}

func (r *borderedContainerRenderer) Destroy() {}

// maxWidthContainer - İçeriğin maksimum genişliğini sınırlandırır
type maxWidthContainer struct {
	widget.BaseWidget
	content   fyne.CanvasObject
	maxWidth  float32
	container *fyne.Container
}

// newMaxWidthContainer creates a container that limits the width of its content
func newMaxWidthContainer(content fyne.CanvasObject, maxWidth float32) *maxWidthContainer {
	mwc := &maxWidthContainer{
		content:  content,
		maxWidth: maxWidth,
	}
	mwc.ExtendBaseWidget(mwc)
	return mwc
}

// CreateRenderer creates the renderer for this widget
func (mwc *maxWidthContainer) CreateRenderer() fyne.WidgetRenderer {
	return &maxWidthRenderer{
		container: mwc,
	}
}

type maxWidthRenderer struct {
	container *maxWidthContainer
}

func (mwr *maxWidthRenderer) Layout(space fyne.Size) {
	// Content'in genişliğini maxWidth kadar sınırla
	contentWidth := space.Width
	if contentWidth > mwr.container.maxWidth {
		contentWidth = mwr.container.maxWidth
	}

	mwr.container.content.Resize(fyne.NewSize(contentWidth, space.Height))
	mwr.container.content.Move(fyne.NewPos(0, 0))
}

func (mwr *maxWidthRenderer) MinSize() fyne.Size {
	// MinSize'ı maxWidth kadar sınırla
	minSize := mwr.container.content.MinSize()
	if minSize.Width > mwr.container.maxWidth {
		minSize.Width = mwr.container.maxWidth
	}
	return minSize
}

func (mwr *maxWidthRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{mwr.container.content}
}

func (mwr *maxWidthRenderer) Refresh() {
	mwr.container.content.Refresh()
}

func (mwr *maxWidthRenderer) Destroy() {}
