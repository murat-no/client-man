package main

const (
	// Application
	DefaultJSONFile = "client_info.json"
	AppID           = "com.clientinfo.manager"
	AppName         = "Client Info Manager"

	// Window
	DefaultWindowWidth  = 800
	DefaultWindowHeight = 600

	// Colors
	DarkBlueRGB = 0x1a2a3a

	// UI
	DefaultTabMinHeight = 300
	TabContentMaxWidth  = 400 // Tab içeriği maksimum genişlik

	// Dialog Titles
	DialogTitleSuccess       = "Başarılı"
	DialogTitleInfo          = "Bilgi"
	DialogTitleDeleteEnv     = "Ortamı Sil"
	DialogTitleSaveData      = "Müşteri Verisi Kaydet"
	DialogTitleOpenData      = "Müşteri Verisi Aç"
	DialogTitleDeleteConfirm = "Sil"

	// Dialog Messages
	DialogMsgFileLoaded               = "Dosya yüklendi!"
	DialogMsgClientAdded              = "Firma eklendi!"
	DialogMsgClientUpdated            = "Güncellendi!"
	DialogMsgClientDeleted            = "Firma silindi!"
	DialogMsgDataExported             = "Müşteri verisi dışa aktarıldı!"
	DialogMsgDataImported             = "Müşteri verisi içe aktarıldı!"
	DialogMsgNoClientsToExport        = "Dışa aktarılacak firma yok!"
	DialogMsgAddClientFirst           = "Önce bir firma eklemelisiniz!"
	DialogMsgInvalidClientSelection   = "geçersiz firma seçimi"
	DialogMsgClientNotFound           = "firma bulunamadı"
	DialogMsgInvalidClientInfo        = "geçersiz firma bilgisi"
	DialogMsgFileNotFound             = "client_info.json dosyası bulunamadı. Yeni dosya oluşturulacak."
	DialogMsgJSONReadError            = "JSON okuma hatası"
	DialogMsgFileReadError            = "dosya okuma hatası"
	DialogMsgDeleteClientConfirm      = "Firma silinecek, emin misiniz?"
	DialogMsgDeleteEnvironmentConfirm = "Ortam silinecek, emin misiniz?"
	DialogMsgExportForClient          = "%d firma müşteri için dışa aktarıldı!"

	// Tab Names
	TabNameCompany      = "Firma"
	TabNameVPN          = "VPN"
	TabNameSystem       = "Sistem"
	TabNameEnvironments = "Ortamlar"

	// Form Labels
	FormLabelAppUsers = "App Users"

	// SSH
	DialogTitleSSH           = "SSH"
	DialogMsgSSHConfig       = "SSH yapılandırması eksik. Sunucu IP ve kullanıcı adı gereklidir."
	DialogMsgSSHFailed       = "SSH açılamadı: %v"
	DialogMsgSSHPasswordCopy = "SSH şifresi panoya kopyalandı.\nTerminal'de Ctrl+V ile yapıştırın."
)
