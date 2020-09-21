package config

type AppConfig struct {
	Config App `toml:"app"`
}

type App struct {
	AppId  string `toml:"appId"`
	LogDir string `toml:"logDir"`
}
