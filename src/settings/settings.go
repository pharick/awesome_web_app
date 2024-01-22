package settings

import (
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type DatabaseSettings struct {
	Host     string `koanf:"HOST" validate:"required"`
	Port     int    `koanf:"PORT" validate:"required"`
	User     string `koanf:"USER" validate:"required"`
	Password string `koanf:"PASSWORD" validate:"required"`
	Database string `koanf:"DATABASE" validate:"required"`
}

type GoogleOAuthSettings struct {
	ClientID     string `koanf:"CLIENT_ID" validate:"required"`
	ClientSecret string `koanf:"CLIENT_SECRET" validate:"required"`
}

type Settings struct {
	BaseUrl       string              `koanf:"BASE_URL" validate:"required"`
	Port          int                 `koanf:"PORT" validate:"required"`
	SessionSecret string              `koanf:"SESSION_SECRET" validate:"required"`
	DB            DatabaseSettings    `koanf:"DB" validate:"required"`
	Google        GoogleOAuthSettings `koanf:"GOOGLE" validate:"required"`
}

func LoadSettings() (settings *Settings, err error) {
	k := koanf.New(".")
	_ = k.Load(
		file.Provider(".env"),
		dotenv.ParserEnv("", "__", func(s string) string {
			return s
		}),
	)
	_ = k.Load(env.Provider("", "__", nil), nil)

	err = k.Unmarshal("", &settings)
	if err != nil {
		return
	}

	validate := validator.New()
	err = validate.Struct(settings)
	return
}
