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

// metaCmd represents the meta command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Meta data of the spdx file",
	Long:  `Meta data of the spdx file`,
	Run: func(cmd *cobra.Command, args []string) {
		dataFile := string(viper.ConfigFileUsed())
		s := &model.SPDX{}
		if err := s.Load(dataFile); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		s.PrintMeta()
		// fmt.Println("Spdx ID:", s.Spdxid)
		// fmt.Println("Spdx version:", s.SpdxVersion)
		// fmt.Println("Spdx creation date:", s.CreationInfo.Created.Format(time.RFC822))
		// if len(s.CreationInfo.Creators) > 0 {
		// 	fmt.Println("created by:", strings.Join(s.CreationInfo.Creators, ", "))
		// }
		// fmt.Println("Project version:", s.Name)
		// fmt.Println("File License(not projects):", s.DataLicense)
		// fmt.Println("Document Namespace:", s.DocumentNamespace)
		// fmt.Println("Document Describes:", strings.Join(s.DocumentDescribes, ", "))
		// fmt.Println("Number of Packages:", len(s.Packages))
		// fmt.Println("Number of Files:", len(s.Files))
		// fmt.Println("Number of Packages:", len(s.Relationships))

	},
}

func init() {
	rootCmd.AddCommand(metaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
