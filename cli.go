/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/apache/calcite-avatica-go/v5"
	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	keywords "github.com/satyakommula96/calcite-cli/keywords"
	"github.com/spf13/cobra"
)

var (
	connectionURL          = "http://localhost:8080"
	serialization          = "protobuf"
	enablePartitionPruning = true
	distributedExecution   = false
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "calcite cli",
		Short: "A calcite CLI prompt to execute queries",
		Run:   runSQLPrompt,
	}

	// Define flags for connection URL and additional parameters
	rootCmd.Flags().StringVar(&connectionURL, "url", connectionURL, "Connection URL")
	rootCmd.Flags().StringVar(&serialization, "serialization", serialization, "Serialization parameter")
	rootCmd.Flags().BoolVar(&enablePartitionPruning, "enablePartitionPruning", enablePartitionPruning, "Enable Partition Pruning")
	rootCmd.Flags().BoolVar(&distributedExecution, "distributedExecution", distributedExecution, "Distributed Execution")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func runSQLPrompt(cmd *cobra.Command, args []string) {
	// Establish a connection to the calcite server
	db, err := sql.Open("avatica", buildConnectionURL())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new table writer
	table := tablewriter.NewWriter(os.Stdout)

	fmt.Println("Enter your queries. Type 'exit' to quit.")
	fmt.Println()

	p := prompt.New(
		executeQueryWrapper(db, table),
		keywords.CustomCompleter,
		prompt.OptionLivePrefix(LivePrefix),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSuggestionBGColor(prompt.White),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionCompletionOnDown(),
		prompt.OptionTitle("Calcite CLI Prompt"),               // Set a title for the prompt
		prompt.OptionInputTextColor(prompt.Cyan),               // Customize input text color
		prompt.OptionDescriptionTextColor(prompt.White),        // Customize description text color
		prompt.OptionSelectedSuggestionTextColor(prompt.Green), // Customize selected suggestion text color
		prompt.OptionSelectedSuggestionBGColor(prompt.Black),   // Customize selected suggestion background color
		prompt.OptionPrefix("calcite \U0001F48E:sql> "),        // Set a custom prefix for the prompt
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(buf *prompt.Buffer) {
				fmt.Println("Exiting calcite CLI Prompt...")
				os.Exit(0)
			},
		}), // Add a custom key binding (e.g., Ctrl+C to exit)
	)

	p.Run()

}

var isMultiline bool

func LivePrefix() (prefix string, useLivePrefix bool) {
	if isMultiline {
		prefix = "... "
		useLivePrefix = true
	} else {
		prefix = "calcite \U0001F48E:sql> "
		useLivePrefix = !isMultiline
	}
	return prefix, useLivePrefix
}

func executeQueryWrapper(db *sql.DB, table *tablewriter.Table) func(string) {
	var multiLineQuery strings.Builder

	return func(query string) {
		// Check for exit command
		if strings.ToLower(query) == "exit" || strings.ToLower(query) == "quit" {
			fmt.Println("Exiting calcite CLI Prompt...")
			os.Exit(0)
		}

		trimmedQuery := strings.TrimSpace(query)

		// Check if it is a multiline query
		if strings.HasSuffix(trimmedQuery, ";") {
			multiLineQuery.WriteString(trimmedQuery)
			executeQuery(db, table, multiLineQuery.String())
			multiLineQuery.Reset()
			isMultiline = false
		} else {
			if !isMultiline {
				multiLineQuery.Reset()
				isMultiline = true
			}
			multiLineQuery.WriteString(trimmedQuery)
			multiLineQuery.WriteString(" ")
		}
	}
}

func executeQuery(db *sql.DB, table *tablewriter.Table, query string) {
	// Execute the query
	start := time.Now()
	cmd := strings.TrimRight(query, ";")
	rows, err := db.Query(cmd)
	if err != nil {
		log.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Println("Error retrieving column names:", err)
		return
	}

	// Create a slice to store the query results
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Clear the table and set new header
	table.ClearRows()
	table.SetHeader(columns)

	// Fetch and print rows
	count := 0
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println("Error retrieving row data:", err)
			continue
		}

		// Prepare row data
		rowData := make([]string, len(columns))
		for i, v := range values {
			if v != nil {
				rowData[i] = fmt.Sprintf("%v", v)
			} else {
				rowData[i] = "NULL"
			}
		}

		// Add row to the table
		table.Append(rowData)
		count++
	}

	duration := time.Since(start)

	// Render the table
	table.Render()

	fmt.Printf("Rows: %d\nExecution Time: %s\n\n", count, duration)
}

func buildConnectionURL() string {
	var params []string

	// Add serialization parameter
	if serialization != "" {
		params = append(params, "serialization="+serialization)
	}

	// Add enablePartitionPruning parameter
	if enablePartitionPruning {
		params = append(params, "enablePartitionPruning=true")
	}

	// Add distributedExecution parameter
	if !distributedExecution {
		params = append(params, "distributedExecution=false")
	}

	// Combine the connection URL and parameters
	url := connectionURL
	if len(params) > 0 {
		url += "?" + strings.Join(params, "&")
	}

	return url
}
