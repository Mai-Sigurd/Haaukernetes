package guacamole

type Guacamole struct {
	Username  string
	Password  string
	BaseUrl   string // TODO vi skal bruge user, pass, server ip, guac port somehow
	AuthToken string
}

type CreateUserAttributes struct {
	Disabled          string `json:"disabled"`
	Expired           string `json:"expired"`
	AccessWindowStart string `json:"access-window-start"`
	AccessWindowEnd   string `json:"access-window-end"`
	ValidFrom         string `json:"valid-from"`
	ValidUntil        string `json:"valid-until"`
	TimeZone          string `json:"timezone"`
}

type UserInfo struct {
	Username   string               `json:"username"`
	Password   string               `json:"password"`
	Attributes CreateUserAttributes `json:"attributes"`
}

type RDPConnection struct {
	ParentIdentifier string        `json:"parentIdentifier"`
	Name             string        `json:"name"`
	Protocol         string        `json:"protocol"`
	Parameters       RDPParameters `json:"parameters"`
	Attributes       RDPAttributes `json:"attributes"`
}

type RDPParameters struct {
	ClipboardEncoding        string `json:"clipboard-encoding"`
	ColorDepth               string `json:"color-depth"`
	Console                  string `json:"console"`
	ConsoleAudio             string `json:"console-audio"`
	CreateDrivePath          bool   `json:"create-drive-path"`
	Cursor                   string `json:"cursor"`
	DestPort                 string `json:"dest-port"`
	DisableAudio             bool   `json:"disable-audio"`
	DisableAuth              bool   `json:"disable-auth"`
	DPI                      string `json:"dpi"`
	DrivePath                string `json:"drive-path"`
	EnableAudio              bool   `json:"enable-audio"`
	EnableAudioInput         bool   `json:"enable-audio-input"`
	EnableDesktopComposition bool   `json:"enable-desktop-composition"`
	EnableDrive              bool   `json:"enable-drive"`
	EnableFontSmoothing      bool   `json:"enable-font-smoothing"`
	EnableFullWindowDrag     bool   `json:"enable-full-window-drag"`
	EnableMenuAnimations     bool   `json:"enable-menu-animations"`
	EnablePrinting           bool   `json:"enable-printing"`
	EnableSFTP               bool   `json:"enable-sftp"`
	EnableTheming            bool   `json:"enable-theming"`
	EnableWallpaper          bool   `json:"enable-wallpaper"`
	GatewayPort              string `json:"gateway-port"`
	Height                   string `json:"height"`
	Hostname                 string `json:"hostname"`
	IgnoreCert               bool   `json:"ignore-cert"`
	Password                 string `json:"password"`
	PreConnectionID          string `json:"preconnection-id"`
	Port                     string `json:"port"`
	ReadOnly                 bool   `json:"read-only"`
	ResizeMethod             string `json:"resize-method"`
	Security                 string `json:"security"`
	ServerLayout             string `json:"server-layout"`
	SFTPPort                 string `json:"sftp-port"`
	SFTPServerAliveInterval  string `json:"sftp-server-alive-interval"`
	SwapRedBlue              bool   `json:"swap-red-blue"`
	Username                 string `json:"username"`
	Width                    string `json:"width"`
}

type RDPAttributes struct {
	FailOverOnly          bool   `json:"failover-only"`
	GuacdEncryption       string `json:"guacd-encryption"`
	GuacdPort             string `json:"guacd-port"`
	MaxConnections        string `json:"max-connections"`
	MaxConnectionsPerUser string `json:"max-connections-per-user"`
	Weight                string `json:"weight"`
}

type AddConnection struct {
	Operation string `json:"op"`
	Path      string `json:"path"`
	Value     string `json:"value"`
}
