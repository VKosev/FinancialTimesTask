/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ft_cli",
	Short: "ft cli evaluates arithmetic expressions",
	Long:  `ft cli evaluates arithmetic expressions. It uses ft api internally for the functionality.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.ReadInConfig()
}
