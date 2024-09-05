package api

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	os.Setenv("SECRET", "dummy-secret") // Set the secret environment variable

	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid Username",
			args: args{
				username: "testuser",
			},
			wantErr: false,
		},
		{
			name: "Empty Username",
			args: args{
				username: "",
			},
			wantErr: false, // A token can still be generated for an empty username
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateJWT(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	os.Setenv("SECRET", "dummy-secret") // Set the secret environment variable

	validToken, _ := GenerateJWT("testuser")

	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		want    *Claims
		wantErr bool
	}{
		{
			name: "Valid Token",
			args: args{
				tokenString: validToken,
			},
			want: &Claims{
				Username: "testuser",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Token",
			args: args{
				tokenString: "invalidTokenString",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Expired Token",
			args: args{
				tokenString: func() string {
					claims := &Claims{
						Username: "expireduser",
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
						},
					}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
					return tokenString
				}(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateJWT(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}