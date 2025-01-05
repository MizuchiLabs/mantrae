// Package api provides handlers for the API
package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/MizuchiLabs/mantrae/test"
)

func setupHandler(method, path string, body string) (*httptest.ResponseRecorder, *http.Request) {
	if err := test.SetupDB(); err != nil {
		return nil, nil
	}
	token, err := util.EncodeUserJWT("test", time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, nil
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)
	if body != "" {
		r = httptest.NewRequest(method, path, bytes.NewBuffer([]byte(body)))
	}
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	Server().ServeHTTP(w, r)

	return w, r
}

func Test_writeJSON(t *testing.T) {
	w := httptest.NewRecorder()

	type args struct {
		w    http.ResponseWriter
		data any
	}
	tests := []struct {
		name string
		args args
	}{
		{"writeJSON - Success", args{w: w, data: map[string]string{"test": "test"}}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			writeJSON(tt.args.w, tt.args.data)
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"Login - Successful", `{"username": "test", "password": "test"}`, http.StatusOK},
		{"Login - Empty Credentials", `{"username": "", "password": ""}`, http.StatusBadRequest},
		{
			"Login - Invalid Credentials",
			`{"username": "test", "password": "wrong"}`,
			http.StatusUnauthorized,
		},
		{
			"Login - Invalid Username",
			`{"username": "wrong", "password": "wrong"}`,
			http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("POST", "/api/login", tt.body)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			Login(w, r)
			if status := w.Code; status != tt.wantStatus {
				t.Errorf(
					"Login() = %v, want %v, error = %v",
					status,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"Valid Token", `{"token": "valid_token"}`, http.StatusOK},
		{"Empty Token", `{"token": ""}`, http.StatusOK},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("POST", "/api/verify", tt.body)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			VerifyToken(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"VerifyToken() status = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
	}{
		{"Get Version - Successful", http.StatusOK},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("GET", "/api/version", "")
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetVersion(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"GetVersion() status = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestGetEvents(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("GET", "/api/events", "")
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetEvents(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"GetEvents() status = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestGetAgents(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		wantStatus int
	}{
		{"Get Agents - Successful", "/api/agent/1", http.StatusOK},
		{"Get Agents - Missing Fields", "/api/agent", http.StatusNotFound},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("GET", tt.endpoint, "")
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetAgents(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"GetAgents() status = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestUpsertAgent(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"Upsert Agent - Successful", `{"profileId": 1, "hostname": "localhost"}`, http.StatusOK},
		{"Upsert Agent - Missing Hostname", `{"profileId": 1}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("PUT", "/api/agent", tt.body)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			UpsertAgent(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"UpsertAgent() status = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestDeleteAgent(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		wantStatus int
	}{
		{"Delete Agent - Successful", "/api/agent/test/hard", http.StatusOK},
		{"Delete Agent - Invalid Agent", "/api/agent/wrong/wrong", http.StatusNotFound},
		{"Delete Agent - Invalid Type", "/api/agent/test/wrong", http.StatusBadRequest},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("DELETE", tt.endpoint, "")
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			DeleteAgent(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"DeleteAgent() status = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestGetProfiles(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
	}{
		{"Get Profiles - Successful", http.StatusOK},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("GET", "/api/profile", "")
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetProfiles(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"GetProfiles() = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		wantStatus int
	}{
		{"Get Profile - Existing ID", "/api/profile/1", http.StatusOK},
		{"Get Profile - Nonexistent ID", "/api/profile/999", http.StatusNotFound},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("GET", tt.endpoint, "")
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetProfile(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"GetProfile() = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestCreateProfile(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			"Create Profile - Successful",
			`{"name": "testprofile", "url": "http://test.com"}`,
			http.StatusOK,
		},
		{
			"Create Profile - Missing Fields",
			`{"name": "Incomplete"}`,
			http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("POST", "/api/profile", tt.body)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			CreateProfile(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"CreateProfile() status = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			"Update Profile - Successful",
			`{"id": 1, "name": "updated", "url": "http://test.com"}`,
			http.StatusOK,
		},
		{"Update Profile - Missing URL", `{"id": 1, "name": "updated"}`, http.StatusBadRequest},
		{"Update Profile - Missing ID", `{"name": "updated"}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("PUT", "/api/profile", tt.body)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			UpdateProfile(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"UpdateProfile() = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestDeleteProfile(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		wantStatus int
	}{
		{"Delete Profile - Existing ID", "/api/profile/1", http.StatusOK},
		{"Delete Profile - Nonexistent ID", "/api/profile/999", http.StatusNotFound},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("DELETE", tt.endpoint, "")
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			DeleteProfile(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf("DeleteProfile() = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
	}{
		{"Get Users - Success", http.StatusOK},
		{"Get Users - No Users", http.StatusOK},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("GET", "/api/user", "")
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetUsers(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"GetUsers() status = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		wantStatus int
	}{
		//{"Get User - Success", "1", http.StatusOK},
		{"Get User - Not Found", "9999", http.StatusNotFound},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("GET", fmt.Sprintf("/api/user/%s", tt.userID), "")
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetUser(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf(
					"GetUser() status = %v, want %v, error = %v",
					got,
					tt.wantStatus,
					w.Body.String(),
				)
			}
		})
	}
}

func TestUpsertUser(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			UpsertUser(tt.args.w, tt.args.r)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			DeleteUser(tt.args.w, tt.args.r)
		})
	}
}

func TestGetProviders(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetProviders(tt.args.w, tt.args.r)
		})
	}
}

func TestDeleteProvider(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			DeleteProvider(tt.args.w, tt.args.r)
		})
	}
}

