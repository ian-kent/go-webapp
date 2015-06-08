package config

import (
	"reflect"
	"strings"

	"github.com/ian-kent/go-webapp/logger"
	"github.com/ian-kent/gofigure"
)

type config struct {
	gofigure      interface{} `order:"env,flag"`
	BindAddr      string      `env:"BIND_ADDR" flag:"bind-addr" flagDesc:"Bind address"`
	SessionSecret []byte      `env:"SESSION_SECRET" flag:"session-secret" flagDesc:"Session secret"`
	SessionName   string      `env:"SESSION_NAME" flag:"session-name" flagDesc:"Session name"`
	LogLevel      string      `env:"LOGLEVEL" flag:"log-level" flagDesc:"Log level"`
}

// Get configures the application and returns the configuration
func Get() (*config, error) {
	cfg := config{
		BindAddr:      ":3123",
		SessionSecret: []byte("0s63d96c9sd23idnkdf098fb62jkdkjfb7982gng"),
		SessionName:   "go-webapp",
		LogLevel:      "info",
	}

	err := gofigure.Gofigure(&cfg)
	if err != nil {
		return nil, err
	}

	cfg.print()

	return &cfg, nil
}

func (c *config) print() {
	logger.Println(nil, "Configuration:")

	s := reflect.ValueOf(c).Elem()
	t := s.Type()

	ml := 0
	for i := 0; i < s.NumField(); i++ {
		if !s.Field(i).CanSet() {
			continue
		}
		if l := len(t.Field(i).Name); l > ml {
			ml = l
		}
	}

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !s.Field(i).CanSet() {
			continue
		}
		logger.Printf(nil, "\t%s%s: %s\n", strings.Repeat(" ", ml-len(t.Field(i).Name)), t.Field(i).Name, f.Interface())
	}
}
