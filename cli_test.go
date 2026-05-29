package main

import (
	"testing"
)

func TestBuildConnectionURL(t *testing.T) {
	tests := []struct {
		name string
		cfg  ConnectionConfig
		want string
	}{
		{
			name: "default connection",
			cfg: ConnectionConfig{
				ConnectionURL: "http://localhost:8080",
			},
			want: "http://localhost:8080",
		},
		{
			name: "with schema",
			cfg: ConnectionConfig{
				ConnectionURL: "http://localhost:8080",
				Schema:        "myschema",
			},
			want: "http://localhost:8080/myschema",
		},
		{
			name: "with auth",
			cfg: ConnectionConfig{
				ConnectionURL: "http://localhost:8080",
				User:          "user1",
				Passwd:        "pass1",
			},
			want: "http://localhost:8080?avaticaPassword=pass1&avaticaUser=user1",
		},
		{
			name: "with all params",
			cfg: ConnectionConfig{
				ConnectionURL: "http://localhost:8080",
				Serialization: "protobuf",
				Schema:        "myschema",
				User:          "user1",
				Passwd:        "pass1",
				MaxRowsTotal:  "1000",
			},
			want: "http://localhost:8080/myschema?avaticaPassword=pass1&avaticaUser=user1&maxRowsTotal=1000&serialization=protobuf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildConnectionURL(tt.cfg)
			if got != tt.want {
				t.Errorf("buildConnectionURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
