package batis

import (
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"os"
	"strings"
)

func init() {

	viper.SetConfigName("batis")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	env := "dev"
	if envVar := os.Getenv("env"); envVar != "" {
		env = strings.ToLower(envVar)
	}
	viper.SetConfigName(env)
	viper.MergeInConfig()

}


