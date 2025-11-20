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
	DialogTitleSuccess       = "Success"
	DialogTitleInfo          = "Information"
	DialogTitleDeleteEnv     = "Delete Environment"
	DialogTitleSaveData      = "Save Customer Data"
	DialogTitleOpenData      = "Open Customer Data"
	DialogTitleDeleteConfirm = "Delete"

	// Dialog Messages
	DialogMsgFileLoaded               = "File loaded!"
	DialogMsgClientAdded              = "Customer added!"
	DialogMsgClientUpdated            = "Updated!"
	DialogMsgClientDeleted            = "Customer deleted!"
	DialogMsgDataExported             = "Customer data exported!"
	DialogMsgDataImported             = "Customer data imported!"
	DialogMsgNoClientsToExport        = "No customers to export!"
	DialogMsgAddClientFirst           = "Please add a customer first!"
	DialogMsgInvalidClientSelection   = "Invalid customer selection"
	DialogMsgClientNotFound           = "Customer not found"
	DialogMsgInvalidClientInfo        = "Invalid customer information"
	DialogMsgFileNotFound             = "client_info.json file not found. A new file will be created."
	DialogMsgJSONReadError            = "JSON read error"
	DialogMsgFileReadError            = "File read error"
	DialogMsgDeleteClientConfirm      = "Customer will be deleted, are you sure?"
	DialogMsgDeleteEnvironmentConfirm = "Environment will be deleted, are you sure?"
	DialogMsgExportForClient          = "%d customers exported for client!"
	// Tab Names
	TabNameCompany      = "Customer"
	TabNameVPN          = "VPN"
	TabNameSystem       = "System"
	TabNameEnvironments = "Apps Environments"

	// Form Labels
	FormLabelAppUsers = "App Users"

	// SSH
	DialogTitleSSH           = "SSH"
	DialogMsgSSHConfig       = "SSH configuration is missing. Server IP and username are required."
	DialogMsgSSHFailed       = "SSH failed to open: %v"
	DialogMsgSSHPasswordCopy = "SSH password copied to clipboard.\nPaste it in the terminal with Ctrl+V."
)
