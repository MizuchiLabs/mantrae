// Package config provides functions for parsing command-line flags and
// setting up the application's default settings.
package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

func TestFlags_Parse(t *testing.T) {
	type fields struct {
		Version  bool
		Port     int
		URL      string
		Username string
		Password string
		Update   bool
		Reset    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Passing all flags",
			fields: fields{
				Version:  false,
				Port:     3000,
				URL:      "http://localhost:8080",
				Username: "admin",
				Password: "password",
				Update:   false,
				Reset:    false,
			},
			wantErr: false,
		},
		{
			name: "Passing only URL",
			fields: fields{
				Version:  false,
				Port:     3000,
				URL:      "http://localhost:8080",
				Username: "",
				Password: "",
				Update:   false,
				Reset:    false,
			},
			wantErr: false,
		},
		{
			name: "Passing only username",
			fields: fields{
				Version:  false,
				Port:     3000,
				URL:      "",
				Username: "admin",
				Password: "",
				Update:   false,
				Reset:    false,
			},
			wantErr: false,
		},
		{
			name: "Passing only password",
			fields: fields{
				Version:  false,
				Port:     3000,
				URL:      "",
				Username: "",
				Password: "password",
				Update:   false,
				Reset:    false,
			},
			wantErr: false,
		},
		{
			name: "Passing reset flag",
			fields: fields{
				Version:  false,
				Port:     3000,
				URL:      "",
				Username: "",
				Password: "",
				Update:   false,
				Reset:    true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := db.InitDB(); err != nil {
				t.Errorf("InitDB() error = %v", err)
			}

			os.Args = []string{
				"-port",
				strconv.Itoa(tt.fields.Port),
				"-url",
				tt.fields.URL,
				"-username",
				tt.fields.Username,
				"-password",
				tt.fields.Password,
			}
			if tt.fields.Update {
				os.Args = append(os.Args, "-update")
			}
			if tt.fields.Reset {
				os.Args = append(os.Args, "-reset")
			}
			var f Flags
			if err := f.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("Flags.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetDefaultAdminUser(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Pass", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := db.InitDB(); err != nil {
				t.Errorf("InitDB() error = %v", err)
			}
			if err := SetDefaultAdminUser(); (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultAdminUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetDefaultProfile(t *testing.T) {
	type args struct {
		url      string
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Passing all arguments",
			args: args{
				url:      "http://localhost:8080",
				username: "admin",
				password: "password",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := db.InitDB(); err != nil {
				t.Errorf("InitDB() error = %v", err)
			}
			if err := SetDefaultProfile(tt.args.url, tt.args.username, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetDefaultSettings(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Pass", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := db.InitDB(); err != nil {
				t.Errorf("InitDB() error = %v", err)
			}
			if err := SetDefaultSettings(); (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResetAdminUser(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Pass",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := db.InitDB(); err != nil {
				t.Errorf("InitDB() error = %v", err)
			}

			if err := ResetAdminUser(); (err != nil) != tt.wantErr {
				t.Errorf("ResetAdminUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
