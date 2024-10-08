package db

import "testing"

func TestCreateProfileParams_Verify(t *testing.T) {
	type fields struct {
		Name     string
		Url      string
		Username string
		Password string
		Tls      bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test CreateProfileParams.Verify",
			fields: fields{
				Name:     "Test",
				Url:      "https://test.com",
				Username: "test",
				Password: "test",
				Tls:      true,
			},
			wantErr: false,
		},
		{
			name: "Test CreateProfileParams.Verify with empty name",
			fields: fields{
				Name:     "",
				Url:      "https://test.com",
				Username: "test",
				Password: "test",
				Tls:      true,
			},
			wantErr: true,
		},
		{
			name: "Test CreateProfileParams.Verify with empty url",
			fields: fields{
				Name:     "Test",
				Url:      "",
				Username: "test",
				Password: "test",
				Tls:      true,
			},
			wantErr: true,
		},
		{
			name: "Test CreateProfileParams.Verify with empty username",
			fields: fields{
				Name:     "Test",
				Url:      "https://test.com",
				Username: "",
				Password: "test",
				Tls:      true,
			},
			wantErr: false,
		},
		{
			name: "Test CreateProfileParams.Verify with empty password",
			fields: fields{
				Name:     "Test",
				Url:      "https://test.com",
				Username: "test",
				Password: "",
				Tls:      true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &CreateProfileParams{
				Name:     tt.fields.Name,
				Url:      tt.fields.Url,
				Username: &tt.fields.Username,
				Password: &tt.fields.Password,
				Tls:      tt.fields.Tls,
			}
			if err := p.Verify(); (err != nil) != tt.wantErr {
				t.Errorf("CreateProfileParams.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateProfileParams_Verify(t *testing.T) {
	type fields struct {
		Name     string
		Url      string
		Username string
		Password string
		Tls      bool
		ID       int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test UpdateProfileParams.Verify",
			fields: fields{
				Name:     "Test",
				Url:      "https://test.com",
				Username: "test",
				Password: "test",
				Tls:      true,
				ID:       1,
			},
			wantErr: false,
		},
		{
			name: "Test UpdateProfileParams.Verify with empty name",
			fields: fields{
				Name:     "",
				Url:      "https://test.com",
				Username: "test",
				Password: "test",
				Tls:      true,
				ID:       1,
			},
			wantErr: true,
		},
		{
			name: "Test UpdateProfileParams.Verify with empty url",
			fields: fields{
				Name:     "Test",
				Url:      "",
				Username: "test",
				Password: "test",
				Tls:      true,
				ID:       1,
			},
			wantErr: true,
		},
		{
			name: "Test UpdateProfileParams.Verify with empty username",
			fields: fields{
				Name:     "Test",
				Url:      "https://test.com",
				Username: "",
				Password: "test",
				Tls:      true,
				ID:       1,
			},
			wantErr: false,
		},
		{
			name: "Test UpdateProfileParams.Verify with empty password",
			fields: fields{
				Name:     "Test",
				Url:      "https://test.com",
				Username: "test",
				Password: "",
				Tls:      true,
				ID:       1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &UpdateProfileParams{
				Name:     tt.fields.Name,
				Url:      tt.fields.Url,
				Username: &tt.fields.Username,
				Password: &tt.fields.Password,
				Tls:      tt.fields.Tls,
				ID:       tt.fields.ID,
			}
			if err := p.Verify(); (err != nil) != tt.wantErr {
				t.Errorf("UpdateProfileParams.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateUserParams_Verify(t *testing.T) {
	type fields struct {
		Username string
		Password string
		Email    string
		Type     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test CreateUserParams.Verify",
			fields: fields{
				Username: "test",
				Password: "test",
				Email:    "test@test.com",
				Type:     "test",
			},
			wantErr: false,
		},
		{
			name: "Test CreateUserParams.Verify with empty username",
			fields: fields{
				Username: "",
				Password: "test",
				Email:    "test@test.com",
				Type:     "test",
			},
			wantErr: true,
		},
		{
			name: "Test CreateUserParams.Verify with empty password",
			fields: fields{
				Username: "test",
				Password: "",
				Email:    "test@test.com",
				Type:     "test",
			},
			wantErr: true,
		},
		{
			name: "Test CreateUserParams.Verify with empty email",
			fields: fields{
				Username: "test",
				Password: "test",
				Email:    "",
				Type:     "test",
			},
			wantErr: false,
		},
		{
			name: "Test CreateUserParams.Verify with empty type",
			fields: fields{
				Username: "test",
				Password: "test",
				Email:    "test@test.com",
				Type:     "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := &CreateUserParams{
				Username: tt.fields.Username,
				Password: tt.fields.Password,
				Email:    &tt.fields.Email,
				Type:     tt.fields.Type,
			}
			if err := u.Verify(); (err != nil) != tt.wantErr {
				t.Errorf("CreateUserParams.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateUserParams_Verify(t *testing.T) {
	type fields struct {
		Username string
		Password string
		Email    string
		Type     string
		ID       int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test UpdateUserParams.Verify",
			fields: fields{
				Username: "test",
				Password: "test",
				Email:    "test@test.com",
				Type:     "test",
				ID:       1,
			},
			wantErr: false,
		},
		{
			name: "Test UpdateUserParams.Verify with empty username",
			fields: fields{
				Username: "",
				Password: "test",
				Email:    "test@test.com",
				Type:     "test",
				ID:       1,
			},
			wantErr: true,
		},
		{
			name: "Test UpdateUserParams.Verify with empty password",
			fields: fields{
				Username: "test",
				Password: "",
				Email:    "test@test.com",
				Type:     "test",
				ID:       1,
			},
			wantErr: true,
		},
		{
			name: "Test UpdateUserParams.Verify with empty email",
			fields: fields{
				Username: "test",
				Password: "test",
				Email:    "",
				Type:     "test",
				ID:       1,
			},
			wantErr: false,
		},
		{
			name: "Test UpdateUserParams.Verify with empty type",
			fields: fields{
				Username: "test",
				Password: "test",
				Email:    "test@test.com",
				Type:     "",
				ID:       1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := &UpdateUserParams{
				Username: tt.fields.Username,
				Password: tt.fields.Password,
				Email:    &tt.fields.Email,
				Type:     tt.fields.Type,
				ID:       tt.fields.ID,
			}
			if err := u.Verify(); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserParams.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateProviderParams_Verify(t *testing.T) {
	type fields struct {
		Name       string
		Type       string
		ExternalIp string
		ApiKey     string
		ApiUrl     string
		IsActive   bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test CreateProviderParams.Verify",
			fields: fields{
				Name:       "Test",
				Type:       "test",
				ExternalIp: "test",
				ApiKey:     "test",
				ApiUrl:     "test",
				IsActive:   true,
			},
			wantErr: false,
		},
		{
			name: "Test CreateProviderParams.Verify with empty name",
			fields: fields{
				Name:       "",
				Type:       "test",
				ExternalIp: "test",
				ApiKey:     "test",
				ApiUrl:     "test",
				IsActive:   true,
			},
			wantErr: true,
		},
		{
			name: "Test CreateProviderParams.Verify with empty type",
			fields: fields{
				Name:       "Test",
				Type:       "",
				ExternalIp: "test",
				ApiKey:     "test",
				ApiUrl:     "test",
				IsActive:   true,
			},
			wantErr: true,
		},
		{
			name: "Test CreateProviderParams.Verify with empty external ip",
			fields: fields{
				Name:       "Test",
				Type:       "test",
				ExternalIp: "",
				ApiKey:     "test",
				ApiUrl:     "test",
				IsActive:   true,
			},
			wantErr: true,
		},
		{
			name: "Test CreateProviderParams.Verify with empty api key",
			fields: fields{
				Name:       "Test",
				Type:       "test",
				ExternalIp: "test",
				ApiKey:     "",
				ApiUrl:     "test",
				IsActive:   true,
			},
			wantErr: true,
		},
		{
			name: "Test CreateProviderParams.Verify with empty api url",
			fields: fields{
				Name:       "Test",
				Type:       "test",
				ExternalIp: "test",
				ApiKey:     "test",
				ApiUrl:     "",
				IsActive:   true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &CreateProviderParams{
				Name:       tt.fields.Name,
				Type:       tt.fields.Type,
				ExternalIp: tt.fields.ExternalIp,
				ApiKey:     tt.fields.ApiKey,
				ApiUrl:     &tt.fields.ApiUrl,
				IsActive:   tt.fields.IsActive,
			}
			if err := p.Verify(); (err != nil) != tt.wantErr {
				t.Errorf("CreateProviderParams.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateProviderParams_Verify(t *testing.T) {
	type fields struct {
		Name       string
		Type       string
		ExternalIp string
		ApiKey     string
		ApiUrl     string
		IsActive   bool
		ID         int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test UpdateProviderParams.Verify",
			fields: fields{
				Name:       "Test",
				Type:       "test",
				ExternalIp: "test",
				ApiKey:     "test",
				ApiUrl:     "test",
				IsActive:   true,
				ID:         1,
			},
			wantErr: false,
		},
		{
			name: "Test UpdateProviderParams.Verify with empty name",
			fields: fields{
				Name:       "",
				Type:       "test",
				ExternalIp: "test",
				ApiKey:     "test",
				ApiUrl:     "test",
				IsActive:   true,
				ID:         1,
			},
			wantErr: true,
		},
		{
			name: "Test UpdateProviderParams.Verify with empty type",
			fields: fields{
				Name:       "Test",
				Type:       "",
				ExternalIp: "test",
				ApiKey:     "test",
				ApiUrl:     "test",
				IsActive:   true,
				ID:         1,
			},
			wantErr: true,
		},
		{
			name: "Test UpdateProviderParams.Verify with empty external ip",
			fields: fields{
				Name:       "Test",
				Type:       "test",
				ExternalIp: "",
				ApiKey:     "test",
				ApiUrl:     "test",
				IsActive:   true,
				ID:         1,
			},
			wantErr: true,
		},
		{
			name: "Test UpdateProviderParams.Verify with empty api key",
			fields: fields{
				Name:       "Test",
				Type:       "test",
				ExternalIp: "test",
				ApiKey:     "",
				ApiUrl:     "test",
				IsActive:   true,
				ID:         1,
			},
			wantErr: true,
		},
		{
			name: "Test UpdateProviderParams.Verify with empty api url",
			fields: fields{
				Name:       "Test",
				Type:       "test",
				ExternalIp: "test",
				ApiKey:     "test",
				ApiUrl:     "",
				IsActive:   true,
				ID:         1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &UpdateProviderParams{
				Name:       tt.fields.Name,
				Type:       tt.fields.Type,
				ExternalIp: tt.fields.ExternalIp,
				ApiKey:     tt.fields.ApiKey,
				ApiUrl:     &tt.fields.ApiUrl,
				IsActive:   tt.fields.IsActive,
				ID:         tt.fields.ID,
			}
			if err := p.Verify(); (err != nil) != tt.wantErr {
				t.Errorf("UpdateProviderParams.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
