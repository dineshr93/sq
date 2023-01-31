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

// pkgsCmd represents the pkgs command
var pkgsCmd = &cobra.Command{
	Use:   "pkgs",
	Short: "Command to list pkgs section",
	Long: `Command to list pkgs section
For Ex: 
To list all packages: ./sq pkgs   
To list first 5     : ./sq pkgs 5`,
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
			s.Printpkgs(i)
		} else {
			if lenPkgs > 0 {
				s.Printpkgs(lenPkgs)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pkgsCmd)
}
