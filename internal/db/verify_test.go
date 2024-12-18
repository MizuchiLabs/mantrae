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
		IsAdmin  bool
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
				IsAdmin:  true,
			},
			wantErr: false,
		},
		{
			name: "Test CreateUserParams.Verify with empty username",
			fields: fields{
				Username: "",
				Password: "test",
				Email:    "test@test.com",
				IsAdmin:  true,
			},
			wantErr: true,
		},
		{
			name: "Test CreateUserParams.Verify with empty password",
			fields: fields{
				Username: "test",
				Password: "",
				Email:    "test@test.com",
				IsAdmin:  true,
			},
			wantErr: true,
		},
		{
			name: "Test CreateUserParams.Verify with empty email",
			fields: fields{
				Username: "test",
				Password: "test",
				Email:    "",
				IsAdmin:  true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := &User{
				Username: tt.fields.Username,
				Password: tt.fields.Password,
				Email:    &tt.fields.Email,
				IsAdmin:  tt.fields.IsAdmin,
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
		IsAdmin  bool
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
				IsAdmin:  true,
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
				IsAdmin:  true,
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
				IsAdmin:  true,
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
				IsAdmin:  true,
				ID:       1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := &User{
				Username: tt.fields.Username,
				Password: tt.fields.Password,
				Email:    &tt.fields.Email,
				IsAdmin:  tt.fields.IsAdmin,
				ID:       tt.fields.ID,
			}
			if err := u.Verify(); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserParams.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
