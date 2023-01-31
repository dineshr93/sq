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

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Command to list files section",
	Long: `Command to list on files section
	
For Ex: 
To list all Files   : ./sq files   
To list first 5     : ./sq files 5`,
	Run: func(cmd *cobra.Command, args []string) {

		dataFile := string(viper.ConfigFileUsed())
		s := &model.SPDX{}
		if err := s.Load(dataFile); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		files := s.Files
		lenFiles := len(files)

		if len(args) > 0 {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			s.PrintFiles(i)
		} else {
			if lenFiles > 0 {
				s.PrintFiles(lenFiles)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(filesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// filesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// filesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
