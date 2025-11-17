package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Renk tanımlamaları
var (
	colorDarkBlue   = color.NRGBA{R: 30, G: 60, B: 100, A: 255} // Ana koyu mavi (içerik arka planı)
	colorLightBlue  = color.NRGBA{R: 40, G: 80, B: 130, A: 255} // Açık mavi
	colorOrange     = color.NRGBA{R: 255, G: 165, B: 0, A: 255} // Turuncu
	colorDarkcyan   = color.NRGBA{R: 0, G: 139, B: 139, A: 255} // Koyu camgöbeği
	colorYellow     = color.NRGBA{R: 255, G: 215, B: 0, A: 255} // Sarı
	colorBrown      = color.NRGBA{R: 80, G: 60, B: 40, A: 255}  // Kahverengi (tab başlıkları)
	colorOlive      = color.NRGBA{R: 85, G: 107, B: 47, A: 255} // Zeytin yeşili
	colorBackground = color.NRGBA{R: 20, G: 20, B: 30, A: 255}  // Genel arka plan
	colorMenuBg     = color.NRGBA{R: 45, G: 45, B: 45, A: 255}  // Menü arka planı
	colorMenuBorder = color.NRGBA{R: 70, G: 70, B: 70, A: 255}  // Menü çerçevesi
	colorMenuHover  = color.NRGBA{R: 70, G: 70, B: 70, A: 255}  // Menü hover rengi
)

// blueTheme - Custom tema for tabs
type blueTheme struct {
	fyne.Theme
}

func (t *blueTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Tab için özel renkler
	switch name {
	case theme.ColorNameHeaderBackground:
		return colorDarkcyan // Tab başlık arka planı
	case theme.ColorNameButton:
		return colorDarkcyan // Tab butonları
	case theme.ColorNameInputBackground:
		return colorDarkcyan // Input/Tab arka planı
	case theme.ColorNameMenuBackground:
		return colorDarkcyan // Menu/Tab arka planı
	case theme.ColorNameOverlayBackground:
		return colorDarkcyan // Overlay arka planı
	case theme.ColorNameBackground:
		return colorBackground // Genel arka plan
	case theme.ColorNameForeground:
		return color.White // Yazı rengi
	case theme.ColorNamePrimary:
		return colorYellow // Seçili tab text rengi - turuncu
	case theme.ColorNameHover:
		return colorDarkcyan // Hover rengi - açık kahverengi
	case theme.ColorNamePressed:
		return colorDarkcyan // Pressed rengi - orta kahverengi
	}
	// Diğer renkler için varsayılan dark theme'i kullan
	return theme.DefaultTheme().Color(name, variant)
}
