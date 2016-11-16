package haproxyconfigparser

type Socket struct {
	Type string
	Addr string
}

type Global struct {
	Stats   []*Socket `json:"stats"`
	Daemon  bool      `json:"daemon"`
	User    string    `json:"user"`
	Group   string    `json:"group"`
	Maxconn int       `json:"maxconn"`
}

type Acl struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Condition []string `json:"condition"`
}

type UseBackendClauses struct {
	ReverseJudge bool       `json:"reserve_judge"`
	Any          [][]string `json:"any"`
}

type UseBackend struct {
	Name      string             `json:"name"`
	Condition *UseBackendClauses `json:"cluses"`
	Backend   *Backend           `json:"backend"`
	Acls      []*Acl             `json:"acls"`
}

type Bind struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Frontend struct {
	Name        string        `json:"name"`
	Mode        string        `json:"mode"`
	Maxconn     int           `json:"maxconn"`
	Bind        Bind          `json:"bind"`
	Acls        []*Acl        `json:"acls"`
	UseBackends []*UseBackend `json:"use_backends"`
}

type Server struct {
	Label   string   `json:"label"`
	Host    string   `json:"host"`
	Port    int      `json:"port"`
	Options []string `json:"options"`
	Enabled bool     `json:"enabled"`
}

type Backend struct {
	Name    string     `json:"name"`
	Options [][]string `json:"options"`
	Servers []*Server  `json:"servers"`
}

type Services struct {
	Global    Global     `json:"global"`
	Frontends []Frontend `json:"frontends"`
	Backends  []Backend  `json:"backends"`
}

func NewServices() *Services {
	return &Services{}
}
