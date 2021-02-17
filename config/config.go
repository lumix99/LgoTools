package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Config struct {
	Name        string
	data        map[string]string
	sectionData map[string]map[string]string
	lock        sync.RWMutex
	autoLock    sync.WaitGroup
	modTime     int64
	opts        Options
}

type Options struct {
	FilePath      string
	AutoReload    bool
	CheckDuration time.Duration
}

const (
	MAIN_PATH       = "config"
	DFAUTORELOAD    = false
	DFCHECKDURATION = time.Minute
)

var conf_map map[string]*Config
var conf_map_lock sync.RWMutex
var main_conf_name string

func init() {
	conf_map = make(map[string]*Config)
	filename := "local.ini"

	env := os.Getenv("SYS_ENV")
	if len(env) > 0 {
		filename = env + ".ini"
	}
	path := filepath.Join(MAIN_PATH, filename)
	c, err := Create(path)
	if err != nil {
		fmt.Println("config init error")
		return
	}
	main_conf_name = c.Name
}

func Create(arg interface{}) (*Config, error) {

	opts := Options{
		AutoReload:    DFAUTORELOAD,
		CheckDuration: DFCHECKDURATION,
	}

	switch s := arg.(type) {
	case string:
		if len(s) == 0 {
			return nil, errors.New("Invalid Option on Config.Create()")
		}
		opts.FilePath = s
	case Options:
		if len(s.FilePath) == 0 {
			return nil, errors.New("Invalid Option on Config.Create()")
		}
		if s.CheckDuration == 0 {
			s.CheckDuration = DFCHECKDURATION
		}
		opts = s
	}

	conf := &Config{
		opts:        opts,
		data:        make(map[string]string),
		sectionData: make(map[string]map[string]string),
	}

	err := conf.new()

	if err != nil {
		return nil, err
	}

	if conf.opts.AutoReload {
		go conf.autoReload()
	}

	return conf, nil
}

func ConfigGet(filename string) *Config {
	if len(filename) == 0 {
		filename = main_conf_name
	}

	if c, ok := conf_map[filename]; ok {
		return c
	} else {
		return nil
	}
}

func Get(args ...string) string {
	c, ok := conf_map[main_conf_name]
	if !ok {
		return ""
	}

	switch len(args) {
	case 1:
		if val, ok := c.Get(args[0]); ok {
			return val
		}
	case 2:
		if val, ok := c.SectionGet(args[0], args[1]); ok {
			return val
		}
	}

	return ""
}

func SectionsGet(sec string) map[string]string {
	c := conf_map[main_conf_name]
	return c.SectionsGet(sec)
}
