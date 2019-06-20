package config

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	Dialect string `required:"true"`
	Host    string `required:"true"`
	Port    string `required:"true"`
	User    string `required:"true"`
	Pass    string `required:"true"`
	Name    string `required:"true"`
	SSLMode string `default:"disable"`
}

func (dbc DBConfig) GetDialect() string {
	return dbc.Dialect
}

func (dbc DBConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbc.Host, dbc.Port, dbc.User, dbc.Name, dbc.Pass, dbc.SSLMode)
}

func LoadDBConfig() (*DBConfig, error) {
	appConfig := &DBConfig{}
	err := envconfig.Process("et_db", appConfig)
	if err != nil {
		return nil, errors.New("unable to load config")
	}

	return appConfig, nil
}

func ConnectToDB() (*gorm.DB, error) {
	dbConfig, err := LoadDBConfig()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dbConfig.GetDialect(), dbConfig.GetConnectionString())
	if err != nil {
		return nil, err
	}

	return db, nil
}
