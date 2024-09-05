package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_statusRecorder_WriteHeader(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected int
	}{
		{
			name:     "Status code 200",
			code:     http.StatusOK,
			expected: http.StatusOK,
		},
		{
			name:     "Status code 404",
			code:     http.StatusNotFound,
			expected: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			wrappedRec := &statusRecorder{
				ResponseWriter: rec,
			}
			wrappedRec.WriteHeader(tt.code)
			if wrappedRec.statusCode != tt.expected {
				t.Errorf(
					"statusRecorder.WriteHeader() = %v, want %v",
					wrappedRec.statusCode,
					tt.expected,
				)
			}
		})
	}
}

func Test_statusRecorder_Flush(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Flush call",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			wrappedRec := &statusRecorder{
				ResponseWriter: rec,
			}
			wrappedRec.Flush()
			// No direct assertion here, we check if Flush does not panic
		})
	}
}

func TestLog(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name string
	}{
		{
			name: "Log middleware",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			Log(handler).ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("Log() = %v, want %v", rec.Code, http.StatusOK)
			}
		})
	}
}

func TestCors(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name string
	}{
		{
			name: "CORS middleware",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodOptions, "/", nil)
			Cors(handler).ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("Cors() = %v, want %v", rec.Code, http.StatusOK)
			}
			if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
				t.Errorf(
					"Cors() Access-Control-Allow-Origin = %v, want %v",
					rec.Header().Get("Access-Control-Allow-Origin"),
					"*",
				)
			}
		})
	}
}

func TestBasicAuth(t *testing.T) {
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := BasicAuth(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BasicAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJWT(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Generate a token for testing
	token, _ := GenerateJWT("testuser")

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "Valid JWT",
			token:      token,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid JWT",
			token:      "invalidToken",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)
			rec := httptest.NewRecorder()
			JWT(handler).ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("JWT() = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestChain(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name        string
		middlewares []func(http.Handler) http.Handler
		wantStatus  int
	}{
		{
			name: "Single middleware",
			middlewares: []func(http.Handler) http.Handler{
				Log,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Multiple middlewares",
			middlewares: []func(http.Handler) http.Handler{
				Log,
				Cors,
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			Chain(tt.middlewares...)(handler).ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("Chain() = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}
