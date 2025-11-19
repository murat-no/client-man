package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CustomComboBox wraps a read-only text display with inline select editor.
type CustomComboBox struct {
	widget.BaseWidget

	text     string
	options  []string
	readOnly bool

	onSave   func(string)
	onWindow func() fyne.Window

	displayLabel *widget.Label
	editSelect   *widget.Select
}

func NewCustomComboBox(text string, options []string, onSave func(string), getWindow func() fyne.Window) *CustomComboBox {
	ccb := &CustomComboBox{
		text:     text,
		options:  options,
		readOnly: true,
		onSave:   onSave,
		onWindow: getWindow,
	}

	ccb.displayLabel = widget.NewLabel(text)

	ccb.ExtendBaseWidget(ccb)
	return ccb
}

func (ccb *CustomComboBox) CreateRenderer() fyne.WidgetRenderer {
	return newCustomComboBoxRenderer(ccb)
}

func (ccb *CustomComboBox) DoubleTapped(_ *fyne.PointEvent) {
	if !ccb.readOnly {
		if ccb.editSelect != nil {
			if canvas := fyne.CurrentApp().Driver().CanvasForObject(ccb.editSelect); canvas != nil {
				canvas.Focus(ccb.editSelect)
			}
		}
		return
	}

	ccb.readOnly = false
	if ccb.editSelect != nil {
		ccb.editSelect.SetSelected(ccb.text)
	}
	ccb.BaseWidget.Refresh()

	fyne.Do(func() {
		if ccb.editSelect != nil {
			if canvas := fyne.CurrentApp().Driver().CanvasForObject(ccb.editSelect); canvas != nil {
				canvas.Focus(ccb.editSelect)
			}
		}
	})
}

func (ccb *CustomComboBox) Tapped(_ *fyne.PointEvent) {}

type customComboBoxRenderer struct {
	comboBox *CustomComboBox

	readOnlyContainer *fyne.Container
	editContainer     *fyne.Container
}

func newCustomComboBoxRenderer(comboBox *CustomComboBox) *customComboBoxRenderer {
	r := &customComboBoxRenderer{
		comboBox: comboBox,
	}

	r.buildReadOnlyUI()
	r.buildEditUI()
	r.updateVisibility()
	return r
}

func (r *customComboBoxRenderer) buildReadOnlyUI() {
	r.comboBox.displayLabel.SetText(r.comboBox.text)
	r.readOnlyContainer = container.NewHBox(r.comboBox.displayLabel)
}

func (r *customComboBoxRenderer) buildEditUI() {
	select_ := widget.NewSelect(r.comboBox.options, func(value string) {
		// OnChanged callback - sadece değeri sakla
		r.comboBox.text = value
	})

	r.comboBox.editSelect = select_
	r.comboBox.editSelect.SetSelected(r.comboBox.text)

	saveBtn := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		r.saveEdit()
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		r.cancelEdit()
	})
	cancelBtn.Importance = widget.LowImportance

	buttonBar := container.NewHBox(saveBtn, cancelBtn)
	r.editContainer = container.NewBorder(nil, nil, nil, buttonBar, r.comboBox.editSelect)
}

func (r *customComboBoxRenderer) saveEdit() {
	if r.comboBox.editSelect == nil {
		return
	}

	r.comboBox.text = r.comboBox.editSelect.Selected
	r.comboBox.readOnly = true
	r.comboBox.displayLabel.SetText(r.comboBox.text)
	r.comboBox.BaseWidget.Refresh()

	if r.comboBox.onSave != nil {
		r.comboBox.onSave(r.comboBox.text)
	}
}

func (r *customComboBoxRenderer) cancelEdit() {
	r.comboBox.readOnly = true
	if r.comboBox.editSelect != nil {
		r.comboBox.editSelect.SetSelected(r.comboBox.text)
	}
	r.comboBox.displayLabel.SetText(r.comboBox.text)
	r.comboBox.BaseWidget.Refresh()
}

func (r *customComboBoxRenderer) updateVisibility() {
	if r.comboBox.readOnly {
		r.editContainer.Hide()
		r.readOnlyContainer.Show()
	} else {
		r.readOnlyContainer.Hide()
		r.editContainer.Show()
	}
}

func (r *customComboBoxRenderer) Layout(size fyne.Size) {
	r.readOnlyContainer.Resize(size)
	r.editContainer.Resize(size)
}

func (r *customComboBoxRenderer) MinSize() fyne.Size {
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

func (r *customComboBoxRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.readOnlyContainer, r.editContainer}
}

func (r *customComboBoxRenderer) Refresh() {
	r.comboBox.displayLabel.SetText(r.comboBox.text)
	r.comboBox.displayLabel.Refresh()
	r.updateVisibility()
}

func (r *customComboBoxRenderer) Destroy() {}

func (s *AppState) createCustomComboBoxItem(label string, text string, options []string, clientIndex int, updateFunc func(*Client, string)) *widget.FormItem {
	comboBox := NewCustomComboBox(text, options, func(newText string) {
		if clientIndex >= 0 && clientIndex < len(s.clients) {
			updateFunc(&s.clients[clientIndex], newText)
			if err := s.saveClients(); err != nil {
				dialog.ShowError(err, s.window)
			}
		}
	}, func() fyne.Window {
		return s.window
	})

	return widget.NewFormItem(label, comboBox)
}
