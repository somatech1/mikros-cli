package settings

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/huh"
	"github.com/creasty/defaults"

	"github.com/mikros-dev/mikros-cli/internal/path"
)

const (
	settingsFilename = "$HOME/.mikros/config.toml"
)

type Settings struct {
	Paths   Path               `toml:"paths"`
	UI      UI                 `toml:"ui"`
	App     Profile            `toml:"app"`
	Profile map[string]Profile `toml:"profile"`
}

type Path struct {
	Plugins Plugins `toml:"plugins"`
}

type Plugins struct {
	Services string `toml:"services" default:"$HOME/.mikros/plugins/services"`
	Features string `toml:"features" default:"$HOME/.mikros/plugins/features"`
}

type Profile struct {
	Project Project `toml:"project"`
}

type Project struct {
	ProtobufMonorepo ProtobufMonorepo `toml:"protobuf_monorepo"`
	Templates        Templates        `toml:"templates"`
}

type ProtobufMonorepo struct {
	RepositoryName string `toml:"repository_name" default:"protobuf-workspace"`
	ProjectName    string `toml:"project_name" default:"services"`
	VcsPath        string `toml:"vcs_path" default:"github.com/your-organization"`
}

type Templates struct {
	Protobuf ProtobufTemplates `toml:"protobuf"`
}

type ProtobufTemplates struct {
	CustomAuthName string `toml:"custom_auth_name" default:"scopes"`
}

type UI struct {
	Theme      string `toml:"theme"`
	Accessible bool   `toml:"accessible"`
}

func Load() (*Settings, error) {
	cfg, err := NewDefault()
	if err != nil {
		return nil, err
	}

	if name, ok := FileExists(); ok {
		if _, err := toml.DecodeFile(name, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func NewDefault() (*Settings, error) {
	cfg := &Settings{}
	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	cfg.Paths.Plugins.Services = os.ExpandEnv(cfg.Paths.Plugins.Services)
	cfg.Paths.Plugins.Features = os.ExpandEnv(cfg.Paths.Plugins.Features)

	return cfg, nil
}

func FileExists() (string, bool) {
	name := os.ExpandEnv(settingsFilename)
	return name, path.FindPath(name)
}

func (s *Settings) Write() error {
	var (
		basePath = os.ExpandEnv(settingsFilename)
	)

	if _, err := path.CreatePath(filepath.Dir(basePath)); err != nil {
		return err
	}

	file, err := os.Create(os.ExpandEnv(settingsFilename))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	en := toml.NewEncoder(file)
	if err := en.Encode(s); err != nil {
		return err
	}

	return nil
}

func (s *Settings) GetTheme() *huh.Theme {
	switch strings.ToLower(s.UI.Theme) {
	case "charm":
		return huh.ThemeCharm()
	case "dracula":
		return huh.ThemeDracula()
	case "catppuccin":
		return huh.ThemeCatppuccin()
	case "base16":
		return huh.ThemeBase16()
	}

	return huh.ThemeBase()
}

func (s *Settings) Hash() (string, error) {
	b, err := toml.Marshal(s)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(b)

	return hex.EncodeToString(h.Sum(nil)), nil
}
