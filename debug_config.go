package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".hypon-api.yaml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	fmt.Printf("Config file: %s\n", viper.ConfigFileUsed())
	fmt.Printf("user.name: '%s'\n", viper.GetString("user.name"))
	fmt.Printf("hypon.api_url: '%s'\n", viper.GetString("hypon.api_url"))
	fmt.Printf("All settings: %+v\n", viper.AllSettings())
}
