package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type MongoConfig struct {
	serverUrl    string
	databaseName string
	userName     string
	password     string
	port         string
	authDatabase string
	Timeout      uint64
}

func (config *MongoConfig) Init() {
	config.serverUrl = os.Getenv("MONGO_HOST")
	config.userName = os.Getenv("MONGO_USER")
	config.password = os.Getenv("MONGO_PWD")
	config.databaseName = os.Getenv("MONGO_DATABASE")
	config.authDatabase = os.Getenv("MONGO_AUTH_DATABASE")
	config.Timeout, _ = strconv.ParseUint(os.Getenv("MONGO_CONNECTION_TIMEOUT"), 10, 4)
}

func (config *MongoConfig) BuildUrl() (string, error) {
	if strings.Compare(config.userName, "") == 0 && strings.Compare(config.password, "") == 0 {

		if strings.Compare(config.serverUrl, "") == 0 {
			return "", errors.New("cowingram: please configure mongo source")
		}

		if strings.Compare(config.port, "") == 0 {
			config.port = "27017"
		}

		return "mongodb://" + config.serverUrl + ":" + config.port, nil
	}

	if strings.Compare(config.authDatabase, "") == 0 {
		config.authDatabase = "admin"
	}

	return "mongodb://" + config.userName + "@" + config.password + "@" +
		config.serverUrl + ":" + config.port + "/authSource=" + config.authDatabase, nil
}
