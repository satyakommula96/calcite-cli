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

package prompt

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/apache/calcite-avatica-go/v5"
	"github.com/c-bata/go-prompt"
	calcitesql "github.com/satyakommula96/calcite-cli/calcitesql"
)

type PromptSession struct {
	db             *sql.DB
	isMultiline    bool
	multiLineQuery strings.Builder
	suggestions    []prompt.Suggest
}

func CreateAndRunPrompt(db *sql.DB) {
	fmt.Println("Welcome! Use SQL to query Apache Calcite.\nUse Ctrl+D, type \"exit\" or \"quit\" to exit.")
	fmt.Println()

	session := &PromptSession{db: db}

	// Initialize with static SQL suggestions
	session.suggestions = append(session.suggestions, sqlSuggestions...)

	// Fetch database-specific tables and columns
	metaSugg := fetchMetadataSuggestions(db)
	session.suggestions = append(session.suggestions, metaSugg...)

	p := prompt.New(
		session.executor,
		session.completer,
		prompt.OptionLivePrefix(session.LivePrefix),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSuggestionBGColor(prompt.White),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionCompletionOnDown(),
		prompt.OptionTitle("Calcite CLI Prompt"),                 // Set a title for the prompt
		prompt.OptionInputTextColor(prompt.Fuchsia),              // Customize input text color
		prompt.OptionDescriptionTextColor(prompt.Black),          // Customize description text color
		prompt.OptionSelectedSuggestionTextColor(prompt.White),   // Customize selected suggestion text color
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray), // Customize selected suggestion background color
		prompt.OptionPrefix("calcite \U0001F48E:sql> "),          // Set a custom prefix for the prompt
	)

	p.Run()
}

func (s *PromptSession) LivePrefix() (prefix string, useLivePrefix bool) {
	if s.isMultiline {
		prefix = "... "
		useLivePrefix = true
	} else {
		prefix = "calcite \U0001F48E:sql> "
		useLivePrefix = !s.isMultiline
	}
	return prefix, useLivePrefix
}

func (s *PromptSession) executor(query string) {
	// Check for exit command
	if strings.ToLower(query) == "exit" || strings.ToLower(query) == "quit" {
		fmt.Println("Exiting calcite CLI Prompt...")
		os.Exit(0)
	}

	trimmedQuery := strings.TrimSpace(query)

	// Check if it is a multiline query
	if strings.HasSuffix(trimmedQuery, ";") {
		s.multiLineQuery.WriteString(trimmedQuery)
		calcitesql.ExecuteQuery(s.db, s.multiLineQuery.String())
		s.multiLineQuery.Reset()
		s.isMultiline = false
	} else {
		if !s.isMultiline {
			s.multiLineQuery.Reset()
			s.isMultiline = true
		}
		s.multiLineQuery.WriteString(trimmedQuery)
		s.multiLineQuery.WriteString(" ")
	}
}

func (s *PromptSession) completer(d prompt.Document) []prompt.Suggest {
	input := d.GetWordBeforeCursor()
	if input == "" {
		return nil
	}
	return prompt.FilterHasPrefix(s.suggestions, input, true)
}

func fetchMetadataSuggestions(db *sql.DB) []prompt.Suggest {
	var suggestions []prompt.Suggest

	// Fetch tables
	tableRows, err := db.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES")
	if err == nil {
		defer tableRows.Close()
		for tableRows.Next() {
			var tableName string
			if err := tableRows.Scan(&tableName); err == nil {
				suggestions = append(suggestions, prompt.Suggest{
					Text:        tableName,
					Description: "Table Name",
				})
			}
		}
	}

	// Fetch columns
	columnRows, err := db.Query("SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS")
	if err == nil {
		defer columnRows.Close()
		for columnRows.Next() {
			var columnName string
			if err := columnRows.Scan(&columnName); err == nil {
				suggestions = append(suggestions, prompt.Suggest{
					Text:        columnName,
					Description: "Column Name",
				})
			}
		}
	}

	return suggestions
}
