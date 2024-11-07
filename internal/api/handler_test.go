// Package api provides handlers for the API
package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/MizuchiLabs/mantrae/test"
)

func setupHandler(method, path string, body string) (*httptest.ResponseRecorder, *http.Request) {
	if err := test.SetupDB(); err != nil {
		return nil, nil
	}
	token, err := util.EncodeUserJWT("test")
	if err != nil {
		return nil, nil
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/profile", bytes.NewBuffer([]byte(body)))
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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
		{
			name: "Test writeJSON",
			args: args{
				w:    w,
				data: map[string]string{"test": "test"},
			},
		},
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
		{
			name:       "Login",
			body:       `{"username": "test", "password": "test"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Login Fail",
			body:       `{"username": "test", "password": "wrong"}`,
			wantStatus: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("POST", "/login", tt.body)
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
		{
			name:       "Valid Token",
			body:       `{"token": "valid_token"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Empty Token",
			body:       `{"token": ""}`,
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("POST", "/api/verify", tt.body)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			VerifyToken(w, r)

			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf("VerifyToken() status = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
	}{
		{
			name:       "Get Version - Successful",
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("GET", "/api/version", "")

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetVersion(w, r)

			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf("GetVersion() status = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestGetEvents(t *testing.T) {
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
			GetEvents(tt.args.w, tt.args.r)
		})
	}
}

func TestGetAgents(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		wantStatus int
	}{
		{
			name:       "Get Agents - Successful",
			endpoint:   "/api/agents/1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Get Agents - Missing Fields",
			endpoint:   "/api/agents",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("GET", tt.endpoint, "")

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetAgents(w, r)
		})
	}
}

func TestUpsertAgent(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "Upsert Agent - Successful",
			body:       `{"profileId": 1, "hostname": "localhost"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Upsert Agent - Missing Hostname",
			body:       `{"profileId": 1}`,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("PUT", "/api/agent/testagent", tt.body)

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			UpsertAgent(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf("UpsertAgent() status = %v, want %v", got, tt.wantStatus)
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
		{
			name:       "Delete Agent - Invalid Agent",
			endpoint:   "/api/agent/test/wrong",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("DELETE", tt.endpoint, "")
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			DeleteAgent(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf("DeleteAgent() status = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestGetAgentToken(t *testing.T) {
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
			GetAgentToken(tt.args.w, tt.args.r)
		})
	}
}

func TestGetProfiles(t *testing.T) {
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
			GetProfiles(tt.args.w, tt.args.r)
		})
	}
}

func TestGetProfile(t *testing.T) {
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
			GetProfile(tt.args.w, tt.args.r)
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
			name:       "Create Profile - Successful",
			body:       `{"name": "testprofile", "url": "http://test.com"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Create Profile - Missing Fields",
			body:       `{"name": "Incomplete"}`,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		w, r := setupHandler("POST", "/api/profile", tt.body)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			CreateProfile(w, r)
			if got := w.Result().StatusCode; got != tt.wantStatus {
				t.Errorf("CreateProfile() status = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
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
			UpdateProfile(tt.args.w, tt.args.r)
		})
	}
}

func TestDeleteProfile(t *testing.T) {
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
			DeleteProfile(tt.args.w, tt.args.r)
		})
	}
}

func TestGetUsers(t *testing.T) {
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
			GetUsers(tt.args.w, tt.args.r)
		})
	}
}

func TestGetUser(t *testing.T) {
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
			GetUser(tt.args.w, tt.args.r)
		})
	}
}

func TestCreateUser(t *testing.T) {
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
			CreateUser(tt.args.w, tt.args.r)
		})
	}
}

func TestUpdateUser(t *testing.T) {
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
			UpdateUser(tt.args.w, tt.args.r)
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

func TestGetProvider(t *testing.T) {
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
			GetProvider(tt.args.w, tt.args.r)
		})
	}
}

func TestCreateProvider(t *testing.T) {
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
			CreateProvider(tt.args.w, tt.args.r)
		})
	}
}

func TestUpdateProvider(t *testing.T) {
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
			UpdateProvider(tt.args.w, tt.args.r)
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
