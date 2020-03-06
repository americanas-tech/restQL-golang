package conf

import (
	"sync"
)

const configFileName = "restql.yml"

var once sync.Once

type Config struct {
	fs    *FileSource
	env   EnvSource
	build string
}

func New(build string) Config {
	var fs FileSource
	once.Do(func() {
		readConfigFile(&fs)
	})

	return Config{fs: &fs, env: EnvSource{}, build: build}
}

func (c Config) File() *FileSource {
	return c.fs
}

func (c Config) Env() EnvSource {
	return c.env
}

func (c Config) Build() string {
	return c.build
}
