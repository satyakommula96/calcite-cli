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
}

func CreateAndRunPrompt(db *sql.DB) {
	fmt.Println("Welcome! Use SQL to query Apache Calcite.\nUse Ctrl+D, type \"exit\" or \"quit\" to exit.")
	fmt.Println()

	session := &PromptSession{db: db}

	p := prompt.New(
		session.executor,
		CustomCompleter,
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
