package conf

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const configFileName = "restql.yml"

var once sync.Once

type Config struct {
	fs  *FileSource
	env EnvSource
}

func New() Config {
	var fs FileSource
	once.Do(func() {
		readConfigFile(&fs)
	})

	return Config{fs: &fs, env: EnvSource{}}
}

func (c Config) File() *FileSource {
	return c.fs
}

type FileSource struct {
	data []byte
}

func (fs *FileSource) Unmarshal(target interface{}) error {
	if len(fs.data) == 0 {
		return errors.New("no file data present, you must have a config file")
	}

	return yaml.Unmarshal(fs.data, target)
}

func (c Config) Env() EnvSource {
	return c.env
}

type EnvSource struct{}

func (e EnvSource) GetString(key string) string {
	return os.Getenv(key)
}

func readConfigFile(fs *FileSource) {
	path := getConfigFilepath()
	if path == "" {
		log.Printf("[WARN] no config file present")
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("[ERROR] could not load file at %s", path)
		return
	}

	fs.data = data
}

func getConfigFilepath() string {
	fileAtRoot := filepath.Join(".", configFileName)
	if doesFileExist(fileAtRoot) {
		return fileAtRoot
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("[DEBUG] failed to find home directory: %v", err)
		return ""
	}

	fileAtHome := filepath.Join(homeDir, configFileName)
	if doesFileExist(fileAtHome) {
		return fileAtHome
	}

	return ""
}

func doesFileExist(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
