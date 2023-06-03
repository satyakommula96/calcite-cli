package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/apache/calcite-avatica-go/v5"
	"github.com/chzyer/readline"
	"github.com/olekukonko/tablewriter"
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

	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "calcite \U0001F48E:sql>  ",
		HistoryFile:            "/tmp/calcite-cli-history",
		InterruptPrompt:        "^C",
		EOFPrompt:              "exit",
		HistorySearchFold:      true,
		DisableAutoSaveHistory: false,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	var cmds []string
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		// Check for exit command
		if strings.ToLower(line) == "exit" || strings.ToLower(line) == "quit" {
			break
		}
		cmds = append(cmds, line)
		if !strings.HasSuffix(line, ";") {
			rl.SetPrompt("...\t>")
			continue
		}
		cmd := strings.Join(cmds, " ")
		cmds = cmds[:0]
		rl.SetPrompt("calcite \U0001F48E:sql>  ")

		executeQuery(db, table, strings.TrimRight(cmd, ";"))
		rl.SaveHistory(cmd)
	}

	fmt.Println("Exiting calcite CLI Prompt...")
}

func executeQuery(db *sql.DB, table *tablewriter.Table, query string) {
	// Execute the query
	start := time.Now()
	rows, err := db.Query(query)
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
