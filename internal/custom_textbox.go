package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CustomTextBox wraps a read-only text display with copy and optional password toggle, plus an inline editor.
type CustomTextBox struct {
	widget.BaseWidget

	text        string
	isPassword  bool
	isMultiLine bool
	readOnly    bool
	hidden      bool

	onSave   func(string)
	onWindow func() fyne.Window

	displayLabel  *widget.Label
	editEntry     *widget.Entry
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
		maxDisplayLen: 30,
	}

	ctb.displayLabel = widget.NewLabel(ctb.getDisplayText())
	if isMultiLine {
		ctb.displayLabel.Wrapping = fyne.TextWrapWord
	}

	ctb.ExtendBaseWidget(ctb)
	return ctb
}

func (ctb *CustomTextBox) getDisplayText() string {
	text := ctb.text
	if !ctb.isMultiLine {
		text = strings.ReplaceAll(text, "\n", " ")
	}

	if ctb.isPassword && ctb.hidden {
		runeCount := len([]rune(text))
		if runeCount == 0 {
			return ""
		}
		return strings.Repeat("•", runeCount)
	}

	if ctb.readOnly && !ctb.isMultiLine {
		runes := []rune(text)
		if len(runes) > ctb.maxDisplayLen {
			return string(runes[:ctb.maxDisplayLen]) + "…"
		}
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
			}
			shortcut := &fyne.ShortcutSelectAll{}
			ctb.editEntry.TypedShortcut(shortcut)
		}
		return
	}

	ctb.readOnly = false
	ctb.hidden = false
	if ctb.editEntry != nil {
		ctb.editEntry.SetText(ctb.text)
	}
	ctb.BaseWidget.Refresh()

	fyne.Do(func() {
		if ctb.editEntry != nil {
			if canvas := fyne.CurrentApp().Driver().CanvasForObject(ctb.editEntry); canvas != nil {
				canvas.Focus(ctb.editEntry)
				shortcut := &fyne.ShortcutSelectAll{}
				ctb.editEntry.TypedShortcut(shortcut)
			}
		}
	})
}

func (ctb *CustomTextBox) Tapped(_ *fyne.PointEvent) {}

type customTextBoxRenderer struct {
	textBox *CustomTextBox

	readOnlyContainer *fyne.Container
	editContainer     *fyne.Container

	copyButton *IconButton
	eyeButton  *IconButton
	copyIcon   fyne.Resource
	copiedIcon fyne.Resource
	eyeIcon    fyne.Resource
	hiddenIcon fyne.Resource
	buttonSize fyne.Size
}

func newCustomTextBoxRenderer(textBox *CustomTextBox) *customTextBoxRenderer {
	r := &customTextBoxRenderer{
		textBox:    textBox,
		copyIcon:   loadIconResource("copy", theme.ContentCopyIcon()),
		copiedIcon: loadIconResource("check", theme.ConfirmIcon()),
		eyeIcon:    loadIconResource("eye", theme.VisibilityIcon()),
		hiddenIcon: loadIconResource("hidden", theme.VisibilityOffIcon()),
		buttonSize: fyne.NewSize(18, 18),
	}

	r.buildReadOnlyUI()
	r.buildEditUI()
	r.updateVisibility()
	return r
}

