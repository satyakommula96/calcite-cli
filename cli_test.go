package main

import (
	"testing"
)

func TestBuildConnectionURL(t *testing.T) {
	tests := []struct {
		name             string
		connectionURL    string
		serialization    string
		schema           string
		connectionParams string
		user             string
		passwd           string
		maxRowsTotal     string
		want             string
	}{
		{
			name:          "default connection",
			connectionURL: "http://localhost:8080",
			want:          "http://localhost:8080",
		},
		{
			name:          "with schema",
			connectionURL: "http://localhost:8080",
			schema:        "myschema",
			want:          "http://localhost:8080/myschema",
		},
		{
			name:          "with auth",
			connectionURL: "http://localhost:8080",
			user:          "user1",
			passwd:        "pass1",
			want:          "http://localhost:8080?avaticaUser=user1&avaticaPassword=pass1",
		},
		{
			name:          "with all params",
			connectionURL: "http://localhost:8080",
			serialization: "protobuf",
			schema:        "myschema",
			user:          "user1",
			passwd:        "pass1",
			maxRowsTotal:  "1000",
			want:          "http://localhost:8080/myschema?serialization=protobuf&avaticaUser=user1&avaticaPassword=pass1&maxRowsTotal=1000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global variables for test
			connectionURL = tt.connectionURL
			serialization = tt.serialization
			schema = tt.schema
			connectionParams = tt.connectionParams
			user = tt.user
			passwd = tt.passwd
			maxRowsTotal = tt.maxRowsTotal

			got := buildConnectionURL()
			if got != tt.want {
				t.Errorf("buildConnectionURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
