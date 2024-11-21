// Copyright © 2023 Dinesh Ravi dineshr93@gmail.com
// SPDX-FileCopyrightText: 2023 Dinesh Ravi
//
// SPDX-License-Identifier: Apache-2.0

package model

import (
    "encoding/csv"
    "html/template"
    "os"
    "github.com/alexeyco/simpletable"
)

// ExportToCSV exports SPDX data to CSV format
func (s *SPDX) ExportToCSV(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write headers
    headers := []string{"Name", "Spdxid", "Version", "License", "Copyright"}
    if err := writer.Write(headers); err != nil {
        return err
    }

    // Write package data
    for _, pkg := range s.Packages {
        row := []string{
            pkg.Name,
            pkg.Spdxid,
            pkg.VersionInfo,
            pkg.LicenseConcluded,
            pkg.CopyrightText,
        }
        if err := writer.Write(row); err != nil {
            return err
        }
    }
    return nil
}

// ExportToHTML exports SPDX data to HTML format
func (s *SPDX) ExportToHTML(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Create preview table
    table := simpletable.New()
    table.Header = &simpletable.Header{
        Cells: []*simpletable.Cell{
            {Text: "Name"},
            {Text: "Spdxid"},
            {Text: "Version"},
            {Text: "License"},
            {Text: "Copyright"},
        },
    }

    // Initialize table body
    table.Body = &simpletable.Body{
        Cells: make([][]*simpletable.Cell, 0),
    }

    for _, pkg := range s.Packages {
        row := []*simpletable.Cell{
            {Text: pkg.Name},
            {Text: pkg.Spdxid},
            {Text: pkg.VersionInfo},
            {Text: pkg.LicenseConcluded},
            {Text: pkg.CopyrightText},
        }
        table.Body.Cells = append(table.Body.Cells, row)
    }

    const htmlTemplate = `
    <!DOCTYPE html>
    <html>
    <head>
        <title>SPDX Export</title>
        <style>
            table { border-collapse: collapse; width: 100%; }
            th, td { border: 1px solid black; padding: 8px; text-align: left; }
            th { background-color: #f2f2f2; }
        </style>
    </head>
    <body>
        <h1>SPDX Data Export</h1>
        <table>
            <tr>
                <th>Name</th>
                <th>Spdxid</th>
                <th>Version</th>
                <th>License</th>
                <th>Copyright</th>
            </tr>
            {{range .Packages}}
            <tr>
                <td>{{.Name}}</td>
                <td>{{.Spdxid}}</td>
                <td>{{.VersionInfo}}</td>
                <td>{{.LicenseConcluded}}</td>
                <td>{{.CopyrightText}}</td>
            </tr>
            {{end}}
        </table>
    </body>
    </html>`

    tmpl, err := template.New("export").Parse(htmlTemplate)
    if err != nil {
        return err
    }

    return tmpl.Execute(file, s)
}