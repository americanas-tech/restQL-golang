package conf

import (
	"github.com/b2wdigital/restQL-golang/internal/eval"
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

func (c Config) File() eval.FileSource {
	return c.fs
}

func (c Config) Env() eval.EnvSource {
	return c.env
}

func (c Config) Build() string {
	return c.build
}
