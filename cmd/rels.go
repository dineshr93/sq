// Copyright © 2023 NAME HERE <EMAIL ADDRESS>
// SPDX-FileCopyrightText: 2023 Dinesh Ravi
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dineshr93/sq/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// relsCmd represents the rels command
var relsCmd = &cobra.Command{
	Use:   "rels",
	Short: "Lists Relationships",
	Long: `Lists Relationships
	
For Ex: 
To list all relationship: ./sq rels   
To list first 5         : ./sq rels 5 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		dataFile := string(viper.ConfigFileUsed())
		s := &model.SPDX{}
		if err := s.Load(dataFile); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		rels := s.Relationships
		lenRels := len(rels)

		if len(args) > 0 {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			s.PrintRelsClarified(i)

		} else {
			if lenRels > 0 {
				s.PrintRelsClarified(lenRels)
				// s.PrintRelsinSPDX(lenRels)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(relsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// relsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// relsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
