// Copyright © 2023 Dinesh Ravi dineshr93@gmail.com
// SPDX-FileCopyrightText: 2023 Dinesh Ravi
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sq",
	Short: "A SBOM Query Module",
	Long: `A SBOM Query CLI (for issue -> https://github.com/dineshr93/sq/issues)
	
	1. List Meta data (sq meta)
	2. List Files (sq files)
	3. List Packages (sq pkgs)
	4. List Relationships (sq rels)
	5. List pkgs and files in Relationships (sq rels dig)`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func findFiles(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}
func isStringInFile(anyFile, text string) bool {
	b, err := ioutil.ReadFile(anyFile)
	if err != nil {
		panic(err)
	}
	s := string(b)
	// //check whether s contains substring text
	return strings.Contains(s, text)
}
func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "YOUR SBOM JSON FILE (default is $CURRENT_DIR/sbom.spdx.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func isvalidSPDXJSONFile(spdxjsonFile string) bool {
	// fmt.Println(filepath.Ext(strings.TrimSpace(spdxjsonFile)))
	if filepath.Ext(strings.TrimSpace(spdxjsonFile)) != ".json" {
		fmt.Println("Error:", spdxjsonFile, " Not a JSON file")
		return false
	}
	// check if valid json if SPDXID Keyword is present in file
	if !isStringInFile(spdxjsonFile, "SPDXID") {
		fmt.Println("Error: Not a Valid SPDX JSON file")
		return false
	}
	return true
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Check if valid json and spdx file and fail if not
		if !isvalidSPDXJSONFile(cfgFile) {
			os.Exit(1)
		}

		viper.SetConfigFile(cfgFile)

	} else {
		// Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(pwd)
		var isSPDXJSONFilefound bool = false
		fmt.Println("--config option not passed So finding first valid spdx json file in current folder")

		for _, s := range findFiles(pwd, ".json") {
			// viper.SetConfigType("json")
			// viper.SetConfigName("sbom.spdx")
			if !isSPDXJSONFilefound && isvalidSPDXJSONFile(s) {
				viper.SetConfigFile(s)
				isSPDXJSONFilefound = true
			}
		}

		if !isSPDXJSONFilefound {
			fmt.Println("No valid spdx json format file found in current directory!")
			os.Exit(1)
		}

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using SBOM file =======>", viper.ConfigFileUsed())
	}
}
