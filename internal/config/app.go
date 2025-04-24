package config

import "github.com/gin-gonic/gin"

type App struct {
	Env      string `env:"ENV" envDefault:"development" json:"env"`
	Timezone string `env:"TIMEZONE" envDefault:"Asia/Ho_Chi_Minh" json:"timezone"`
	LogPath  string `env:"LOG_PATH" envDefault:"logger" json:"logPath"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"" json:"logLevel"`
	Version  string `env:"VERSION" json:"version"`
}

func (app App) IsProduction() bool {
	return app.Env == "production"
}

func (app App) IsDevelopment() bool {
	return app.Env == "development"
}

func (app App) IsTest() bool {
	return app.Env == "test"
}

func (app App) GetMode() string {
	if app.IsProduction() {
		return gin.ReleaseMode
	}

	return gin.DebugMode
}
