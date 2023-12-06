package config

import "fmt"

type Config struct {
	ServerStartURL string
	ReturnURL      string
}

func NewConfig() *Config {
	return &Config{
		ServerStartURL: "",
		ReturnURL:      "",
	}
}

func (c *Config) SetServerStartURL(s string) {
	c.ServerStartURL = s
}

func (c *Config) SetReturnURL(s string) {
	c.ReturnURL = s
}

func (c *Config) Print() {
	fmt.Println(c.ServerStartURL)
	fmt.Println(c.ReturnURL)
}
