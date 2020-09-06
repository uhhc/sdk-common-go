package db

import (
	"errors"
	"fmt"

	// use mysql library
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	"github.com/uhhc/sdk-common-go/log"
)

// Config is the config for db connection
type Config struct {
	Engine   string
	User     string
	Password string
	DBName   string
	Host     string
	Port     string
	Charset  string
}

// DbClient is the struct of db client
type DbClient struct {
	*gorm.DB
	logger log.Logger
}

// NewDB to get a db instance
func NewDB(logger log.Logger, config *Config) (*DbClient, error) {
	var (
		db                                                  *gorm.DB
		err                                                 error
		engine, user, password, dbName, host, port, charset string
		timeout                                             uint32
	)

	if config == nil {
		engine = viper.GetString("DB_ENGINE")
		user = viper.GetString("DB_USER")
		password = viper.GetString("DB_PASSWORD")
		dbName = viper.GetString("DB_NAME")
		host = viper.GetString("DB_HOST")
		port = viper.GetString("DB_PORT")
		charset = viper.GetString("DB_CHARSET")
	} else {
		engine = config.Engine
		user = config.User
		password = config.Password
		dbName = config.DBName
		host = config.Host
		port = config.Port
		charset = config.Charset
	}
	timeout = viper.GetUint32("DB_CONN_TIMEOUT")
	if timeout == 0 {
		timeout = 10
	}

	if charset == "" {
		charset = "utf8mb4"
	}

	if engine == "mysql" {
		// See https://gorm.io/docs/connecting_to_the_database.html
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=%ds", user, password, host, port, dbName, charset, timeout)
		// fmt.Printf("dsn: %+v\n", dsn)
		db, err = gorm.Open(engine, dsn)
	} else {
		db = nil
		err = errors.New(engine + " is an unsupported database engine")
	}

	return &DbClient{
		DB:     db,
		logger: logger,
	}, err
}
