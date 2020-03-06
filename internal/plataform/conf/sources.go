package conf

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type FileSource struct {
	data []byte
}

func (fs *FileSource) Unmarshal(target interface{}) error {
	if len(fs.data) == 0 {
		return errors.New("no file data present, you must have a config file")
	}

	return yaml.Unmarshal(fs.data, target)
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

type EnvSource struct{}

func (e EnvSource) GetString(key string) string {
	return os.Getenv(key)
}
