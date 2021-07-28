package subscribe

const (
	SSRServer = iota
	SSRPort
	SSRProtocol
	SSRCipher
	SSROBFS
	SSRSuffix
)

type Vmess struct {
	Add  string      `json:"add"`
	Aid  interface{} `json:"aid"`
	Host string      `json:"host"`
	ID   string      `json:"id"`
	Net  string      `json:"net"`
	Path string      `json:"path"`
	Port interface{} `json:"port"`
	PS   string      `json:"ps"`
	TLS  string      `json:"tls"`
	Type string      `json:"type"`
	V    interface{} `json:"v"`
}

type ClashVmess struct {
	Name           string            `json:"name,omitempty"`
	Type           string            `json:"type,omitempty"`
	Server         string            `json:"server,omitempty"`
	Port           interface{}       `json:"port,omitempty"`
	UUID           string            `json:"uuid,omitempty"`
	AlterID        interface{}       `json:"alterId,omitempty"`
	Cipher         string            `json:"cipher,omitempty"`
	TLS            bool              `json:"tls,omitempty"`
	Network        string            `json:"network,omitempty"`
	WSPATH         string            `json:"ws-path,omitempty"`
	WSHeaders      map[string]string `json:"ws-headers,omitempty"`
	SkipCertVerify bool              `json:"skip-cert-verify,omitempty"`
}

type ClashRSSR struct {
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	Server        string      `json:"server"`
	Port          interface{} `json:"port"`
	Password      string      `json:"password"`
	Cipher        string      `json:"cipher"`
	Protocol      string      `json:"protocol"`
	ProtocolParam string      `json:"protocol-param"`
	OBFS          string      `json:"obfs"`
	OBFSParam     string      `json:"obfs-param"`
}

type ClashSS struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Server     string      `json:"server"`
	Port       interface{} `json:"port"`
	Password   string      `json:"password"`
	Cipher     string      `json:"cipher"`
	Plugin     string      `json:"plugin"`
	PluginOpts PluginOpts  `json:"plugin-opts"`
}

type PluginOpts struct {
	Mode string `json:"mode"`
	Host string `json:"host"`
}

type SSD struct {
	Airport      string  `json:"airport"`
	Port         int     `json:"port"`
	Encryption   string  `json:"encryption"`
	Password     string  `json:"password"`
	TrafficUsed  float64 `json:"traffic_used"`
	TrafficTotal float64 `json:"traffic_total"`
	Expiry       string  `json:"expiry"`
	URL          string  `json:"url"`
	Servers      []struct {
		ID            int     `json:"id"`
		Server        string  `json:"server"`
		Ratio         float64 `json:"ratio"`
		Remarks       string  `json:"remarks"`
		Port          string  `json:"port"`
		Encryption    string  `json:"encryption"`
		Password      string  `json:"password"`
		Plugin        string  `json:"plugin"`
		PluginOptions string  `json:"plugin_options"`
	} `json:"servers"`
}

type Clash struct {
	Port      int `yaml:"port"`
	SocksPort int `yaml:"socks-port"`
	// RedirPort          int                      `yaml:"redir-port"`
	// Authentication     []string                 `yaml:"authentication"`
	AllowLan           bool   `yaml:"allow-lan"`
	Mode               string `yaml:"mode"`
	LogLevel           string `yaml:"log-level"`
	ExternalController string `yaml:"external-controller"`
	// ExternalUI         string                   `yaml:"external-ui"`
	// Secret             string                   `yaml:"secret"`
	// Experimental       map[string]interface{} 	`yaml:"experimental"`
	Proxy             []map[string]interface{} `yaml:"proxies"`
	ProxyGroup        []map[string]interface{} `yaml:"proxy-groups"`
	Rule              []string                 `yaml:"rules"`
	CFWByPass         []string                 `yaml:"cfw-bypass"`
	CFWLatencyTimeout int                      `yaml:"cfw-latency-timeout"`
}

type Trojan struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Server   string      `json:"server"`
	Password string      `json:"password"`
	Sni      string      `json:"sni,omitempty"`
	Port     interface{} `json:"port"`
}