func TestGetEntryPoints(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetEntryPoints(tt.args.w, tt.args.r)
		})
	}
}

func TestGetRouters(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetRouters(tt.args.w, tt.args.r)
		})
	}
}

func TestUpsertRouter(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			UpsertRouter(tt.args.w, tt.args.r)
		})
	}
}

func TestDeleteRouter(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			DeleteRouter(tt.args.w, tt.args.r)
		})
	}
}

func TestGetServices(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetServices(tt.args.w, tt.args.r)
		})
	}
}

func TestUpsertService(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			UpsertService(tt.args.w, tt.args.r)
		})
	}
}

func TestDeleteService(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			DeleteService(tt.args.w, tt.args.r)
		})
	}
}

func TestGetMiddlewares(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetMiddlewares(tt.args.w, tt.args.r)
		})
	}
}

func TestUpsertMiddleware(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			UpsertMiddleware(tt.args.w, tt.args.r)
		})
	}
}

func TestDeleteMiddleware(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			DeleteMiddleware(tt.args.w, tt.args.r)
		})
	}
}

func TestGetSettings(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetSettings(tt.args.w, tt.args.r)
		})
	}
}

func TestGetSetting(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetSetting(tt.args.w, tt.args.r)
		})
	}
}

func TestUpdateSetting(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			UpdateSetting(tt.args.w, tt.args.r)
		})
	}
}

func TestDownloadBackup(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			DownloadBackup(tt.args.w, tt.args.r)
		})
	}
}

func TestUploadBackup(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			UploadBackup(tt.args.w, tt.args.r)
		})
	}
}

func TestDeleteRouterDNS(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			DeleteRouterDNS(tt.args.w, tt.args.r)
		})
	}
}

func TestGetMiddlewarePlugins(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetMiddlewarePlugins(tt.args.w, tt.args.r)
		})
	}
}

func TestGetTraefikOverview(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetTraefikOverview(tt.args.w, tt.args.r)
		})
	}
}

func TestGetTraefikConfig(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetTraefikConfig(tt.args.w, tt.args.r)
		})
	}
}

func TestGetPublicIP(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetPublicIP(tt.args.w, tt.args.r)
		})
	}
}
