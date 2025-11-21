package main

// VPNInfo holds VPN connection details
type VPNInfo struct {
	App           string `json:"app"`
	Host          string `json:"host"`
	User          string `json:"user"`
	Password      string `json:"password"`
	TwoFATokenApp string `json:"two_fa_token_app"`
	Notes         string `json:"not"`
}

// ClientData holds system-specific information
type ClientData struct {
	JiraURI       string   `json:"jira_uri"`
	JiraUser      string   `json:"jira_user"`
	JiraPassword  string   `json:"jira_password"`
	User          string   `json:"user"`
	PasswordReset string   `json:"pass_reset"`
	RDC           []string `json:"rdc"`
	Hosts         []string `json:"hosts"`
	Notes         string   `json:"not"`
}

// AppInfo holds application environment details
type AppInfo struct {
	Type          string   `json:"type"`
	Name          string   `json:"name"`
	User          string   `json:"user"`
	Password      string   `json:"pass"`
	DBServerIP    string   `json:"db_server_ip"`
	TNS           string   `json:"tns"`
	AppServerIP   string   `json:"app_server_ip"`
	AppServerURI  string   `json:"app_server_uri"`
	AppServerUser string   `json:"app_server_user"`
	AppServerPass string   `json:"app_server_pass"`
	WeblogicPass  string   `json:"weblogic_pass"`
	AppURI        string   `json:"app_uri"`
	AppUsers      []string `json:"app_users"`
	SSHParams     string   `json:"ssh_params"`
	Notes         string   `json:"not"`
}

// Client represents a single client with all their information
type Client struct {
	Company    string     `json:"company"`
	EBSVersion string     `json:"ebs_version"`
	VPN        VPNInfo    `json:"vpn"`
	Data       ClientData `json:"data"`
	Apps       []AppInfo  `json:"apps"`
	Notes      string     `json:"not"`
}
