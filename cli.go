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
	"net/url"
	"strings"

	avatica "github.com/apache/calcite-avatica-go/v5"
	prompt "github.com/satyakommula96/calcite-cli/prompt"
	"github.com/spf13/cobra"
)

type ConnectionConfig struct {
	ConnectionURL    string
	Serialization    string
	Schema           string
	ConnectionParams string
	User             string
	Passwd           string
	MaxRowsTotal     string
	CustomParams     string
}

var config = ConnectionConfig{
	ConnectionURL: "http://localhost:8080",
	Serialization: "protobuf",
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "calcite cli",
		Short: "A calcite CLI prompt to execute queries",
		Run:   runSQLPrompt,
	}

	// Define flags for connection URL and additional parameters
	rootCmd.Flags().StringVar(&config.ConnectionURL, "url", config.ConnectionURL, "Connection URL")
	rootCmd.Flags().StringVar(&config.Serialization, "serialization", "", "Serialization parameter")
	rootCmd.Flags().StringVar(&config.ConnectionParams, "params", "", "Extra parameters for avatica connection (ex: \"parameter1=value&...parameterN=value\")")
	rootCmd.Flags().StringVarP(&config.Schema, "schema", "s", "", "The schema path sets the default schema to use for this connection.")
	rootCmd.Flags().StringVarP(&config.User, "username", "u", "", "The user to use when authenticating against Avatica")
	rootCmd.Flags().StringVarP(&config.Passwd, "password", "p", "", "The password to use when authenticating against Avatica")
	rootCmd.MarkFlagsRequiredTogether("username", "password")
	rootCmd.Flags().StringVarP(&config.MaxRowsTotal, "maxRowsTotal", "m", "", "The maxRowsTotal parameter sets the maximum number of rows to return for a given query")
	rootCmd.Flags().StringVar(&config.CustomParams, "extra_params", "", "Custom connection parameters for avatica connection (ex: \"parameter1=value;...parameterN=value\")")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func runSQLPrompt(cmd *cobra.Command, args []string) {
	// Establish a connection to the calcite server
	db := establishConnection(config)
	defer db.Close()

	// Create and run the SQL prompt
	prompt.CreateAndRunPrompt(db)
}

func establishConnection(cfg ConnectionConfig) *sql.DB {
	dsn := buildConnectionURL(cfg)
	fmt.Println("Connecting to ", dsn)

	// Prepare the info map
	info := make(map[string]string)
	if cfg.CustomParams != "" {
		pairs := strings.Split(cfg.CustomParams, ";")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				info[kv[0]] = kv[1]
			}
		}
	}

	// Create a new connector
	connector := avatica.NewConnector(dsn).(*avatica.Connector)

	// Set the info map in the connector
	connector.Info = info

	// Open the database using the connector
	db := sql.OpenDB(connector)
	fmt.Println("Connected")
	return db
}
func buildConnectionURL(cfg ConnectionConfig) string {
	u, err := url.Parse(cfg.ConnectionURL)
	if err != nil {
		log.Fatalf("Invalid connection URL: %v", err)
	}

	if cfg.Schema != "" {
		if !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
		u.Path += cfg.Schema
	}

	q := u.Query()

	// Add serialization parameter by default protobuf
	if cfg.Serialization != "" {
		q.Set("serialization", cfg.Serialization)
	}

	// Add username and password as parameter
	if cfg.User != "" {
		q.Set("avaticaUser", cfg.User)
		q.Set("avaticaPassword", cfg.Passwd)
	}

	if cfg.MaxRowsTotal != "" {
		q.Set("maxRowsTotal", cfg.MaxRowsTotal)
	}

	// Add connection parameters
	if cfg.ConnectionParams != "" {
		extraQ, err := url.ParseQuery(strings.ReplaceAll(cfg.ConnectionParams, ";", "&"))
		if err == nil {
			for k, v := range extraQ {
				for _, val := range v {
					q.Add(k, val)
				}
			}
		} else {
			log.Printf("Failed to parse connection params: %v", err)
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}
