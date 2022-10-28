package client

type ClientList struct {
	EMoneyServices Server
}

type Server struct {
	TLS     bool   `yaml:"tls"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}
