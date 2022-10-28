package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"

	"github.com/vins7/top-up-services/config/client"
	"github.com/vins7/top-up-services/config/db"
	"github.com/vins7/top-up-services/config/server"
)

type config struct {
	Server   server.ServerList
	Database db.DatabaseList
	Client   client.ClientList
}

var cfg config

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func init() {

	viper.AddConfigPath(basepath + "/server")
	viper.SetConfigType("yaml")
	viper.SetConfigName("server.yml")
	err := viper.MergeInConfig()
	if err != nil {
		panic(fmt.Errorf("Cannot load server client config: %v", err))
	}

	viper.AddConfigPath(basepath + "/db")
	viper.SetConfigType("yaml")
	viper.SetConfigName("mysql.yml")
	err = viper.MergeInConfig()
	if err != nil {
		panic(fmt.Errorf("Cannot read database config: %v", err))
	}

	viper.AddConfigPath(basepath + "/client")
	viper.SetConfigType("yaml")
	viper.SetConfigName("grpc.yml")
	err = viper.MergeInConfig()
	if err != nil {
		panic(fmt.Errorf("Cannot read database config: %v", err))
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}

	viper.Unmarshal(&cfg)

	fmt.Println("============================")
	fmt.Println(Stringify(cfg))
	fmt.Println("============================")
}

func GetConfig() *config {
	return &cfg
}

func Stringify(data interface{}) string {
	dataByte, _ := json.Marshal(data)
	return string(dataByte)
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(env) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}
