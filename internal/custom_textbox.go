package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CustomTextBox - ReadOnly/Edit toggle, Copy button, Password support, Multiline support
// Layout: [text (copy)(eye) ...] - buttons right next to text, ellipsis at the end
type CustomTextBox struct {
	widget.BaseWidget
	text        string
	isPassword  bool
	isMultiLine bool
	readOnly    bool
	hidden      bool
	onSave      func(string)
	onWindow    func() fyne.Window

	displayLabel  *canvas.Text
	editEntry     *widget.Entry
	container     *fyne.Container
	isEditing     bool
	maxDisplayLen int
}

func NewCustomTextBox(text string, isPassword bool, isMultiLine bool, onSave func(string), getWindow func() fyne.Window) *CustomTextBox {
	ctb := &CustomTextBox{
		text:          text,
		isPassword:    isPassword,
		isMultiLine:   isMultiLine,
		readOnly:      true,
		hidden:        isPassword,
		onSave:        onSave,
		onWindow:      getWindow,
		isEditing:     false,
		maxDisplayLen: 30,
	}

	ctb.displayLabel = canvas.NewText(ctb.getDisplayText(), theme.ForegroundColor())
	ctb.displayLabel.TextSize = theme.TextSize()

	ctb.container = container.NewStack(ctb.displayLabel)

	ctb.ExtendBaseWidget(ctb)
	return ctb
}

func (ctb *CustomTextBox) getDisplayText() string {
	text := ctb.text

	if ctb.isPassword && ctb.hidden {
		text = strings.Repeat("•", len(text))
	}

	if ctb.readOnly && !ctb.isMultiLine && len(text) > ctb.maxDisplayLen {
		text = text[:ctb.maxDisplayLen] + "…"
	}

	return text
}

func (ctb *CustomTextBox) CreateRenderer() fyne.WidgetRenderer {
	return newCustomTextBoxRenderer(ctb)
}

func (ctb *CustomTextBox) DoubleTapped(_ *fyne.PointEvent) {
	if !ctb.readOnly {
		if ctb.editEntry != nil {
			if canvas := fyne.CurrentApp().Driver().CanvasForObject(ctb.editEntry); canvas != nil {
				canvas.Focus(ctb.editEntry)
				selectAll := &fyne.ShortcutSelectAll{}
				ctb.editEntry.TypedShortcut(selectAll)
			}
		}
		return
	}

	ctb.readOnly = false
	ctb.isEditing = true

	// BaseWidget.Refresh() çağır ki Fyne canvas'ı invalidate etsin
	ctb.BaseWidget.Refresh()

	// Canvas'ı manuel olarak refresh et
	if canvas := fyne.CurrentApp().Driver().CanvasForObject(ctb); canvas != nil {
		canvas.Refresh(ctb)
	}

	if ctb.editEntry != nil {
		if canvas := fyne.CurrentApp().Driver().CanvasForObject(ctb.editEntry); canvas != nil {
			canvas.Focus(ctb.editEntry)
			selectAll := &fyne.ShortcutSelectAll{}
			ctb.editEntry.TypedShortcut(selectAll)
		}
	}
}

func (ctb *CustomTextBox) Tapped(_ *fyne.PointEvent) {}

// readOnlyTextBoxContainer - özel layout container
type readOnlyTextBoxContainer struct {
	widget.BaseWidget
	textBox   *CustomTextBox
	textLabel *canvas.Text
	buttons   fyne.CanvasObject
}

func newReadOnlyTextBoxContainer(textBox *CustomTextBox, buttons fyne.CanvasObject) *readOnlyTextBoxContainer {
	c := &readOnlyTextBoxContainer{
		textBox:   textBox,
		textLabel: textBox.displayLabel,
		buttons:   buttons,
	}
	c.ExtendBaseWidget(c)
	return c
}

func (c *readOnlyTextBoxContainer) MinSize() fyne.Size {
	textSize := c.textLabel.MinSize()
	buttonsSize := c.buttons.MinSize()

	// Height'ı button ve text'in max'ı yap
	height := textSize.Height
	if buttonsSize.Height > height {
		height = buttonsSize.Height
	}

	return fyne.NewSize(textSize.Width+buttonsSize.Width+10, height+4)
}

