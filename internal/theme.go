package main

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Renk tanımlamaları
var (
	colorDarkBlue   = color.NRGBA{R: 30, G: 60, B: 100, A: 255}   // Ana koyu mavi (içerik arka planı)
	colorLightBlue  = color.NRGBA{R: 40, G: 80, B: 130, A: 255}   // Açık mavi
	colorOrange     = color.NRGBA{R: 255, G: 165, B: 0, A: 255}   // Turuncu
	colorDarkcyan   = color.NRGBA{R: 0, G: 139, B: 139, A: 255}   // Koyu camgöbeği
	colorYellow     = color.NRGBA{R: 255, G: 215, B: 0, A: 255}   // Sarı
	colorBrown      = color.NRGBA{R: 80, G: 60, B: 40, A: 255}    // Kahverengi (tab başlıkları)
	colorOlive      = color.NRGBA{R: 85, G: 107, B: 47, A: 255}   // Zeytin yeşili
	colorDarkGray   = color.NRGBA{R: 136, G: 148, B: 172, A: 255} // Genel arka plan
	colorMenuBg     = color.NRGBA{R: 45, G: 45, B: 45, A: 255}    // Menü arka planı
	colorMenuBorder = color.NRGBA{R: 70, G: 70, B: 70, A: 255}    // Menü çerçevesi
	colorMenuHover  = color.NRGBA{R: 70, G: 70, B: 70, A: 255}    // Menü hover rengi
	colorBackground = color.NRGBA{R: 20, G: 20, B: 30, A: 255}    // Genel arka plan
	colorListSelect = colorOrange
	//colorTabBG      = colorLightBlue
	colorButton = colorDarkGray

	// Ortam tipi çerçeve renkleri
	colorAppTypeProd  = color.NRGBA{R: 220, G: 38, B: 38, A: 255} // Kırmızı - PROD
	colorAppTypeStage = color.NRGBA{R: 255, G: 165, B: 0, A: 255} // Turuncu - PREP/UAT
	colorAppTypeOther = color.NRGBA{R: 34, G: 197, B: 94, A: 255} // Yeşil - Diğerleri (DEV, TEST vb.)

	// Badge renkleri
	colorBadgeGreen = color.NRGBA{R: 34, G: 197, B: 94, A: 255}  // Yeşil - VPN gibi durumlar için
	colorBadgeBlue  = color.NRGBA{R: 59, G: 130, B: 246, A: 255} // Mavi - EBS versiyonu, RDC/Host sayısı gibi bilgilendirme için

	// Separator ve border renkleri
	colorSeparator = color.NRGBA{R: 117, G: 140, B: 163, A: 255} // #758CA3 - Başlık altı çizgi ve çerçeve rengi
)

// getAppTypeBorderColor ortam tipine göre çerçeve rengini döndürür
func getAppTypeBorderColor(appType string) color.Color {
	// Büyük/küçük harf duyarsız karşılaştırma için
	appTypeLower := strings.ToLower(strings.TrimSpace(appType))

	switch appTypeLower {
	case "prod", "production":
		return colorAppTypeProd // Kırmızı
	case "prep", "uat":
		return colorAppTypeStage // Turuncu
	default:
		return colorAppTypeOther // Yeşil (dev, test, staging vb.)
	}
}

// blueTheme - Custom tema for tabs
type blueTheme struct {
	fyne.Theme
}

func (t *blueTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Tab için özel renkler
	switch name {
	case theme.ColorNameHeaderBackground:
		return colorListSelect // Tab başlık arka planı
	case theme.ColorNameButton:
		return colorButton // Tab butonları
	case theme.ColorNameInputBackground:
		return colorBackground // Input/Tab arka planı
	case theme.ColorNameMenuBackground:
		return colorMenuBg // Menu/Tab arka planı
	case theme.ColorNameOverlayBackground:
		return colorMenuBg // Overlay arka planı
	case theme.ColorNameBackground:
		return colorBackground // Genel arka plan
	case theme.ColorNameForeground:
		return color.White // Yazı rengi
	case theme.ColorNamePrimary:
		return colorDarkGray
	case theme.ColorNameSeparator:
		return colorSeparator // Separator rengi #758CA3
	case theme.ColorNameShadow:
		return colorSeparator // Card kenarlık rengi #758CA3
	case theme.ColorNameHover:
		return colorDarkcyan // Hover rengi - açık kahverengi
	case theme.ColorNamePressed:
		return colorListSelect // Pressed rengi - orta kahverengi
	case theme.ColorNameHyperlink:
		return colorYellow // Pressed rengi - orta kahverengi
	}
	// Diğer renkler için varsayılan dark theme'i kullan
	return theme.DefaultTheme().Color(name, variant)
}
