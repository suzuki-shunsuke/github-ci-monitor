package controller

import (
	"errors"
	"fmt"
	"os"

	"github.com/suzuki-shunsuke/go-findconfig/findconfig"
	"gopkg.in/yaml.v2"
)

type ExistFile func(string) bool

type Reader struct {
	ExistFile ExistFile
}

func (reader Reader) read(p string) (Params, error) {
	params := Params{}
	f, err := os.Open(p)
	if err != nil {
		return params, fmt.Errorf("open a config file "+p+": %w", err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(true)
	if err := decoder.Decode(&params); err != nil {
		return params, fmt.Errorf("parse config as YAML: %w", err)
	}
	return params, nil
}

var ErrNotFound error = errors.New("configuration file isn't found")

func (reader Reader) FindAndRead(cfgPath, wd string) (Params, string, error) {
	params := Params{}
	if cfgPath == "" {
		p := findconfig.Find(wd, reader.ExistFile, ".github-ci-monitor.yml", ".github-ci-monitor.yaml")
		if p == "" {
			return params, "", ErrNotFound
		}
		cfgPath = p
	}
	cfg, err := reader.read(cfgPath)
	return cfg, cfgPath, err
}
