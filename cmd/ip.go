/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dineshr93/sq/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Display copyright and licenses",
	Long:  `Display copyright and licenses`,
	Run: func(cmd *cobra.Command, args []string) {
		dataFile := string(viper.ConfigFileUsed())
		s := &model.SPDX{}
		if err := s.Load(dataFile); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		pkgs := s.Packages
		lenPkgs := len(pkgs)

		if len(args) > 0 {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			s.PrintpkgsIP(i)
		} else {
			if lenPkgs > 0 {
				s.PrintpkgsIP(lenPkgs)
			}
		}
	},
}

func init() {
	pkgsCmd.AddCommand(ipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
