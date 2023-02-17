/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dineshr93/sq/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// digCmd represents the dig command
var digCmd = &cobra.Command{
	Use:   "dig",
	Short: "dig the relationship and show actual packages and files",
	Long: `dig the relationship and show actual packages and files
	
For Ex: 
To list all relationship: ./sq rels dig`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("digging")
		dataFile := string(viper.ConfigFileUsed())
		s := &model.SPDX{}
		if err := s.Load(dataFile); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		rels := s.Relationships
		lenRels := len(rels)

		if len(args) > 0 {
			// fmt.Sprintln(len(args))
			// i, err := strconv.Atoi(args[0])
			// if err != nil {
			// 	panic(err)
			// }
			// s.PrintDigRels(i)
			s.PrintDigRels()
		} else {
			// fmt.Sprintln(lenRels)
			if lenRels > 0 {
				s.PrintDigRels()
			}
		}
	},
}

func init() {
	relsCmd.AddCommand(digCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// digCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// digCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
