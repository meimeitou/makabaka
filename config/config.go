package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/meimeitou/makabaka/db"
	"github.com/meimeitou/makabaka/model"
	"github.com/sirupsen/logrus"
)

var DBSet = make(dBSet)

type dBSet map[string]*db.Conn

func InitDBSet(c *Config, logger *logrus.Logger) error {
	for _, item := range c.Proxy {
		conn, err := item.Storage.Config.Open(logger)
		if err != nil {
			return err
		}
		if c.AutoMigrate {
			if err := conn.AutoMigrate(&model.Apis{}); err != nil {
				logger.Error(err)
				logger.Warnf("db migrate error, skipping proxy %s!", item.Name)
				continue
			}
		}
		logger.Infof("init proxy %s success!", item.Name)
		DBSet[item.Name] = conn
	}
	return nil
}

func (d dBSet) GetDB(key string) (*db.Conn, error) {
	if db, ok := d[key]; ok {
		return db, nil
	}
	return nil, errors.New("proxy db key not exist!")
}

func ReadConfigFile(configFile string) (*Config, error) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", configFile, err)
	}
	var c Config
	if err := yaml.Unmarshal(configData, &c); err != nil {
		return nil, fmt.Errorf("error parse config file %s: %v", configFile, err)
	}
	Cfg = &c
	return &c, nil
}
