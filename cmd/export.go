// Copyright © 2023 Dinesh Ravi dineshr93@gmail.com
// SPDX-FileCopyrightText: 2023 Dinesh Ravi
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "github.com/dineshr93/sq/model"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var exportCmd = &cobra.Command{
    Use:   "export",
    Short: "Export SPDX data to different formats",
    Long: `Export SPDX data in different formats
    
For Ex: 
To export as CSV    : ./sq export -f csv -o output.csv
To export as HTML   : ./sq export -f html -o output.html`,
    Run: func(cmd *cobra.Command, args []string) {
        dataFile := string(viper.ConfigFileUsed())
        s := &model.SPDX{}
        if err := s.Load(dataFile); err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            os.Exit(1)
        }

        format, _ := cmd.Flags().GetString("format")
        output, _ := cmd.Flags().GetString("output")

        // Create Reports directory if it doesn't exist
        reportsDir := "Reports"
        if err := os.MkdirAll(reportsDir, 0755); err != nil {
            fmt.Fprintln(os.Stderr, "Error creating Reports directory:", err)
            os.Exit(1)
        }

        // Construct full path with Reports directory
        output = filepath.Join(reportsDir, output)

        if !isValidOutputFile(output, format) {
            fmt.Fprintf(os.Stderr, "Error: Output file extension doesn't match format %s\n", format)
            os.Exit(1)
        }

        switch format {
        case "csv":
            if err := s.ExportToCSV(output); err != nil {
                fmt.Fprintln(os.Stderr, err.Error())
                os.Exit(1)
            }
            fmt.Printf("Successfully exported to CSV: %s\n", output)
        case "html":
            if err := s.ExportToHTML(output); err != nil {
                fmt.Fprintln(os.Stderr, err.Error())
                os.Exit(1)
            }
            fmt.Printf("Successfully exported to HTML: %s\n", output)
        default:
            fmt.Fprintf(os.Stderr, "unsupported format: %s\n", format)
            os.Exit(1)
        }
    },
}

func init() {
    rootCmd.AddCommand(exportCmd)
    exportCmd.Flags().StringP("format", "f", "csv", "Export format (csv/html)")
    exportCmd.Flags().StringP("output", "o", "export.csv", "Output file path")
}

func isValidOutputFile(outputFile, format string) bool {
    ext := filepath.Ext(strings.TrimSpace(outputFile))
    switch format {
    case "csv":
        return ext == ".csv"
    case "html":
        return ext == ".html"
    default:
        return false
    }
}