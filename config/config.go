package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var defaultConf = []byte(`
# example: debug, release, test
CORE_MODE=debug
CORE_PORT=8080
`)

var Conf ConfENV
var once sync.Once

type ConfENV struct {
	Core  SectionCore
	Log   SectionLog
	DB    SectionDB
	Redis SectionRedis
}

type SectionCore struct {
	Mode string
	Port string
}

type SectionLog struct {
	Format string
	Output string
	Level  string
}

type SectionDB struct {
	Host         string
	Port         string
	Database     string
	Username     string
	Password     string
	MaxIdleConns int
	MaxOpenConns int
}

type SectionRedis struct {
	Hosts string
}

func InitConf(confPath string) error {
	var err error
	once.Do(func() {
		Conf, err = LoadConf(confPath)
	})
	return err
}

// LoadConf load config from file and read in environment variables that match
func LoadConf(confPath string) (ConfENV, error) {
	var conf ConfENV

	viper.SetConfigType("env")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if confPath != "" {
		content, err := os.ReadFile(confPath)

		if err != nil {
			return conf, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
	} else {
		// Search config in home directory with name ".gorush" (without extension).
		viper.AddConfigPath(".")
		// viper.SetConfigName("")
		viper.SetConfigFile(".env")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			// load default config
			if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
				return conf, err
			}
		}
	}

	conf.Core.Mode = viper.GetString("core_mode")
	conf.Core.Port = viper.GetString("core_port")
	if len(conf.Core.Port) == 0 {
		conf.Core.Port = "8080"
	}

	conf.Log.Format = viper.GetString("log_format")
	conf.Log.Level = viper.GetString("log_level")
	conf.Log.Output = viper.GetString("log_output")

	conf.DB.Host = viper.GetString("db_host")
	conf.DB.Port = viper.GetString("db_port")
	conf.DB.Username = viper.GetString("db_username")
	conf.DB.Password = viper.GetString("db_password")
	conf.DB.Database = viper.GetString("db_database")
	conf.DB.MaxIdleConns = viper.GetInt("db_max_idle_conns")
	conf.DB.MaxOpenConns = viper.GetInt("db_max_open_conns")

	conf.Redis.Hosts = viper.GetString("redis_hosts")

	return conf, nil
}
