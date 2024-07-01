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
	Core   SectionCore
	Log    SectionLog
	SQLite SectionSQLite
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

type SectionSQLite struct {
	Database string
	MaxConn  int
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

	conf.SQLite.Database = viper.GetString("sqlite_database")
	conf.SQLite.MaxConn = viper.GetInt("sqlite_db_max_conn")

	return conf, nil
}
