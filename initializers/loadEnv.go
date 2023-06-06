package initializers

import "github.com/spf13/viper"

type Config struct {
	JWT_SECRET         string `mapstructure:"JWT_SECRET"`
	REFRESH_JWT_SECRET string `mapstructure:"REFRESH_JWT_SECRET"`
	DATABASE_URL       string `mapstructure:"DATABASE_URL"`
	REDIS_URL          string `mapstructure:"REDIS_URL"`
	REDIS_PASS         string `mapstructure:"REDIS_PASS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
