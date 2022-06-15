package config

import (
	"github.com/afonsojota/go-afonsojota-toolkit/goafonsojotapass"
	"github.com/afonsojota/go-afonsojota-toolkit/goutils/logger"
	"github.com/spf13/viper"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

var Settings = getConnection()

func getConnection() mysql.ConnectionURL {
	Init()
	if InLocal() {
		return mysql.ConnectionURL{
			Database: viper.GetString("database.schema"),
			Host:     viper.GetString("database.host"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.pass"),
		}
	}

	return mysql.ConnectionURL{
		Database: viper.GetString("database.schema"),
		Host:     goafonsojotapass.GetEnv(viper.GetString("database.host")),
		User:     viper.GetString("database.user"),
		Password: goafonsojotapass.GetEnv(viper.GetString("database.pass")),
	}
}

func GetDatabase() sqlbuilder.Database {
	logger.Info("Connecting to database")
	sess, err := mysql.Open(Settings)
	logger.Info("Database connected")
	if err != nil {
		logger.Error("Database Error", err)
	}
	return sess
}