func (r *customTextBoxRenderer) buildReadOnlyUI() {
	r.textBox.displayLabel.SetText(r.textBox.getDisplayText())

	r.copyButton = NewIconButtonSimple(
		r.copyIcon,
		"",
		r.buttonSize,
		"Panoya kopyala",
		func() {
			window := r.textBox.onWindow()
			if window == nil {
				return
			}

			window.Clipboard().SetContent(r.textBox.text)
			r.copyButton.icon = r.copiedIcon
			r.copyButton.Refresh()

			go func() {
				time.Sleep(2 * time.Second)
				fyne.Do(func() {
					r.copyButton.icon = r.copyIcon
					r.copyButton.Refresh()
				})
			}()
		},
	)

	buttonObjects := []fyne.CanvasObject{r.copyButton}

	if r.textBox.isPassword {
		r.eyeButton = NewIconButtonSimple(
			r.eyeIcon,
			"",
			r.buttonSize,
			"Göster/Gizle",
			func() {
				r.textBox.hidden = !r.textBox.hidden
				r.textBox.displayLabel.SetText(r.textBox.getDisplayText())
				r.textBox.displayLabel.Refresh()
				r.updateEyeIcon()
			},
		)
		buttonObjects = append([]fyne.CanvasObject{r.eyeButton}, buttonObjects...)
	}

	buttons := container.NewHBox(buttonObjects...)
	r.readOnlyContainer = container.NewBorder(nil, nil, nil, buttons, r.textBox.displayLabel)
}

func (r *customTextBoxRenderer) buildEditUI() {
	if r.textBox.isMultiLine {
		entry := widget.NewMultiLineEntry()
		entry.Wrapping = fyne.TextWrapWord
		r.textBox.editEntry = entry
	} else {
		entry := widget.NewEntry()
		entry.OnSubmitted = func(value string) {
			r.saveEdit()
		}
		r.textBox.editEntry = entry
	}

	r.textBox.editEntry.SetText(r.textBox.text)
	r.textBox.editEntry.Password = r.textBox.isPassword

	saveBtn := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		r.saveEdit()
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		r.cancelEdit()
	})
	cancelBtn.Importance = widget.LowImportance

	buttonBar := container.NewHBox(saveBtn, cancelBtn)
	r.editContainer = container.NewBorder(nil, nil, nil, buttonBar, r.textBox.editEntry)
}

func (r *customTextBoxRenderer) saveEdit() {
	if r.textBox.editEntry == nil {
		return
	}

	r.textBox.text = r.textBox.editEntry.Text
	r.textBox.readOnly = true
	r.textBox.hidden = r.textBox.isPassword
	r.textBox.displayLabel.SetText(r.textBox.getDisplayText())
	r.textBox.BaseWidget.Refresh()

	if r.textBox.onSave != nil {
		r.textBox.onSave(r.textBox.text)
	}
}

func (r *customTextBoxRenderer) cancelEdit() {
	r.textBox.readOnly = true
	r.textBox.hidden = r.textBox.isPassword
	if r.textBox.editEntry != nil {
		r.textBox.editEntry.SetText(r.textBox.text)
	}
	r.textBox.displayLabel.SetText(r.textBox.getDisplayText())
	r.textBox.BaseWidget.Refresh()
}

func (r *customTextBoxRenderer) updateEyeIcon() {
	if r.eyeButton == nil {
		return
	}
	if r.textBox.hidden {
		r.eyeButton.icon = r.eyeIcon
	} else {
		r.eyeButton.icon = r.hiddenIcon
	}
	r.eyeButton.Refresh()
}

func (r *customTextBoxRenderer) updateVisibility() {
	if r.textBox.readOnly {
		r.editContainer.Hide()
		r.readOnlyContainer.Show()
	} else {
		r.readOnlyContainer.Hide()
		r.editContainer.Show()
	}
}

func (r *customTextBoxRenderer) Layout(size fyne.Size) {
	r.readOnlyContainer.Resize(size)
	r.editContainer.Resize(size)
}

func (r *customTextBoxRenderer) MinSize() fyne.Size {
	ro := r.readOnlyContainer.MinSize()
	ed := r.editContainer.MinSize()

	width := ro.Width
	if ed.Width > width {
		width = ed.Width
	}

	height := ro.Height
	if ed.Height > height {
		height = ed.Height
	}

	return fyne.NewSize(width, height)
}

func (r *customTextBoxRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.readOnlyContainer, r.editContainer}
}

func (r *customTextBoxRenderer) Refresh() {
	r.textBox.displayLabel.SetText(r.textBox.getDisplayText())
	r.textBox.displayLabel.Refresh()
	r.updateEyeIcon()
	r.updateVisibility()
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
