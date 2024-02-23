package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HelloText   string `yaml:"hello_text"`
	SuccessText string `yaml:"success_text"`
	TimeBan     int    `yaml:"time_ban_ms"`
}

func NewConfig(fileName string) (*Config, error) {
	p := Config{}
	b, err := os.ReadFile(fileName)
	if err != nil {
		return &p, err
	}
	err = yaml.Unmarshal(b, &p)
	if err != nil {
		return &p, err
	}
	return &p, nil
}

func DefaultConfig() *Config {
	return &Config{
		HelloText:   "Добро пожаловать, @%s!\n\nПодтвердите, что вы не бот, у вас есть 2 попытки и 5 минут.\n\n",
		SuccessText: "Вы успешно прошли капчу, спасибо.",
		TimeBan:     300000,
	}
}