func (c *readOnlyTextBoxContainer) CreateRenderer() fyne.WidgetRenderer {
	return &readOnlyTextBoxRenderer{
		container: c,
	}
}

type readOnlyTextBoxRenderer struct {
	container *readOnlyTextBoxContainer
}

func (r *readOnlyTextBoxRenderer) Layout(space fyne.Size) {
	buttonsSize := r.container.buttons.MinSize()
	availableWidth := space.Width - buttonsSize.Width - 5

	if availableWidth < 50 {
		availableWidth = 50
	}

	r.container.textLabel.Text = r.container.textBox.getDisplayText()
	r.container.textLabel.TextSize = theme.TextSize()

	textSize := r.container.textLabel.MinSize()
	// Bottom align - 5px yukarı (aşağıdan 5px az)
	textY := space.Height - textSize.Height + 5
	if textY < 0 {
		textY = 0
	}

	r.container.textLabel.Resize(fyne.NewSize(availableWidth, textSize.Height))
	r.container.textLabel.Move(fyne.NewPos(0, textY))

	// Buttons da bottom align - 5px yukarı
	buttonsY := space.Height - buttonsSize.Height + 5
	if buttonsY < 0 {
		buttonsY = 0
	}
	r.container.buttons.Resize(buttonsSize)
	r.container.buttons.Move(fyne.NewPos(availableWidth+5, buttonsY))
}

func (r *readOnlyTextBoxRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *readOnlyTextBoxRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container.textLabel, r.container.buttons}
}

func (r *readOnlyTextBoxRenderer) Refresh() {
	r.container.textLabel.Text = r.container.textBox.getDisplayText()
	r.container.textLabel.Refresh()
	r.container.buttons.Refresh()
}

func (r *readOnlyTextBoxRenderer) Destroy() {}

type customTextBoxRenderer struct {
	textBox     *CustomTextBox
	readOnlyUI  *readOnlyTextBoxContainer
	editUI      *fyne.Container
	currentUI   fyne.CanvasObject
	initialized bool
}

func newCustomTextBoxRenderer(textBox *CustomTextBox) *customTextBoxRenderer {
	r := &customTextBoxRenderer{
		textBox:     textBox,
		initialized: false,
	}
	r.initializeUIs()
	return r
}

func (r *customTextBoxRenderer) initializeUIs() {
	if r.initialized {
		return
	}

	// ReadOnly UI
	var copyBtn *IconButton
	copyBtn = NewIconButton(
		theme.ContentCopyIcon(),
		"",
		fyne.NewSize(12, 12),
		"Kopyala",
		func() {
			window := r.textBox.onWindow()
			if window != nil {
				window.Clipboard().SetContent(r.textBox.text)

				// Button'ı "Copied" text'i ile göster
				originalText := copyBtn.text
				copyBtn.text = "✓"
				copyBtn.icon = theme.ConfirmIcon()

				// Widget'i refresh et
				fyne.Do(func() {
					copyBtn.BaseWidget.Refresh()
				})

				// 2 saniye sonra geri dön
				go func() {
					time.Sleep(2 * time.Second)
					fyne.Do(func() {
						copyBtn.text = originalText
						copyBtn.icon = theme.ContentCopyIcon()
						copyBtn.BaseWidget.Refresh()
					})
				}()
			}
		},
		nil, nil,
	)

	var eyeBtn *IconButton
	if r.textBox.isPassword {
		eyeBtn = NewIconButton(
			theme.VisibilityIcon(),
			"",
			fyne.NewSize(12, 12),
			"Göster/Gizle",
			func() {
				r.textBox.hidden = !r.textBox.hidden
				r.textBox.displayLabel.Text = r.textBox.getDisplayText()
				r.textBox.displayLabel.Refresh()
			},
			nil, nil,
		)
	}

	var buttonsBox fyne.CanvasObject
	if r.textBox.isPassword {
		buttonsBox = container.NewHBox(eyeBtn, copyBtn)
	} else {
		buttonsBox = copyBtn
	}

	r.readOnlyUI = newReadOnlyTextBoxContainer(r.textBox, buttonsBox)

	// Edit UI
	if r.textBox.isMultiLine {
		r.textBox.editEntry = widget.NewMultiLineEntry()
		r.textBox.editEntry.Wrapping = fyne.TextWrapWord
	} else {
		r.textBox.editEntry = widget.NewEntry()
	}
	r.textBox.editEntry.SetText(r.textBox.text)

	saveBtn := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		r.textBox.text = r.textBox.editEntry.Text
		r.textBox.readOnly = true
		r.textBox.displayLabel.Text = r.textBox.getDisplayText()
		r.textBox.BaseWidget.Refresh()
		if canvas := fyne.CurrentApp().Driver().CanvasForObject(r.textBox); canvas != nil {
			canvas.Refresh(r.textBox)
		}
		if r.textBox.onSave != nil {
			r.textBox.onSave(r.textBox.text)
		}
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		r.textBox.readOnly = true
		r.textBox.BaseWidget.Refresh()
		if canvas := fyne.CurrentApp().Driver().CanvasForObject(r.textBox); canvas != nil {
			canvas.Refresh(r.textBox)
		}
	})
	cancelBtn.Importance = widget.LowImportance

	buttons := container.NewHBox(saveBtn, cancelBtn)

	r.editUI = container.NewBorder(
		nil, nil, nil, buttons,
		r.textBox.editEntry,
	)

	r.initialized = true
	r.updateCurrentUI()
}

func (r *customTextBoxRenderer) updateCurrentUI() {
	if r.textBox.readOnly {
		r.currentUI = r.readOnlyUI
	} else {
		r.currentUI = r.editUI
	}
}

func (r *customTextBoxRenderer) Layout(space fyne.Size) {
	if r.currentUI != nil {
		r.currentUI.Move(fyne.NewPos(0, 0))
		r.currentUI.Resize(space)

		// Eğer readOnly ise, renderer'ını layout et
		if r.textBox.readOnly && r.readOnlyUI != nil {
			renderer := r.readOnlyUI.CreateRenderer()
			if renderer != nil {
				renderer.Layout(space)
			}
		}
	}
}

func (r *customTextBoxRenderer) MinSize() fyne.Size {
	if r.currentUI != nil {
		return r.currentUI.MinSize()
	}
	return fyne.NewSize(100, 30)
}

func (r *customTextBoxRenderer) Objects() []fyne.CanvasObject {
	if r.currentUI != nil {
		// Eğer readOnly ise, readOnlyUI'nin objects'lerini döndür
		if r.textBox.readOnly && r.readOnlyUI != nil {
			renderer := r.readOnlyUI.CreateRenderer()
			if renderer != nil {
				return renderer.Objects()
			}
		}
		// Eğer edit mode ise, editUI'nin objects'lerini döndür
		if !r.textBox.readOnly && r.editUI != nil {
			return r.editUI.Objects
		}
	}
	return []fyne.CanvasObject{}
}

func (r *customTextBoxRenderer) Refresh() {
	// Edit mode'a geçilirse, editEntry'yi create et
	if !r.textBox.readOnly && r.textBox.editEntry == nil {
		if r.textBox.isMultiLine {
			r.textBox.editEntry = widget.NewMultiLineEntry()
			r.textBox.editEntry.Wrapping = fyne.TextWrapWord
		} else {
			r.textBox.editEntry = widget.NewEntry()
		}
		r.textBox.editEntry.SetText(r.textBox.text)
	}

	r.updateCurrentUI()

	if r.textBox.readOnly {
		r.textBox.displayLabel.Text = r.textBox.getDisplayText()
	} else if r.textBox.editEntry != nil {
		r.textBox.editEntry.SetText(r.textBox.text)
	}
}

func (r *customTextBoxRenderer) Destroy() {}

func (s *AppState) createCustomTextBoxItem(label string, text string, isPassword bool, isMultiLine bool, clientIndex int, updateFunc func(*Client, string)) *widget.FormItem {
	if isEncrypted(text) {
		decrypted, err := decryptString(text)
		if err == nil {
			text = decrypted
		}
	}

	textBox := NewCustomTextBox(text, isPassword, isMultiLine, func(newText string) {
		if clientIndex >= 0 && clientIndex < len(s.clients) {
			updateFunc(&s.clients[clientIndex], newText)
			if err := s.saveClients(); err != nil {
				dialog.ShowError(err, s.window)
			}
		}
	}, func() fyne.Window {
		return s.window
	})

	return widget.NewFormItem(label, textBox)
}
